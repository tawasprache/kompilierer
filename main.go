// +build !wasm

package main

import (
	"Tawa/codegenerierung"
	"Tawa/interpreter"
	"Tawa/parser"
	"Tawa/typisierung"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	llvm := flag.Bool("als-llvm", false, "compilieren nach llvm")
	interp := flag.Bool("als-interpreter", false, "interpreter")
	parse := flag.Bool("parser", false, "parser drücken")

	flag.Parse()

	if *parse {
		println(parser.Parser.String())
		return
	}
	datei, feh := ioutil.ReadFile(flag.Arg(0))
	if feh != nil {
		panic(feh)
	}
	es := parser.Datei{}
	feh = parser.Parser.ParseBytes(flag.Arg(0), datei, &es)
	if feh != nil {
		panic(feh)
	}
	es.Vorverarbeiten()
	v := typisierung.NeuVollKontext()
	v.Push()
	err := typisierung.Typisierung(v, &es)
	if err != nil {
		if v, ok := err.(*typisierung.Fehler); ok {
			s := string(datei)
			fmt.Printf("%s:%d:\n", flag.Arg(0), v.Line)
			println(strings.Split(s, "\n")[v.Line-1])
		}
		println(err.Error())
		os.Exit(1)
	}

	if *llvm {
		println(codegenerierung.Codegen(&es))
		return
	} else if *interp {
		vk := interpreter.NeuVollKontext()
		vk.Push()

		ok, err := interpreter.Interpret(es, "main", vk)
		if err != nil {
			println(err.Error())
			os.Exit(1)
		}

		println(ok.AlsString())
		return
	}

	flag.PrintDefaults()
}
