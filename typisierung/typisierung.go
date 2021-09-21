package typisierung

import (
	"Tawa/parser"
	"Tawa/typen"
	"errors"
	"fmt"
	"reflect"

	"github.com/alecthomas/participle/v2/lexer"
)

type ErrMismatch struct {
	A, B typen.Typ
	Pos  lexer.Position
}

func (e *ErrMismatch) Error() string {
	return fmt.Sprintf("%s: %s != %s", e.Pos, e.A, e.B)
}

func mismatchErr(a, b typen.Typ, pos lexer.Position) error {
	return &ErrMismatch{
		A:   a,
		B:   b,
		Pos: pos,
	}
}

func lenMismatchErr(pos lexer.Position) error {
	return fmt.Errorf("fehler[%s]: %+w", pos, errLenMismatch_)
}

func nichtGefundenErr(pos lexer.Position) error {
	return fmt.Errorf("fehler[%s]: %+w", pos, errNichtGefunden_)
}

func typNichtGefundenErr(pos lexer.Position) error {
	return fmt.Errorf("fehler[%s]: %+w", pos, errTypNichtGefunden_)
}

var errMismatch_ = errors.New("mismatch")
var errLenMismatch_ = errors.New("len mismatch")
var errNichtGefunden_ = errors.New("nicht gefunden")
var errTypNichtGefunden_ = errors.New("typ nicht gefunden")

func gleich(a typen.Typ, b typen.Typ) bool {
	return reflect.DeepEqual(a, b)
}

func gleichErr(a typen.Typ, b typen.Typ, pos lexer.Position) error {
	if !gleich(a, b) {
		return mismatchErr(a, b, pos)
	}
	return nil
}

// typeof e == a
func checkExpression(ktx *kontext, expr *parser.Expression, gegenArt typen.Typ) error {
	if expr.Ganz != nil {
		return gleichErr(gegenArt, typen.Integer{}, expr.Pos)
	} else if expr.Block != nil {
		ruck, err := synthExpression(ktx, expr)
		if err != nil {
			return err
		}
		return gleichErr(ruck, gegenArt, expr.Pos)
	} else if expr.Funktionsaufruf != nil {
		ruck, err := synthExpression(ktx, expr)
		if err != nil {
			return err
		}
		return gleichErr(ruck, gegenArt, expr.Pos)
	} else if expr.Variable != nil {
		ruck, err := synthExpression(ktx, expr)
		if err != nil {
			return err
		}
		return gleichErr(ruck, gegenArt, expr.Pos)
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
		funktionArt, ok := ktx.sucheFnTyp(expr.Funktionsaufruf.Name)
		if !ok {
			return nil, nichtGefundenErr(expr.Pos)
		}
		funktion, ok := ktx.sucheFn(expr.Funktionsaufruf.Name)
		if !ok {
			return nil, nichtGefundenErr(expr.Pos)
		}

		return synthApplication(ktx, funktionArt, funktion, expr, expr.Funktionsaufruf.Argumente)
	} else if expr.Variable != nil {
		varArt, ok := ktx.sucheVar(*expr.Variable)
		if !ok {
			return nil, nichtGefundenErr(expr.Pos)
		}

		return varArt, nil
	}
	panic("feh synthExpression")
}

func synthApplication(ktx *kontext, funkTyp typen.Funktion, funkTemp parser.Funktion, e *parser.Expression, arg []parser.Expression) (typen.Typ, error) {
	if len(funkTyp.Eingabe) != len(arg) {
		return nil, lenMismatchErr(e.Pos)
	}

	vars := map[string]typen.Typ{}

	for idx := range funkTyp.Eingabe {
		eingabe := arg[idx]

		b, f := synthExpression(ktx, &eingabe)
		if f != nil {
			return nil, f
		}

		a := funkTyp.Eingabe[idx]

		switch inT := a.(type) {
		case typen.Typevar:
			if _, ok := vars[inT.Name]; !ok {
				vars[inT.Name] = b
			}
			a = vars[inT.Name]
		}

		if !gleich(a, b) {
			return nil, mismatchErr(a, b, eingabe.Pos)
		}
	}

	ret := funkTyp.Ausgabe
	switch retT := ret.(type) {
	case typen.Typevar:
		if _, ok := vars[retT.Name]; !ok {
			panic("unreachable")
		}
		ret = vars[retT.Name]
	}

	if len(vars) > 0 {
		t, err := typVonFunktionMono(ktx, &funkTemp, vars)
		if err != nil {
			return nil, err
		}

		e.Funktionsaufruf.MonomorphisierteTyp = t
		e.Funktionsaufruf.MonomorphisierteTypargumenten = vars

		ktx.neuScope()

		for idx, v := range funkTemp.Funktionsargumente {
			ktx.head().vars[v.Name] = t.Eingabe[idx]
		}

		err = checkExpression(ktx, &funkTemp.Expression, t.Ausgabe)
		if err != nil {
			return nil, err
		}

		ktx.loescheScope()

		return t.Ausgabe, nil
	}

	return ret, nil
}

func typVonFunktionMono(ktx *kontext, it *parser.Funktion, vars map[string]typen.Typ) (typen.Funktion, error) {
	ktx.neuScope()
	defer ktx.loescheScope()

	for k, v := range vars {
		ktx.head().typs[k] = v
	}

	var ausgabe typen.Typ
	var ok bool

	if it.Resultatart == nil {
		ausgabe = typen.Nichts{}
	} else {
		ausgabe, ok = ktx.typVonParser(it.Resultatart)
		if !ok {
			return typen.Funktion{}, typNichtGefundenErr(it.Resultatart.Pos)
		}
	}

	if ausgabe == nil {
		panic("nil ausgabe")
	}

	var eingaben []typen.Typ

	for _, it := range it.Funktionsargumente {
		eingabe, ok := ktx.typVonParser(&it.Art)
		if !ok {
			return typen.Funktion{}, typNichtGefundenErr(it.Art.Pos)
		}
		eingaben = append(eingaben, eingabe)
	}

	return typen.Funktion{
		Eingabe: eingaben,
		Ausgabe: ausgabe,
	}, nil
}

func typVonFunktion(ktx *kontext, it *parser.Funktion) (typen.Funktion, error) {
	ktx.neuScope()
	defer ktx.loescheScope()

	for _, arg := range it.Typargumenten {
		ktx.head().typs[arg] = typen.Typevar{Name: arg}
	}

	var ausgabe typen.Typ
	var ok bool

	if it.Resultatart == nil {
		ausgabe = typen.Nichts{}
	} else {
		ausgabe, ok = ktx.typVonParser(it.Resultatart)
		if !ok {
			return typen.Funktion{}, typNichtGefundenErr(it.Resultatart.Pos)
		}
	}

	if ausgabe == nil {
		panic("nil ausgabe")
	}

	var eingaben []typen.Typ

	for _, it := range it.Funktionsargumente {
		eingabe, ok := ktx.typVonParser(&it.Art)
		if !ok {
			return typen.Funktion{}, typNichtGefundenErr(it.Art.Pos)
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

		ktx.head().fnTyps[it.Name] = kind
		ktx.head().fns[it.Name] = it
		d.Funktionen[idx].CodeTyp = kind
	}

	for _, it := range d.Funktionen {
		if len(it.Typargumenten) > 0 {
			continue
		}

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
