package parser

import "github.com/alecthomas/participle/v2"

type Datei struct {
	Paket        string     `"paket" @String`
	Importierten []string   `(@String "ist" "importiert")*`
	Funktionen   []Funktion `@@*`
}

type Art struct {
	Normal *string `@Ident`
}

type Funktion struct {
	Name                string `"funk" @Ident "("`
	Funktionsargumenten []struct {
		Name string `@Ident ":"`
		Art  Art    `@@`
	} `@@* ")"`
	Resultatart *Art       `(":" @@)?`
	Expression  Expression `@@`
}

type Expression struct {
	Bedingung *struct {
		Falls     Expression  `"falls" @@`
		Werden    Expression  `@@`
		WennNicht *Expression `("wenn" "nicht" @@)?`
	} `@@ |`
	Variable *string      `@Ident |`
	Block    []Expression `("{" @@* "}")`
}

var (
	Parser = participle.MustBuild(&Datei{}, participle.UseLookahead(2))
)
