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

type Symbolkette struct {
	pos

	Name []*Ident
}

var _ Node = &Symbolkette{}

func (s *Symbolkette) String() string {
	var sb strings.Builder
	for idx, es := range s.Name {
		sb.WriteString(es.Name)
		if idx+1 < len(s.Name) {
			sb.WriteString("::")
		}
	}
	return sb.String()
}

func symbolketteVonParser(p parser.Symbolkette) *Symbolkette {
	idents := []*Ident{}
	for _, es := range p.Symbolen {
		idents = append(idents, identVonParser(es))
	}
	return &Symbolkette{
		pos: pos{
			anfang: p.Pos,
			ende:   p.EndPos,
		},
		Name: idents,
	}
}
