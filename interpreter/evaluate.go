package interpreter

import (
	"Tawa/parser"
	"Tawa/typisierung"
	"fmt"
	"strings"
)

func Evaluate(prog, filename string) string {
	es := parser.Datei{}
	feh := parser.Parser.ParseString(filename, prog, &es)
	if feh != nil {
		return feh.Error()
	}
	es.Vorverarbeiten()
	v := typisierung.NeuVollKontext()
	v.Push()
	err := typisierung.Typisierung(v, &es)
	if err != nil {
		s := strings.Builder{}
		if v, ok := err.(*typisierung.Fehler); ok {
			fmt.Fprintf(&s, "%s:%d:\n", filename, v.Line)
			fmt.Fprintln(&s, strings.Split(prog, "\n")[v.Line-1])
		}
		fmt.Fprintln(&s, err.Error())
		return s.String()
	}

	vk := NeuVollKontext()
	vk.Push()

	a, err := Interpret(es, "main", vk)
	if err != nil {
		return err.Error()
	}

	return a.AlsString()
}
