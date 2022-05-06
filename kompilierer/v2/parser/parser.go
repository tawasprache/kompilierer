package parser

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type Modul struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Name               Ident              `"paket" @@ EOS`
	Verwendungen       []Verwendung       `(@@ EOS)*`
	Moduldeklarationen []Moduldeklaration `(@@ EOS)*`
}

type Verwendung struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Paket string `"verwende" (@String | @RawString)`
	Als   *Ident `("als" @@)?`
}

type Moduldeklaration struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Typdeklaration      *Typdeklaration      `  @@`
	Funktiondeklaration *Funktiondeklaration `| @@`
}

type Typ struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Expression *Expression `@@`
	// Typvariable    *Typvariable    `  ("§"@Ident)`
	// Typfunktion    *Typfunktion    `| (@@)`
	// Typkonstruktor *Typkonstruktor `(@@)`
}

// type Typfunktion struct {
// 	Argumenten  []Typ `"funk" "(" ( @@ ( "," @@ )* )? ")"`
// 	Rückgabetyp *Typ  `(":" @@)?`
// }
// type Typvariable string
// type Typkonstruktor struct {
// 	Pos    lexer.Position
// 	EndPos lexer.Position

// 	Name Symbolkette `@@`
// 	// Generischeargumenten []Typ       `("[" ( @@ ( "," @@ )* )? "]")?`
// }

type Ident struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Name string `@Ident`
}

type Argument struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Namen []Ident `@@ ("," @@)*`
	Typ   Typ     `":" @@`
}

type Argumentliste struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Argumenten []Argument `( @@ ( "," @@ )* )?`
}

type Funktiondeklaration struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Name        Ident         `"funk" @@`
	Argumenten  Argumentliste `"(" @@ ")"`
	Rückgabetyp *Typ          `@@?`

	Inhalt Block `":" @@`
}

type Block struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Anweisungen []Anweisung `"beginne" EOS? (@@ EOS)* "beende"`
}

type Anweisung struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Gib *Gib `  @@`
	Sei *Sei `| @@`
	Ist *Ist `| @@`
}

type Gib struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Wert *Expression `"gib" @@?`
}

type Sei struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Name Ident      `"sei" @@`
	Wert Expression `":=" @@`
}

type Ist struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Name Ident      `@@`
	Wert Expression `":=" @@`
}

type Typdeklaration struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Name               Ident              `"typ" @@ "="`
	Verbunddeklaration Verbunddeklaration `@@`
}

type Verbunddeklaration struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Felden []Verbunddeklarationsfeld `"verbund" EOS? (@@ EOS)*`
	Fallen []Verbunddeklarationsfall `("mit" EOS? "fälle" EOS? (@@ EOS)*)? "beende"`
}

type Verbunddeklarationsfeld struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Name Ident `@@ ":"`
	Typ  Typ   `@@`
}

type Verbunddeklarationsfall struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Name   Ident                     `@@ ":"`
	Felden []Verbunddeklarationsfeld `"(" (@@ EOS)* ")"`
}

var (
	optionen = []participle.Option{
		participle.UseLookahead(4),
		participle.Lexer(&lexFac{}),
		participle.Elide("Comment"),
	}
	Parser         = participle.MustBuild(&Modul{}, optionen...)
	TerminalParser = participle.MustBuild(&Terminal{}, optionen...)
)
