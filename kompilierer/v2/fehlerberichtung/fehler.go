package fehlerberichtung

import (
	"Tawa/kompilierer/v2/ast"
	"fmt"

	"github.com/alecthomas/participle/v2/lexer"
)

type Fehlercode int

const (
	KeinRückgabeErwartet Fehlercode = iota
	ArithmetikSeitenNichtGleichTyp
	GleichheitSeitenNichtGleichTyp
	VergleichSeitenNichtGleichTyp
	NichtErwarteTyp
	TypNichtGefunden
	VarNichtGefunden
	FeldNichtGefunden
	MüssenFallSein
	NichtGenugArgumente
	ZuVieleArgumente
	IstKeinTyp
	IstKeinVariable
	IstKeinFall
	HatKeineFälle
)

type Fehler struct {
	Fehlercode Fehlercode
	Anfang     lexer.Position
	Ende       lexer.Position
}

func String(f Fehlercode) string {
	switch f {
	case KeinRückgabeErwartet:
		return "Ich erwarte kein Rückgabe"
	case ArithmetikSeitenNichtGleichTyp:
		return "Beide Seiten eines Arithmetischer Expression müssen vom gleichen Typen sein"
	case GleichheitSeitenNichtGleichTyp:
		return "Beide Seiten eines Vergleischsexpression müssen vom gleichen Typen sein"
	case VergleichSeitenNichtGleichTyp:
		return "Beide Seiten eines Vergleischsexpression müssen vom gleichen Typen sein"
	case NichtErwarteTyp:
		return "Das typ hatte ich nicht erwartet"
	case TypNichtGefunden:
		return "Typ kann nicht gefunden"
	case VarNichtGefunden:
		return "Var kann nicht gefunden"
	case IstKeinTyp:
		return "Es ist kein Typ"
	case IstKeinVariable:
		return "Es ist kein Variable"
	case FeldNichtGefunden:
		return "Feld nicht gefunden"
	case IstKeinFall:
		return "Ist kein Fall"
	case HatKeineFälle:
		return "Typ hat keine Fälle"
	default:
		panic("a")
	}
}

func (f *Fehler) Error() string {
	return fmt.Sprintf("%s: %s", f.Anfang, String(f.Fehlercode))
}

func Neu(f Fehlercode, n ast.Node) *Fehler {
	return &Fehler{
		Fehlercode: f,
		Anfang:     n.Anfang(),
		Ende:       n.Ende(),
	}
}
