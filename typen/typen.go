package typen

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/alecthomas/repr"
)

type Typargumenten struct {
	Argumenten *[]Typvariable
}

func (t Typargumenten) Typvariablen() []Typvariable {
	if *t.Argumenten == nil {
		*t.Argumenten = []Typvariable{}
	}
	return *t.Argumenten
}

func (t Typargumenten) SetTypvariablen(tv []Typvariable) {
	*t.Argumenten = tv
}

type keiner struct {
}

func (k keiner) Typvariablen() []Typvariable {
	return []Typvariable{}
}

func (k keiner) SetTypvariablen(t []Typvariable) {

}

type Art interface {
	String() string
	IstGleich(Art) bool
	KannVon(Art) bool
	KannNach(Art) bool
	Typvariablen() []Typvariable
	SetTypvariablen(t []Typvariable)
}

type kannVonNicht struct{}

func (kannVonNicht) KannVon(Art) bool  { return false }
func (kannVonNicht) KannNach(Art) bool { return false }

type Primitiv struct {
	Name string
	kannVonNicht
	keiner
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
	Typargumenten
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

type Nichts struct {
	kannVonNicht
	keiner
}

func (n Nichts) String() string {
	return "nichts"
}
func (n Nichts) IstGleich(a Art) bool {
	_, ok := a.(Nichts)
	return ok
}

type Logik struct {
	kannVonNicht
	keiner
}

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
	Typargumenten
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
	keiner
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
	keiner
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

type Typvariable struct {
	Name string
	Idx  int
	keiner
}

func (n Typvariable) String() string {
	return fmt.Sprintf("%s", n.Name)
}
func (n Typvariable) IstGleich(a Art) bool {
	return false
}
func (n Typvariable) KannVon(a Art) bool {
	return false
}
func (n Typvariable) KannNach(a Art) bool {
	return false
}

type Entweder struct {
	Fallen map[string]Art
	Typargumenten
}

func (n Entweder) String() string {
	var a strings.Builder
	a.WriteString("entweder ")

	i := false
	for k, v := range n.Fallen {
		if i {
			a.WriteString(" oder ")
		}
		if v == nil {
			a.WriteString(fmt.Sprintf("%s", k))
		} else {
			a.WriteString(fmt.Sprintf("%s von %s", k, v))
		}
		i = true
	}

	return a.String()
}

func (n Entweder) IstGleich(a Art) bool {
	z, ok := a.(Entweder)
	if !ok {
		return false
	}

	repr.Println(z.Fallen)
	repr.Println(n.Fallen)
	repr.Println(z.Fallen)

	return reflect.DeepEqual(n.Fallen, z.Fallen)
}
func (n Entweder) KannVon(a Art) bool {
	return n.IstGleich(a)
}
func (n Entweder) KannNach(a Art) bool {
	return n.IstGleich(a)
}
