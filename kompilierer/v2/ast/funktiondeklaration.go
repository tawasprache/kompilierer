package ast

import (
	"Tawa/kompilierer/v2/parser"
)

type Argument struct {
	pos

	Name *Ident
	Typ  Typnutzung
}

type Argumentliste struct {
	pos

	Argumente []*Argument
}

func argumentlisteVonParser(p parser.Argumentliste) *Argumentliste {
	var a []*Argument
	for _, arg := range p.Argumenten {
		for _, name := range arg.Namen {
			a = append(a, &Argument{
				pos: pos{name.Pos, arg.EndPos},

				Name: identVonParser(name),
				Typ:  typnutzungVonParser(arg.Typ),
			})
		}
	}

	return &Argumentliste{
		pos:       pos{p.Pos, p.EndPos},
		Argumente: a,
	}
}

type Funktiondeklaration struct {
	pos

	Name        *Ident
	Inhalt      *Block
	Argumenten  *Argumentliste
	Rückgabetyp Typnutzung
}

var _ Deklaration = &Funktiondeklaration{}

func funktionDeklarationVonParser(p parser.Funktiondeklaration) *Funktiondeklaration {
	k := &Funktiondeklaration{
		pos: pos{
			anfang: p.Pos,
			ende:   p.EndPos,
		},
		Name:       identVonParser(p.Name),
		Inhalt:     blockVonParser(p.Inhalt),
		Argumenten: argumentlisteVonParser(p.Argumenten),
	}
	if p.Rückgabetyp != nil {
		k.Rückgabetyp = typnutzungVonParser(*p.Rückgabetyp)
	}
	return k
}
