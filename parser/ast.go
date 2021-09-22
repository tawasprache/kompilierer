package parser

import (
	"Tawa/typen"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type Datei struct {
	Paket string `"paket" @Ident`
	//	Importierten     []string         `(@String "ist" "importiert")*`
	//	Typdeklarationen []Typdeklaration `@@*`
	Funktionen []Funktion `@@*`
}

type Art struct {
	Name          string `@Ident`
	Typargumenten []*Art `("[" @@* "]")?`

	Pos    lexer.Position
	EndPos lexer.Position
}

type Funktionsargument struct {
	Name string `@Ident ":"`
	Art  Art    `@@`
}

type Funktion struct {
	Name               string              `"funk" @Ident `
	Typargumenten      []string            `("[" @Ident* "]")?`
	Funktionsargumente []Funktionsargument `"(" ( @@ ( "," @@ )* )? ")"`
	Resultatart        *Art                `(":" @@)?`
	Expression         Expression          `("=" ">")? @@`

	CodeTyp                       typen.Typ
	MonomorphisierteTypargumenten map[string]typen.Typ

	// Pos    lexer.Position
	// EndPos lexer.Position
}

// type Fieldexpression struct {
// 	Expr  Expression
// 	Field string

// 	FieldIndex int
// }

// type Zuweisungsexpression struct {
// 	Links  Expression
// 	Rechts Expression
// }

type Zuweisungsoperator struct {
	Wert Expression `"=" @@`
}

type Postfix struct {
	//	Fieldoperator      *Fieldoperator      `@@ |`
	Zuweisungsoperator *Zuweisungsoperator `@@`

	Pos    lexer.Position
	EndPos lexer.Position
}

type Block struct {
	Expressionen []*Expression `"tu" @@* "beende"`
}

type Funktionsaufruf struct {
	Name      string       `@Ident`
	Argumente []Expression `"(" ( @@ ( "," @@ )* )? ")"`

	MonomorphisierteTyp           typen.Typ
	MonomorphisierteTypargumenten map[string]typen.Typ
}

type Definierung struct {
	Name string      `"lass" @Ident`
	Art  *Art        `":" @@?`
	Wert *Expression `"=" @@`
}

type Zuweisungsexpression struct {
	Links  Expression
	Rechts Expression
}

type Expression struct {
	Ganz            *int64           `@Int |`
	Funktionsaufruf *Funktionsaufruf `@@ |`
	Definierung     *Definierung     `@@ |`
	Block           *Block           `@@ |`
	Variable        *string          `(?! "beende") @Ident`

	Zuweisungsexpression *Zuweisungsexpression

	// Bedingung        *Bedingung        `(@@ |`
	// Zuweisung        *Zuweisung        `@@ |`
	// Logik            *Logik            `@@ |`
	// Cast             *Cast             `@@ |`
	// Integer          *Integer          `@@ |`
	// Löschen          *Löschen          `@@ |`
	// Neu              *Neu              `@@ |`
	// Stack            *Stack            `@@ |`
	// Dereferenzierung *Dereferenzierung `@@ |`
	// Fieldexpression      *Fieldexpression

	Postfix []*Postfix `(@@*)?`

	Pos    lexer.Position
	EndPos lexer.Position
}

func (d Expression) Vorverarbeiten() Expression {
	if d.Block != nil {
		for idx, it := range d.Block.Expressionen {
			*d.Block.Expressionen[idx] = it.Vorverarbeiten()
		}
	}

	if len(d.Postfix) == 0 {
		return d
	}

	head := d.Postfix[0]
	elm := d.Postfix[1:]
	d.Postfix = []*Postfix{}

	if head.Zuweisungsoperator != nil {
		return Expression{
			Zuweisungsexpression: &Zuweisungsexpression{
				Links:  d,
				Rechts: head.Zuweisungsoperator.Wert,
			},
			Pos:     d.Pos,
			EndPos:  head.EndPos,
			Postfix: elm,
		}.Vorverarbeiten()
	} else {
		panic("e")
	}
}

func (d *Datei) Vorverarbeiten() {
	for idx, it := range d.Funktionen {
		e := it.Expression

		d.Funktionen[idx].Expression = e.Vorverarbeiten()
	}
}

var (
	Parser = participle.MustBuild(&Datei{}, participle.UseLookahead(4))
)

func VonStringX(filename, content string) (r Datei) {
	err := Parser.ParseString(filename, content, &r)
	if err != nil {
		panic(err)
	}
	r.Vorverarbeiten()
	return
}
