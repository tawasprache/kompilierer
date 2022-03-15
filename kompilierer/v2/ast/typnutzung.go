package ast

import "Tawa/kompilierer/v2/parser"

type Typnutzung interface {
	Node
}

type Typkonstruktor struct {
	pos

	Symbolkette *Symbolkette
}

var _ Typnutzung = &Typkonstruktor{}

func typnutzungVonParser(p parser.Typ) Typnutzung {
	if p.Typkonstruktor != nil {
		return &Typkonstruktor{
			pos: pos{
				anfang: p.Pos,
				ende:   p.EndPos,
			},
			Symbolkette: symbolketteVonParser(p.Typkonstruktor.Name),
		}
	} else {
		panic("e")
	}
}
