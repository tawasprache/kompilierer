package ast

import "Tawa/kompilierer/v2/parser"

type Typdeklaration struct {
	pos

	Felden []Feld
}

var _ Deklaration = Typdeklaration{}

func typDeklarationVonParser(p parser.Typdeklaration) Typdeklaration {
	var s []Feld
	for _, feld := range p.Verbunddeklaration.Felden {
		s = append(s, feldVonParser(feld))
	}

	return Typdeklaration{
		pos: pos{
			anfang: p.Pos,
			ende:   p.EndPos,
		},
		Felden: s,
	}
}
