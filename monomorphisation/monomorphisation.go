package monomorphisation

import (
	"Tawa/parser"
	"Tawa/typen"
	"strings"

	"github.com/jinzhu/copier"
)

func mangleName(s string, t typen.Funktion) string {
	var sb strings.Builder
	sb.WriteString(s)
	for _, it := range t.Eingabe {
		sb.WriteString("_")
		sb.WriteString(it.MangledString())
	}
	sb.WriteString("_")
	sb.WriteString(t.Ausgabe.MangledString())

	return sb.String()
}

func monomorphise(e *parser.Expression, fns map[string]parser.Funktion, genFunctions map[string]parser.Funktion) {
	if e.Block != nil {
		for _, it := range e.Block.Expressionen {
			monomorphise(it, fns, genFunctions)
		}
	} else if e.Funktionsaufruf != nil {
		if _, ok := genFunctions[e.Funktionsaufruf.Name]; !ok {
			return
		}

		name := mangleName(e.Funktionsaufruf.Name, e.Funktionsaufruf.MonomorphisierteTyp.(typen.Funktion))
		if _, ok := fns[name]; ok {
			return
		}

		gen := genFunctions[e.Funktionsaufruf.Name]
		gen.Name = name
		gen.CodeTyp = e.Funktionsaufruf.MonomorphisierteTyp
		gen.MonomorphisierteTypargumenten = e.Funktionsaufruf.MonomorphisierteTypargumenten
		monomorphise(&gen.Expression, fns, genFunctions)

		fns[name] = gen

		e.Funktionsaufruf.Name = name
	} else if e.Definierung != nil {
		monomorphise(e.Definierung.Wert, fns, genFunctions)
	} else if e.Zuweisungsexpression != nil {
		monomorphise(e.Definierung.Wert, fns, genFunctions)
	}
}

func Monomorphise(d parser.Datei) (r parser.Datei) {
	copier.Copy(&r, &d)

	genFunctions := map[string]parser.Funktion{}
	functions := map[string]parser.Funktion{}

	for _, it := range r.Funktionen {
		if len(it.Typargumenten) > 0 {
			genFunctions[it.Name] = it
		} else {
			functions[it.Name] = it
		}
	}

	for _, it := range r.Funktionen {
		if len(it.Typargumenten) > 0 {
			continue
		}

		monomorphise(&it.Expression, functions, genFunctions)
		functions[it.Name] = it
	}

	r.Funktionen = []parser.Funktion{}
	for _, it := range functions {
		r.Funktionen = append(r.Funktionen, it)
	}

	return r
}
