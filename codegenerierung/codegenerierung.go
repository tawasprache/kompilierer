package codegenerierung

import (
	"Tawa/parser"
	"Tawa/typen"

	"github.com/alecthomas/repr"
	"github.com/pontaoski/gccjit"
)

func gccjitTypVonTypen(p typen.Typ, ctx *gccjit.Context) gccjit.IType {
	switch w := p.(type) {
	case typen.Logik:
		return ctx.GetType(gccjit.TypeBool)
	case typen.Integer:
		return ctx.GetType(gccjit.TypeInt)
	case typen.Nichts:
		return ctx.GetType(gccjit.TypeVoid)
	default:
		_ = w
		panic("a " + repr.String(p))
	}
}

func codegenPrefunktionen(c *kontext, d *parser.Datei, ctx *gccjit.Context) {
	for _, funk := range d.Funktionen {
		typ := funk.CodeTyp.(typen.Funktion)

		var params []gccjit.Param
		for idx, es := range funk.Funktionsargumente {
			params = append(params, gccjit.Param{
				Type: gccjitTypVonTypen(typ.Eingabe[idx], ctx),
				Name: es.Name,
			})
		}

		c.funktionen[funk.Name] = ctx.CreateFunction(funk.Name, gccjit.Exported, gccjitTypVonTypen(typ.Ausgabe, ctx), params, false)
	}
}

func codegenExpression(c *kontext, e *parser.Expression, ctx *gccjit.Context, fn *gccjit.Function, b **gccjit.Block) *wert {
	if e == nil {
		return nil
	}

	if e.Ganz != nil {
		typ := ctx.GetType(gccjit.TypeInt)
		return nurrechts(ctx.IntRValue(int(*e.Ganz), typ), typ)
	} else if e.Funktionsaufruf != nil {
		fn := c.funktionen[e.Funktionsaufruf.Name]
		var w []gccjit.IRValue
		for _, it := range e.Funktionsaufruf.Argumente {
			w = append(w, codegenExpression(c, &it, ctx, fn, b).rvalue)
		}
		return nurrechts(ctx.NewCall(fn, w), fn.ReturnType())
	} else if e.Block != nil {
		var last *wert

		c.pushScope()
		for _, statement := range e.Block.Expressionen {
			last = codegenExpression(c, statement, ctx, fn, b)
		}
		c.popScope()

		return last
	}

	repr.Println(e)
	panic("ee")
}

func codegenFunktion(c *kontext, d *parser.Funktion, ctx *gccjit.Context) {
	fn := c.funktionen[d.Name]

	blk := fn.NewBlock("zuerst")

	c.pushScope()
	for idx, arg := range d.Funktionsargumente {
		it := fn.GetParam(idx)

		c.top()[arg.Name] = links(it, it.Kind())
	}
	ret := codegenExpression(c, &d.Expression, ctx, fn, &blk)
	c.popScope()

	retTyp := d.CodeTyp.(typen.Funktion).Ausgabe
	_, nichts := retTyp.(typen.Nichts)

	if nichts {
		fn.Blocks()[len(fn.Blocks())-1].EndWithReturnVoid()
	} else {
		fn.Blocks()[len(fn.Blocks())-1].EndWithReturnValue(ret.rvalue)
	}

	fn.DumpToDot("test.dot")
}

func codegen(c *kontext, d *parser.Datei) *gccjit.Context {
	ctx := gccjit.NewContext()

	c.neu = ctx.CreateFunction("malloc", gccjit.Imported, ctx.GetType(gccjit.TypeVoidPtr), []gccjit.Param{{Type: ctx.GetType(gccjit.TypeSizeT), Name: "size"}}, false)
	c.löschen = ctx.CreateFunction("free", gccjit.Imported, ctx.GetType(gccjit.TypeVoid), []gccjit.Param{{Type: ctx.GetType(gccjit.TypeVoidPtr), Name: "ptr"}}, false)

	codegenPrefunktionen(c, d, ctx)
	for _, fn := range d.Funktionen {
		codegenFunktion(c, &fn, ctx)
	}

	return ctx
}

func CodegenZuDatei(d *parser.Datei, output string) {
	c := &kontext{
		namen:      []map[string]*wert{{}},
		funktionen: map[string]*gccjit.Function{},
		typen:      map[string]gccjit.IType{},
	}

	it := codegen(c, d)

	it.CompileToFile(gccjit.Executable, output)
}
