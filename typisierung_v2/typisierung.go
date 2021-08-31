package typisierungv2

import (
	"Tawa/parser"
	"errors"
	"fmt"
	"reflect"

	"github.com/alecthomas/repr"
)

type kontext struct {
	fns map[string]funktion
}

type typ interface{ istTyp() }

type istTypImpl struct{}

func (i istTypImpl) istTyp() {}

type funktion struct {
	eingabe []typ
	ausgabe typ

	istTypImpl
}

type nichts struct {
	istTypImpl
}

type integer struct {
	istTypImpl
}

type logik struct {
	istTypImpl
}

type entweder struct {
	fallen map[string]typ

	istTypImpl
}

type generischerTyp struct {
	von        typ
	argumenten []string

	istTypImpl
}

func (g generischerTyp) voll(args map[string]typ) typ {
	if len(args) != len(g.argumenten) {
		panic("ne enough")
	}
	for _, it := range g.argumenten {
		_, ok := args[it]
		if !ok {
			panic("ne enough")
		}
	}

	switch a := g.von.(type) {
	case funktion:
		b := a

		b.eingabe = make([]typ, len(a.eingabe))
		copy(b.eingabe, a.eingabe)

		for idx := range b.eingabe {
			v := b.eingabe[idx]
			if vv, ok := v.(kvar); ok {
				b.eingabe[idx] = args[vv.n]
			}
		}
		if vv, ok := b.ausgabe.(kvar); ok {
			b.ausgabe = args[vv.n]
		}

		return b
	case entweder:
		b := a
		b.fallen = map[string]typ{}
		for k, v := range a.fallen {
			b.fallen[k] = v
		}
		for it := range b.fallen {
			v := b.fallen[it]
			if vv, ok := v.(kvar); ok {
				b.fallen[it] = args[vv.n]
			}
		}
		return b
	default:
		return a
	}
}

type kvar struct {
	n string

	istTypImpl
}

var errMismatch = errors.New("mismatch")
var errLenMismatch = errors.New("len mismatch")
var errNichtGefunden = errors.New("nicht gefunden")

func gleich(a typ, b typ) bool {
	repr.Println(a)
	repr.Println(b)

	return reflect.DeepEqual(a, b)
}

func neuKontext() *kontext {
	return &kontext{
		fns: map[string]funktion{},
	}
}

// typeof e == a
func checkExpression(ktx *kontext, expr *parser.Expression, gegenArt typ) error {
	if expr.Integer != nil {
		if !gleich(gegenArt, integer{}) {
			return errMismatch
		}

		return nil
	} else if expr.Logik != nil {
		if !gleich(gegenArt, logik{}) {
			return errMismatch
		}

		return nil
	} else if expr.Funktionsaufruf != nil {
		funktionArt, ok := ktx.fns[expr.Funktionsaufruf.Name]
		if !ok {
			return errNichtGefunden
		}

		rückgabeArt, err := synthApplication(ktx, funktionArt, expr.Funktionsaufruf.Argumente)
		if err != nil {
			return err
		}

		fmt.Printf("got return type %s\n", repr.String(rückgabeArt))

		if !gleich(rückgabeArt, gegenArt) {
			return errMismatch
		}

		return nil
	}

	panic("feh checkExpression")
}

func synthExpression(ktx *kontext, expr *parser.Expression) (typ, error) {
	if expr.Integer != nil {
		return integer{}, nil
	} else if expr.Logik != nil {
		return logik{}, nil
	} else if expr.Funktionsaufruf != nil {
		funktionArt, ok := ktx.fns[expr.Funktionsaufruf.Name]
		if !ok {
			return nil, errNichtGefunden
		}

		return synthApplication(ktx, funktionArt, expr.Funktionsaufruf.Argumente)
	}

	panic("feh synthExpression")
}

func synthApplication(ktx *kontext, funk funktion, arg []parser.Expression) (typ, error) {
	if len(funk.eingabe) != len(arg) {
		return nil, errLenMismatch
	}

	vars := map[string]typ{}

	for idx := range funk.eingabe {
		a := funk.eingabe[idx]
		eingabe := arg[idx]

		b, f := synthExpression(ktx, &eingabe)
		if f != nil {
			return nil, f
		}

		switch inT := a.(type) {
		case kvar:
			if _, ok := vars[inT.n]; !ok {
				vars[inT.n] = b
			}
			a = vars[inT.n]
		case entweder:
			for fall, fallTyp := range inT.fallen {
				if fallKVar, ok := fallTyp.(kvar); ok {
					if _, ok := vars[fallKVar.n]; !ok {
						vars[fallKVar.n] = b.(entweder).fallen[fall]
					}

					inT.fallen[fall] = vars[fallKVar.n]
				}
			}
		}

		repr.Println(a)
		repr.Println(b)

		if !gleich(a, b) {
			return nil, errMismatch
		}
	}

	ret := funk.ausgabe
	switch retT := ret.(type) {
	case kvar:
		if _, ok := vars[retT.n]; !ok {
			panic("unreachable")
		}

		ret = vars[retT.n]
	case entweder:
		for fall, fallTyp := range retT.fallen {
			if fallKVar, ok := fallTyp.(kvar); ok {
				if _, ok := vars[fallKVar.n]; !ok {
					panic("unreachable")
				}

				retT.fallen[fall] = vars[fallKVar.n]
			}
		}
	}

	return ret, nil
}
