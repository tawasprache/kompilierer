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

//go:embed testdata/Wahr.tawa
var wahr string

func TestTypisierungWahr(t *testing.T) {
	var v parser.Modul
	feh := parser.Parser.ParseString("Wahr.tawa", wahr, &v)
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
	if len(fehler) != 0 {
		t.Fatalf("fehler: %+v", fehler)
	}
}

//go:embed testdata/Strukt.tawa
var strukt string

func TestTypisierungStrukt(t *testing.T) {
	var v parser.Modul
	feh := parser.Parser.ParseString("Strukt.tawa", strukt, &v)
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
	if len(fehler) != 0 {
		t.Fatalf("fehler: %+v", fehler)
	}
}

//go:embed testdata/SchlechtStrukt.tawa
var schlechtStrukt string

func TestTypisierungSchlechtStrukt(t *testing.T) {
	var v parser.Modul
	feh := parser.Parser.ParseString("SchlechtStrukt.tawa", schlechtStrukt, &v)
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
