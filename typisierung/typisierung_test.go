package typisierung_test

import (
	"Tawa/parser"
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
