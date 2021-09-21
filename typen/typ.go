package typen

import "strings"

type Typ interface {
	istTyp()
	TypeString() string
	MangledString() string
}

type istTypImpl struct{}

func (i istTypImpl) istTyp() {}

type Funktion struct {
	Eingabe []Typ
	Ausgabe Typ

	istTypImpl
}

func (f Funktion) TypeString() string {
	var s strings.Builder

	var ss []string
	for _, it := range f.Eingabe {
		ss = append(ss, it.TypeString())
	}

	s.WriteString("funk(")
	s.WriteString(strings.Join(ss, ", "))
	s.WriteString(")")

	if f.Ausgabe != nil {
		s.WriteString(": ")
		s.WriteString(f.Ausgabe.TypeString())
	}

	return s.String()
}

func (f Funktion) MangledString() string {
	var s strings.Builder

	var ss []string
	for _, it := range f.Eingabe {
		ss = append(ss, it.MangledString())
	}

	s.WriteString("funk___")
	s.WriteString(strings.Join(ss, "__"))
	s.WriteString("___")

	if f.Ausgabe != nil {
		s.WriteString("____")
		s.WriteString(f.Ausgabe.MangledString())
	}

	return s.String()
}

type Integer struct {
	istTypImpl
}

func (i Integer) TypeString() string    { return "ganz" }
func (i Integer) MangledString() string { return "ganz" }

type Logik struct {
	istTypImpl
}

func (l Logik) TypeString() string    { return "logik" }
func (l Logik) MangledString() string { return "logik" }

type Nichts struct {
	istTypImpl
}

func (n Nichts) TypeString() string    { return "nichts" }
func (n Nichts) MangledString() string { return "nichts" }

type Typevar struct {
	Name string

	istTypImpl
}

func (t Typevar) TypeString() string { return t.Name }
func (t Typevar) MangledString() string {
	panic("you shouldn't be attempting to mangle something with a typevar still in it")
}
