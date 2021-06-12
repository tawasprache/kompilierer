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
		return types.I8
	case typen.Neutyp:
		return llvmTypVonTypen(w.Von)
	default:
		panic("a")
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

func codegenExpression(c *ctx, e *parser.Expression, b **ir.Block) value.Value {
	if e == nil {
		return nil
	}

	if e.Block != nil {
		var last value.Value

		c.pushScope()
		for _, statement := range e.Block.Expr {
			last = codegenExpression(c, &statement, b)
		}
		c.popScope()

		return last
	} else if e.Definierung != nil {
		val := codegenExpression(c, &e.Definierung.Wert, b)

		alloca := (*b).NewAlloca(val.Type())
		(*b).NewStore(val, alloca)

		c.top()[e.Definierung.Variable] = LLVMMutableValue{Value: alloca}

		return val
	} else if e.Variable != nil {
		switch v := c.lookup(*e.Variable).(type) {
		case LLVMValue:
			return v.Value
		case LLVMMutableValue:
			return (*b).NewLoad(v.Value.Type().(*types.PointerType).ElemType, v.Value)
		default:
			panic("aa")
		}
	} else if e.Bedingung != nil {
		c.i++

		condVal := codegenExpression(c, &e.Bedingung.Wenn, b)

		fn := (*b).Parent

		thenBloc := fn.NewBlock("then" + strconv.Itoa(c.i))
		thenValue := codegenExpression(c, &e.Bedingung.Werden, &thenBloc)

		condCmp := (*b).NewICmp(enum.IPredNE, condVal, constant.False)
		elseBloc := fn.NewBlock("else" + strconv.Itoa(c.i))

		if e.Bedingung.Sonst == nil {
			thenBloc.NewBr(elseBloc)
			(*b).NewCondBr(condCmp, thenBloc, elseBloc)
			*b = elseBloc
			return nil
		}

		elseValue := codegenExpression(c, e.Bedingung.Sonst, &elseBloc)

		mergeBloc := fn.NewBlock("ifcont" + strconv.Itoa(c.i))
		phi := mergeBloc.NewPhi(ir.NewIncoming(thenValue, thenBloc), ir.NewIncoming(elseValue, elseBloc))

		// time to add the conditional now that we built the blocks
		(*b).NewCondBr(condCmp, thenBloc, elseBloc)

		// now we chain the branches to the merge block
		thenBloc.NewBr(mergeBloc)
		elseBloc.NewBr(mergeBloc)

		*b = mergeBloc

		return phi
	} else if e.Integer != nil {
		return constant.NewInt(types.I32, e.Integer.Value)
	} else if e.Logik != nil {
		if e.Logik.Wert == "Wahr" {
			return constant.True
		} else {
			return constant.False
		}
	} else if e.Funktionsaufruf != nil {
		fn := c.lookup(e.Funktionsaufruf.Name).(LLVMValue).Value

		args := []value.Value{}
		for _, arg := range e.Funktionsaufruf.Argumente {
			args = append(args, codegenExpression(c, &arg, b))
		}

		return (*b).NewCall(fn, args...)
	} else if e.Cast != nil {
		return codegenExpression(c, &e.Cast.Von, b)
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
		funk.Blocks[len(funk.Blocks)-1].NewRet(ret)
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

	it := codegen(c, d)

	return it.String()
}