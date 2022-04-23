package ast

import "Tawa/kompilierer/v2/parser"

type Typnutzung interface {
	Expression
}

func typnutzungVonParser(p parser.Typ) Typnutzung {
	return expressionVonParser(*p.Expression)
}
