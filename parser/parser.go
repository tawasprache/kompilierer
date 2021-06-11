package parser

import (
	"Tawa/typen"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type Datei struct {
	Paket        string     `"paket" @String`
	Importierten []string   `(@String "ist" "importiert")*`
	Funktionen   []Funktion `@@*`
}

type Art struct {
	Normal *string `@Ident`

	Pos    lexer.Position
	EndPos lexer.Position
}

type Funktionsargument struct {
	Name string `@Ident ":"`
	Art  Art    `@@`
}

type Funktion struct {
	Name                string              `"funk" @Ident "("`
	Funktionsargumenten []Funktionsargument `( @@ ( "," @@ )* )? ")"`
	Resultatart         *Art                `(":" @@)?`
	Expression          Expression          `@@`
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
type Expression struct {
	Bedingung   *Bedingung   `@@ |`
	Definierung *Definierung `@@ |`
	Zuweisung   *Zuweisung   `@@ |`
	Variable    *string      `@Ident |`
	Block       []Expression `("{" @@* "}")`

	Pos    lexer.Position
	EndPos lexer.Position

	Art typen.Art
}

var (
	Parser = participle.MustBuild(&Datei{}, participle.UseLookahead(2))
)
