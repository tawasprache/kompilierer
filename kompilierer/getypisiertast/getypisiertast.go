package getypisiertast

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2/lexer"
)

type Modul struct {
	Name         string
	ZeigeAlles   bool
	Nativcode    map[string]string
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
	Dokumentation     string
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
	String() string
}

type Typvariable struct {
	Name string
}

func (t Typvariable) istTyp() {}

func (t Typvariable) String() string {
	return "§" + t.Name
}

type Nativ struct {
	Code map[string]string

	LTyp ITyp
	LPos Span
}

func (v Nativ) istExpression() {}
func (v Nativ) Typ() ITyp      { return v.LTyp }
func (v Nativ) Pos() Span      { return v.LPos }

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

type Typfunktion struct {
	Argumenten  []ITyp
	Rückgabetyp ITyp
}

func (t Typfunktion) istTyp() {}
func (t Typfunktion) String() string {
	var s []string
	for _, it := range t.Argumenten {
		s = append(s, it.String())
	}
	return "funk(" + strings.Join(s, ", ") + ")"
}

type Nichtunifiziert struct {
}

func (t Nichtunifiziert) istTyp() {}
func (t Nichtunifiziert) String() string {
	return ""
}

type SymbolURL struct {
	Paket string
	Name  string
}

func (s SymbolURL) String() string {
	return strings.ReplaceAll(s.Paket, "/", ":") + ":" + s.Name
}

type Typ struct {
	Dokumentation        string
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
	Pos() Span
}

type Sei struct {
	Name string
	Wert Expression

	In   Expression
	LPos Span
}

func (Sei) istExpression() {}
func (s Sei) Typ() ITyp {
	return s.In.Typ()
}
func (s Sei) Pos() Span {
	return s.LPos
}

type Span struct {
	Von lexer.Position
	Zu  lexer.Position
}

func NeuSpan(von, zu lexer.Position) Span {
	return Span{
		Von: von,
		Zu:  zu,
	}
}

type Ganzzahl struct {
	Wert int
	LPos Span
}

func (g Ganzzahl) Pos() Span      { return g.LPos }
func (g Ganzzahl) istExpression() {}
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
var TypEinheit = Typnutzung{
	SymbolURL: SymbolURL{
		Paket: "Tawa/Eingebaut",
		Name:  "Einheit",
	},
}

func TypLeiste(i ITyp) ITyp {
	return Typnutzung{
		SymbolURL: SymbolURL{
			Paket: "Tawa/Leiste",
			Name:  "Leiste",
		},
		Generischeargumenten: []ITyp{i},
	}
}

type Zeichenkette struct {
	Wert string
	LPos Span
}

func (g Zeichenkette) Pos() Span      { return g.LPos }
func (g Zeichenkette) istExpression() {}
func (g Zeichenkette) Typ() ITyp {
	return TypZeichenkette
}

type Variable struct {
	Name string
	ITyp ITyp
	LPos Span
}

func (v Variable) istExpression() {}
func (v Variable) Typ() ITyp      { return v.ITyp }
func (v Variable) Pos() Span      { return v.LPos }

type Funktionsaufruf struct {
	Funktion    SymbolURL
	Argumenten  []Expression
	Rückgabetyp ITyp

	LPos Span
}

func (v Funktionsaufruf) istExpression() {}
func (v Funktionsaufruf) Typ() ITyp      { return v.Rückgabetyp }
func (v Funktionsaufruf) Pos() Span      { return v.LPos }

type FunktionErsteKlasseAufruf struct {
	Funktion    Expression
	Argumenten  []Expression
	Rückgabetyp ITyp

	LPos Span
}

func (v FunktionErsteKlasseAufruf) istExpression() {}
func (v FunktionErsteKlasseAufruf) Typ() ITyp      { return v.Rückgabetyp }
func (v FunktionErsteKlasseAufruf) Pos() Span      { return v.LPos }

type Variantaufruf struct {
	Variant        SymbolURL
	Konstruktor    string
	Argumenten     []Expression
	Strukturfelden []Strukturfeld
	Varianttyp     ITyp

	LPos Span
}

func (v Variantaufruf) istExpression() {}
func (v Variantaufruf) Typ() ITyp      { return v.Varianttyp }
func (v Variantaufruf) Pos() Span      { return v.LPos }

type Strukturaktualisierung struct {
	Wert   Expression
	Felden []Strukturaktualisierungsfeld

	LPos Span
}

func (s Strukturaktualisierung) istExpression() {}
func (v Strukturaktualisierung) Typ() ITyp      { return v.Wert.Typ() }
func (s Strukturaktualisierung) Pos() Span      { return s.LPos }

type Strukturaktualisierungsfeld struct {
	Name string
	Wert Expression
}

type Pattern struct {
	Wert    Expression
	Mustern []Muster

	LTyp ITyp
	LPos Span
}

func (v Pattern) istExpression() {}
func (v Pattern) Typ() ITyp      { return v.LTyp }
func (v Pattern) Pos() Span      { return v.LPos }

type Leiste struct {
	Werte []Expression

	LTyp ITyp
	LPos Span
}

func (v Leiste) istExpression() {}
func (v Leiste) Typ() ITyp      { return v.LTyp }
func (v Leiste) Pos() Span      { return v.LPos }

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
	BinOpVerketten
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
	LPos Span
}

func (v ValBinaryOperator) istExpression() {}
func (v ValBinaryOperator) Typ() ITyp      { return v.LTyp }
func (v ValBinaryOperator) Pos() Span      { return v.LPos }

type LogikBinaryOperator struct {
	Links  Expression
	Rechts Expression
	Art    LogikBinOp

	LPos Span
}

func (v LogikBinaryOperator) istExpression() {}
func (v LogikBinaryOperator) Pos() Span      { return v.LPos }
func (v LogikBinaryOperator) Typ() ITyp {
	return TypLogik
}

type Feldzugriff struct {
	Links Expression
	Feld  string

	LTyp ITyp
	LPos Span
}

func (f Feldzugriff) istExpression() {}
func (f Feldzugriff) Typ() ITyp      { return f.LTyp }
func (f Feldzugriff) Pos() Span      { return f.LPos }

type Strukturfeld struct {
	Name string
	Wert Expression
}

type Funktionsliteral struct {
	Formvariabeln []Formvariable
	Expression    Expression

	LTyp Typfunktion
	LPos Span
}

func (f Funktionsliteral) istExpression() {}
func (f Funktionsliteral) Typ() ITyp      { return f.LTyp }
func (f Funktionsliteral) Pos() Span      { return f.LPos }
