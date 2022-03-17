package ast

import (
	"Tawa/kompilierer/v2/parser"
	"strings"
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

func symbolketteVonParser(p parser.Symbolkette) *Ident {
	idents := []string{}
	for _, es := range p.Symbolen {
		idents = append(idents, es.Name)
	}
	return &Ident{
		pos: pos{
			anfang: p.Pos,
			ende:   p.EndPos,
		},
		Name: strings.Join(idents, "::"),
	}
}
