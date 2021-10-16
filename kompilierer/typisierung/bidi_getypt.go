package typisierung

import (
	"Tawa/kompilierer/fehlerberichtung"
	"Tawa/kompilierer/getypisiertast"

	"github.com/alecthomas/repr"
)

// assert typeof e == a
func checkGetypisiertExpression(l *lokalekontext, s *scopes, expr getypisiertast.Expression, gegenTyp getypisiertast.ITyp) (getypisiertast.Expression, error) {
	ruck, err := synthGetypisiertExpression(l, s, expr)
	if err != nil {
		return nil, err
	}

	if !TypGleich(ruck.Typ(), gegenTyp) {
		return nil, fehlerberichtung.GleichErr(expr.Pos(), "check", ruck.Typ(), gegenTyp)
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
			return nil, fehlerberichtung.NeuFehler(e.Pos(), "variable »%s« nicht gefunden", e.Name)
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
				return nil, fehlerberichtung.NeuFehler(e.Pos(), " »%s« hat kein varianten", k.SymbolURL)
			}

			if len(typDekl.Varianten) != len(e.Mustern) {
				return nil, fehlerberichtung.NeuFehler(e.Pos(), "len nicht gleich")
			}

			suche := func(s string, p getypisiertast.Span) (getypisiertast.Variant, error) {
				for _, it := range typDekl.Varianten {
					if it.Name == s {
						return it, nil
					}
				}
				return getypisiertast.Variant{}, fehlerberichtung.NeuFehler(p, "variant »%s« existiert nicht", s)
			}

			varis := map[string]getypisiertast.Variant{}

			for _, it := range e.Mustern {
				vari, feh := suche(it.Konstruktor, it.Expression.Pos())
				if feh != nil {
					return nil, feh
				}

				if len(vari.Datenfelden) != len(it.Variablen) {
					return nil, fehlerberichtung.NeuFehler(it.Expression.Pos(), "variante »%s« hat %d variablen, aber du nutzt %d", vari.Name, len(vari.Datenfelden), len(it.Variablen))
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
					if !TypGleich(kind, expr.Typ()) {
						return nil, fehlerberichtung.NeuFehler(e.Pos(), "arme sind nicht gleich, erwartete %s und sah %s", kind, expr.Typ())
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

		if e.Art == getypisiertast.BinOpVerketten {
			if !TypGleich(links.Typ(), rechts.Typ()) {
				return nil, fehlerberichtung.GleichErr(e.Pos(), "verketten", links.Typ(), rechts.Typ())
			}

			linksIstLeiste := TypGleich(links.Typ(), getypisiertast.TypLeiste(getypisiertast.Typvariable{Name: "L"}))
			rechtsIstLeiste := TypGleich(rechts.Typ(), getypisiertast.TypLeiste(getypisiertast.Typvariable{Name: "L"}))

			if !linksIstLeiste {
				return nil, fehlerberichtung.NeuFehler(links.Pos(), "ist kein leiste")
			}
			if !rechtsIstLeiste {
				return nil, fehlerberichtung.NeuFehler(rechts.Pos(), "ist kein leiste")
			}

			return getypisiertast.ValBinaryOperator{
				Links:  links,
				Rechts: rechts,
				Art:    e.Art,
				LTyp:   links.Typ(),
				LPos:   e.LPos,
			}, nil
		}

		linksIstGanz := TypGleich(links.Typ(), getypisiertast.TypGanz)
		rechtsIstGanz := TypGleich(rechts.Typ(), getypisiertast.TypGanz)

		if !linksIstGanz {
			return nil, fehlerberichtung.GleichErr(links.Pos(), "term", links.Typ(), getypisiertast.TypGanz)
		}
		if !rechtsIstGanz {
			return nil, fehlerberichtung.GleichErr(rechts.Pos(), "term", rechts.Typ(), getypisiertast.TypGanz)
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

		if !TypGleich(links.Typ(), rechts.Typ()) {
			return nil, fehlerberichtung.GleichErr(e.Pos(), "vergleich", links.Typ(), rechts.Typ())
		}

		switch e.Art {
		case getypisiertast.BinOpWeniger:
			fallthrough
		case getypisiertast.BinOpWenigerGleich:
			fallthrough
		case getypisiertast.BinOpGrößer:
			fallthrough
		case getypisiertast.BinOpGrößerGleich:
			if !TypGleich(links.Typ(), getypisiertast.TypGanz) {
				return nil, fehlerberichtung.GleichErr(e.Pos(), "vergleich", links.Typ(), getypisiertast.TypGanz)
			}
		}

		return getypisiertast.LogikBinaryOperator{
			Links:  links,
			Rechts: rechts,
			Art:    e.Art,
			LPos:   e.Pos(),
		}, nil
	case getypisiertast.Strukturaktualisierung:
		links, feh := synthGetypisiertExpression(l, s, e.Wert)
		if feh != nil {
			return nil, feh
		}
		var felden []getypisiertast.Strukturaktualisierungsfeld
		switch k := links.Typ().(type) {
		case getypisiertast.Typnutzung:
			typ, feh := l.typDekl(k.SymbolURL, e.Pos())
			if feh != nil {
				panic(feh)
			}
			for idx := range k.Generischeargumenten {
				typ = substituteTypdekl(typ, getypisiertast.Typvariable{Name: typ.Generischeargumenten[idx]}, k.Generischeargumenten[idx])
			}
			if len(typ.Datenfelden) == 0 {
				return nil, fehlerberichtung.NeuFehler(e.Pos(), "%s ist kein struktur", typ.SymbolURL)
			}
		äußere:
			for _, usrFeld := range e.Felden {
				for _, typFeld := range typ.Datenfelden {
					if usrFeld.Name != typFeld.Name {
						continue
					}

					b, f := synthGetypisiertExpression(l, s, usrFeld.Wert)
					if f != nil {
						return nil, f
					}

					if !TypGleich(typFeld.Typ, b.Typ()) {
						return nil, fehlerberichtung.GleichErr(usrFeld.Wert.Pos(), "strukturfeld", typFeld.Typ, b.Typ())
					}

					felden = append(felden, getypisiertast.Strukturaktualisierungsfeld{
						Name: usrFeld.Name,
						Wert: b,
					})

					continue äußere
				}
				return nil, fehlerberichtung.NeuFehler(usrFeld.Wert.Pos(), "%s ist kein feld von %s", usrFeld.Name, typ.SymbolURL)
			}
			return getypisiertast.Strukturaktualisierung{
				Wert:   links,
				Felden: felden,
				LPos:   e.LPos,
			}, nil
		case getypisiertast.Typvariable:
			panic("idk")
		}
	case getypisiertast.Feldzugriff:
		links, feh := synthGetypisiertExpression(l, s, e.Links)
		if feh != nil {
			return nil, feh
		}
		switch k := links.Typ().(type) {
		case getypisiertast.Typnutzung:
			typ, feh := l.typDekl(k.SymbolURL, e.Pos())
			if feh != nil {
				panic(feh)
			}
			for idx := range k.Generischeargumenten {
				typ = substituteTypdekl(typ, getypisiertast.Typvariable{Name: typ.Generischeargumenten[idx]}, k.Generischeargumenten[idx])
			}

			for _, feld := range typ.Datenfelden {
				if feld.Name != e.Feld {
					continue
				}
				return getypisiertast.Feldzugriff{
					Links: links,
					Feld:  e.Feld,

					LTyp: feld.Typ,
					LPos: e.Pos(),
				}, nil
			}

			return nil, fehlerberichtung.NeuFehler(e.Pos(), "%s ist kein feld von %s", e.Feld, typ.SymbolURL)
		case getypisiertast.Typvariable:
			panic("idk")
		}
	case getypisiertast.Leiste:
		var (
			werte []getypisiertast.Expression
			typ   getypisiertast.ITyp
		)
		for _, it := range e.Werte {
			if typ == nil {
				wert, feh := synthGetypisiertExpression(l, s, it)
				if feh != nil {
					return nil, feh
				}
				werte = append(werte, wert)
				typ = wert.Typ()
			} else {
				wert, feh := checkGetypisiertExpression(l, s, it, typ)
				if feh != nil {
					return nil, feh
				}
				werte = append(werte, wert)
			}
		}
		if typ == nil {
			panic("e")
		}
		return getypisiertast.Leiste{
			Werte: werte,
			LTyp:  getypisiertast.TypLeiste(typ),
			LPos:  e.LPos,
		}, nil
	case getypisiertast.Nativ:
		return e, nil
	case getypisiertast.Sei:
		var neuer getypisiertast.Sei
		var feh error

		neuer.Name = e.Name
		neuer.Wert, feh = synthGetypisiertExpression(l, s, e.Wert)
		if feh != nil {
			return nil, feh
		}
		neuer.LPos = e.LPos

		s.neuScope()
		s.head().vars[neuer.Name] = neuer.Wert.Typ()
		neuer.In, feh = synthGetypisiertExpression(l, s, e.In)
		if feh != nil {
			return nil, feh
		}
		s.loescheScope()

		return neuer, nil
	}
	panic("unreachable " + repr.String(expr))
}

func substituteVars(pos getypisiertast.Span, vars map[string]getypisiertast.ITyp, substitutions map[string]getypisiertast.ITyp) error {
	for k, v := range substitutions {
		if _, ok := vars[k]; !ok {
			vars[k] = v
		} else {
			if TypGleich(vars[k], v) {
				return nil
			}
			return fehlerberichtung.NeuFehler(pos, "this wants %s to be %s, but %s is already %s", k, v, k, vars[k])
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
	if !ok && len(typDekl.Datenfelden) == 0 {
		panic("!ok")
	}
	if len(typDekl.Datenfelden) > len(aufruf.Strukturfelden) {
		return nil, fehlerberichtung.NeuFehler(aufruf.Pos(), "nicht genug felden")
	}

	aufruf = copy(aufruf).(getypisiertast.Variantaufruf)

	if len(typ.Datenfelden) != len(aufruf.Argumenten) {
		return nil, fehlerberichtung.NeuFehler(aufruf.Pos(), "len nicht gleich")
	}

	vars := map[string]getypisiertast.ITyp{}

	var felden []getypisiertast.Strukturfeld

äußere:
	for _, typFeld := range typDekl.Datenfelden {
		for _, usrFeld := range aufruf.Strukturfelden {
			if usrFeld.Name != typFeld.Name {
				continue
			}

			b, f := synthGetypisiertExpression(l, s, usrFeld.Wert)
			if f != nil {
				return nil, f
			}

			a := typFeld.Typ

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

			if !TypGleich(a, b.Typ()) {
				return nil, fehlerberichtung.GleichErr(aufruf.Pos(), "variante", a, b.Typ())
			}

			felden = append(felden, getypisiertast.Strukturfeld{
				Name: usrFeld.Name,
				Wert: b,
			})

			continue äußere
		}
		return nil, fehlerberichtung.NeuFehler(aufruf.LPos, "du hast %s vergessen", typFeld.Name)
	}

äußere2:
	for _, usrFeld := range aufruf.Strukturfelden {
		for _, typFeld := range typDekl.Datenfelden {
			if usrFeld.Name == typFeld.Name {
				continue äußere2
			}
		}
		return nil, fehlerberichtung.NeuFehler(aufruf.LPos, "%s ist kein feld von %s", usrFeld.Name, typDekl.SymbolURL)
	}

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

		if !TypGleich(a, b.Typ()) {
			return nil, fehlerberichtung.GleichErr(aufruf.Pos(), "variante", a, b.Typ())
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
		Variant:        typDekl.SymbolURL,
		Konstruktor:    typ.Name,
		Argumenten:     bs,
		Strukturfelden: felden,
		Varianttyp:     es,

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
		return nil, fehlerberichtung.NeuFehler(aufruf.Pos(), "len nicht gleich")
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

		if !TypGleich(a, b.Typ()) {
			return nil, fehlerberichtung.GleichErr(eingabe.Pos(), "funktionsaufruf", a, b.Typ())
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
