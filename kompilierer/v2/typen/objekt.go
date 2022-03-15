package typen

import "github.com/alecthomas/participle/v2/lexer"

type Objekt interface {
	Name() string
	Paket() string
	Pos() lexer.Position

	Sichtbarkeitsbereich() *Sichtbarkeitsbereich
}

type objekt struct {
	sichtbarkeitsbereich *Sichtbarkeitsbereich
	name                 string
	paket                string
	pos                  lexer.Position
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

func (o *objekt) Sichtbarkeitsbereich() *Sichtbarkeitsbereich {
	return o.sichtbarkeitsbereich
}

type Typname struct {
	objekt
}

type Funktion struct {
	objekt

	Argumenten  []*Typname
	Rückgabetyp *Typname
}

type Variable struct {
	objekt
}
