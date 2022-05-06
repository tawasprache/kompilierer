package ast

import "Tawa/kompilierer/v2/parser"

type Verwendung struct {
	pos

	Paket string
	Als   *Ident
}

var _ Node = &Verwendung{}

func verwendungVonParser(p parser.Verwendung) *Verwendung {
	v := &Verwendung{
		pos: pos{p.Pos, p.EndPos},

		Paket: p.Paket,
	}
	if p.Als != nil {
		v.Als = identVonParser(*p.Als)
	}

	return v
}
