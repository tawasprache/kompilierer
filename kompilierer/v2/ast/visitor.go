package ast

import (
	"github.com/alecthomas/repr"
)

type Visitor interface {
	Visit(node Node) (w Visitor)
	EndVisit(node Node)
}

func Walk(v Visitor, n Node) {
	if n == nil {
		panic("nil node")
	}
	defer v.EndVisit(n)
	if v = v.Visit(n); v == nil {
		return
	}

	switch node := n.(type) {
	case *BinaryExpression:
		Walk(v, node.Links)
		Walk(v, node.Rechts)
	case *Datei:
		for _, it := range node.Deklarationen {
			Walk(v, it)
		}
	case *Funktiondeklaration:
		Walk(v, node.Name)
		Walk(v, node.Argumenten)
		if node.Rückgabetyp != nil {
			Walk(v, node.Rückgabetyp)
		}
		Walk(v, node.Inhalt)
	case *Typdeklaration:
		Walk(v, node.Name)
		for _, feld := range node.Felden {
			Walk(v, feld)
		}
	case *Feld:
		Walk(v, node.Name)
		Walk(v, node.Typ)
	case *Typkonstruktor:
		Walk(v, node.Ident)
	case *Block:
		for _, anweisung := range node.Anweisungen {
			Walk(v, anweisung)
		}
	case *Sei:
		Walk(v, node.Variable)
		Walk(v, node.Wert)
	case *Ist:
		Walk(v, node.Variable)
		Walk(v, node.Wert)
	case *Gib:
		if node.Wert != nil {
			Walk(v, node.Wert)
		}
	case *IdentExpression:
		Walk(v, node.Ident)
	case *SelektorExpression:
		Walk(v, node.Objekt)
		Walk(v, node.Feld)
	case *Argumentliste:
		for _, arg := range node.Argumente {
			Walk(v, arg)
		}
	case *Argument:
		Walk(v, node.Name)
		Walk(v, node.Typ)
	case *StrukturwertExpression:
		Walk(v, node.Name)
		for _, arg := range node.Argumente {
			Walk(v, arg)
		}
		for _, feld := range node.Felden {
			Walk(v, feld)
		}
	case *StrukturwertFeld:
		Walk(v, node.Name)
		Walk(v, node.Wert)
	case *MusterabgleichExpression:
		Walk(v, node.Wert)
		for _, muster := range node.Mustern {
			Walk(v, muster)
		}
	case *Muster:
		Walk(v, node.Pattern)
		Walk(v, node.Expression)
	case *Pattern:
		Walk(v, node.Name)
		for _, vari := range node.Variabeln {
			Walk(v, vari)
		}
	case *Ident, *GanzzahlExpression, *ZeichenketteExpression:
		return
	default:
		panic("eep " + repr.String(node))
	}
}

type inspector func(Node) bool

func (f inspector) EndVisit(node Node) {
}

func (f inspector) Visit(node Node) Visitor {
	if f(node) {
		return f
	}
	return nil
}

func Inspect(node Node, f func(Node) bool) {
	Walk(inspector(f), node)
}
