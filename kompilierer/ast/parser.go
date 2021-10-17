package ast

import (
	"strings"
	"text/scanner"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type Modul struct {
	Package        string         `"paket" @Ident`
	Zeigen         Zeigen         `@@`
	Nativauftakt   *Nativauftakt  `@@?`
	Importierungen []Importierung `@@*`
	Deklarationen  []Deklaration  `@@*`
}

type Nativauftakt struct {
	Code []Nativencode `"¤" (@@*)`
}

type Zeigen struct {
	Symbolen *[]string `  "zeigt" "(" (@Ident)* ")"`
	Alles    *string   `| "zeigt" @"alles"`
	Nichts   *string   `| "zeigt" @"nichts"`
}

type Importierung struct {
	Pos lexer.Position

	Import      Symbolkette  `"import" @@`
	Als         *Symbolkette `("als" @@)?`
	Zeigen      []string     `(("zeigende" "(" @Ident* ")" )`
	ZeigenAlles *string      `| ("zeigende" @"alles"))?`
}

type Symbolkette struct {
	Symbolen []string `@Ident ( ":" @Ident )*`
}

func (s Symbolkette) String() string {
	return strings.Join(s.Symbolen, ":")
}

type Deklaration struct {
	Tokens []lexer.Token

	Funktiondeklaration *Funktiondeklaration `(@@ |`
	Typdeklarationen    *Typdeklarationen    `@@)`
}

func (d Deklaration) Comments() string {
	var s strings.Builder
	for _, it := range d.Tokens {
		if it.Type != scanner.Comment {
			return strings.TrimSpace(s.String())
		}

		if strings.HasPrefix(it.Value, "/*") {
			s.WriteString(it.Value[2 : len(it.Value)-2])
		} else {
			s.WriteString(it.Value[2:])
			s.WriteString("\n")
		}
	}
	return strings.TrimSpace(s.String())
}

type Funktiondeklaration struct {
	Name                 string         `"funk" @Ident`
	Generischeargumenten []string       `("[" ( @Ident ( "," @Ident )* )? "]")?`
	Formvariabeln        []Formvariable `"(" ( @@ ( "," @@ )* )? ")"`
	Rückgabetyp          *Typ           `(":" @@)?`
	Expression           Expression     `(("=>") | ("->"))? @@`
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
	Pos    lexer.Position
	EndPos lexer.Position

	Typvariable    *Typvariable    `  ("§"@Ident)`
	Typfunktion    *Typfunktion    `| (@@)`
	Typkonstruktor *Typkonstruktor `| (@@)`
}

type Typfunktion struct {
	Argumenten  []Typ `"funk" "(" ( @@ ( "," @@ )* )? ")"`
	Rückgabetyp *Typ  `(":" @@)?`
}
type Typvariable string
type Typkonstruktor struct {
	Name                 Symbolkette `@@`
	Generischeargumenten []Typ       `("[" ( @@ ( "," @@ )* )? "]")?`
}

var (
	Parser                 = participle.MustBuild(&Modul{}, participle.UseLookahead(4), participle.Lexer(&lexFac{}), participle.Elide("Comment"))
	TerminalParser         = participle.MustBuild(&Terminal{}, participle.UseLookahead(4), participle.Lexer(&lexFac{}), participle.Elide("Comment"))
	argLeisteParser        = participle.MustBuild(&Argumentleiste{}, participle.UseLookahead(4), participle.Lexer(&lexFac{}), participle.Elide("Comment"))
	typFunktionParser      = participle.MustBuild(&Typfunktion{}, participle.UseLookahead(4), participle.Lexer(&lexFac{}), participle.Elide("Comment"))
	funktionsLiteralParser = participle.MustBuild(&Funktionsliteral{}, participle.UseLookahead(4), participle.Lexer(&lexFac{}), participle.Elide("Comment"))
)

func VonStringX(filename, content string) (r Modul) {
	err := Parser.ParseString(filename, content, &r)
	if err != nil {
		panic(err)
	}
	return
}
