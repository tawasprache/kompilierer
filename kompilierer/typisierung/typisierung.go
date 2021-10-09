package typisierung

import "Tawa/kompilierer/getypisiertast"

func Typiere(k *Kontext, m getypisiertast.Modul, modulePrefix string) (getypisiertast.Modul, error) {
	m = copy(m).(getypisiertast.Modul)

	l := &lokalekontext{
		k:                k,
		modul:            &m,
		inModul:          modulePrefix + "/" + m.Name,
		lokaleFunktionen: map[string]getypisiertast.Funktionssignatur{},
		importieren:      m.Dependencies,
	}

	for _, it := range m.Funktionen {
		l.lokaleFunktionen[it.SymbolURL.Name] = it.Funktionssignatur
	}

	var serr error

	for idx, it := range m.Funktionen {
		s := scopes{}
		s.neuScope()

		for _, it := range it.Funktionssignatur.Formvariabeln {
			s.head().vars[it.Name] = it.Typ
		}

		istEinheit := false
		v, ok := it.Funktionssignatur.Rückgabetyp.(getypisiertast.Typnutzung)
		if ok && v.SymbolURL.Name == "Einheit" && v.SymbolURL.Paket == "Tawa/Eingebaut" {
			istEinheit = true
		}

		var rückgabe getypisiertast.Expression
		var err error

		if istEinheit {
			rückgabe, err = synthGetypisiertExpression(l, &s, it.Expression)
		} else {
			rückgabe, err = checkGetypisiertExpression(l, &s, it.Expression, it.Funktionssignatur.Rückgabetyp)
		}

		serr = fehlerVerketten(serr, err)
		m.Funktionen[idx].Expression = rückgabe
	}

	return m, serr
}
