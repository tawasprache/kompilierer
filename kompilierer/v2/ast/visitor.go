package ast

import "reflect"

type Visitor interface {
	Visit(node Node) (w Visitor)
	EndVisit(node Node)
}

func Walk(v Visitor, n Node) {
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
	case *Ident, *GanzzahlExpression, *ZeichenketteExpression:
		return
	default:
		panic("eep " + reflect.TypeOf(node).Name())
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
