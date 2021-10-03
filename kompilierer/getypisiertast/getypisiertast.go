package getypisiertast

import (
	"github.com/alecthomas/participle/v2/lexer"
)

type Modul struct {
	Name         string
	Dependencies []string
	Typen        []Typ
	Funktionen   []Funktion
}

type Funktion struct {
	SymbolURL         SymbolURL
	Funktionssignatur Funktionssignatur
	Expression        Expression
}

type Funktionssignatur struct {
	Generischeargumenten []string
	Formvariabeln        []Formvariable
	Rückgabetyp          ITyp
}

type Formvariable struct {
	Name string
	Typ  ITyp
}

type ITyp interface {
	istTyp()
}

type Typvariable struct {
	Name string
}

func (t Typvariable) istTyp() {}

type Typnutzung struct {
	SymbolURL            SymbolURL
	Generischeargumenten []ITyp
}

func (t Typnutzung) istTyp() {}

type Nichtunifiziert struct {
}

func (t Nichtunifiziert) istTyp() {}

type SymbolURL struct {
	Paket string
	Name  string
}

type Typ struct {
	SymbolURL            SymbolURL
	Generischeargumenten []string
	Datenfelden          []Datenfeld
	Varianten            []Variant
}

type Datenfeld struct {
	Name string
	Typ  ITyp
}

type Variant struct {
	Name        string
	Datenfelden []Datenfeld
}

type Expression interface {
	istExpression()
	Typ() ITyp
	Pos() lexer.Position
}

type Ganzzahl struct {
	Wert int
	LPos lexer.Position
}

func (g Ganzzahl) Pos() lexer.Position { return g.LPos }
func (g Ganzzahl) istExpression()      {}
func (g Ganzzahl) Typ() ITyp {
	return Typnutzung{
		SymbolURL: SymbolURL{
			Paket: "Tawa/Eingebaut",
			Name:  "Ganz",
		},
	}
}

type Zeichenkette struct {
	Wert string
	LPos lexer.Position
}

func (g Zeichenkette) Pos() lexer.Position { return g.LPos }
func (g Zeichenkette) istExpression()      {}
func (g Zeichenkette) Typ() ITyp {
	return Typnutzung{
		SymbolURL: SymbolURL{
			Paket: "Tawa/Eingebaut",
			Name:  "Zeichenkette",
		},
	}
}

type Variable struct {
	Name string
	ITyp ITyp
	LPos lexer.Position
}

func (v Variable) istExpression()      {}
func (v Variable) Typ() ITyp           { return v.ITyp }
func (v Variable) Pos() lexer.Position { return v.LPos }

type Funktionsaufruf struct {
	Funktion    SymbolURL
	Argumenten  []Expression
	Rückgabetyp ITyp

	LPos lexer.Position
}

func (v Funktionsaufruf) istExpression()      {}
func (v Funktionsaufruf) Typ() ITyp           { return v.Rückgabetyp }
func (v Funktionsaufruf) Pos() lexer.Position { return v.LPos }

type Variantaufruf struct {
	Variant     SymbolURL
	Konstruktor string
	Argumenten  []Expression
	Varianttyp  ITyp

	LPos lexer.Position
}

func (v Variantaufruf) istExpression()      {}
func (v Variantaufruf) Typ() ITyp           { return v.Varianttyp }
func (v Variantaufruf) Pos() lexer.Position { return v.LPos }

type Pattern struct {
	Wert    Expression
	Mustern []Muster

	LTyp ITyp
	LPos lexer.Position
}

func (v Pattern) istExpression()      {}
func (v Pattern) Typ() ITyp           { return v.LTyp }
func (v Pattern) Pos() lexer.Position { return v.LPos }

type Muster struct {
	Variante    SymbolURL
	Konstruktor string

	Variablen []Mustervariable

	Expression Expression
}

type Mustervariable struct {
	Variante    SymbolURL
	Konstruktor string
	VonFeld     int

	Name string
}
