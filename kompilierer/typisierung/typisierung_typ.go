package typisierung

import (
	"Tawa/kompilierer/ast"
	"Tawa/kompilierer/getypisiertast"
)

func typ(l *lokalekontext, t ast.Typ, generischeargumenten []string) (getypisiertast.ITyp, error) {
	if t.Typkonstruktor != nil {
		haupt, err := l.sucheTyp(t.Typkonstruktor.Name, t.Pos)
		if err != nil {
			return nil, err
		}

		var ts []getypisiertast.ITyp
		for _, it := range t.Typkonstruktor.Generischeargumenten {
			t, err := typ(l, it, generischeargumenten)
			if err != nil {
				return nil, err
			}
			ts = append(ts, t)
		}

		return getypisiertast.Typnutzung{
			SymbolURL:            haupt,
			Generischeargumenten: ts,
		}, nil
	} else if t.Typvariable != nil {
		s := string(*t.Typvariable)
		for _, it := range generischeargumenten {
			if it == s {
				return getypisiertast.Typvariable{Name: string(*t.Typvariable)}, nil
			}
		}
		return nil, neuFehler(t.Pos, "typvariable »%s« nicht gefunden", s)
	}
	panic("unreachable")
}

func typDeklZu(l *lokalekontext, t ast.Typdeklarationen) (getypisiertast.Typ, error) {
	g := getypisiertast.Typ{}
	g.SymbolURL = getypisiertast.SymbolURL{
		Paket: l.inModul,
		Name:  t.Name,
	}
	g.Generischeargumenten = t.Generischeargumenten

	for _, it := range t.Datenfelden {
		feld := getypisiertast.Datenfeld{
			Name: it.Name,
		}
		typ, err := typ(l, it.Typ, t.Generischeargumenten)
		if err != nil {
			return getypisiertast.Typ{}, err
		}
		feld.Typ = typ
		g.Datenfelden = append(g.Datenfelden, feld)
	}
	for _, it := range t.Varianten {
		vari := getypisiertast.Variant{
			Name: it.Name,
		}
		for _, feld := range it.Datenfelden {
			typ, err := typ(l, feld.Typ, t.Generischeargumenten)
			if err != nil {
				return getypisiertast.Typ{}, err
			}
			vari.Datenfelden = append(vari.Datenfelden, getypisiertast.Datenfeld{
				Name: feld.Name,
				Typ:  typ,
			})
		}
		g.Varianten = append(g.Varianten, vari)
	}

	return g, nil
}
