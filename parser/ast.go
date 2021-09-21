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
	Typargumenten []*Art `("(" @@* ")")?`
}

type Funktionsargument struct {
	Name string `@Ident ":"`
	Art  Art    `@@`
}

type Funktion struct {
	Name               string              `"funk" @Ident "("`
	Funktionsargumente []Funktionsargument `( @@ ( "," @@ )* )? ")"`
	Resultatart        *Art                `(":" @@)?`
	Expression         Expression          `@@`

	CodeTyp typen.Typ

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

// type Postfix struct {
// 	Fieldoperator      *Fieldoperator      `@@ |`
// 	Zuweisungsoperator *Zuweisungsoperator `@@`

// 	Pos    lexer.Position
// 	EndPos lexer.Position
// }

type Block struct {
	Expressionen []*Expression `"tu" @@* "beende"`
}

type Funktionsaufruf struct {
	Name      string       `@Ident`
	Argumente []Expression `"(" ( @@ ( "," @@ )* )? ")"`
}

type Expression struct {
	Ganz            *int64           `@Int |`
	Funktionsaufruf *Funktionsaufruf `@@ |`
	Block           *Block           `@@`

	// Bedingung        *Bedingung        `(@@ |`
	// Definierung      *Definierung      `@@ |`
	// Zuweisung        *Zuweisung        `@@ |`
	// Logik            *Logik            `@@ |`
	// Cast             *Cast             `@@ |`
	// Integer          *Integer          `@@ |`
	// Löschen          *Löschen          `@@ |`
	// Neu              *Neu              `@@ |`
	// Stack            *Stack            `@@ |`
	// Dereferenzierung *Dereferenzierung `@@ |`
	// Variable         *string           `@Ident |`
	// Block            *Block            `@@)`
	// Fieldexpression      *Fieldexpression
	// Zuweisungsexpression *Zuweisungsexpression

	// Postfix []*Postfix `(@@*)?`

	Pos    lexer.Position
	EndPos lexer.Position
}

// func (d Expression) Vorverarbeiten() Expression {
// 	if d.Block != nil {
// 		for idx, it := range d.Block.Expr {
// 			d.Block.Expr[idx] = it.Vorverarbeiten()
// 		}
// 	}

// 	if len(d.Postfix) == 0 {
// 		return d
// 	}

// 	head := d.Postfix[0]
// 	elm := d.Postfix[1:]
// 	d.Postfix = []*Postfix{}

// 	if head.Fieldoperator != nil {
// 		return Expression{
// 			Fieldexpression: &Fieldexpression{
// 				Expr:  d,
// 				Field: head.Fieldoperator.Field,
// 			},
// 			Pos:     d.Pos,
// 			EndPos:  head.EndPos,
// 			Postfix: elm,
// 		}.Vorverarbeiten()
// 	} else if head.Zuweisungsoperator != nil {
// 		return Expression{
// 			Zuweisungsexpression: &Zuweisungsexpression{
// 				Links:  d,
// 				Rechts: head.Zuweisungsoperator.Wert,
// 			},
// 			Pos:     d.Pos,
// 			EndPos:  head.EndPos,
// 			Postfix: elm,
// 		}.Vorverarbeiten()
// 	} else {
// 		panic("e")
// 	}
// }

// func (d *Datei) Vorverarbeiten() {
// 	for idx, it := range d.Funktionen {
// 		e := it.Expression

// 		d.Funktionen[idx].Expression = e.Vorverarbeiten()
// 	}
// }

var (
	Parser = participle.MustBuild(&Datei{}, participle.UseLookahead(4))
)

func VonStringX(filename, content string) (r Datei) {
	err := Parser.ParseString(filename, content, &r)
	if err != nil {
		panic(err)
	}
	return
}
