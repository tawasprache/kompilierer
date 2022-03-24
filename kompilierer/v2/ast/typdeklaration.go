package ast

import "Tawa/kompilierer/v2/parser"

type Fall struct {
	pos

	Name   *Ident
	Felden []*Feld
}

type Typdeklaration struct {
	pos

	Name   *Ident
	Felden []*Feld
	Fälle  []*Fall
}

var _ Deklaration = &Typdeklaration{}

func typDeklarationVonParser(p parser.Typdeklaration) *Typdeklaration {
	var s []*Feld
	var f []*Fall
	for _, feld := range p.Verbunddeklaration.Felden {
		s = append(s, feldVonParser(feld))
	}
	for _, fall := range p.Verbunddeklaration.Fallen {
		f = append(f, &Fall{
			pos:  pos{fall.Pos, fall.EndPos},
			Name: identVonParser(fall.Name),
			Felden: func() (r []*Feld) {
				for _, es := range fall.Felden {
					r = append(r, feldVonParser(es))
				}
				return r
			}(),
		})
	}

	return &Typdeklaration{
		pos: pos{
			anfang: p.Pos,
			ende:   p.EndPos,
		},
		Name:   identVonParser(p.Name),
		Felden: s,
		Fälle:  f,
	}
}
