package typisierung

import (
	"Tawa/kompilierer/ast"
	"Tawa/kompilierer/getypisiertast"

	"github.com/alecthomas/participle/v2/lexer"
)

type lokalekontext struct {
	k *Kontext

	modul            *getypisiertast.Modul
	inModul          string
	importieren      []getypisiertast.Dependency
	lokaleFunktionen map[string]getypisiertast.Funktionssignatur
}

type scopes struct {
	scopes []*scope
}

func (k *scopes) head() *scope {
	return k.scopes[len(k.scopes)-1]
}

func (k *scopes) neuScope() *scope {
	k.scopes = append(k.scopes, &scope{
		vars: map[string]getypisiertast.ITyp{},
	})
	return k.scopes[len(k.scopes)-1]
}

func (k *scopes) loescheScope() {
	k.scopes = k.scopes[:len(k.scopes)-1]
}

func (k *scopes) suche(n string) (getypisiertast.ITyp, bool) {
	for i := range k.scopes {
		ding := k.scopes[len(k.scopes)-1-i]
		if v, ok := ding.vars[n]; ok {
			return v, true
		}
	}
	return nil, false
}

type scope struct {
	vars map[string]getypisiertast.ITyp
}

func (l *lokalekontext) auflöseTyp(n ast.Symbolkette, pos lexer.Position) (getypisiertast.SymbolURL, error) {
	switch art := l.tabelleErstellen().namen[n.String()].(type) {
	case typEintrag:
		return art.SymURL, nil
	case error:
		return getypisiertast.SymbolURL{}, neuFehler(pos, "%s", art.Error())
	default:
		return getypisiertast.SymbolURL{}, neuFehler(pos, "typ »%s« nicht gefunden", n)
	}
}

func (l *lokalekontext) typDekl(url getypisiertast.SymbolURL, pos lexer.Position) (t getypisiertast.Typ, e error) {
	defer func() {
		if e == nil {
			t = copy(t).(getypisiertast.Typ)
		}
	}()
	for _, it := range l.modul.Typen {
		if it.SymbolURL == url {
			return it, nil
		}
	}
	for _, it := range l.importieren {
		modul := l.k.Module[it.Paket]
		for _, it := range modul.Typen {
			if it.SymbolURL == url {
				return it, nil
			}
		}
	}
	return getypisiertast.Typ{}, neuFehler(pos, "typ »%s« nicht gefunden", url)
}

func (l *lokalekontext) auflöseVariant(n ast.Symbolkette, pos lexer.Position) (t getypisiertast.Typ, v getypisiertast.Variant, s getypisiertast.SymbolURL, e error) {
	defer func() {
		if e == nil {
			t = copy(t).(getypisiertast.Typ)
			v = copy(v).(getypisiertast.Variant)
		}
	}()

	switch art := l.tabelleErstellen().namen[n.String()].(type) {
	case variantEintrag:
		return art.Typ, art.Variant, art.SymURL, nil
	case error:
		return getypisiertast.Typ{}, getypisiertast.Variant{}, getypisiertast.SymbolURL{}, neuFehler(pos, "%s", art.Error())
	default:
		return getypisiertast.Typ{}, getypisiertast.Variant{}, getypisiertast.SymbolURL{}, neuFehler(pos, "variant »%s« nicht gefunden", n)
	}
}

func (l *lokalekontext) auflöseFunkSig(n ast.Symbolkette, pos lexer.Position) (s getypisiertast.Funktionssignatur, sym getypisiertast.SymbolURL, e error) {
	defer func() {
		if e == nil {
			s = copy(s).(getypisiertast.Funktionssignatur)
			sym = copy(sym).(getypisiertast.SymbolURL)
		}
	}()

	switch art := l.tabelleErstellen().namen[n.String()].(type) {
	case funktionEintrag:
		return art.Sig, art.SymURL, nil
	case error:
		return getypisiertast.Funktionssignatur{}, getypisiertast.SymbolURL{}, neuFehler(pos, "%s", art.Error())
	default:
		return getypisiertast.Funktionssignatur{}, getypisiertast.SymbolURL{}, neuFehler(pos, "funktion »%s« nicht gefunden", n)
	}
}

func (l *lokalekontext) funktionsrumpf(url getypisiertast.SymbolURL, pos lexer.Position) (f getypisiertast.Funktion, e error) {
	defer func() {
		if e == nil {
			f = copy(f).(getypisiertast.Funktion)
		}
	}()
	for _, it := range l.modul.Funktionen {
		if it.SymbolURL == url {
			return it, nil
		}
	}
	for _, it := range l.importieren {
		modul := l.k.Module[it.Paket]
		for _, it := range modul.Funktionen {
			if it.SymbolURL == url {
				return it, nil
			}
		}
	}
	return getypisiertast.Funktion{}, neuFehler(pos, "funktion »%s« muss definiert sein vor nutzung", url.Name)
}
