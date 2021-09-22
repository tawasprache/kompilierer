package typisierung_test

import (
	"Tawa/parser"
	"Tawa/typen"
	"Tawa/typisierung"
	_ "embed"
	"testing"
)

//go:embed data/nicht-gefunden.tawa
var nichtGefunden string

//go:embed data/gut.tawa
var gut string

//go:embed data/anruf.tawa
var anruf string

//go:embed data/mismatch.tawa
var mismatch string

//go:embed data/nichts.tawa
var nichts string

//go:embed data/self.tawa
var self string

//go:embed data/generisch.tawa
var generisch string

//go:embed data/schlechte-generisch.tawa
var schlechteGenerisch string

//go:embed data/gute-generisch.tawa
var guteGenerisch string

//go:embed data/var-nicht-gefunden.tawa
var varNichtGefunden string

//go:embed data/variabeln.tawa
var variabeln string

func TestNichtGefunden(t *testing.T) {
	a := parser.VonStringX("data/nicht-gefunden.tawa", nichtGefunden)
	err := typisierung.PrüfDatei(&a)
	if err == nil {
		t.Fatalf("expected an error, got none")
	}
}

func TestGut(t *testing.T) {
	a := parser.VonStringX("data/gut.tawa", gut)
	err := typisierung.PrüfDatei(&a)
	if err != nil {
		t.Fatalf("expected no errors, got one: %+s", err)
	}
}

func TestAnruf(t *testing.T) {
	a := parser.VonStringX("data/anruf.tawa", anruf)
	err := typisierung.PrüfDatei(&a)
	if err != nil {
		t.Fatalf("expected no errors, got one: %+s", err)
	}
}

func TestMismatch(t *testing.T) {
	a := parser.VonStringX("data/mismatch.tawa", mismatch)
	err := typisierung.PrüfDatei(&a)
	if err == nil {
		t.Fatalf("expected an error, didn't get one")
	}
}

func TestNichts(t *testing.T) {
	a := parser.VonStringX("data/nichts.tawa", nichts)
	err := typisierung.PrüfDatei(&a)
	if err != nil {
		t.Fatalf("expected no errors, got one: %+s", err)
	}
}

func TestSelf(t *testing.T) {
	a := parser.VonStringX("data/self.tawa", self)
	err := typisierung.PrüfDatei(&a)
	if err != nil {
		t.Fatalf("expected no errors, got one: %+s", err)
	}
}

func TestGenerisch(t *testing.T) {
	a := parser.VonStringX("data/generisch.tawa", generisch)
	err := typisierung.PrüfDatei(&a)
	if err != nil {
		t.Fatalf("expected no errors, got one: %+s", err)
	}
}

func TestSchlechtGenerisch(t *testing.T) {
	a := parser.VonStringX("data/schlechte-generisch.tawa", schlechteGenerisch)
	err := typisierung.PrüfDatei(&a)

	if err == nil {
		t.Fatalf("expected an error, didn't get one")
	}

	if err.(*typisierung.ErrMismatch).A != (typen.Logik{}) || err.(*typisierung.ErrMismatch).B != (typen.Integer{}) {
		t.Fatalf("expected A = logik and B = ganz")
	}
}

func TestGuteGenerisch(t *testing.T) {
	a := parser.VonStringX("data/gute-generisch.tawa", guteGenerisch)
	err := typisierung.PrüfDatei(&a)
	if err != nil {
		t.Fatalf("expected no errors, got one: %+s", err)
	}
}

func TestVarNichtGefunden(t *testing.T) {
	a := parser.VonStringX("data/var-nicht-gefunden.tawa", varNichtGefunden)
	err := typisierung.PrüfDatei(&a)
	if err == nil {
		t.Fatalf("expected an error, didn't get one")
	}
}

func TestVariabeln(t *testing.T) {
	a := parser.VonStringX("data/variabeln.tawa", variabeln)
	err := typisierung.PrüfDatei(&a)
	if err != nil {
		t.Fatalf("expected no errors, got one: %s", err)
	}
}
