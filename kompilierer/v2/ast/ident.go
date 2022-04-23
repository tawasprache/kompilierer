package ast

import (
	"Tawa/kompilierer/v2/parser"
)

type Ident struct {
	pos

	Name string
}

func (i Ident) String() string {
	return i.Name
}

var _ Node = &Ident{}

func identVonParser(p parser.Ident) *Ident {
	return &Ident{
		pos: pos{
			anfang: p.Pos,
			ende:   p.EndPos,
		},
		Name: p.Name,
	}
}
