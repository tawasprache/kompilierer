package parser

import "github.com/alecthomas/participle/v2/lexer"

type Expression struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Terminal *Terminal

	// oder

	Links  *Expression
	Op     *BinaryOperator
	Rechts *Expression

	// oder

	Objekt   *Expression
	Selektor *Ident
}

type Terminal struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Ganzzahl       *int            `  @Int`
	Musterabgleich *Musterabgleich `| @@`
	Variable       *Symbolkette    `| @@`
	Zeichenkette   *string         `| (@String | @RawString)`
	Strukturwert   *Strukturwert   `| @@`
}

type Strukturwert struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Name           Symbolkette    `"#" @@`
	Argumente      []Expression   `("(" ( @@ ( "," @@ )* )? ")")?`
	Strukturfelden []Strukturfeld `("{" ( @@ ( "," @@ )* )? "}")?`
}

type Strukturfeld struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Name Ident      `@@`
	Wert Expression `"=" @@`
}

type Musterabgleich struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Wert    Expression `"musterabgleich" @@ "mit" EOS?`
	Mustern []Muster   `@@* "beende"`
}

type Muster struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Pattern    Pattern    `@@`
	Expression Expression `"=>" @@ EOS`
}

type Pattern struct {
	Pos    lexer.Position
	EndPos lexer.Position

	Name      Symbolkette `"#" @@`
	Variabeln []Ident     `("(" ( @@ ( "," @@ )* )? ")")?`
}
