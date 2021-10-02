package typisierung

import (
	"Tawa/kompilierer/ast"
	"Tawa/kompilierer/getypisiertast"
	"reflect"

	"github.com/alecthomas/participle/v2/lexer"
)

func gleich(a getypisiertast.ITyp, b getypisiertast.ITyp) bool {
	_, lhsIstTypvar := a.(getypisiertast.Typvariable)
	_, rhsIstTypvar := b.(getypisiertast.Typvariable)
	de := reflect.DeepEqual(a, b)

	return !(!de && !lhsIstTypvar && !rhsIstTypvar)
}

func gleichErr(p lexer.Position, a getypisiertast.ITyp, b getypisiertast.ITyp) error {
	return neuFehler(p, "nicht gleich: »%s« »%s«", a, b)
}

// assert typeof e == a
func checkExpression(l *lokalekontext, s *scopes, expr ast.Expression, gegenTyp getypisiertast.ITyp) (getypisiertast.Expression, error) {
	ruck, err := synthExpression(l, s, expr)
	if err != nil {
		return nil, err
	}

	if !gleich(ruck.Typ(), gegenTyp) {
		return nil, gleichErr(expr.Pos, ruck.Typ(), gegenTyp)
	}

	return ruck, nil
}

// typeof e
func synthExpression(l *lokalekontext, s *scopes, expr ast.Expression) (getypisiertast.Expression, error) {
	if expr.Ganzzahl != nil {
		return getypisiertast.Ganzzahl{Wert: *expr.Ganzzahl}, nil
	} else if expr.Variable != nil {
		typ, ok := s.suche(*expr.Variable)
		if !ok {
			return nil, neuFehler(expr.Pos, "variable »%s« nicht gefunden", *expr.Variable)
		}
		return getypisiertast.Variable{
			Name: *expr.Variable,
			ITyp: typ,
		}, nil
	} else if expr.Funktionsaufruf != nil {
		sig, url, err := l.sucheFunkSig(expr.Funktionsaufruf.Name, expr.Pos)
		if err != nil {
			return nil, err
		}

		return synthApplication(l, s, url, sig, expr, expr.Funktionsaufruf.Argumente)
	} else if expr.Variantaufruf != nil {
		sig, vari, url, err := l.sucheVariant(expr.Variantaufruf.Name, expr.Pos)
		if err != nil {
			return nil, err
		}

		return synthVariantApplication(l, s, sig, vari, url, expr, expr.Variantaufruf.Argumente)
	} else if expr.Passt != nil {
		wert, feh := synthExpression(l, s, expr.Passt.Wert)
		if feh != nil {
			return nil, feh
		}
		switch k := wert.Typ().(type) {
		case getypisiertast.Typvariable:
			panic("idk")
		case getypisiertast.Typnutzung:
			typDekl, feh := l.sucheTypDekl(k.SymbolURL, expr.Passt.Wert.Pos)
			if feh != nil {
				return nil, feh
			}

			for idx, tvar := range typDekl.Generischeargumenten {
				typDekl = substituteTypdekl(typDekl, getypisiertast.Typvariable{Name: tvar}, k.Generischeargumenten[idx])
			}

			if len(typDekl.Varianten) < 2 {
				return nil, neuFehler(expr.Passt.Wert.Pos, " »%s« hat kein varianten", k.SymbolURL)
			}

			if len(typDekl.Varianten) != len(expr.Passt.Mustern) {
				return nil, neuFehler(expr.Passt.Wert.Pos, "len nicht gleich")
			}

			suche := func(s string, p lexer.Position) (getypisiertast.Variant, error) {
				for _, it := range typDekl.Varianten {
					if it.Name == s {
						return it, nil
					}
				}
				return getypisiertast.Variant{}, neuFehler(p, "variant »%s« existiert nicht", s)
			}

			varis := map[string]getypisiertast.Variant{}

			for _, it := range expr.Passt.Mustern {
				vari, feh := suche(it.Pattern.Name, it.Pattern.Pos)
				if feh != nil {
					return nil, feh
				}

				if len(vari.Datenfelden) != len(it.Pattern.Variabeln) {
					return nil, neuFehler(it.Pattern.Pos, "variante »%s« hat %d variablen, aber du nutzt %d", vari.Name, len(vari.Datenfelden), len(it.Pattern.Variabeln))
				}

				varis[it.Pattern.Name] = vari
			}

			var mustern []getypisiertast.Muster
			var kind getypisiertast.ITyp

			for _, it := range expr.Passt.Mustern {
				s.neuScope()

				vari := varis[it.Pattern.Name]
				for idx, feld := range it.Pattern.Variabeln {
					s.head().vars[feld] = vari.Datenfelden[idx].Typ
				}

				expr, feh := synthExpression(l, s, it.Expression)
				if feh != nil {
					return nil, feh
				}

				var mvars []getypisiertast.Mustervariable

				for idx, it := range it.Pattern.Variabeln {
					mvars = append(mvars, getypisiertast.Mustervariable{
						Name:    it,
						VonFeld: vari.Datenfelden[idx].Name,
						ITyp:    vari.Datenfelden[idx].Typ,
					})
				}

				mustern = append(mustern, getypisiertast.Muster{
					Name:       it.Pattern.Name,
					Variablen:  mvars,
					Expression: expr,
				})

				if kind == nil {
					kind = expr.Typ()
				} else {
					if !gleich(kind, expr.Typ()) {
						return nil, neuFehler(it.Expression.Pos, "arme sind nicht gleich, erwartete %s und sah %s", kind, expr.Typ())
					}
				}

				s.loescheScope()
			}

			return getypisiertast.Pattern{
				Wert:    wert,
				Mustern: mustern,

				LTyp: kind,
				LPos: expr.Pos,
			}, nil
		}
	}
	panic("unreachable")
}

func unify(lhs getypisiertast.ITyp, rhs getypisiertast.ITyp) map[string]getypisiertast.ITyp {
	if rhv, ok := rhs.(getypisiertast.Typvariable); ok {
		return map[string]getypisiertast.ITyp{
			rhv.Name: lhs,
		}
	} else {
		lhv := lhs.(getypisiertast.Typnutzung)
		rhv := rhs.(getypisiertast.Typnutzung)

		if len(lhv.Generischeargumenten) != len(rhv.Generischeargumenten) {
			panic("idk")
		}

		r := map[string]getypisiertast.ITyp{}

		for idx := range lhv.Generischeargumenten {
			lhs2 := lhv.Generischeargumenten[idx]
			rhs2 := rhv.Generischeargumenten[idx]

			for k, v := range unify(lhs2, rhs2) {
				r[k] = v
			}
		}
	}
	panic("a")
}

func substituteList(typ []getypisiertast.ITyp, suche getypisiertast.ITyp, ersetzen getypisiertast.ITyp) (r []getypisiertast.ITyp) {
	for _, it := range typ {
		r = append(r, substitute(it, suche, ersetzen))
	}
	return r
}

func substituteTypdekl(typ getypisiertast.Typ, suche getypisiertast.ITyp, ersetzen getypisiertast.ITyp) (r getypisiertast.Typ) {
	r.Generischeargumenten = append(r.Generischeargumenten, typ.Generischeargumenten...)
	r.SymbolURL = typ.SymbolURL
	for _, it := range typ.Datenfelden {
		r.Datenfelden = append(r.Datenfelden, getypisiertast.Datenfeld{
			Name: it.Name,
			Typ:  substitute(it.Typ, suche, ersetzen),
		})
	}
	for _, vari := range typ.Varianten {
		var rVari getypisiertast.Variant
		rVari.Name = vari.Name

		for _, feld := range vari.Datenfelden {
			rVari.Datenfelden = append(rVari.Datenfelden, getypisiertast.Datenfeld{
				Name: feld.Name,
				Typ:  substitute(feld.Typ, suche, ersetzen),
			})
		}

		r.Varianten = append(r.Varianten, rVari)
	}
	return r
}

func substitute(typ getypisiertast.ITyp, suche getypisiertast.ITyp, ersetzen getypisiertast.ITyp) getypisiertast.ITyp {
	switch v := typ.(type) {
	case getypisiertast.Typvariable:
		if gleich(v, suche) {
			return ersetzen
		}
		return v
	case getypisiertast.Typnutzung:
		return getypisiertast.Typnutzung{
			SymbolURL:            v.SymbolURL,
			Generischeargumenten: substituteList(v.Generischeargumenten, suche, ersetzen),
		}
	default:
		panic("e")
	}
}

func substituteExpression(expr getypisiertast.Expression, suche getypisiertast.ITyp, ersetzen getypisiertast.ITyp) getypisiertast.Expression {
	switch v := expr.(type) {
	case getypisiertast.Ganzzahl:
		return v
	case getypisiertast.Funktionsaufruf:
		for idx, arg := range v.Argumenten {
			v.Argumenten[idx] = substituteExpression(arg, suche, ersetzen)
		}
		if gleich(v.Rückgabetyp, suche) {
			v.Rückgabetyp = ersetzen
		}
		return v
	case getypisiertast.Variable:
		if gleich(v.ITyp, suche) {
			v.ITyp = ersetzen
		}
		return v
	default:
		panic("e")
	}
}

func synthVariantApplication(l *lokalekontext, s *scopes, typ getypisiertast.Typ, vari getypisiertast.Variant, url getypisiertast.SymbolURL, e ast.Expression, arg []ast.Expression) (getypisiertast.Expression, error) {
	if len(vari.Datenfelden) != len(arg) {
		return nil, neuFehler(e.Pos, "len nicht gleich")
	}

	vars := map[string]getypisiertast.ITyp{}
	var bs []getypisiertast.Expression

	for idx := range vari.Datenfelden {
		eingabe := arg[idx]

		b, f := synthExpression(l, s, eingabe)
		if f != nil {
			return nil, f
		}

		a := vari.Datenfelden[idx].Typ

		es := unify(b.Typ(), a)
		for k, v := range es {
			if _, ok := vars[k]; !ok {
				vars[k] = v
			}
			a = vars[k]
		}

		if !gleich(a, b.Typ()) {
			return nil, gleichErr(e.Pos, a, b.Typ())
		}

		bs = append(bs, b)
	}

	ret := typ
	for k, v := range vars {
		ret = substituteTypdekl(ret, getypisiertast.Typvariable{Name: k}, v)
	}

	var es getypisiertast.Typnutzung
	es.SymbolURL = url
	for _, k := range typ.Generischeargumenten {
		if v, ok := vars[k]; ok {
			es.Generischeargumenten = append(es.Generischeargumenten, v)
		} else {
			// TODO: sus
			es.Generischeargumenten = append(es.Generischeargumenten, getypisiertast.Typvariable{Name: k})
		}
	}

	return getypisiertast.Variantaufruf{
		Variant:    url,
		Argumenten: bs,
		Varianttyp: es,

		LPos: e.Pos,
	}, nil
}

func synthApplication(l *lokalekontext, s *scopes, url getypisiertast.SymbolURL, sig getypisiertast.Funktionssignatur, e ast.Expression, arg []ast.Expression) (getypisiertast.Expression, error) {
	if len(sig.Formvariabeln) != len(arg) {
		return nil, neuFehler(e.Pos, "len nicht gleich")
	}

	vars := map[string]getypisiertast.ITyp{}
	var bs []getypisiertast.Expression

	for idx := range sig.Formvariabeln {
		eingabe := arg[idx]

		b, f := synthExpression(l, s, eingabe)
		if f != nil {
			return nil, f
		}

		a := sig.Formvariabeln[idx].Typ

		es := unify(b.Typ(), a)
		for k, v := range es {
			if _, ok := vars[k]; !ok {
				vars[k] = v
			}
			a = vars[k]
		}

		if !gleich(a, b.Typ()) {
			return nil, gleichErr(e.Pos, a, b.Typ())
		}

		bs = append(bs, b)
	}

	ret := sig.Rückgabetyp
	for k, v := range vars {
		ret = substitute(ret, getypisiertast.Typvariable{Name: k}, v)
	}

	if len(vars) > 0 {
		rumpf, feh := l.sucheFunktionsrumpf(url, e.Pos)
		if feh != nil {
			return nil, feh
		}
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

	return getypisiertast.Funktionsaufruf{
		Funktion:    url,
		Argumenten:  bs,
		Rückgabetyp: ret,
	}, nil
}
