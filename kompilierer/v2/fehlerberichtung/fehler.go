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
	NichtGefunden
	IstKeinTyp
	IstKeinVariable
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
	case NichtGefunden:
		return "Es wird nicht gefunden"
	case IstKeinTyp:
		return "Es ist kein Typ"
	case IstKeinVariable:
		return "Es ist kein Variable"
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
