package codegenerierung

import (
	"Tawa/parser"
	"Tawa/typen"
	"strconv"

	"github.com/alecthomas/repr"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type vollwert struct {
	links  value.Value
	rechts value.Value
}

func vrechts(v value.Value) *vollwert {
	return &vollwert{rechts: v}
}

func llvmTypVonTypen(p typen.Art) types.Type {
	switch w := p.(type) {
	case typen.Funktion:
		params := []types.Type{}
		for _, es := range w.Argumente {
			params = append(params, llvmTypVonTypen(es))
		}
		return types.NewFunc(llvmTypVonTypen(w.Returntyp), params...)
	case typen.Primitiv:
		switch w.Name {
		case "ganz":
			return types.I32
		case "g8":
			return types.I8
		case "g16":
			return types.I16
		case "g32":
			return types.I32
		case "g64":
			return types.I64
		case "vzlosganz":
			return types.I32
		case "vzlosg8":
			return types.I8
		case "vzlosg16":
			return types.I16
		case "vzlosg32":
			return types.I32
		case "vzlosg64":
			return types.I64
		default:
			panic("a")
		}
	case typen.Nichts:
		return types.Void
	case typen.Logik:
		return types.I1
	case typen.Neutyp:
		return llvmTypVonTypen(w.Von)
	case typen.Struktur:
		a := types.NewStruct()
		for _, it := range w.Fields {
			a.Fields = append(a.Fields, llvmTypVonTypen(it.Typ))
		}
		return a
	case typen.Zeiger:
		return types.NewPointer(llvmTypVonTypen(w.Auf))
	default:
		panic("a " + repr.String(p))
	}
}

func codegenPrefunktionen(c *ctx, d *parser.Datei, module *ir.Module) {
	for _, funk := range d.Funktionen {
		typ := funk.Art.(typen.Funktion)

		var params []*ir.Param
		for idx, es := range funk.Funktionsargumente {
			params = append(params, ir.NewParam(es.Name, llvmTypVonTypen(typ.Argumente[idx])))
		}

		fn := module.NewFunc(funk.Name, llvmTypVonTypen(typ.Returntyp), params...)

		c.top()[funk.Name] = LLVMValue{Value: fn}
	}
}

func codegenExpression(c *ctx, e *parser.Expression, b **ir.Block) *vollwert {
	if e == nil {
		return nil
	}

	if e.Block != nil {
		var last value.Value

		c.pushScope()
		for _, statement := range e.Block.Expr {
			last = codegenExpression(c, &statement, b).rechts
		}
		c.popScope()

		return vrechts(last)
	} else if e.Definierung != nil {
		val := codegenExpression(c, &e.Definierung.Wert, b)

		alloca := (*b).NewAlloca(val.rechts.Type())
		(*b).NewStore(val.rechts, alloca)

		c.top()[e.Definierung.Variable] = LLVMMutableValue{Value: alloca}

		return val
	} else if e.Variable != nil {
		switch v := c.lookup(*e.Variable).(type) {
		case LLVMValue:
			return vrechts(v.Value)
		case LLVMMutableValue:
			return &vollwert{
				links:  v.Value,
				rechts: (*b).NewLoad(v.Value.Type().(*types.PointerType).ElemType, v.Value),
			}
		default:
			panic("aa")
		}
	} else if e.Bedingung != nil {
		c.i++

		condVal := codegenExpression(c, &e.Bedingung.Wenn, b)

		fn := (*b).Parent

		thenBloc := fn.NewBlock("then" + strconv.Itoa(c.i))
		thenValue := codegenExpression(c, &e.Bedingung.Werden, &thenBloc)

		condCmp := (*b).NewICmp(enum.IPredNE, condVal.rechts, constant.False)
		elseBloc := fn.NewBlock("else" + strconv.Itoa(c.i))

		if e.Bedingung.Sonst == nil {
			thenBloc.NewBr(elseBloc)
			(*b).NewCondBr(condCmp, thenBloc, elseBloc)
			*b = elseBloc
			return nil
		}

		elseValue := codegenExpression(c, e.Bedingung.Sonst, &elseBloc)

		mergeBloc := fn.NewBlock("ifcont" + strconv.Itoa(c.i))
		phi := mergeBloc.NewPhi(ir.NewIncoming(thenValue.rechts, thenBloc), ir.NewIncoming(elseValue.rechts, elseBloc))

		// time to add the conditional now that we built the blocks
		(*b).NewCondBr(condCmp, thenBloc, elseBloc)

		// now we chain the branches to the merge block
		thenBloc.NewBr(mergeBloc)
		elseBloc.NewBr(mergeBloc)

		*b = mergeBloc

		return vrechts(phi)
	} else if e.Integer != nil {
		return vrechts(constant.NewInt(types.I32, e.Integer.Value))
	} else if e.Logik != nil {
		if e.Logik.Wert == "Wahr" {
			return vrechts(constant.True)
		} else {
			return vrechts(constant.False)
		}
	} else if e.Funktionsaufruf != nil {
		fn := c.lookup(e.Funktionsaufruf.Name).(LLVMValue).Value

		args := []value.Value{}
		for _, arg := range e.Funktionsaufruf.Argumente {
			args = append(args, codegenExpression(c, &arg, b).rechts)
		}

		return vrechts((*b).NewCall(fn, args...))
	} else if e.Cast != nil {
		return codegenExpression(c, &e.Cast.Von, b)
	} else if e.Löschen != nil {
		val := codegenExpression(c, &e.Löschen.Expr, b)

		fn := (*b).Parent.Parent.NewFunc("free", types.Void, ir.NewParam("addr", types.NewInt(64)))

		(*b).NewCall(fn, (*b).NewPtrToInt(val.rechts, types.I64))

		return val
	} else if e.Neu != nil {
		wert := codegenExpression(c, e.Neu.Expression, b)

		ptr := types.NewPointer(wert.rechts.Type())

		gep := (*b).NewGetElementPtr(wert.rechts.Type(), constant.NewNull(ptr), constant.NewInt(types.I32, 1))
		grosse := (*b).NewPtrToInt(gep, types.I64)

		fn := (*b).Parent.Parent.NewFunc("malloc", types.NewPointer(types.NewInt(8)), ir.NewParam("size", types.NewInt(64)))

		ret := (*b).NewCall(fn, grosse)

		ptrret := (*b).NewBitCast(ret, types.NewPointer(wert.rechts.Type()))

		(*b).NewStore(wert.rechts, ptrret)

		return vrechts(ptrret)
	} else if e.Stack != nil {
		kind := c.lookup(e.Stack.Initialisierung.Name).(TypenTyp).Art.(typen.Struktur)

		typ := llvmTypVonTypen(kind)
		alloca := (*b).NewAlloca(typ)

		for _, it := range e.Stack.Initialisierung.Fields {
			wert := codegenExpression(c, &it.Wert, b)

			a := -1
			for pos, f := range kind.Fields {
				if f.Name == it.Name {
					a = pos
					break
				}
			}

			ptr := (*b).NewGetElementPtr(typ, alloca, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, int64(a)))

			if v, ok := wert.rechts.(*ir.InstAlloca); ok {
				loaded := (*b).NewLoad(v.ElemType, wert.rechts)
				_ = (*b).NewStore(loaded, wert.rechts)
			} else {
				_ = (*b).NewStore(wert.rechts, ptr)
			}

		}

		return vrechts(alloca)
	} else if e.Dereferenzierung != nil {
		w := codegenExpression(c, &e.Dereferenzierung.Expr, b)

		l := (*b).NewLoad(w.rechts.Type().(*types.PointerType).ElemType, w.rechts)

		return vrechts(l)
	} else if e.Zuweisungsexpression != nil {
		links := codegenExpression(c, &e.Zuweisungsexpression.Links, b)
		rechts := codegenExpression(c, &e.Zuweisungsexpression.Rechts, b)

		es := (*b).NewStore(rechts.rechts, links.links)

		return vrechts(es.Src)
	} else if e.Fieldexpression != nil {
		expr := codegenExpression(c, &e.Fieldexpression.Expr, b)
		idx := e.Fieldexpression.FieldIndex

		t := expr.rechts.(value.Value).Type().(*types.PointerType).ElemType.(*types.StructType)

		ptr := (*b).NewGetElementPtr(expr.rechts.Type().(*types.PointerType).ElemType, expr.rechts, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, int64(idx)))
		a := (*b).NewLoad(t.Fields[idx], ptr)

		if types.IsStruct(a.Type()) {
			return &vollwert{links: ptr, rechts: ptr}
		} else {
			return &vollwert{links: ptr, rechts: a}
		}
	}

	repr.Println(e)
	panic("ee")
}

func codegenFunktion(c *ctx, d *parser.Funktion, module *ir.Module) {
	wert := c.lookup(d.Name).(LLVMValue)

	funk := wert.Value.(*ir.Func)
	block := funk.NewBlock("entry")

	c.pushScope()
	for idx, arg := range d.Funktionsargumente {
		param := funk.Params[idx]

		alloca := block.NewAlloca(param.Typ)
		block.NewStore(param, alloca)

		c.top()[arg.Name] = LLVMMutableValue{Value: alloca}
	}
	ret := codegenExpression(c, &d.Expression, &block)
	c.popScope()

	if types.IsVoid(funk.Sig.RetType) {
		funk.Blocks[len(funk.Blocks)-1].NewRet(nil)
	} else {
		if types.IsStruct(funk.Sig.RetType) {
			es := funk.Blocks[len(funk.Blocks)-1].NewLoad(funk.Sig.RetType, ret.rechts)
			funk.Blocks[len(funk.Blocks)-1].NewRet(es)
		} else {
			funk.Blocks[len(funk.Blocks)-1].NewRet(ret.rechts)
		}
	}
}

func codegen(c *ctx, d *parser.Datei) *ir.Module {
	module := ir.NewModule()

	codegenPrefunktionen(c, d, module)
	for _, fn := range d.Funktionen {
		codegenFunktion(c, &fn, module)
		c.i = 0
	}

	return module
}

func Codegen(d *parser.Datei) string {
	c := &ctx{
		names:           []map[string]namedThing{{}},
		stringConstants: map[string]value.Value{},
	}

	for _, it := range d.Typdeklarationen {
		c.names[0][it.Name] = TypenTyp{Art: it.CodeArt}
	}

	it := codegen(c, d)

	return it.String()
}
