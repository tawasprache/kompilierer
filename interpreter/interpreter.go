package interpreter

import (
	"Tawa/parser"
	"Tawa/typen"

	"github.com/alecthomas/repr"
)

func interpretExpression(es parser.Datei, fn parser.Expression, mit *VollKontext) (wert, error) {
	if fn.Bedingung != nil {
		cond, err := interpretExpression(es, fn.Bedingung.Wenn, mit)
		if err != nil {
			return nil, err
		}
		if cond.(logik).val {
			return interpretExpression(es, fn.Bedingung.Wenn, mit)
		} else if fn.Bedingung.Sonst != nil {
			return interpretExpression(es, *fn.Bedingung.Sonst, mit)
		}
		return nil, nil
	} else if fn.Definierung != nil {
		w, err := interpretExpression(es, fn.Definierung.Wert, mit)
		if err != nil {
			return nil, err
		}
		mit.Top().Variabeln[fn.Definierung.Variable] = w
		return w, nil
	} else if fn.Zuweisung != nil {
		w, err := interpretExpression(es, fn.Zuweisung.Wert, mit)
		if err != nil {
			return nil, err
		}
		mit.Top().Variabeln[fn.Zuweisung.Variable] = w
		return w, nil
	} else if fn.Funktionsaufruf != nil {
		var exprs []wert
		for _, it := range fn.Funktionsaufruf.Argumente {
			w, err := interpretExpression(es, it, mit)
			if err != nil {
				return nil, err
			}
			exprs = append(exprs, w)
		}

		mit.Push()
		defer mit.Pop()

		for idx, it := range mit.Funktionen[fn.Funktionsaufruf.Name].Funktionsargumente {
			mit.Top().Variabeln[it.Name] = exprs[idx]
		}

		return interpretExpression(es, mit.Funktionen[fn.Funktionsaufruf.Name].Expression, mit)
	} else if fn.Logik != nil {
		return logik{val: fn.Logik.Wert == "Wahr"}, nil
	} else if fn.Cast != nil {
		return interpretExpression(es, fn.Cast.Von, mit)
	} else if fn.Integer != nil {
		return ganz{val: fn.Integer.Value}, nil
	} else if fn.Löschen != nil {
		w, err := interpretExpression(es, *&fn.Löschen.Expr, mit)
		if err != nil {
			return nil, err
		}

		*(w.(ptr).w) = nil
		return nil, nil
	} else if fn.Neu != nil {
		w, err := interpretExpression(es, *fn.Neu.Expression, mit)
		if err != nil {
			return nil, err
		}

		a := new(*wert)
		*a = &w

		return ptr{w: a}, nil
	} else if fn.Stack != nil {
		a, ok := mit.LookupArt(fn.Stack.Initialisierung.Name)
		if !ok {
			panic("not found")
		}
		b := a.(typen.Neutyp).Von.(typen.Struktur)

		userExprs := map[string]wert{}
		for _, it := range fn.Stack.Initialisierung.Fields {
			w, err := interpretExpression(es, it.Wert, mit)
			if err != nil {
				return nil, err
			}
			userExprs[it.Name] = w
		}

		exprs := map[string]wert{}
		for _, it := range b.Fields {
			if v, ok := userExprs[it.Name]; ok {
				exprs[it.Name] = v
			} else {
				exprs[it.Name] = nil
			}
		}

		return struktur{fields: exprs}, nil
	} else if fn.Variable != nil {
		v, ok := mit.LookupVariable(*fn.Variable)
		if !ok {
			panic("var nicht gefunden")
		}
		return v, nil
	} else if fn.Block != nil {
		mit.Push()
		defer mit.Pop()
		v := wert(nil)
		for _, it := range fn.Block.Expr {
			va, err := interpretExpression(es, it, mit)
			if err != nil {
				return nil, err
			}
			v = va
		}
		return v, nil
	} else if fn.Dereferenzierung != nil {
		v, feh := interpretExpression(es, fn.Dereferenzierung.Expr, mit)
		if feh != nil {
			return nil, feh
		}
		return **v.(ptr).w, nil
	}

	panic("e")
}

func artVonParser(v *VollKontext, a *parser.Art) (typen.Art, error) {
	if a == nil {
		return typen.Nichts{}, nil
	}

	if a.Normal != nil {
		we, ok := v.LookupArt(*a.Normal)
		if !ok {
			panic("a")
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
	}

	panic("a " + repr.String(a))
}

func Interpret(es parser.Datei, von string, mit *VollKontext) (wert, error) {
	for _, es := range es.Typdeklarationen {
		typ, feh := artVonParser(mit, &es.Art)
		if feh != nil {
			return nil, feh
		}
		mit.Top().Arten[es.Name] = typen.Neutyp{
			Name: es.Name,
			Von:  typ,
		}
	}
	for _, it := range es.Funktionen {
		mit.Funktionen[it.Name] = it
	}
	for _, it := range es.Funktionen {
		if it.Name == von {
			return interpretExpression(es, it.Expression, mit)
		}
	}
	panic("missing foo")
}
