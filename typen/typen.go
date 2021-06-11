package typen

import "strings"

type Art interface {
	String() string
	IstGleich(Art) bool
}

type Primitiv struct {
	Name string
}

func (p Primitiv) IstGleich(a Art) bool {
	v, ok := a.(Primitiv)
	if !ok {
		return false
	}

	return p.Name == v.Name
}

func (p Primitiv) String() string {
	return p.Name
}

type Funktion struct {
	Argumente []Art
	Returntyp Art
}

func (p Funktion) IstGleich(a Art) bool {
	v, ok := a.(Funktion)
	if !ok {
		return false
	}

	if !v.Returntyp.IstGleich(p.Returntyp) {
		return false
	}
	if len(v.Argumente) != len(p.Argumente) {
		return false
	}

	for idx := range p.Argumente {
		if !v.Argumente[idx].IstGleich(p.Argumente[idx]) {
			return false
		}
	}

	return true
}

func (p Funktion) String() string {
	var s strings.Builder
	s.WriteString("funk(")
	for idx, es := range p.Argumente {
		s.WriteString(es.String())
		if idx != len(p.Argumente)-1 {
			s.WriteString(", ")
		}
	}
	s.WriteString(")")
	return s.String()
}

type Nichts struct{}

func (n Nichts) String() string {
	return "nichts"
}
func (n Nichts) IstGleich(a Art) bool {
	_, ok := a.(Nichts)
	return ok
}

type Logik struct{}

func (n Logik) String() string {
	return "logik"
}
func (n Logik) IstGleich(a Art) bool {
	_, ok := a.(Logik)
	return ok
}
