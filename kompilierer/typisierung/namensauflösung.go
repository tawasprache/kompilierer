package typisierung

import (
	"Tawa/kompilierer/ast"
	"Tawa/kompilierer/getypisiertast"

	"github.com/alecthomas/participle/v2/lexer"
)

func exprNamensauflösung(k *Kontext, s *scopes, l *lokalekontext, astExpr ast.Expression) (getypisiertast.Expression, error) {
	var feh error

	if astExpr.Ganzzahl != nil {
		return getypisiertast.Ganzzahl{
			Wert: *astExpr.Ganzzahl,
			LPos: astExpr.Pos,
		}, nil
	} else if astExpr.Passt != nil {
		var (
			wert    getypisiertast.Expression
			mustern []getypisiertast.Muster

			lTyp getypisiertast.ITyp = getypisiertast.Nichtunifiziert{}
			lPos lexer.Position      = astExpr.Pos
		)

		wert, feh = exprNamensauflösung(k, s, l, astExpr.Passt.Wert)
		if feh != nil {
			return nil, feh
		}

		for _, it := range astExpr.Passt.Mustern {
			_, _, sym, feh := l.auflöseVariant(it.Pattern.Name, it.Pattern.Pos)
			if feh != nil {
				return nil, feh
			}

			muster := getypisiertast.Muster{
				Variante:    sym,
				Konstruktor: it.Pattern.Name,
			}

			s.neuScope()
			defer s.loescheScope()

			for idx, musterVari := range it.Pattern.Variabeln {
				s.head().vars[musterVari] = getypisiertast.Nichtunifiziert{}
				muster.Variablen = append(muster.Variablen, getypisiertast.Mustervariable{
					Variante:    sym,
					Konstruktor: it.Pattern.Name,
					VonFeld:     idx,
					Name:        musterVari,
				})
			}

			musterExpr, feh := exprNamensauflösung(k, s, l, it.Expression)
			if feh != nil {
				return nil, feh
			}

			muster.Expression = musterExpr

			mustern = append(mustern, muster)
		}

		return getypisiertast.Pattern{
			Wert:    wert,
			Mustern: mustern,

			LTyp: lTyp,
			LPos: lPos,
		}, nil
	} else if astExpr.Variantaufruf != nil {
		var (
			variant     getypisiertast.SymbolURL
			konstruktor string = astExpr.Variantaufruf.Name
			argumenten  []getypisiertast.Expression
			varianttyp  getypisiertast.ITyp = getypisiertast.Nichtunifiziert{}

			lPos lexer.Position = astExpr.Pos
		)

		_, _, variant, feh := l.auflöseVariant(astExpr.Variantaufruf.Name, astExpr.Pos)
		if feh != nil {
			return nil, feh
		}

		for _, it := range astExpr.Variantaufruf.Argumente {
			variExpr, feh := exprNamensauflösung(k, s, l, it)
			if feh != nil {
				return nil, feh
			}
			argumenten = append(argumenten, variExpr)
		}

		return getypisiertast.Variantaufruf{
			Variant:     variant,
			Konstruktor: konstruktor,
			Argumenten:  argumenten,
			Varianttyp:  varianttyp,

			LPos: lPos,
		}, nil
	} else if astExpr.Funktionsaufruf != nil {
		var (
			funktion    getypisiertast.SymbolURL
			argumenten  []getypisiertast.Expression
			rückgabetyp getypisiertast.ITyp = getypisiertast.Nichtunifiziert{}
			lPos        lexer.Position      = astExpr.Pos
		)

		_, funktion, feh = l.auflöseFunkSig(astExpr.Funktionsaufruf.Name, astExpr.Pos)
		if feh != nil {
			return nil, feh
		}

		for _, arg := range astExpr.Funktionsaufruf.Argumente {
			garg, feh := exprNamensauflösung(k, s, l, arg)
			if feh != nil {
				return nil, feh
			}
			argumenten = append(argumenten, garg)
		}

		return getypisiertast.Funktionsaufruf{
			Funktion:    funktion,
			Argumenten:  argumenten,
			Rückgabetyp: rückgabetyp,
			LPos:        lPos,
		}, nil
	} else if astExpr.Variable != nil {
		_, gefunden := s.suche(*astExpr.Variable)

		if !gefunden {
			return nil, neuFehler(astExpr.Pos, "variable »%s« nicht gefunden", *astExpr.Variable)
		}

		return getypisiertast.Variable{
			Name: *astExpr.Variable,
			ITyp: getypisiertast.Nichtunifiziert{},
			LPos: astExpr.Pos,
		}, nil
	}

	panic("unhandled case")
}

func typNamensauflösung(k *Kontext, l *lokalekontext, t ast.Typdeklarationen) (getypisiertast.Typ, error) {
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

func typ(l *lokalekontext, t ast.Typ, generischeargumenten []string) (getypisiertast.ITyp, error) {
	if t.Typkonstruktor != nil {
		haupt, err := l.auflöseTyp(t.Typkonstruktor.Name, t.Pos)
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

func funkZuSignatur(l *lokalekontext, t ast.Funktiondeklaration) (getypisiertast.Funktionssignatur, error) {
	r := getypisiertast.Funktionssignatur{}
	r.Generischeargumenten = t.Generischeargumenten

	if t.Rückgabetyp == nil {
		r.Rückgabetyp = getypisiertast.Typnutzung{
			SymbolURL: getypisiertast.SymbolURL{
				Paket: "Tawa/Eingebaut",
				Name:  "Einheit",
			},
		}
	} else {
		t, err := typ(l, *t.Rückgabetyp, t.Generischeargumenten)
		if err != nil {
			return getypisiertast.Funktionssignatur{}, nil
		}
		r.Rückgabetyp = t
	}

	for _, it := range t.Formvariabeln {
		t, err := typ(l, it.Typ, t.Generischeargumenten)
		if err != nil {
			return getypisiertast.Funktionssignatur{}, nil
		}
		r.Formvariabeln = append(r.Formvariabeln, getypisiertast.Formvariable{
			Name: it.Name,
			Typ:  t,
		})
	}

	return r, nil
}

func Auflösenamen(k *Kontext, m ast.Modul, modulePrefix string) (getypisiertast.Modul, error) {
	modul := getypisiertast.Modul{}
	l := &lokalekontext{
		k:                k,
		modul:            &modul,
		inModul:          modulePrefix + "/" + m.Package,
		lokaleFunktionen: map[string]getypisiertast.Funktionssignatur{},
		importieren:      []string{"Tawa/Eingebaut"},
	}
	modul.Name = modulePrefix + "/" + m.Package
	if modul.Name == "Tawa/Eingebaut" {
		l.importieren = []string{}
	}
	for _, it := range m.Importierungen {
		l.importieren = append(l.importieren, it.Import)
		// TODO: geimportierte paket zum kontext hinzufuegen
	}
	modul.Dependencies = append(modul.Dependencies, l.importieren...)

	for _, it := range m.Deklarationen {
		if it.Typdeklarationen != nil {
			typ, feh := typNamensauflösung(k, l, *it.Typdeklarationen)
			if feh != nil {
				return getypisiertast.Modul{}, feh
			}
			modul.Typen = append(modul.Typen, typ)
		}
	}

	for _, it := range m.Deklarationen {
		if it.Funktiondeklaration != nil {
			sig, feh := funkZuSignatur(l, *it.Funktiondeklaration)
			if feh != nil {
				return getypisiertast.Modul{}, feh
			}
			l.lokaleFunktionen[it.Funktiondeklaration.Name] = sig
		}
	}

	for _, it := range m.Deklarationen {
		if it.Funktiondeklaration != nil {
			s := scopes{}
			s.neuScope()

			for _, it := range it.Funktiondeklaration.Formvariabeln {
				s.head().vars[it.Name] = getypisiertast.Nichtunifiziert{}
			}

			expr, feh := exprNamensauflösung(k, &s, l, it.Funktiondeklaration.Expression)
			if feh != nil {
				return getypisiertast.Modul{}, feh
			}

			modul.Funktionen = append(modul.Funktionen, getypisiertast.Funktion{
				SymbolURL: getypisiertast.SymbolURL{
					Paket: l.inModul,
					Name:  it.Funktiondeklaration.Name,
				},
				Funktionssignatur: l.lokaleFunktionen[it.Funktiondeklaration.Name],
				Expression:        expr,
			})
		}
	}

	return modul, nil
}
