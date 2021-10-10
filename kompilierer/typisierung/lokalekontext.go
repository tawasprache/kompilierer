package typisierung

import (
	"Tawa/kompilierer/ast"
	"Tawa/kompilierer/fehlerberichtung"
	"Tawa/kompilierer/getypisiertast"
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

func (l *lokalekontext) auflöseTyp(n ast.Symbolkette, pos getypisiertast.Span) (getypisiertast.Typ, getypisiertast.SymbolURL, error) {
	switch art := l.tabelleErstellen().namen[n.String()].(type) {
	case typEintrag:
		return art.Typ, art.SymURL, nil
	case error:
		return getypisiertast.Typ{}, getypisiertast.SymbolURL{}, fehlerberichtung.NeuFehler(pos, "%s", art.Error())
	default:

		return getypisiertast.Typ{}, getypisiertast.SymbolURL{}, fehlerberichtung.NeuFehler(pos, "typ »%s« nicht gefunden", n)
	}
}

func (l *lokalekontext) typDekl(url getypisiertast.SymbolURL, pos getypisiertast.Span) (t getypisiertast.Typ, e error) {
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
	return getypisiertast.Typ{}, fehlerberichtung.NeuFehler(pos, "typ »%s« nicht gefunden", url)
}

func (l *lokalekontext) auflöseVariant(n ast.Symbolkette, pos getypisiertast.Span) (t getypisiertast.Typ, v getypisiertast.Variant, s getypisiertast.SymbolURL, e error) {
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
		return getypisiertast.Typ{}, getypisiertast.Variant{}, getypisiertast.SymbolURL{}, fehlerberichtung.NeuFehler(pos, "%s", art.Error())
	default:
		return getypisiertast.Typ{}, getypisiertast.Variant{}, getypisiertast.SymbolURL{}, fehlerberichtung.NeuFehler(pos, "variant »%s« nicht gefunden", n)
	}
}

func (l *lokalekontext) auflöseFunkSig(n ast.Symbolkette, pos getypisiertast.Span) (s getypisiertast.Funktionssignatur, sym getypisiertast.SymbolURL, e error) {
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
		return getypisiertast.Funktionssignatur{}, getypisiertast.SymbolURL{}, fehlerberichtung.NeuFehler(pos, "%s", art.Error())
	default:
		return getypisiertast.Funktionssignatur{}, getypisiertast.SymbolURL{}, fehlerberichtung.NeuFehler(pos, "funktion »%s« nicht gefunden", n)
	}
}

func (l *lokalekontext) funktionsrumpf(url getypisiertast.SymbolURL, pos getypisiertast.Span) (f getypisiertast.Funktion, e error) {
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
	return getypisiertast.Funktion{}, fehlerberichtung.NeuFehler(pos, "funktion »%s« muss definiert sein vor nutzung", url.Name)
}
