package typisierung

import (
	"Tawa/kompilierer/ast"
	"Tawa/kompilierer/getypisiertast"
)

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

func getyptFunkZu(l *lokalekontext, t getypisiertast.Funktion) error {
	s := scopes{}
	s.neuScope()

	for _, it := range t.Funktionssignatur.Formvariabeln {
		s.head().vars[it.Name] = it.Typ
	}

	istEinheit := false
	v, ok := t.Funktionssignatur.Rückgabetyp.(getypisiertast.Typnutzung)
	if ok && v.SymbolURL.Name == "Einheit" && v.SymbolURL.Paket == "Tawa/Eingebaut" {
		istEinheit = true
	}

	var err error

	if istEinheit {
		_, err = synthGetypisiertExpression(l, &s, t.Expression)
	} else {
		err = checkGetypisiertExpression(l, &s, t.Expression, t.Funktionssignatur.Rückgabetyp)
	}

	return err
}

func funkZu(l *lokalekontext, t ast.Funktiondeklaration) (getypisiertast.Funktion, error) {
	s := scopes{}
	s.neuScope()

	sig := l.lokaleFunktionen[t.Name]

	for _, it := range sig.Formvariabeln {
		s.head().vars[it.Name] = it.Typ
	}

	istEinheit := false
	v, ok := sig.Rückgabetyp.(getypisiertast.Typnutzung)
	if ok && v.SymbolURL.Name == "Einheit" && v.SymbolURL.Paket == "Tawa/Eingebaut" {
		istEinheit = true
	}

	var expr getypisiertast.Expression
	var err error

	if istEinheit {
		expr, err = synthExpression(l, &s, t.Expression)
		if err != nil {
			return getypisiertast.Funktion{}, err
		}
	} else {
		expr, err = checkExpression(l, &s, t.Expression, sig.Rückgabetyp)
		if err != nil {
			return getypisiertast.Funktion{}, err
		}
	}

	return getypisiertast.Funktion{
		SymbolURL: getypisiertast.SymbolURL{
			Paket: l.inModul,
			Name:  t.Name,
		},
		Funktionssignatur: sig,
		Expression:        expr,
	}, nil
}
