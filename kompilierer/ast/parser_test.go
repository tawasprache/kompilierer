package ast

import "testing"

func TestFunk(t *testing.T) {
	var v Typfunktion
	feh := typFunktionParser.ParseString("mald", "funk()", &v)
	if feh != nil {
		t.Fatalf("error: %s", feh)
	}
}

func TestFunktionsLiteral(t *testing.T) {
	var v Funktionsliteral
	feh := funktionsLiteralParser.ParseString("mald", "\\(x: Logik) => x", &v)
	if feh != nil {
		t.Fatalf("error: %s", feh)
	}

	feh = funktionsLiteralParser.ParseString("mald", "\\(x) => x", &v)
	if feh != nil {
		t.Fatalf("error: %s", feh)
	}
}
