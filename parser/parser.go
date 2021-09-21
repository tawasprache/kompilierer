package parser

import (
	"Tawa/typen"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type Datei struct {
	Paket            string           `"paket" @String`
	Importierten     []string         `(@String "ist" "importiert")*`
	Typdeklarationen []Typdeklaration `@@*`
	Funktionen       []Funktion       `@@*`
}

type Strukturfield struct {
	Name string `@Ident`
	Art  Art    `@@`
}

type Struktur struct {
	Fields []Strukturfield `"struktur" "(" @@* ")"`
}

type Zeiger struct {
	Auf *Art `"zeiger" "auf" @@`
}

type Entwedersfall struct {
	Name string `@Ident`
	Von  *Art   `("von" @@)?`
}

type Entweder struct {
	Fallen []Entwedersfall `"entweder" ( @@ ( "oder" @@ )* )?`
}

type Normalart struct {
	Name          string `@Ident`
	Typargumenten []Art  `("[" ( @@ ( "," @@ )* )? "]")?`
}

type Art struct {
	Struktur *Struktur  `@@ |`
	Zeiger   *Zeiger    `@@ |`
	Entweder *Entweder  `@@ |`
	Normal   *Normalart `@@`

	Pos    lexer.Position
	EndPos lexer.Position
}

type Funktionsargument struct {
	Name string `@Ident ":"`
	Art  Art    `@@`

	TypenArt typen.Art
}

type Funktion struct {
	Name               string              `"funk" @Ident "("`
	Funktionsargumente []Funktionsargument `( @@ ( "," @@ )* )? ")"`
	Resultatart        *Art                `(":" @@)?`
	Expression         Expression          `@@`

	Art    typen.Art
	Pos    lexer.Position
	EndPos lexer.Position
}

type Typdeklaration struct {
	Name          string   `"typ" @Ident`
	Typargumenten []string `((?! "ist") @Ident)* "ist"`
	Art           Art      `@@`

	CodeArt typen.Art
}

type Bedingung struct {
	Wenn   Expression  `"wenn" @@`
	Werden Expression  `@@`
	Sonst  *Expression `("sonst" @@)?`

	Art typen.Art
}

type Definierung struct {
	Variable string     `@Ident`
	Art      *Art       `":" @@?`
	Wert     Expression `"=" @@`
}

type Zuweisung struct {
	Variable string     `@Ident`
	Wert     Expression `"=" @@`
}

type Funktionsaufruf struct {
	Name      string       `@Ident`
	Argumente []Expression `"(" ( @@ ( "," @@ )* )? ")"`
}

type Block struct {
	Expr []Expression `("(" @@* ")")`
}

type Integer struct {
	Value int64 `@Int`
}

type Logik struct {
	Wert string `@("Wahr" | "Falsch")`
}

type Cast struct {
	Von  Expression `"cast" @@`
	Nach Art        `"nach" @@`
}

type Strukturinitialisierungsfield struct {
	Name string     `@Ident "ist"`
	Wert Expression `@@`
}

type Strukturinitialisierung struct {
	Name   string                          `@Ident "("`
	Fields []Strukturinitialisierungsfield `@@* ")"`

	Pos    lexer.Position
	EndPos lexer.Position
}

type Neu struct {
	Expression *Expression `"neu" @@`
}

type Stack struct {
	Initialisierung Strukturinitialisierung `"stack" @@`
}

type Löschen struct {
	Expr Expression `"lösche" @@`
}

type Dereferenzierung struct {
	Expr Expression `"deref" @@`
}

type Fieldoperator struct {
	Field string `"." @Ident`
}

type Zuweisungsoperator struct {
	Wert Expression `"=" @@`
}

type Fieldexpression struct {
	Expr  Expression
	Field string

	FieldIndex int
}

type Zuweisungsexpression struct {
	Links  Expression
	Rechts Expression
}

type Postfix struct {
	Fieldoperator      *Fieldoperator      `@@ |`
	Zuweisungsoperator *Zuweisungsoperator `@@`

	Pos    lexer.Position
	EndPos lexer.Position
}

type Expression struct {
	Bedingung            *Bedingung        `(@@ |`
	Definierung          *Definierung      `@@ |`
	Zuweisung            *Zuweisung        `@@ |`
	Funktionsaufruf      *Funktionsaufruf  `@@ |`
	Logik                *Logik            `@@ |`
	Cast                 *Cast             `@@ |`
	Integer              *Integer          `@@ |`
	Löschen              *Löschen          `@@ |`
	Neu                  *Neu              `@@ |`
	Stack                *Stack            `@@ |`
	Dereferenzierung     *Dereferenzierung `@@ |`
	Variable             *string           `@Ident |`
	Block                *Block            `@@)`
	Fieldexpression      *Fieldexpression
	Zuweisungsexpression *Zuweisungsexpression

	Postfix []*Postfix `(@@*)?`

	Pos    lexer.Position
	EndPos lexer.Position

	Art typen.Art
}

func (d Expression) Vorverarbeiten() Expression {
	if d.Block != nil {
		for idx, it := range d.Block.Expr {
			d.Block.Expr[idx] = it.Vorverarbeiten()
		}
	}

	if len(d.Postfix) == 0 {
		return d
	}

	head := d.Postfix[0]
	elm := d.Postfix[1:]
	d.Postfix = []*Postfix{}

	if head.Fieldoperator != nil {
		return Expression{
			Fieldexpression: &Fieldexpression{
				Expr:  d,
				Field: head.Fieldoperator.Field,
			},
			Pos:     d.Pos,
			EndPos:  head.EndPos,
			Postfix: elm,
		}.Vorverarbeiten()
	} else if head.Zuweisungsoperator != nil {
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
