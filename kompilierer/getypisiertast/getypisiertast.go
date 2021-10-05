package getypisiertast

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2/lexer"
)

type Modul struct {
	Name         string
	ZeigeAlles   bool
	Zeigen       map[string]struct{}
	Dependencies []Dependency
	Typen        []Typ
	Funktionen   []Funktion
}

type Dependency struct {
	Paket      string
	Als        string
	ZeigeAlles bool
	Zeigen     []string
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
func (t Typnutzung) String() string {
	var s []string
	for _, it := range t.Generischeargumenten {
		s = append(s, fmt.Sprint(it))
	}
	if len(t.Generischeargumenten) > 0 {
		return fmt.Sprintf("%s[%s]", t.SymbolURL.String(), strings.Join(s, ","))
	} else {
		return t.SymbolURL.String()
	}
}

type Nichtunifiziert struct {
}

func (t Nichtunifiziert) istTyp() {}

type SymbolURL struct {
	Paket string
	Name  string
}

func (s SymbolURL) String() string {
	return strings.ReplaceAll(s.Paket, "/", ":") + ":" + s.Name
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
	return TypGanz
}

var TypGanz = Typnutzung{
	SymbolURL: SymbolURL{
		Paket: "Tawa/Eingebaut",
		Name:  "Ganz",
	},
}
var TypLogik = Typnutzung{
	SymbolURL: SymbolURL{
		Paket: "Tawa/Eingebaut",
		Name:  "Logik",
	},
}
var TypZeichenkette = Typnutzung{
	SymbolURL: SymbolURL{
		Paket: "Tawa/Eingebaut",
		Name:  "Zeichenkette",
	},
}

type Zeichenkette struct {
	Wert string
	LPos lexer.Position
}

func (g Zeichenkette) Pos() lexer.Position { return g.LPos }
func (g Zeichenkette) istExpression()      {}
func (g Zeichenkette) Typ() ITyp {
	return TypZeichenkette
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
	Variant        SymbolURL
	Konstruktor    string
	Argumenten     []Expression
	Strukturfelden []Strukturfeld
	Varianttyp     ITyp

	LPos lexer.Position
}

func (v Variantaufruf) istExpression()      {}
func (v Variantaufruf) Typ() ITyp           { return v.Varianttyp }
func (v Variantaufruf) Pos() lexer.Position { return v.LPos }

type Strukturaktualisierung struct {
	Wert   Expression
	Felden []Strukturaktualisierungsfeld

	LPos lexer.Position
}

func (s Strukturaktualisierung) istExpression()      {}
func (v Strukturaktualisierung) Typ() ITyp           { return v.Wert.Typ() }
func (s Strukturaktualisierung) Pos() lexer.Position { return s.LPos }

type Strukturaktualisierungsfeld struct {
	Name string
	Wert Expression
}

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

type ValBinOp int

const (
	_ ValBinOp = iota
	BinOpAdd
	BinOpSub
	BinOpMul
	BinOpDiv
	BinOpPow
	BinOpMod
)

type LogikBinOp int

const (
	_ LogikBinOp = iota
	BinOpGleich
	BinOpNichtGleich
	BinOpWeniger
	BinOpWenigerGleich
	BinOpGrößer
	BinOpGrößerGleich
)

type ValBinaryOperator struct {
	Links  Expression
	Rechts Expression
	Art    ValBinOp

	LTyp ITyp
	LPos lexer.Position
}

func (v ValBinaryOperator) istExpression()      {}
func (v ValBinaryOperator) Typ() ITyp           { return v.LTyp }
func (v ValBinaryOperator) Pos() lexer.Position { return v.LPos }

type LogikBinaryOperator struct {
	Links  Expression
	Rechts Expression
	Art    LogikBinOp

	LPos lexer.Position
}

func (v LogikBinaryOperator) istExpression()      {}
func (v LogikBinaryOperator) Pos() lexer.Position { return v.LPos }
func (v LogikBinaryOperator) Typ() ITyp {
	return TypLogik
}

type Feldzugriff struct {
	Links Expression
	Feld  string

	LTyp ITyp
	LPos lexer.Position
}

func (f Feldzugriff) istExpression()      {}
func (f Feldzugriff) Typ() ITyp           { return f.LTyp }
func (f Feldzugriff) Pos() lexer.Position { return f.LPos }

type Strukturfeld struct {
	Name string
	Wert Expression
}
