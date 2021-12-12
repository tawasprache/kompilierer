package ast

import "testing"

func TestFunk(t *testing.T) {
	var v Typfunktion
	feh := typFunktionParser.ParseString("mald", "funk()", &v)
	if feh != nil {
		t.Fatalf("error: %s", feh)
	}
}

func TestParens(t *testing.T) {
	var v Terminal
	feh := TerminalParser.ParseString("mald", "sei x = 1 == (2 == 3) in 1", &v)
	if feh != nil {
		t.Fatalf("error: %s", feh)
	}
}

func TestFirstClassCall(t *testing.T) {
	var v Terminal
	feh := TerminalParser.ParseString("mald", "sei x = a.(2 == 3) in 1", &v)
	if feh != nil {
		t.Fatalf("error: %s", feh)
	}
}

func TestNonassoc(t *testing.T) {
	var v Terminal
	feh := TerminalParser.ParseString("mald", "sei x = a == b == c in 1", &v)
	if feh == nil {
		t.Fatalf("wanted error, got none")
	}
}

func TestDeclarations(t *testing.T) {
	var d Deklaration
	feh := deklParser.ParseString("mald", `/* mald */ funk A() => 1`, &d)
	if feh != nil {
		t.Fatalf("error: %s", feh)
	}
	if d.Comments() != "mald" {
		t.Fatalf("expected: %s, got: %s", "mald", d.Comments())
	}

	feh = deklParser.ParseString("mald",
		`// mald
funk A() => 1`, &d)
	if feh != nil {
		t.Fatalf("error: %s", feh)
	}
	if d.Comments() != "mald" {
		t.Fatalf("expected: %s, got: %s", "mald", d.Comments())
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
