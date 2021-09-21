package typisierung

import (
	"Tawa/parser"
	"Tawa/typen"
	"errors"
	"reflect"
)

var errMismatch = errors.New("mismatch")
var errLenMismatch = errors.New("len mismatch")
var errNichtGefunden = errors.New("nicht gefunden")
var errTypNichtGefunden = errors.New("typ nicht gefunden")

func gleich(a typen.Typ, b typen.Typ) bool {
	return reflect.DeepEqual(a, b)
}

func gleichErr(a typen.Typ, b typen.Typ) error {
	if !gleich(a, b) {
		return errMismatch
	}
	return nil
}

// typeof e == a
func checkExpression(ktx *kontext, expr *parser.Expression, gegenArt typen.Typ) error {
	if expr.Ganz != nil {
		return gleichErr(gegenArt, typen.Integer{})
	} else if expr.Block != nil {
		ruck, err := synthExpression(ktx, expr)
		if err != nil {
			return err
		}
		return gleichErr(ruck, gegenArt)
	} else if expr.Funktionsaufruf != nil {
		ruck, err := synthExpression(ktx, expr)
		if err != nil {
			return err
		}
		return gleichErr(ruck, gegenArt)
	} else if expr.Variable != nil {
		ruck, err := synthExpression(ktx, expr)
		if err != nil {
			return err
		}
		return gleichErr(ruck, gegenArt)
	}
	panic("feh checkExpression")
}

func synthExpression(ktx *kontext, expr *parser.Expression) (typen.Typ, error) {
	if expr.Ganz != nil {
		return typen.Integer{}, nil
	} else if expr.Block != nil {
		ktx.neuScope()
		var ruck typen.Typ = typen.Nichts{}
		for _, it := range expr.Block.Expressionen {
			t, err := synthExpression(ktx, it)
			if err != nil {
				return nil, err
			}
			ruck = t
		}
		ktx.loescheScope()

		return ruck, nil
	} else if expr.Funktionsaufruf != nil {
		funktionArt, ok := ktx.sucheFn(expr.Funktionsaufruf.Name)
		if !ok {
			return nil, errNichtGefunden
		}

		return synthApplication(ktx, funktionArt, expr.Funktionsaufruf.Argumente)
	} else if expr.Variable != nil {
		varArt, ok := ktx.sucheVar(*expr.Variable)
		if !ok {
			return nil, errNichtGefunden
		}

		return varArt, nil
	}
	panic("feh synthExpression")
}

func synthApplication(ktx *kontext, funk typen.Funktion, arg []parser.Expression) (typen.Typ, error) {
	if len(funk.Eingabe) != len(arg) {
		return nil, errLenMismatch
	}

	for idx := range funk.Eingabe {
		eingabe := arg[idx]

		b, f := synthExpression(ktx, &eingabe)
		if f != nil {
			return nil, f
		}

		a := funk.Eingabe[idx]
		if !gleich(a, b) {
			return nil, errMismatch
		}
	}

	ret := funk.Ausgabe
	return ret, nil
}

func typVonFunktion(ktx *kontext, it *parser.Funktion) (typen.Funktion, error) {
	var ausgabe typen.Typ
	var ok bool

	if it.Resultatart == nil {
		ausgabe = typen.Nichts{}
	} else {
		ausgabe, ok = ktx.typVonParser(it.Resultatart)
		if !ok {
			return typen.Funktion{}, errTypNichtGefunden
		}
	}

	if ausgabe == nil {
		panic("nil ausgabe")
	}

	var eingaben []typen.Typ

	for _, it := range it.Funktionsargumente {
		eingabe, ok := ktx.typVonParser(&it.Art)
		if !ok {
			return typen.Funktion{}, errTypNichtGefunden
		}
		eingaben = append(eingaben, eingabe)
	}

	return typen.Funktion{
		Eingabe: eingaben,
		Ausgabe: ausgabe,
	}, nil
}

func PrüfDatei(d *parser.Datei) error {
	ktx := neuKontext()

	for idx, it := range d.Funktionen {
		kind, err := typVonFunktion(ktx, &it)
		if err != nil {
			return err
		}

		ktx.head().fns[it.Name] = kind
		d.Funktionen[idx].CodeTyp = kind
	}

	for _, it := range d.Funktionen {
		ktx.neuScope()

		kind := it.CodeTyp.(typen.Funktion)

		for idx, v := range it.Funktionsargumente {
			ktx.head().vars[v.Name] = kind.Eingabe[idx]
		}

		err := checkExpression(ktx, &it.Expression, kind.Ausgabe)
		if err != nil {
			return err
		}

		ktx.loescheScope()
	}

	return nil
}
