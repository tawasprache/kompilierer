package ast

import "Tawa/kompilierer/v2/parser"

type Expression interface {
	Node
	istExpression()
}

type istExpressionImpl struct {
}

func (istExpressionImpl) istExpression() {}

func terminalVonParser(p parser.Terminal) Expression {
	if p.Variable != nil {
		return &IdentExpression{
			pos:   pos{p.Pos, p.EndPos},
			Ident: symbolketteVonParser(*p.Variable),
		}
	} else if p.Ganzzahl != nil {
		return &GanzzahlExpression{
			pos:  pos{p.Pos, p.EndPos},
			Wert: *p.Ganzzahl,
		}
	} else if p.Zeichenkette != nil {
		return &ZeichenketteExpression{
			pos:  pos{p.Pos, p.EndPos},
			Wert: *p.Zeichenkette,
		}
	} else {
		panic("e")
	}
}

type GanzzahlExpression struct {
	pos

	Wert int

	istExpressionImpl
}

type ZeichenketteExpression struct {
	pos

	Wert string

	istExpressionImpl
}

type IdentExpression struct {
	pos

	Ident *Ident

	istExpressionImpl
}

var _ Expression = &IdentExpression{}

type SelektorExpression struct {
	pos

	Objekt Expression
	Feld   *Ident

	istExpressionImpl
}

type BinaryExpression struct {
	pos

	Links    Expression
	Operator parser.BinaryOperator
	Rechts   Expression

	istExpressionImpl
}

var _ Expression = &BinaryExpression{}

func binaryVonParser(
	p parser.Expression,
	links parser.Expression,
	op parser.BinaryOperator,
	rechts parser.Expression,
) Expression {
	return &BinaryExpression{
		pos: pos{p.Pos, p.EndPos},

		Links:    expressionVonParser(links),
		Operator: op,
		Rechts:   expressionVonParser(rechts),
	}
}

func expressionVonParser(p parser.Expression) Expression {
	if p.Terminal != nil {
		return terminalVonParser(*p.Terminal)
	} else if p.Links != nil && p.Rechts != nil && p.Op != nil {
		return binaryVonParser(p, *p.Links, *p.Op, *p.Rechts)
	} else if p.Objekt != nil && p.Selektor != nil {
		return &SelektorExpression{
			pos:    pos{p.Pos, p.EndPos},
			Objekt: expressionVonParser(*p.Objekt),
			Feld:   identVonParser(*p.Selektor),
		}
	} else {
		panic("e")
	}
}
