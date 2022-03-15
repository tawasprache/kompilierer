package ast

import "Tawa/kompilierer/v2/parser"

type Funktiondeklaration struct {
	pos

	Name   Ident
	Inhalt Block
}

var _ Deklaration = Funktiondeklaration{}

func funktionDeklarationVonParser(p parser.Funktiondeklaration) Funktiondeklaration {
	return Funktiondeklaration{
		pos: pos{
			anfang: p.Pos,
			ende:   p.EndPos,
		},
		Name:   identVonParser(p.Name),
		Inhalt: blockVonParser(p.Inhalt),
	}
}
