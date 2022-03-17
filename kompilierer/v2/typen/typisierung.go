package typen

import (
	"Tawa/kompilierer/v2/ast"
	"Tawa/kompilierer/v2/fehlerberichtung"
	"Tawa/kompilierer/v2/parser"
)

type typisierung struct {
	fehler []error
	ktx    *Kontext
	s      *Sichtbarkeitsbereich
	f      *Funktion
}

func (t *typisierung) checkGetypisiertExpression(expr ast.Expression, mit Typ) Typ {
	switch expr := expr.(type) {
	default:
		ruck := t.synthGetypisiertExpression(expr)
		if ruck == nil {
			panic("e")
		}
		if !Gleich(mit, ruck) {
			t.fehler = append(t.fehler, fehlerberichtung.Neu(fehlerberichtung.NichtErwarteTyp, expr))
		}

		return mit
	}
}

func (t *typisierung) synthGetypisiertExpression(expr ast.Expression) Typ {
	switch expr := expr.(type) {
	case *ast.BinaryExpression:
		switch expr.Operator {
		case parser.BinOpAdd, parser.BinOpSub, parser.BinOpMul, parser.BinOpDiv, parser.BinOpPow, parser.BinOpMod:
			links := t.synthGetypisiertExpression(expr.Links)
			rechts := t.synthGetypisiertExpression(expr.Rechts)
			if !Gleich(links, rechts) {
				t.fehler = append(t.fehler, fehlerberichtung.Neu(fehlerberichtung.ArithmetikSeitenNichtGleichTyp, expr))
			}
			return links
		case parser.BinOpVerketten:
			zeichenkette := Welt.namen["Zeichenkette"].Typ()
			t.checkGetypisiertExpression(expr.Links, zeichenkette)
			t.checkGetypisiertExpression(expr.Rechts, zeichenkette)
			return zeichenkette
		case parser.BinOpGleich, parser.BinOpNichtGleich:
			links := t.synthGetypisiertExpression(expr.Links)
			rechts := t.synthGetypisiertExpression(expr.Rechts)
			if !Gleich(links, rechts) {
				t.fehler = append(t.fehler, fehlerberichtung.Neu(fehlerberichtung.GleichheitSeitenNichtGleichTyp, expr))
			}
			return Welt.namen["Wahrheitswert"].Typ()
		case parser.BinOpWeniger, parser.BinOpWenigerGleich, parser.BinOpGrößer, parser.BinOpGrößerGleich:
			links := t.synthGetypisiertExpression(expr.Links)
			rechts := t.synthGetypisiertExpression(expr.Rechts)
			if !Gleich(links, rechts) {
				t.fehler = append(t.fehler, fehlerberichtung.Neu(fehlerberichtung.VergleichSeitenNichtGleichTyp, expr))
			}
			return Welt.namen["Wahrheitswert"].Typ()
		default:
			panic("e")
		}
	case *ast.GanzzahlExpression:
		return Welt.namen["Ganzzahl"].Typ()
	case *ast.ZeichenketteExpression:
		return Welt.namen["Zeichenkette"].Typ()
	case *ast.IdentExpression:
		_, obj := t.s.Suchen(expr.Ident.String())
		if obj == nil {
			panic("nicht gefunden")
		}
		return obj.Typ()
	case *ast.SelektorExpression:
		panic("selektor")
	default:
		panic("e")
	}
}

func (t *typisierung) Visit(n ast.Node) ast.Visitor {
	if v, ok := t.ktx.Sichtbarkeitsbereichen[n]; ok {
		t.s = v
	}

	switch x := n.(type) {
	case *ast.Funktiondeklaration:
		t.f = t.ktx.Defs[x.Name].(*Funktion)
	case ast.Expression:
		t.synthGetypisiertExpression(x)
		return nil
	case *ast.Gib:
		if t.f.typ.(*Signature).Rückgabetyp == nil && x.Wert != nil {
			panic("kein gib")
		} else if t.f.typ.(*Signature).Rückgabetyp == nil && x.Wert == nil {
			return nil
		}
		t.checkGetypisiertExpression(x.Wert, t.f.typ.(*Signature).Rückgabetyp.Typ())
		return nil
	case *ast.Ist:
		_, varTyp := t.s.Suchen(x.Variable.String())
		t.checkGetypisiertExpression(x.Wert, varTyp.Typ())
		return nil
	case *ast.Sei:
		_, va := t.s.Suchen(x.Variable.String())
		va.(*Variable).typ = t.synthGetypisiertExpression(x.Wert)
		return nil
	}

	return t
}

func (t *typisierung) EndVisit(n ast.Node) {
	if v, ok := t.ktx.Sichtbarkeitsbereichen[n]; ok {
		t.s = v.übergeordneterSichtbarkeitsbereich
	}

	switch n.(type) {
	case *ast.Funktiondeklaration:
		t.f = nil
	}
}

func Typisierung(a *ast.Datei, ktx *Kontext) []error {
	k := &typisierung{nil, ktx, nil, nil}
	ast.Walk(k, a)
	return k.fehler
}
