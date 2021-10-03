package typisierung

import (
	"Tawa/kompilierer/getypisiertast"

	"github.com/alecthomas/participle/v2/lexer"
	"github.com/alecthomas/repr"
)

// assert typeof e == a
func checkGetypisiertExpression(l *lokalekontext, s *scopes, expr getypisiertast.Expression, gegenTyp getypisiertast.ITyp) (getypisiertast.Expression, error) {
	ruck, err := synthGetypisiertExpression(l, s, expr)
	if err != nil {
		return nil, err
	}

	if !gleich(ruck.Typ(), gegenTyp) {
		return nil, gleichErr(expr.Pos(), "check", ruck.Typ(), gegenTyp)
	}

	return ruck, nil
}

// typeof e
func synthGetypisiertExpression(l *lokalekontext, s *scopes, expr getypisiertast.Expression) (getypisiertast.Expression, error) {
	switch e := expr.(type) {
	case getypisiertast.Ganzzahl:
		return e, nil
	case getypisiertast.Zeichenkette:
		return e, nil
	case getypisiertast.Funktionsaufruf:
		rumpf, feh := l.funktionsrumpf(e.Funktion, expr.Pos())
		if feh != nil {
			return nil, feh
		}
		return synthGetypisiertApplication(l, s, rumpf, e)
	case getypisiertast.Variable:
		typ, ok := s.suche(e.Name)
		if !ok {
			return nil, neuFehler(e.Pos(), "variable »%s« nicht gefunden", e.Name)
		}
		return getypisiertast.Variable{
			Name: e.Name,
			ITyp: typ,
		}, nil
	case getypisiertast.Variantaufruf:
		return synthGetypisiertVariantApplication(l, s, e)
	case getypisiertast.Pattern:
		wert, feh := synthGetypisiertExpression(l, s, e.Wert)
		if feh != nil {
			return nil, feh
		}
		switch k := wert.Typ().(type) {
		case getypisiertast.Typvariable:
			return getypisiertast.Pattern{
				Wert:    wert,
				Mustern: e.Mustern,
				LTyp:    e.LTyp,
				LPos:    e.LPos,
			}, nil
		case getypisiertast.Typnutzung:
			typDekl, feh := l.typDekl(k.SymbolURL, e.Pos())
			if feh != nil {
				return nil, feh
			}

			for idx, tvar := range typDekl.Generischeargumenten {
				typDekl = substituteTypdekl(typDekl, getypisiertast.Typvariable{Name: tvar}, k.Generischeargumenten[idx])
			}

			if len(typDekl.Varianten) < 2 {
				return nil, neuFehler(e.Pos(), " »%s« hat kein varianten", k.SymbolURL)
			}

			if len(typDekl.Varianten) != len(e.Mustern) {
				return nil, neuFehler(e.Pos(), "len nicht gleich")
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

			for _, it := range e.Mustern {
				vari, feh := suche(it.Konstruktor, it.Expression.Pos())
				if feh != nil {
					return nil, feh
				}

				if len(vari.Datenfelden) != len(it.Variablen) {
					return nil, neuFehler(it.Expression.Pos(), "variante »%s« hat %d variablen, aber du nutzt %d", vari.Name, len(vari.Datenfelden), len(it.Variablen))
				}

				varis[it.Konstruktor] = vari
			}

			var mustern []getypisiertast.Muster
			var kind getypisiertast.ITyp

			for _, it := range e.Mustern {
				s.neuScope()

				vari := varis[it.Konstruktor]
				for _, feld := range it.Variablen {
					s.head().vars[feld.Name] = vari.Datenfelden[feld.VonFeld].Typ
				}

				expr, feh := synthGetypisiertExpression(l, s, it.Expression)
				if feh != nil {
					return nil, feh
				}

				var mvars []getypisiertast.Mustervariable

				for _, it := range it.Variablen {
					mvars = append(mvars, getypisiertast.Mustervariable{
						Variante:    k.SymbolURL,
						Name:        it.Name,
						Konstruktor: it.Konstruktor,
						VonFeld:     it.VonFeld,
					})
				}

				mustern = append(mustern, getypisiertast.Muster{
					Konstruktor: it.Konstruktor,
					Variablen:   mvars,
					Expression:  expr,
				})

				if kind == nil {
					kind = expr.Typ()
				} else {
					if !gleich(kind, expr.Typ()) {
						repr.Println(kind)
						repr.Println(expr.Typ())
						return nil, neuFehler(e.Pos(), "arme sind nicht gleich, erwartete %s und sah %s", kind, expr.Typ())
					}
				}

				s.loescheScope()
			}

			return getypisiertast.Pattern{
				Wert:    wert,
				Mustern: mustern,

				LTyp: kind,
				LPos: e.Pos(),
			}, nil
		}
	case getypisiertast.ValBinaryOperator:
		links, feh := synthGetypisiertExpression(l, s, e.Links)
		if feh != nil {
			return nil, feh
		}
		rechts, feh := synthGetypisiertExpression(l, s, e.Rechts)
		if feh != nil {
			return nil, feh
		}

		linksIstGanz := gleich(links.Typ(), getypisiertast.TypGanz)
		rechtsIstGanz := gleich(rechts.Typ(), getypisiertast.TypGanz)

		if !linksIstGanz {
			return nil, gleichErr(links.Pos(), "term", links.Typ(), getypisiertast.TypGanz)
		}
		if !rechtsIstGanz {
			return nil, gleichErr(rechts.Pos(), "term", links.Typ(), getypisiertast.TypGanz)
		}

		return getypisiertast.ValBinaryOperator{
			Links:  links,
			Rechts: rechts,
			Art:    e.Art,
			LTyp:   getypisiertast.TypGanz,
			LPos:   e.LPos,
		}, nil
	case getypisiertast.LogikBinaryOperator:
		links, feh := synthGetypisiertExpression(l, s, e.Links)
		if feh != nil {
			return nil, feh
		}
		rechts, feh := synthGetypisiertExpression(l, s, e.Rechts)
		if feh != nil {
			return nil, feh
		}
		if !gleich(links.Typ(), rechts.Typ()) {
			return nil, gleichErr(e.Pos(), "vergleich", links.Typ(), rechts.Typ())
		}
		return getypisiertast.LogikBinaryOperator{
			Links:  links,
			Rechts: rechts,
			Art:    e.Art,
			LPos:   e.Pos(),
		}, nil
	}
	panic("unreachable " + repr.String(expr))
}

func substituteVars(pos lexer.Position, vars map[string]getypisiertast.ITyp, substitutions map[string]getypisiertast.ITyp) error {
	for k, v := range substitutions {
		if _, ok := vars[k]; !ok {
			vars[k] = v
		} else {
			if gleich(vars[k], v) {
				return nil
			}
			return neuFehler(pos, "this wants %s to be %s, but %s is already %s", k, v, k, vars[k])
		}
	}
	return nil
}

func synthGetypisiertVariantApplication(l *lokalekontext, s *scopes, aufruf getypisiertast.Variantaufruf) (getypisiertast.Expression, error) {
	typDekl, feh := l.typDekl(aufruf.Variant, aufruf.Pos())
	if feh != nil {
		panic(feh)
	}

	var typ getypisiertast.Variant
	var ok bool

	for _, it := range typDekl.Varianten {
		if it.Name == aufruf.Konstruktor {
			typ = it
			ok = true
			break
		}
	}
	if !ok {
		panic("!ok")
	}

	aufruf = copy(aufruf).(getypisiertast.Variantaufruf)

	if len(typ.Datenfelden) != len(aufruf.Argumenten) {
		return nil, neuFehler(aufruf.Pos(), "len nicht gleich")
	}

	vars := map[string]getypisiertast.ITyp{}
	var bs []getypisiertast.Expression

	for idx := range typ.Datenfelden {
		eingabe := aufruf.Argumenten[idx]

		b, f := synthGetypisiertExpression(l, s, eingabe)
		if f != nil {
			return nil, f
		}

		a := typ.Datenfelden[idx].Typ

		es, feh := unify(b.Typ(), a)
		if feh != nil {
			return nil, feh
		}
		feh = substituteVars(b.Pos(), vars, es)
		if feh != nil {
			return nil, feh
		}
		for k, v := range vars {
			a = substitute(a, getypisiertast.Typvariable{Name: k}, v)
		}

		if !gleich(a, b.Typ()) {
			return nil, gleichErr(aufruf.Pos(), "variante", a, b.Typ())
		}

		bs = append(bs, b)
	}

	ret := typDekl
	for k, v := range vars {
		ret = substituteTypdekl(ret, getypisiertast.Typvariable{Name: k}, v)
	}

	var es getypisiertast.Typnutzung
	es.SymbolURL = typDekl.SymbolURL
	for _, k := range typDekl.Generischeargumenten {
		if v, ok := vars[k]; ok {
			es.Generischeargumenten = append(es.Generischeargumenten, v)
		} else {
			// TODO: sus
			es.Generischeargumenten = append(es.Generischeargumenten, getypisiertast.Typvariable{Name: k})
		}
	}

	return getypisiertast.Variantaufruf{
		Variant:     typDekl.SymbolURL,
		Konstruktor: typ.Name,
		Argumenten:  bs,
		Varianttyp:  es,

		LPos: aufruf.Pos(),
	}, nil
}

func synthGetypisiertApplication(l *lokalekontext, s *scopes, funktion getypisiertast.Funktion, aufruf getypisiertast.Funktionsaufruf) (getypisiertast.Expression, error) {
	funktion = copy(funktion).(getypisiertast.Funktion)
	aufruf = copy(aufruf).(getypisiertast.Funktionsaufruf)

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

		es, feh := unify(b.Typ(), a)
		if feh != nil {
			return nil, feh
		}
		feh = substituteVars(b.Pos(), vars, es)
		if feh != nil {
			return nil, feh
		}
		for k, v := range vars {
			a = substitute(a, getypisiertast.Typvariable{Name: k}, v)
		}

		if !gleich(a, b.Typ()) {
			return nil, gleichErr(eingabe.Pos(), "funktionsaufruf", a, b.Typ())
		}

		arg[idx] = b
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

	return getypisiertast.Funktionsaufruf{
		Funktion:    aufruf.Funktion,
		Argumenten:  aufruf.Argumenten,
		Rückgabetyp: ret,
		LPos:        aufruf.Pos(),
	}, nil
}
