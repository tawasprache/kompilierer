package monomorphisation_test

import (
	"Tawa/monomorphisation"
	"Tawa/parser"
	"Tawa/typisierung"
	_ "embed"
	"testing"
)

//go:embed data/ident-a.tawa
var identA string

func TestIdent(t *testing.T) {
	astA := parser.VonStringX("data/ident-a.tawa", identA)

	err := typisierung.PrüfDatei(&astA)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	procA := monomorphisation.Monomorphise(astA)
	identGanzGanzGefunden := false
	for _, it := range procA.Funktionen {
		if it.Name == "Ident" {
			t.Fatalf("expected generic ident to be deleted")
		} else if it.Name == "Ident_ganz_ganz" {
			identGanzGanzGefunden = true
		}
	}

	if !identGanzGanzGefunden {
		t.Fatalf("expected ident_ganz_ganz to exist")
	}
}
