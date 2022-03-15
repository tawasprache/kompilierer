package ast

import "Tawa/kompilierer/v2/parser"

type Anweisung interface {
	Node
}

func gibVonParser(p parser.Gib) *Gib {
	return &Gib{
		pos:  pos{p.Pos, p.EndPos},
		Wert: expressionVonParser(p.Wert),
	}
}

func istVonParser(p parser.Ist) *Ist {
	return &Ist{
		pos:      pos{p.Pos, p.EndPos},
		Variable: identVonParser(p.Name),
		Wert:     expressionVonParser(p.Wert),
	}
}

func seiVonParser(p parser.Sei) *Sei {
	return &Sei{
		pos:      pos{p.Pos, p.EndPos},
		Variable: identVonParser(p.Name),
		Wert:     expressionVonParser(p.Wert),
	}
}

func anweisungVonParser(p parser.Anweisung) Anweisung {
	if p.Gib != nil {
		return gibVonParser(*p.Gib)
	} else if p.Ist != nil {
		return istVonParser(*p.Ist)
	} else if p.Sei != nil {
		return seiVonParser(*p.Sei)
	} else {
		panic("TODO")
	}
}

type Gib struct {
	pos

	Wert Expression
}

var _ Anweisung = &Gib{}

type Ist struct {
	pos

	Variable *Ident
	Wert     Expression
}

var _ Anweisung = &Ist{}

type Sei struct {
	pos

	Variable *Ident
	Wert     Expression
}

var _ Anweisung = &Sei{}
