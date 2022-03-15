package ast

import "Tawa/kompilierer/v2/parser"

type Ident struct {
	pos

	Name string
}

var _ Node = Ident{}

func identVonParser(p parser.Ident) Ident {
	return Ident{
		pos: pos{
			anfang: p.Pos,
			ende:   p.EndPos,
		},
		Name: p.Name,
	}
}

type Symbolkette struct {
	pos

	Name []Ident
}

var _ Node = Symbolkette{}

func symbolketteVonParser(p parser.Symbolkette) Symbolkette {
	idents := []Ident{}
	for _, es := range p.Symbolen {
		idents = append(idents, identVonParser(es))
	}
	return Symbolkette{
		pos: pos{
			anfang: p.Pos,
			ende:   p.EndPos,
		},
		Name: idents,
	}
}
