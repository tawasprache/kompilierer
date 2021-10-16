package llvm

import (
	"Tawa/kompilierer/codegenierung"
	"Tawa/kompilierer/getypisiertast"

	"tinygo.org/x/go-llvm"
)

func init() {
	codegenierung.UnterbauRegistrieren("llvm", &llvmUnterbau{})
}

type llvmUnterbau struct {
	c llvm.Context
}

func (l *llvmUnterbau) Pregen(o codegenierung.Optionen) error {
	l.c = llvm.NewContext()

	return nil
}
func (l *llvmUnterbau) CodegenModul(o codegenierung.Optionen, m getypisiertast.Modul) error {
	ctx := l.c
	_ = ctx
	modu := ctx.NewModule(m.Name)

	t := llvm.FunctionType(ctx.VoidType(), []llvm.Type{}, false)

	for _, v := range m.Funktionen {
		fn := llvm.AddFunction(modu, v.SymbolURL.String(), t)
		block := modu.Context().AddBasicBlock(fn, "entry")

		ir := modu.Context().NewBuilder()
		ir.SetInsertPointAtEnd(block)
		defer ir.Dispose()

		ir.CreateRetVoid()
	}

	println("=== verification ===")
	llvm.VerifyModule(modu, llvm.PrintMessageAction)
	println("=== dumping ===")
	modu.Dump()
	println()

	return nil
}
func (l *llvmUnterbau) Postgen(o codegenierung.Optionen) error {
	return nil
}
