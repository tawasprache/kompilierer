package typisierung

import (
	"github.com/alecthomas/repr"

	"Tawa/parser"
	"Tawa/typen"
)

func artVonParser(v *VollKontext, a *parser.Art) (typen.Art, error) {
	if a == nil {
		return typen.Nichts{}, nil
	}

	if a.Normal != nil {
		we, ok := v.LookupArt(*a.Normal)
		if !ok {
			return nil, NeuFehler(a.Pos, "Typ »%s« nicht definiert", *a.Normal)
		}
		return we, nil
	} else if a.Struktur != nil {
		s := typen.Struktur{}

		for _, f := range a.Struktur.Fields {
			t, feh := artVonParser(v, &f.Art)
			if feh != nil {
				return nil, feh
			}

			s.Fields = append(s.Fields, typen.Strukturfield{
				Name: f.Name,
				Typ:  t,
			})
		}

		return s, nil
	} else if a.Zeiger != nil {
		t, feh := artVonParser(v, a.Zeiger.Auf)
		if feh != nil {
			return nil, feh
		}

		return typen.Zeiger{Auf: t}, nil
	} else if a.Entweder != nil {
		f := typen.Entweder{Fallen: map[string]typen.Art{}}

		for _, it := range a.Entweder.Fallen {
			if it.Von == nil {
				f.Fallen[it.Name] = nil
			} else {
				t, feh := artVonParser(v, it.Von)
				if feh != nil {
					return nil, feh
				}

				f.Fallen[it.Name] = t
			}
		}

		return f, nil
	}

	panic("a " + repr.String(a))
}

func artVonStrukt(v *VollKontext, e *parser.Strukturinitialisierung) (a typen.Art, err error) {
	we, ok := v.LookupArt(e.Name)
	if !ok {
		return nil, NeuFehler(e.Pos, "typ »%s« nicht deklariert", e.Name)
	}
	nt, ok := we.(typen.Neutyp)
	if !ok {
		return nil, NeuFehler(e.Pos, "typ »%s« ist kein Struktur", e.Name)
	}
	strukt, ok := nt.Von.(typen.Struktur)
	if !ok {
		return nil, NeuFehler(e.Pos, "typ »%s« ist kein Struktur", e.Name)
	}

outer:
	for _, it := range e.Fields {
		for _, f := range strukt.Fields {
			if f.Name == it.Name {
				v, err := artVonExpression(v, &it.Wert)
				if err != nil {
					return nil, err
				}

				if !v.IstGleich(f.Typ) {
					return nil, NeuFehler(e.Pos, "field »%s« ist nicht »%s«, sondern »%s«", it.Name, v, f.Typ)
				}

				continue outer
			}
		}
		return nil, NeuFehler(e.Pos, "»%s« ist kein field von »%s«", it.Name, e.Name)
	}

	return we, nil
}

func artVonExpression(v *VollKontext, e *parser.Expression) (a typen.Art, err error) {
	defer func() {
		e.Art = a
	}()

	if e.Bedingung != nil {
		lArt, err := artVonExpression(v, &e.Bedingung.Wenn)
		if err != nil {
			return nil, err
		}

		if !lArt.IstGleich(logikArt) {
			return nil, NeuFehler(e.Pos, "nicht-logikwert verwendent als bedingung")
		}

		rArt, err := artVonExpression(v, &e.Bedingung.Werden)
		if err != nil {
			return nil, err
		}

		if e.Bedingung.Sonst != nil {
			sonstArt, err := artVonExpression(v, e.Bedingung.Sonst)
			if err != nil {
				return nil, err
			}

			if !rArt.IstGleich(sonstArt) {
				return nil, NeuFehler(e.Bedingung.Sonst.Pos, "alle branchen sind nicht gleich")
			}

			e.Bedingung.Art = rArt
			return rArt, nil
		}

		return typen.Nichts{}, nil
	} else if e.Definierung != nil {
		if _, ok := v.LookupVariable(e.Definierung.Variable); ok {
			return nil, NeuFehler(e.Pos, "redefinition von »%s«", e.Definierung.Variable)
		}
		lArt, err := artVonExpression(v, &e.Definierung.Wert)
		if err != nil {
			return nil, err
		}
		v.KontextStack.Top().Variabeln[e.Definierung.Variable] = lArt

		return lArt, nil
	} else if e.Zuweisung != nil {
		va, ok := v.LookupVariable(e.Zuweisung.Variable)
		if !ok {
			return nil, NeuFehler(e.Pos, "»%s« nicht deklariert", e.Definierung.Variable)
		}
		typ, err := artVonExpression(v, &e.Zuweisung.Wert)
		if err != nil {
			return nil, err
		}
		if !va.IstGleich(typ) {
			return nil, NeuFehler(e.Pos, "»%s« kann nicht als »%s« in Zuweisung genutzt werden", typ, va)
		}
		return va, nil
	} else if e.Variable != nil {
		typ, ok := v.LookupVariable(*e.Variable)
		if !ok {
			return nil, NeuFehler(e.Pos, "»%s« nicht deklariert", *e.Variable)
		}
		return typ, nil
	} else if e.Block != nil {
		var art typen.Art = typen.Nichts{}
		var err error

		v.Push()
		for _, it := range e.Block.Expr {
			art, err = artVonExpression(v, &it)
			if err != nil {
				return nil, err
			}
		}
		v.Pop()

		return art, nil
	} else if e.Funktionsaufruf != nil {
		f, ok := v.Funktionen[e.Funktionsaufruf.Name]
		if !ok {
			return nil, NeuFehler(e.Pos, "funktion »%s« nicht deklariert", e.Funktionsaufruf.Name)
		}

		if len(f.Argumente) != len(e.Funktionsaufruf.Argumente) {
			if len(f.Argumente) > len(e.Funktionsaufruf.Argumente) {
				return nil, NeuFehler(e.Pos, "zu viele Argumente für Funktion »%s«", f)
			} else if len(f.Argumente) < len(e.Funktionsaufruf.Argumente) {
				return nil, NeuFehler(e.Pos, "zu wenige Argumente für Funktion »%s«", f)
			}
		}

		for idx := range f.Argumente {
			typ, feh := artVonExpression(v, &e.Funktionsaufruf.Argumente[idx])
			if feh != nil {
				return nil, feh
			}
			if !typ.IstGleich(f.Argumente[idx]) {
				return nil, NeuFehler(e.Funktionsaufruf.Argumente[idx].Pos, "»%s« kann nicht als »%s« in Funktionsaufruf genutzt werden", typ, f.Argumente[idx])
			}
		}

		return f.Returntyp, nil
	} else if e.Integer != nil {
		return ganzArt, nil
	} else if e.Logik != nil {
		return logikArt, nil
	} else if e.Cast != nil {
		von, feh := artVonExpression(v, &e.Cast.Von)
		if feh != nil {
			return nil, feh
		}
		nach, feh := artVonParser(v, &e.Cast.Nach)
		if feh != nil {
			return nil, feh
		}
		if !(von.KannNach(nach) || nach.KannVon(von)) {
			return nil, NeuFehler(e.Pos, "»%s« kann nicht nach »%s« umgewandelt werden", von, nach)
		}
		return nach, nil
	} else if e.Löschen != nil {
		_, err := artVonExpression(v, &e.Löschen.Expr)
		if err != nil {
			return nil, err
		}
		return nichtsArt, nil
	} else if e.Stack != nil {
		return artVonStrukt(v, &e.Stack.Initialisierung)
	} else if e.Neu != nil {
		es, err := artVonExpression(v, e.Neu.Expression)
		if err != nil {
			return nil, err
		}
		return typen.Zeiger{Auf: es}, nil
	} else if e.Dereferenzierung != nil {
		es, err := artVonExpression(v, &e.Dereferenzierung.Expr)
		if err != nil {
			return nil, err
		}

		v, ok := es.(typen.Zeiger)
		if !ok {
			return nil, NeuFehler(e.Pos, "»%s« ist kein Zeiger", v)
		}

		return v.Auf, nil
	} else if e.Zuweisungsexpression != nil {
		es, err := artVonExpression(v, &e.Zuweisungsexpression.Links)
		if err != nil {
			return nil, err
		}

		r, err := artVonExpression(v, &e.Zuweisungsexpression.Rechts)
		if err != nil {
			return nil, err
		}

		if !es.IstGleich(r) {
			return nil, NeuFehler(e.Pos, "»%s« kann nicht als »%s« in Zuweisung genutzt werden", r, es)
		}

		return r, nil
	} else if e.Fieldexpression != nil {
		es, err := artVonExpression(v, &e.Fieldexpression.Expr)
		if err != nil {
			return nil, err
		}

		nt, ok := es.(typen.Neutyp)
		if !ok {
			return nil, NeuFehler(e.Pos, "»%s« ist kein Struktur", es)
		}
		v, ok := nt.Von.(typen.Struktur)
		if !ok {
			return nil, NeuFehler(e.Pos, "»%s« ist kein Struktur", es)
		}

		var f typen.Art
		var i int
		for idx, it := range v.Fields {
			if it.Name == e.Fieldexpression.Field {
				f = it.Typ
				i = idx
				break
			}
		}
		e.Fieldexpression.FieldIndex = i

		if f == nil {
			return nil, NeuFehler(e.Pos, "»%s« ist kein Field von »%s«", e.Fieldexpression.Field, v)
		}

		return f, nil
	} else {
		panic("e " + repr.String(e))
	}

	panic("a " + repr.String(e))
}

func Typisierung(v *VollKontext, d *parser.Datei) error {
	for idx, es := range d.Typdeklarationen {
		ctx := v.Push()
		for i, it := range es.Typargumenten {
			ctx.Arten[it] = typen.Typvariable{Name: it, Idx: i}
		}
		typ, feh := artVonParser(v, &es.Art)
		v.Pop()
		if feh != nil {
			return feh
		}
		d.Typdeklarationen[idx].CodeArt = typ
		v.Top().Arten[es.Name] = typen.Neutyp{
			Name: es.Name,
			Von:  typ,
		}
	}
	for idx, es := range d.Funktionen {
		t := typen.Funktion{}
		typ, feh := artVonParser(v, es.Resultatart)
		if feh != nil {
			return feh
		}
		t.Returntyp = typ
		for _, arg := range es.Funktionsargumente {
			typ, feh := artVonParser(v, &arg.Art)
			if feh != nil {
				return feh
			}
			t.Argumente = append(t.Argumente, typ)
		}

		v.Funktionen[es.Name] = t
		d.Funktionen[idx].Art = t
	}
	for _, fnk := range d.Funktionen {
		v.Push()

		for _, es := range fnk.Funktionsargumente {
			art, feh := artVonParser(v, &es.Art)
			if feh != nil {
				return feh
			}
			v.Top().Variabeln[es.Name] = art
		}

		art, feh := artVonExpression(v, &fnk.Expression)
		if feh != nil {
			return feh
		}

		if !art.IstGleich(fnk.Art.(typen.Funktion).Returntyp) && !fnk.Art.(typen.Funktion).Returntyp.IstGleich(typen.Nichts{}) {
			if art.IstGleich(typen.Nichts{}) {
				return NeuFehler(fnk.Pos, "kein Return-Anweisung mit Wert, in Funktion mit Rückgabetyp »%s«", fnk.Art.(typen.Funktion).Returntyp)
			} else {
				return NeuFehler(fnk.Expression.Pos, "ungültige Umwandlung von »%s« in »%s«", art, fnk.Art.(typen.Funktion).Returntyp)
			}
		}

		v.Pop()
	}
	return nil
}
