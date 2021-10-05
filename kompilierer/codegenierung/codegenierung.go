package codegenierung

import (
	"Tawa/kompilierer/getypisiertast"
	"fmt"
	"strings"
)

type Optionen struct {
	Outpath     string
	JSOutfile   string
	HTMLOutfile string
	Entry       string
}

type Unterbau interface {
	Pregen(o Optionen) error
	CodegenModul(o Optionen, m getypisiertast.Modul) error
	Postgen(o Optionen) error
}

var unterbauen = map[string]Unterbau{}

func GetUnterbau(s string) Unterbau {
	return unterbauen[s]
}

func UnterbauRegistrieren(name string, u Unterbau) {
	unterbauen[name] = u
}

type Filebuilder struct {
	Einzug int

	strings.Builder
}

func (f *Filebuilder) Add(format string, a ...interface{}) {
	f.WriteString(strings.Repeat("\t", f.Einzug))
	f.WriteString(fmt.Sprintf(format, a...))
	f.WriteRune('\n')
}

func (f *Filebuilder) AddE(format string, a ...interface{}) {
	f.WriteString(strings.Repeat("\t", f.Einzug))
	f.WriteString(fmt.Sprintf(format, a...))
}

func (f *Filebuilder) AddK(format string, a ...interface{}) {
	f.WriteString(fmt.Sprintf(format, a...))
}

func (f *Filebuilder) AddNL() {
	f.WriteRune('\n')
}

func (f *Filebuilder) AddI(format string, a ...interface{}) {
	f.Add(format, a...)
	f.Einzug++
}

func (f *Filebuilder) AddD(format string, a ...interface{}) {
	f.Einzug--
	f.Add(format, a...)
}
