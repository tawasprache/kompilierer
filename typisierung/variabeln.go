package typisierung

import (
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
	}
	panic("a")
}

func artVonExpression(v *VollKontext, e *parser.Expression) (typen.Art, error) {
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
		}

		return rArt, nil
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
		var art typen.Art
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
	}
	panic("a")
}

func Typisierung(v *VollKontext, d *parser.Datei) error {
	for _, es := range d.Funktionen {
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

		_ = art

		v.Pop()
	}
	return nil
}
