package typisierung

import (
	"Tawa/kompilierer/getypisiertast"

	"github.com/alecthomas/participle/v2/lexer"
)

type lokalekontext struct {
	k *Kontext

	modul            *getypisiertast.Modul
	inModul          string
	importieren      []string
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

func (l *lokalekontext) sucheTyp(n string, pos lexer.Position) (getypisiertast.SymbolURL, error) {
	for _, it := range l.modul.Typen {
		if it.SymbolURL.Name == n {
			return it.SymbolURL, nil
		}
	}
	for _, it := range l.importieren {
		modul := l.k.Module[it]
		for _, it := range modul.Typen {
			if it.SymbolURL.Name == n {
				return it.SymbolURL, nil
			}
		}
	}
	return getypisiertast.SymbolURL{}, neuFehler(pos, "typ »%s« nicht gefunden", n)
}

func (l *lokalekontext) sucheTypDekl(url getypisiertast.SymbolURL, pos lexer.Position) (t getypisiertast.Typ, e error) {
	defer func() {
		t = copy(t).(getypisiertast.Typ)
	}()
	for _, it := range l.modul.Typen {
		if it.SymbolURL == url {
			return it, nil
		}
	}
	for _, it := range l.importieren {
		modul := l.k.Module[it]
		for _, it := range modul.Typen {
			if it.SymbolURL == url {
				return it, nil
			}
		}
	}
	return getypisiertast.Typ{}, neuFehler(pos, "typ »%s« nicht gefunden", url)
}

func (l *lokalekontext) sucheVariant(n string, pos lexer.Position) (t getypisiertast.Typ, v getypisiertast.Variant, s getypisiertast.SymbolURL, e error) {
	defer func() {
		t = copy(t).(getypisiertast.Typ)
		v = copy(v).(getypisiertast.Variant)
	}()
	for _, typ := range l.modul.Typen {
		for _, vari := range typ.Varianten {
			if vari.Name == n {
				return typ, vari, typ.SymbolURL, nil
			}
		}
	}
	for _, it := range l.importieren {
		modul := l.k.Module[it]
		for _, typ := range modul.Typen {
			for _, vari := range typ.Varianten {
				if vari.Name == n {
					return typ, vari, typ.SymbolURL, nil
				}
			}
		}
	}
	return getypisiertast.Typ{}, getypisiertast.Variant{}, getypisiertast.SymbolURL{}, neuFehler(pos, "variant »%s« nicht gefunden", n)
}

func (l *lokalekontext) sucheFunkSig(n string, pos lexer.Position) (s getypisiertast.Funktionssignatur, sym getypisiertast.SymbolURL, e error) {
	defer func() {
		s = copy(s).(getypisiertast.Funktionssignatur)
		sym = copy(sym).(getypisiertast.SymbolURL)
	}()
	for name, fn := range l.lokaleFunktionen {
		if name == n {
			return fn, getypisiertast.SymbolURL{
				Paket: l.inModul,
				Name:  name,
			}, nil
		}
	}
	for _, it := range l.importieren {
		modul := l.k.Module[it]
		for _, it := range modul.Funktionen {
			if it.SymbolURL.Name == n {
				return it.Funktionssignatur, it.SymbolURL, nil
			}
		}
	}
	return getypisiertast.Funktionssignatur{}, getypisiertast.SymbolURL{}, neuFehler(pos, "funktion »%s« nicht gefunden", n)
}

func (l *lokalekontext) sucheFunktionsrumpf(url getypisiertast.SymbolURL, pos lexer.Position) (f getypisiertast.Funktion, e error) {
	defer func() {
		f = copy(f).(getypisiertast.Funktion)
	}()
	for _, it := range l.modul.Funktionen {
		if it.SymbolURL == url {
			return it, nil
		}
	}
	for _, it := range l.importieren {
		modul := l.k.Module[it]
		for _, it := range modul.Funktionen {
			if it.SymbolURL == url {
				return it, nil
			}
		}
	}
	return getypisiertast.Funktion{}, neuFehler(pos, "funktion »%s« muss definiert sein vor nutzung", url.Name)
}
