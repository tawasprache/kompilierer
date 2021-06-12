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

type Art struct {
	Normal *string `@Ident`

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
	Name string `"typ" @Ident "ist"`
	Art  Art    `@@`
}

type Bedingung struct {
	Wenn   Expression  `"wenn" @@`
	Werden Expression  `@@`
	Sonst  *Expression `("sonst" @@)?`
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
	Expr []Expression `("{" @@* "}")`
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

type Expression struct {
	Bedingung       *Bedingung       `@@ |`
	Definierung     *Definierung     `@@ |`
	Zuweisung       *Zuweisung       `@@ |`
	Funktionsaufruf *Funktionsaufruf `@@ |`
	Logik           *Logik           `@@ |`
	Cast            *Cast            `@@ |`
	Integer         *Integer         `@@ |`
	Variable        *string          `@Ident |`
	Block           *Block           `@@`

	Pos    lexer.Position
	EndPos lexer.Position

	Art typen.Art
}

var (
	Parser = participle.MustBuild(&Datei{}, participle.UseLookahead(2))
)
