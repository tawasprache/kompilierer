package ast

import "github.com/alecthomas/participle/v2/lexer"

type Expression struct {
	Pos lexer.Position

	Terminal *Terminal

	// oder

	Links  *Expression
	Op     *BinaryOperator
	Rechts *Expression
}

type Terminal struct {
	Pos lexer.Position

	Ganzzahl        *int             `  @Int`
	Zeichenkette    *string          `| @String`
	Passt           *Passt           `| @@`
	Variantaufruf   *Variantaufruf   `| @@`
	Funktionsaufruf *Funktionsaufruf `| @@`
	Variable        *string          `| @Ident`
}

type Funktionsaufruf struct {
	Name      string       `@Ident`
	Argumente []Expression `"(" ( @@ ( "," @@ )* )? ")"`
}

type Variantaufruf struct {
	Name      string       `"#" @Ident`
	Argumente []Expression `("(" ( @@ ( "," @@ )* )? ")")?`
}

type Passt struct {
	Wert    Expression `"passt" @@ "zu"`
	Mustern []Muster   `@@* "beende"`
}

type Muster struct {
	Pattern    Pattern    `@@`
	Expression Expression `"=>" @@ "."`
}

type Pattern struct {
	Pos lexer.Position

	Name      string   `"#" @Ident`
	Variabeln []string `("(" ( @Ident ( "," @Ident )* )? ")")?`
}
