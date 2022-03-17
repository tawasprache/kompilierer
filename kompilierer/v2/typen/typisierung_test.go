package typen

import (
	"Tawa/kompilierer/v2/ast"
	"Tawa/kompilierer/v2/parser"
	_ "embed"
	"testing"
)

func TestTypisierung(t *testing.T) {
	var v parser.Modul
	feh := parser.Parser.ParseString("Ganzzahl.tawa", ganzzahl, &v)
	if feh != nil {
		t.Fatalf("error: %s", feh)
	}

	datei := ast.VonParser(v)
	k := NeuKontext()
	fehler := Namenaufslösung(datei, k)
	if len(fehler) > 0 {
		t.Fatalf("fehler: %+v", fehler)
	}
	fehler = Typisierung(datei, k)
	if len(fehler) > 0 {
		t.Fatalf("fehler: %+v", fehler)
	}
}

//go:embed testdata/SchlechtZwei.tawa
var schlechtZwei string

func TestTypisierungZwei(t *testing.T) {
	var v parser.Modul
	feh := parser.Parser.ParseString("SchlechtZwei.tawa", schlechtZwei, &v)
	if feh != nil {
		t.Fatalf("error: %s", feh)
	}

	datei := ast.VonParser(v)
	k := NeuKontext()
	fehler := Namenaufslösung(datei, k)
	if len(fehler) > 0 {
		t.Fatalf("fehler: %+v", fehler)
	}
	fehler = Typisierung(datei, k)
	if len(fehler) != 2 {
		t.Fatalf("fehler: %+v", fehler)
	}
}
