package ast

import "github.com/alecthomas/participle/v2/lexer"

type Expression struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Terminal *Terminal

	// oder

	Links  *Expression
	Op     *BinaryOperator
	Rechts *Expression
}

type Terminal struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Ganzzahl               *int                   `  @Int`
	Zeichenkette           *string                `| @String`
	Sei                    *Sei                   `| @@`
	Leiste                 *Leiste                `| @@`
	Nativ                  *Nativ                 `| @@`
	Passt                  *Passt                 `| @@`
	Variantaufruf          *Variantaufruf         `| @@`
	Funktionsaufruf        *Funktionsaufruf       `| @@`
	Strukturaktualisierung *Stukturaktualisierung `| @@`
	Variable               *string                `| @Ident`
}

type Sei struct {
	Variable string     `"sei" @Ident`
	Wert     Expression `"=" @@`
	In       Expression `"in" @@`
}

type Leiste struct {
	Expressionen []Expression `"[" @@ ( "," @@ )* "]"`
}

type Nativ struct {
	Typ  Typ           `"¤" @@`
	Code []Nativencode `@@*`
}

type Nativencode struct {
	Language string `"|" @Ident ":"`
	Code     string `(@String | @RawString)`
}

type Funktionsaufruf struct {
	Name      Symbolkette  `@@`
	Argumente []Expression `"(" ( @@ ( "," @@ )* )? ")"`
}

type Variantaufruf struct {
	Name           Symbolkette    `"#" @@`
	Argumente      []Expression   `("(" ( @@ ( "," @@ )* )? ")")?`
	Strukturfelden []Strukturfeld `("{" ( @@ ( "," @@ )* )? "}")?`
}

type Stukturaktualisierung struct {
	Struktur Expression     `"{" @@ "|"`
	Felden   []Strukturfeld `( @@ ( "," @@ )* ) "}"`
}

type Strukturfeld struct {
	Name string     `@Ident`
	Wert Expression `"=" @@`
}

type Passt struct {
	Wert    Expression `"passt" @@ "zu"`
	Mustern []Muster   `@@* "beende"`
}

type Muster struct {
	Pattern    Pattern    `@@`
	Expression Expression `"=>" @@`
}

type Pattern struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Name      Symbolkette `"#" @@`
	Variabeln []string    `("(" ( @Ident ( "," @Ident )* )? ")")?`
}
