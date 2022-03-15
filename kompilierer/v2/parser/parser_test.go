package parser

import (
	_ "embed"
	"testing"
)

//go:embed testdata/Hallo.tawa
var hallo string

func TestHallo(t *testing.T) {
	var v Modul
	feh := Parser.ParseString("Hallo.tawa", hallo, &v)
	if feh != nil {
		t.Fatalf("error: %s", feh)
	}
}

//go:embed testdata/Funktionen.tawa
var funktionen string

func TestFunktionen(t *testing.T) {
	var v Modul
	feh := Parser.ParseString("Funktionen.tawa", funktionen, &v)
	if feh != nil {
		t.Fatalf("error: %s", feh)
	}
}
