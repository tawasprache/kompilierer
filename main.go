package main

import (
	"Tawa/parser"
	"Tawa/typisierung"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	datei, feh := ioutil.ReadFile(os.Args[1])
	if feh != nil {
		panic(feh)
	}
	es := parser.Datei{}
	feh = parser.Parser.ParseBytes(os.Args[1], datei, &es)
	if feh != nil {
		panic(feh)
	}
	v := typisierung.VollKontext{}
	v.Push()
	err := typisierung.VariabelnDefinierung(&v, &es)
	if err != nil {
		if v, ok := err.(*typisierung.Fehler); ok {
			s := string(datei)
			fmt.Printf("%s:%d:\n", os.Args[1], v.Line)
			println(strings.Split(s, "\n")[v.Line-1])
		}
		println(err.Error())
		os.Exit(1)
	}
}
