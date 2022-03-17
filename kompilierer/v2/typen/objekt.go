package typen

import "github.com/alecthomas/participle/v2/lexer"

type Objekt interface {
	Name() string
	Paket() string
	Pos() lexer.Position
	Typ() Typ

	Sichtbarkeitsbereich() *Sichtbarkeitsbereich
}

type objekt struct {
	sichtbarkeitsbereich *Sichtbarkeitsbereich
	name                 string
	paket                string
	pos                  lexer.Position
	typ                  Typ
}

func (o *objekt) Name() string {
	return o.name
}

func (o *objekt) Paket() string {
	return o.paket
}

func (o *objekt) Pos() lexer.Position {
	return o.pos
}

func (o *objekt) Typ() Typ {
	return o.typ
}

func (o *objekt) Sichtbarkeitsbereich() *Sichtbarkeitsbereich {
	return o.sichtbarkeitsbereich
}

type Funktion struct {
	objekt
}

type Variable struct {
	objekt
}
