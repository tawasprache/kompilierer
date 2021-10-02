package typisierung

import "Tawa/kompilierer/getypisiertast"

// assert typeof e == a
func checkGetypisiertExpression(l *lokalekontext, s *scopes, expr getypisiertast.Expression, gegenTyp getypisiertast.ITyp) error {
	ruck, err := synthGetypisiertExpression(l, s, expr)
	if err != nil {
		return err
	}

	if !gleich(ruck, gegenTyp) {
		return gleichErr(expr.Pos(), ruck, gegenTyp)
	}

	return nil
}

// typeof e
func synthGetypisiertExpression(l *lokalekontext, s *scopes, expr getypisiertast.Expression) (getypisiertast.ITyp, error) {
	switch e := expr.(type) {
	case getypisiertast.Ganzzahl:
		return e.Typ(), nil
	case getypisiertast.Funktionsaufruf:
		rumpf, feh := l.sucheFunktionsrumpf(e.Funktion, expr.Pos())
		if feh != nil {
			return nil, feh
		}
		return synthGetypisiertApplication(l, s, rumpf, e)
	case getypisiertast.Variable:
		return e.Typ(), nil
	case *getypisiertast.Variantaufruf:
		panic("unimplemented getypisiertast.Variantaufruf")
	}
	panic("unreachable")
}

func synthGetypisiertApplication(l *lokalekontext, s *scopes, funktion getypisiertast.Funktion, aufruf getypisiertast.Funktionsaufruf) (getypisiertast.ITyp, error) {
	sig := funktion.Funktionssignatur
	sigArg := sig.Formvariabeln
	arg := aufruf.Argumenten

	if len(sigArg) != len(arg) {
		return nil, neuFehler(aufruf.Pos(), "len nicht gleich")
	}

	vars := map[string]getypisiertast.ITyp{}

	for idx := range funktion.Funktionssignatur.Formvariabeln {
		eingabe := arg[idx]

		b, f := synthGetypisiertExpression(l, s, eingabe)
		if f != nil {
			return nil, f
		}

		a := sigArg[idx].Typ

		es := unify(b, a)
		for k, v := range es {
			if _, ok := vars[k]; !ok {
				vars[k] = v
			}
			a = vars[k]
		}

		if !gleich(a, b) {
			return nil, gleichErr(aufruf.Pos(), a, b)
		}
	}

	ret := sig.Rückgabetyp
	for k, v := range vars {
		ret = substitute(ret, getypisiertast.Typvariable{Name: k}, v)
	}

	if len(vars) > 0 {
		var feh error

		rumpf := funktion
		for k, tvar := range vars {
			for idx, fvar := range rumpf.Funktionssignatur.Formvariabeln {
				rumpf.Funktionssignatur.Formvariabeln[idx].Typ = substitute(fvar.Typ, getypisiertast.Typvariable{Name: k}, tvar)
			}
			rumpf.Funktionssignatur.Rückgabetyp = substitute(rumpf.Funktionssignatur.Rückgabetyp, getypisiertast.Typvariable{Name: k}, tvar)
			rumpf.Expression = substituteExpression(rumpf.Expression, getypisiertast.Typvariable{Name: k}, tvar)
		}
		feh = getyptFunkZu(l, rumpf)
		if feh != nil {
			return nil, feh
		}
	}

	return ret, nil
}
