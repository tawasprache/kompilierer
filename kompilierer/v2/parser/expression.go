package parser

import "github.com/alecthomas/participle/v2/lexer"

type Expression struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Terminal *Terminal

	// oder

	Links  *Expression
	Op     *BinaryOperator
	Rechts *Expression
}

type Terminal struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Ganzzahl *int    `  @Int`
	Variable *string `| @Ident`
}
