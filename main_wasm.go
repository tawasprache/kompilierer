// +build wasm

package main

import "syscall/js"
import "Tawa/interpreter"

func main() {
	js.Global().Set("evalTawa", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return interpreter.Evaluate(args[0].String(), "online.tawa")
	}))

	c := make(chan struct{}, 0)
	<-c
}
