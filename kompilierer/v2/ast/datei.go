package ast

import "Tawa/kompilierer/v2/parser"

type Datei struct {
	pos

	Deklarationen []Deklaration
}

func VonParser(p parser.Modul) Datei {
	datei := Datei{
		pos: pos{
			anfang: p.Pos,
			ende:   p.EndPos,
		},
	}

	for _, es := range p.Moduldeklarationen {
		if es.Funktiondeklaration != nil {
			datei.Deklarationen = append(datei.Deklarationen, funktionDeklarationVonParser(*es.Funktiondeklaration))
		} else if es.Typdeklaration != nil {
			datei.Deklarationen = append(datei.Deklarationen, typDeklarationVonParser(*es.Typdeklaration))
		} else {
			panic("e")
		}
	}

	return datei
}
