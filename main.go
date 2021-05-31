package main

import (
	"Tawa/parser"
	"io/ioutil"
	"os"

	"github.com/alecthomas/repr"
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
	repr.Println(es)
}
