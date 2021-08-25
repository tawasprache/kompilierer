package main

import "Tawa/libgccjit"

func main() {
	ctx := libgccjit.NewContext()

	printer := ctx.CreateFunction("printf", libgccjit.Imported, ctx.GetType(libgccjit.TypeVoid), []libgccjit.Param{
		{
			Name: "format",
			Type: ctx.GetType(libgccjit.TypeConstCharPtr),
		},
	}, false)

	fn := ctx.CreateFunction("main", libgccjit.Exported, ctx.GetType(libgccjit.TypeVoid), []libgccjit.Param{}, false)
	hello := ctx.StringRValue("hello world!\n")
	blk := fn.NewBlock()
	call := ctx.NewCall(printer, []libgccjit.IRValue{hello})
	blk.AddEval(call)
	blk.EndWithReturnVoid()

	ctx.CompileToFile(libgccjit.Executable, "test")
}
