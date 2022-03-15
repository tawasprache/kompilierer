package ast

import "Tawa/kompilierer/v2/parser"

type Feld struct {
	pos

	Name Ident
	Typ  Typnutzung
}

type Feldliste struct {
	pos

	Felden []Feld
}

var _ Node = Feld{}
var _ Node = Feldliste{}

func feldVonParser(p parser.Verbunddeklarationsfeld) Feld {
	return Feld{
		pos: pos{
			anfang: p.Pos,
			ende:   p.EndPos,
		},
		Name: identVonParser(p.Name),
		Typ:  typnutzungVonParser(p.Typ),
	}
}
