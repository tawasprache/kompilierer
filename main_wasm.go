// +build wasm

package main

import "syscall/js"
import "Tawa/kompilierer/ast"
import "Tawa/kompilierer/codegenierung"
import "Tawa/kompilierer/typisierung"

func main() {
	js.Global().Set("tawaZuJS", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		dat := ast.Modul{}
		err = ast.Parser.Parse(c.Args().First(), fi, &dat)
		if err != nil {
			return err
		}

		ktx := typisierung.NeuKontext()
		genannt, err := typisierung.Auflösenamen(ktx, dat, "User")
		if err != nil {
			return err
		}

		getypt, err := typisierung.Typiere(ktx, genannt, "User")
		if err != nil {
			return err
		}
	}))

	c := make(chan struct{}, 0)
	<-c
}
