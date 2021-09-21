package typisierung

import (
	"Tawa/parser"
	"Tawa/typen"
)

type kontext struct {
	scopes []*scope
}

func neuKontext() *kontext {
	a := &kontext{}
	a.neuScope()
	a.head().typs = map[string]typen.Typ{
		"ganz":  typen.Integer{},
		"logik": typen.Logik{},
	}
	return a
}

func (k *kontext) head() *scope {
	return k.scopes[len(k.scopes)-1]
}

func (k *kontext) neuScope() *scope {
	k.scopes = append(k.scopes, &scope{
		fnTyps: map[string]typen.Funktion{},
		fns:    map[string]parser.Funktion{},
		typs:   map[string]typen.Typ{},
		vars:   map[string]typen.Typ{},
	})
	return k.scopes[len(k.scopes)-1]
}

func (k *kontext) loescheScope() {
	k.scopes = k.scopes[:len(k.scopes)-1]
}

func (k *kontext) sucheFnTyp(n string) (typen.Funktion, bool) {
	for i := range k.scopes {
		ding := k.scopes[len(k.scopes)-1-i]
		if v, ok := ding.fnTyps[n]; ok {
			return v, true
		}
	}
	return typen.Funktion{}, false
}

func (k *kontext) sucheFn(n string) (parser.Funktion, bool) {
	for i := range k.scopes {
		ding := k.scopes[len(k.scopes)-1-i]
		if v, ok := ding.fns[n]; ok {
			return v, true
		}
	}
	return parser.Funktion{}, false
}

func (k *kontext) sucheVar(n string) (typen.Typ, bool) {
	for i := range k.scopes {
		ding := k.scopes[len(k.scopes)-1-i]
		if v, ok := ding.vars[n]; ok {
			return v, true
		}
	}
	return typen.Funktion{}, false
}

func (k *kontext) sucheTyps(n string) (typen.Typ, bool) {
	for i := range k.scopes {
		ding := k.scopes[len(k.scopes)-1-i]
		if v, ok := ding.typs[n]; ok {
			return v, true
		}
	}
	return nil, false
}

func (k *kontext) typVonParser(p *parser.Art) (typen.Typ, bool) {
	a, ok := k.sucheTyps(p.Name)
	if !ok {
		return nil, false
	}

	return a, true
}

type scope struct {
	fnTyps map[string]typen.Funktion
	fns    map[string]parser.Funktion
	typs   map[string]typen.Typ
	vars   map[string]typen.Typ
}
