package ast

import "Tawa/kompilierer/v2/parser"

type Anweisung interface {
	Node
}

func anweisungVonParser(p parser.Anweisung) Anweisung {
	if p.Gib != nil {
		panic("TODO Gib")
	} else if p.Ist != nil {
		panic("TODO Ist")
	} else if p.Sei != nil {
		panic("TODO Sei")
	} else {
		panic("TODO")
	}
}
