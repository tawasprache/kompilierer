package codegenerierungaufgccjit

import (
	"Tawa/libgccjit"
	"Tawa/parser"
	"Tawa/typen"

	"github.com/alecthomas/repr"
)

func gccjitTypVonTypen(p typen.Art, ctx *libgccjit.Context) libgccjit.IType {
	switch w := p.(type) {
	case typen.Nichts:
		return ctx.GetType(libgccjit.TypeVoid)
	case typen.Logik:
		return ctx.GetType(libgccjit.TypeBool)
	case typen.Struktur:
		var f []*libgccjit.Field
		for _, it := range w.Fields {
			f = append(f, ctx.GetField(it.Name, gccjitTypVonTypen(it.Typ, ctx)))
		}
		return ctx.GetStructType(w.String(), f)
	case typen.Neutyp:
		return gccjitTypVonTypen(w.Von, ctx)
	case typen.Primitiv:
		switch w.Name {
		case "ganz":
			return ctx.GetType(libgccjit.TypeInt)
		case "g8":
			return ctx.GetIntType(1, true)
		case "g16":
			return ctx.GetIntType(2, true)
		case "g32":
			return ctx.GetIntType(4, true)
		case "g64":
			return ctx.GetIntType(8, true)
		case "vzlosganz":
			return ctx.GetType(libgccjit.TypeUnsignedInt)
		case "vzlosg8":
			return ctx.GetIntType(1, false)
		case "vzlosg16":
			return ctx.GetIntType(2, false)
		case "vzlosg32":
			return ctx.GetIntType(4, false)
		case "vzlosg64":
			return ctx.GetIntType(8, false)
		default:
			panic("a")
		}
	default:
		_ = w
		panic("a " + repr.String(p))
	}
}

func codegenPrefunktionen(c *kontext, d *parser.Datei, ctx *libgccjit.Context) {
	for _, funk := range d.Funktionen {
		typ := funk.Art.(typen.Funktion)

		var params []libgccjit.Param
		for idx, es := range funk.Funktionsargumente {
			params = append(params, libgccjit.Param{gccjitTypVonTypen(typ.Argumente[idx], ctx), es.Name})
		}

		c.funktionen[funk.Name] = ctx.CreateFunction(funk.Name, libgccjit.Exported, gccjitTypVonTypen(typ.Returntyp, ctx), params, false)
	}
}

func codegenExpression(c *kontext, e *parser.Expression, ctx *libgccjit.Context, fn *libgccjit.Function, b **libgccjit.Block) *wert {
	if e == nil {
		return nil
	}

	if e.Block != nil {
		var last *wert

		c.pushScope()
		for _, statement := range e.Block.Expr {
			last = codegenExpression(c, &statement, ctx, fn, b)
		}
		c.popScope()

		return last
	} else if e.Definierung != nil {
		val := codegenExpression(c, &e.Definierung.Wert, ctx, fn, b)

		loc := fn.NewLocal(c.namemiti(e.Definierung.Variable), val.typ)

		c.top()[e.Definierung.Variable] = links(loc, val.typ)

		return val
	} else if e.Variable != nil {
		return c.lookup(*e.Variable)
	} else if e.Bedingung != nil {
		a := e.Bedingung.Art

		w := codegenExpression(c, &e.Bedingung.Wenn, ctx, fn, b)

		if a == nil {
			b1 := fn.NewBlock("bedingung wenn")
			codegenExpression(c, &e.Bedingung.Werden, ctx, fn, &b1)

			b2 := fn.NewBlock("bedingung sonst nichts")

			(*b).EndWithConditional(w.rvalue, b1, b2)
			b1.EndWithJump(b2)

			*b = b2

			return nil
		} else {
			t := gccjitTypVonTypen(a, ctx)

			l := fn.NewLocal(c.name(), t)

			b1 := fn.NewBlock("bedingung wenn wert")
			w1 := codegenExpression(c, &e.Bedingung.Werden, ctx, fn, &b1)
			b1.AddAssign(l, w1.rvalue)

			b2 := fn.NewBlock("bedingung sonst wert")
			w2 := codegenExpression(c, &e.Bedingung.Werden, ctx, fn, &b2)
			b2.AddAssign(l, w2.rvalue)

			b3 := fn.NewBlock("nach bedingung")

			b1.EndWithJump(b3)
			b2.EndWithJump(b3)

			(*b).EndWithConditional(w.rvalue, b1, b2)

			*b = b3

			return links(l, w.typ)
		}
	} else if e.Integer != nil {
		typ := ctx.GetType(libgccjit.TypeInt)
		return nurrechts(ctx.IntRValue(int(e.Integer.Value), typ), typ)
	} else if e.Logik != nil {
		typ := ctx.GetType(libgccjit.TypeBool)
		if e.Logik.Wert == "Wahr" {
			return nurrechts(ctx.IntRValue(1, typ), typ)
		} else {
			return nurrechts(ctx.IntRValue(0, typ), typ)
		}
	} else if e.Funktionsaufruf != nil {
		fn := c.funktionen[e.Funktionsaufruf.Name]
		var w []libgccjit.IRValue
		for _, it := range e.Funktionsaufruf.Argumente {
			w = append(w, codegenExpression(c, &it, ctx, fn, b).rvalue)
		}
		return nurrechts(ctx.NewCall(fn, w), fn.ReturnType())
	} else if e.Cast != nil {
		return codegenExpression(c, &e.Cast.Von, ctx, fn, b)
	} else if e.Löschen != nil {
		expr := codegenExpression(c, &e.Löschen.Expr, ctx, fn, b)

		call := ctx.NewCall(c.löschen, []libgccjit.IRValue{expr.rvalue})

		(*b).AddEval(call)

		return nil
	} else if e.Neu != nil {
		expr := codegenExpression(c, e.Neu.Expression, ctx, fn, b)

		ptrType := expr.typ.AsPointer()

		call := ctx.NewCall(c.neu, []libgccjit.IRValue{ctx.SizeOf(expr.typ)})
		cast := ctx.NewCast(call, ptrType)
		deref := cast.Dereference()

		(*b).AddAssign(deref, expr.rvalue)

		return nurrechts(cast, ptrType)
	} else if e.Stack != nil {
		typ := c.typen[e.Stack.Initialisierung.Name].(*libgccjit.StructType)
		va := fn.NewLocal(c.name(), c.typen[e.Stack.Initialisierung.Name])

		for _, es := range e.Stack.Initialisierung.Fields {
			lval := libgccjit.AccessLValueField(va, typ.FieldsMap()[es.Name])
			rval := codegenExpression(c, &es.Wert, ctx, fn, b)
			(*b).AddAssign(lval, rval.rvalue)
		}

		return links(va, typ)
	} else if e.Dereferenzierung != nil {
		e := codegenExpression(c, &e.Dereferenzierung.Expr, ctx, fn, b)
		es := e.rvalue.RValue().Dereference()
		return links(es, es.Kind())
	} else if e.Zuweisungsexpression != nil {
		lval := codegenExpression(c, &e.Zuweisungsexpression.Links, ctx, fn, b)
		rval := codegenExpression(c, &e.Zuweisungsexpression.Rechts, ctx, fn, b)
		(*b).AddAssign(lval.lvalue, rval.rvalue)

		return rval
	} else if e.Fieldexpression != nil {
		strct := codegenExpression(c, &e.Fieldexpression.Expr, ctx, fn, b)
		field := strct.typ.(*libgccjit.StructType).FieldsMap()[e.Fieldexpression.Field]

		return links(libgccjit.AccessLValueField(strct.lvalue, field), field.Type())
	}

	repr.Println(e)
	panic("ee")
}

func codegenFunktion(c *kontext, d *parser.Funktion, ctx *libgccjit.Context) {
	fn := c.funktionen[d.Name]

	blk := fn.NewBlock("zuerst")

	c.pushScope()
	for idx, arg := range d.Funktionsargumente {
		it := fn.GetParam(idx)

		c.top()[arg.Name] = links(it, it.Kind())
	}
	ret := codegenExpression(c, &d.Expression, ctx, fn, &blk)
	c.popScope()

	retTyp := d.Art.(typen.Funktion).Returntyp
	_, nichts := retTyp.(typen.Nichts)

	if nichts {
		fn.Blocks()[len(fn.Blocks())-1].EndWithReturnVoid()
	} else {
		fn.Blocks()[len(fn.Blocks())-1].EndWithReturnValue(ret.rvalue)
	}

	fn.DumpToDot("test.dot")
}

func codegen(c *kontext, d *parser.Datei) *libgccjit.Context {
	ctx := libgccjit.NewContext()

	c.neu = ctx.CreateFunction("malloc", libgccjit.Imported, ctx.GetType(libgccjit.TypeVoidPtr), []libgccjit.Param{{Type: ctx.GetType(libgccjit.TypeSizeT), Name: "size"}}, false)
	c.löschen = ctx.CreateFunction("free", libgccjit.Imported, ctx.GetType(libgccjit.TypeVoid), []libgccjit.Param{{Type: ctx.GetType(libgccjit.TypeVoidPtr), Name: "ptr"}}, false)

	for _, es := range d.Typdeklarationen {
		c.typen[es.Name] = gccjitTypVonTypen(es.CodeArt, ctx)
	}

	codegenPrefunktionen(c, d, ctx)
	for _, fn := range d.Funktionen {
		codegenFunktion(c, &fn, ctx)
	}

	return ctx
}

func CodegenZuDatei(d *parser.Datei, output string) {
	c := &kontext{
		namen:      []map[string]*wert{{}},
		funktionen: map[string]*libgccjit.Function{},
		typen:      map[string]libgccjit.IType{},
	}

	it := codegen(c, d)

	it.CompileToFile(libgccjit.Executable, output)
}
