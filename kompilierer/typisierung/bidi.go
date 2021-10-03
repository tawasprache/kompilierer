package typisierung

import (
	"Tawa/kompilierer/getypisiertast"
	"fmt"

	"github.com/alecthomas/participle/v2/lexer"
	"github.com/alecthomas/repr"
)

func gleich(a getypisiertast.ITyp, b getypisiertast.ITyp) bool {
	lhv, lhsIstTypvar := a.(getypisiertast.Typvariable)
	rhv, rhsIstTypvar := b.(getypisiertast.Typvariable)

	if lhsIstTypvar && !rhsIstTypvar {
		return true
	} else if !lhsIstTypvar && rhsIstTypvar {
		return true
	} else if lhsIstTypvar && rhsIstTypvar {
		return lhv.Name == rhv.Name
	}

	nlhv := a.(getypisiertast.Typnutzung)
	nrhv := b.(getypisiertast.Typnutzung)

	if nlhv.SymbolURL != nrhv.SymbolURL {
		return false
	}

	if len(nlhv.Generischeargumenten) != len(nrhv.Generischeargumenten) {
		return false
	}

	for idx := range nlhv.Generischeargumenten {
		lh := nlhv.Generischeargumenten[idx]
		rh := nrhv.Generischeargumenten[idx]

		if !gleich(lh, rh) {
			return false
		}
	}

	return true
}

func gleichErr(p lexer.Position, art string, a getypisiertast.ITyp, b getypisiertast.ITyp) error {
	return neuFehler(p, "%s nicht gleich: »%s« »%s«", art, a, b)
}

func unify(lhs getypisiertast.ITyp, rhs getypisiertast.ITyp) (map[string]getypisiertast.ITyp, error) {
	if rhv, ok := rhs.(getypisiertast.Typvariable); ok {
		return map[string]getypisiertast.ITyp{
			rhv.Name: lhs,
		}, nil
	} else if lhv, ok := rhs.(getypisiertast.Typvariable); ok {
		return map[string]getypisiertast.ITyp{
			lhv.Name: rhs,
		}, nil
	} else {
		lhv := lhs.(getypisiertast.Typnutzung)
		rhv := rhs.(getypisiertast.Typnutzung)

		if len(lhv.Generischeargumenten) != len(rhv.Generischeargumenten) {
			return nil, fmt.Errorf("%s != %s", lhv, rhv)
		}

		r := map[string]getypisiertast.ITyp{}

		for idx := range lhv.Generischeargumenten {
			lhs2 := lhv.Generischeargumenten[idx]
			rhs2 := rhv.Generischeargumenten[idx]

			unified, feh := unify(lhs2, rhs2)
			if feh != nil {
				return nil, feh
			}

			for k, v := range unified {
				r[k] = v
			}
		}

		return r, nil
	}
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
	case getypisiertast.Nichtunifiziert:
		println("warning: ununified type being substituted")
		return v
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
		v.Rückgabetyp = substitute(v.Rückgabetyp, suche, ersetzen)
		return v
	case getypisiertast.Variable:
		v.ITyp = substitute(v.ITyp, suche, ersetzen)
		return v
	case getypisiertast.Variantaufruf:
		for idx, arg := range v.Argumenten {
			v.Argumenten[idx] = substituteExpression(arg, suche, ersetzen)
		}
		v.Varianttyp = substitute(v.Varianttyp, suche, ersetzen)
		return v
	case getypisiertast.Pattern:
		v.Wert = substituteExpression(v.Wert, suche, ersetzen)
		for idx, muster := range v.Mustern {
			v.Mustern[idx].Expression = substituteExpression(muster.Expression, suche, ersetzen)
		}
		return v
	default:
		panic("e " + repr.String(v))
	}
}