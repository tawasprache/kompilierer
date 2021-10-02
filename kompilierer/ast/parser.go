package ast

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type Modul struct {
	Package        string         `"paket" @Ident`
	Importierungen []Importierung `@@*`
	Deklarationen  []Deklaration  `@@*`
}

type Importierung struct {
	Pos lexer.Position

	Import string `"import" @String`
}

type Deklaration struct {
	Funktiondeklaration *Funktiondeklaration `(@@ |`
	Typdeklarationen    *Typdeklarationen    `@@)`
}

type Funktiondeklaration struct {
	Name                 string         `"funk" @Ident`
	Generischeargumenten []string       `("[" ( @Ident ( "," @Ident )* )? "]")?`
	Formvariabeln        []Formvariable `"(" ( @@ ( "," @@ )* )? ")"`
	Rückgabetyp          *Typ           `(":" @@)?`
	Expression           Expression     `(("=" ">") | ("-" ">"))? @@`
}

type Typdeklarationen struct {
	Name                 string      `"typ" @Ident`
	Generischeargumenten []string    `("[" ( @Ident ( "," @Ident )* )? "]")? "ist"`
	Datenfelden          []Datenfeld `@@*`
	Varianten            []Variante  `@@* "beende"`
}

type Formvariable struct {
	Name string `@Ident`
	Typ  Typ    `":" @@`
}

type Datenfeld struct {
	Name string `@Ident`
	Typ  Typ    `":" @@`
}

type Variante struct {
	Name        string      `"|" @Ident`
	Datenfelden []Datenfeld `("(" ( @@ ( "," @@ )* )? ")")?`
}

type Typ struct {
	Pos lexer.Position

	Typvariable    *Typvariable    `("§"@Ident) |`
	Typkonstruktor *Typkonstruktor `(@@)`
}

type Typvariable string
type Typkonstruktor struct {
	Name                 string `@Ident`
	Generischeargumenten []Typ  `("[" ( @@ ( "," @@ )* )? "]")?`
}

var (
	Parser = participle.MustBuild(&Modul{}, participle.UseLookahead(4))
)

func VonStringX(filename, content string) (r Modul) {
	err := Parser.ParseString(filename, content, &r)
	if err != nil {
		panic(err)
	}
	return
}
