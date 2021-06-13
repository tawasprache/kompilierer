package typen

import (
	"fmt"
	"strings"
)

type Art interface {
	String() string
	IstGleich(Art) bool
	KannVon(Art) bool
	KannNach(Art) bool
}

type kannVonNicht struct{}

func (kannVonNicht) KannVon(Art) bool  { return false }
func (kannVonNicht) KannNach(Art) bool { return false }

type Primitiv struct {
	Name string
	kannVonNicht
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
	kannVonNicht
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

type Nichts struct{ kannVonNicht }

func (n Nichts) String() string {
	return "nichts"
}
func (n Nichts) IstGleich(a Art) bool {
	_, ok := a.(Nichts)
	return ok
}

type Logik struct{ kannVonNicht }

func (n Logik) String() string {
	return "logik"
}
func (n Logik) IstGleich(a Art) bool {
	_, ok := a.(Logik)
	return ok
}

type Neutyp struct {
	Name string
	Von  Art
}

func (n Neutyp) String() string {
	return n.Name
}
func (n Neutyp) IstGleich(a Art) bool {
	v, ok := a.(Neutyp)
	return ok && v.Name == n.Name
}
func (n Neutyp) KannVon(a Art) bool {
	switch v := n.Von.(type) {
	case Primitiv:
		t, ok := a.(Primitiv)
		if !ok {
			return false
		}
		return v.Name == t.Name
	case Logik:
		_, ok := a.(Logik)
		return ok
	}
	return false
}
func (n Neutyp) KannNach(a Art) bool {
	switch v := n.Von.(type) {
	case Primitiv:
		t, ok := a.(Primitiv)
		if !ok {
			return false
		}
		return v.Name == t.Name
	case Logik:
		_, ok := a.(Logik)
		return ok
	}
	return false
}

type Strukturfield struct {
	Name string
	Typ  Art
}
type Struktur struct {
	Fields []Strukturfield
}

func (n Struktur) String() string {
	var s strings.Builder
	s.WriteString("struktur {")
	for _, f := range n.Fields {
		s.WriteString(fmt.Sprintf("%s %s", f.Name, f.Typ))
	}
	s.WriteString("}")
	return s.String()
}
func (n Struktur) IstGleich(a Art) bool {
	an, ok := a.(Struktur)
	if !ok {
		return false
	}

	if len(an.Fields) != len(n.Fields) {
		return false
	}

	for idx := range an.Fields {
		if n.Fields[idx].Name != an.Fields[idx].Name {
			return false
		}
		if !n.Fields[idx].Typ.IstGleich(an.Fields[idx].Typ) {
			return false
		}
	}

	return true
}
func (n Struktur) KannVon(a Art) bool {
	return false
}
func (n Struktur) KannNach(a Art) bool {
	return false
}

type Zeiger struct {
	Auf Art
}

func (n Zeiger) String() string {
	return fmt.Sprintf("zeiger auf %s", n.Auf)
}
func (n Zeiger) IstGleich(a Art) bool {
	z, ok := a.(Zeiger)
	if !ok {
		return false
	}

	return z.Auf.IstGleich(n.Auf)
}
func (n Zeiger) KannVon(a Art) bool {
	return false
}
func (n Zeiger) KannNach(a Art) bool {
	return false
}
