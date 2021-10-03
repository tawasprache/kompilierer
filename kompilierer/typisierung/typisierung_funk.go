package typisierung

import (
	"Tawa/kompilierer/getypisiertast"
)

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
		_, err = checkGetypisiertExpression(l, &s, t.Expression, t.Funktionssignatur.Rückgabetyp)
	}

	return err
}
