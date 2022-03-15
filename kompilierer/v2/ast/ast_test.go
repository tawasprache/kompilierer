package ast

import (
	"Tawa/kompilierer/v2/parser"
	_ "embed"
	"testing"
)

//go:embed testdata/Hallo.tawa
var hallo string

func TestConvert(t *testing.T) {
	var v parser.Modul
	feh := parser.Parser.ParseString("Hallo.tawa", hallo, &v)
	if feh != nil {
		t.Fatalf("error: %s", feh)
	}

	VonParser(v)
}
