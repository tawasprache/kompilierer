package typen

import (
	"Tawa/kompilierer/v2/ast"
	"Tawa/kompilierer/v2/parser"
	_ "embed"
	"testing"
)

//go:embed testdata/Ganzzahl.tawa
var ganzzahl string

func TestNamensaufloesung(t *testing.T) {
	var v parser.Modul
	feh := parser.Parser.ParseString("Ganzzahl.tawa", ganzzahl, &v)
	if feh != nil {
		t.Fatalf("error: %s", feh)
	}

	datei := ast.VonParser(v)
	k := NeuKontext()
	fehler := Namenaufslösung(datei, k)
	if len(fehler) > 0 {
		t.Fail()
	}
}
