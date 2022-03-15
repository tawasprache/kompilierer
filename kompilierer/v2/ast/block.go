package ast

import "Tawa/kompilierer/v2/parser"

type Block struct {
	pos

	Anweisungen []Anweisung
}

var _ Node = Block{}

func blockVonParser(p parser.Block) Block {
	var a []Anweisung

	for _, es := range p.Anweisungen {
		a = append(a, anweisungVonParser(es))
	}

	return Block{
		pos: pos{
			anfang: p.Pos,
			ende:   p.EndPos,
		},
		Anweisungen: a,
	}
}
