package ast_test

import (
	"Tawa/kompilierer/v2/ast"
	"Tawa/kompilierer/v2/parser"
	"testing"
)

func TestInspect(t *testing.T) {
	var v parser.Modul
	feh := parser.Parser.ParseString("Hallo.tawa", hallo, &v)
	if feh != nil {
		t.Fatalf("error: %s", feh)
	}

	k := ast.VonParser(v)
	ast.Inspect(k, func(n ast.Node) bool {
		return true
	})
}
