package ast

import "github.com/alecthomas/participle/v2/lexer"

type Node interface {
	Anfang() lexer.Position
	Ende() lexer.Position
}

type Deklaration interface {
	Node
}

type pos struct {
	anfang lexer.Position
	ende   lexer.Position
}

func (p *pos) Anfang() lexer.Position {
	return p.anfang
}

func (p *pos) Ende() lexer.Position {
	return p.ende
}
