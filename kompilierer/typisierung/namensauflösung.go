package typisierung

import (
	"Tawa/kompilierer/ast"
	"Tawa/kompilierer/fehlerberichtung"
	"Tawa/kompilierer/getypisiertast"
	"Tawa/standardbibliothek"
	"errors"
	"strings"

	"github.com/alecthomas/repr"
)

func exprNamensauflösung(k *Kontext, s *scopes, l *lokalekontext, astExpr ast.Expression) (getypisiertast.Expression, error) {
	var feh error

	if astExpr.Terminal != nil {
		terminal := *astExpr.Terminal

		if terminal.Ganzzahl != nil {
			return getypisiertast.Ganzzahl{
				Wert: *terminal.Ganzzahl,
				LPos: getypisiertast.NeuSpan(terminal.Pos, terminal.EndPos),
			}, nil
		} else if terminal.Passt != nil {
			var (
				wert    getypisiertast.Expression
				mustern []getypisiertast.Muster

				lTyp getypisiertast.ITyp = getypisiertast.Nichtunifiziert{}
				lPos getypisiertast.Span = getypisiertast.NeuSpan(terminal.Pos, terminal.EndPos)
			)

			wert, feh = exprNamensauflösung(k, s, l, terminal.Passt.Wert)
			if feh != nil {
				return nil, feh
			}

			for _, it := range terminal.Passt.Mustern {
				_, _, sym, feh := l.auflöseVariant(it.Pattern.Name, getypisiertast.NeuSpan(it.Pattern.Pos, it.Pattern.EndPos))
				if feh != nil {
					return nil, feh
				}

				muster := getypisiertast.Muster{
					Variante:    sym,
					Konstruktor: it.Pattern.Name.Symbolen[len(it.Pattern.Name.Symbolen)-1],
				}

				s.neuScope()
				defer s.loescheScope()

				for idx, musterVari := range it.Pattern.Variabeln {
					s.head().vars[musterVari] = getypisiertast.Nichtunifiziert{}
					muster.Variablen = append(muster.Variablen, getypisiertast.Mustervariable{
						Variante:    sym,
						Konstruktor: it.Pattern.Name.Symbolen[len(it.Pattern.Name.Symbolen)-1],
						VonFeld:     idx,
						Name:        musterVari,
					})
				}

				musterExpr, feh := exprNamensauflösung(k, s, l, it.Expression)
				if feh != nil {
					return nil, feh
				}

				muster.Expression = musterExpr

				mustern = append(mustern, muster)
			}

			return getypisiertast.Pattern{
				Wert:    wert,
				Mustern: mustern,

				LTyp: lTyp,
				LPos: lPos,
			}, nil
		} else if terminal.Variantaufruf != nil {
			var (
				variant        getypisiertast.SymbolURL
				konstruktor    string = terminal.Variantaufruf.Name.Symbolen[len(terminal.Variantaufruf.Name.Symbolen)-1]
				argumenten     []getypisiertast.Expression
				strukturfelden []getypisiertast.Strukturfeld
				varianttyp     getypisiertast.ITyp = getypisiertast.Nichtunifiziert{}

				lPos getypisiertast.Span = getypisiertast.NeuSpan(terminal.Pos, terminal.EndPos)
			)

			_, _, variant, feh := l.auflöseVariant(terminal.Variantaufruf.Name, getypisiertast.NeuSpan(terminal.Pos, terminal.EndPos))
			if feh != nil && len(terminal.Variantaufruf.Argumente) == 0 && len(terminal.Variantaufruf.Strukturfelden) > 0 {
				_, variant, feh = l.auflöseTyp(terminal.Variantaufruf.Name, getypisiertast.NeuSpan(terminal.Pos, terminal.EndPos))
				if feh != nil {
					return nil, feh
				}
			} else if feh != nil {
				return nil, feh
			}

			for _, it := range terminal.Variantaufruf.Argumente {
				variExpr, feh := exprNamensauflösung(k, s, l, it)
				if feh != nil {
					return nil, feh
				}
				argumenten = append(argumenten, variExpr)
			}

			for _, it := range terminal.Variantaufruf.Strukturfelden {
				variExpr, feh := exprNamensauflösung(k, s, l, it.Wert)
				if feh != nil {
					return nil, feh
				}
				strukturfelden = append(strukturfelden, getypisiertast.Strukturfeld{
					Name: it.Name,
					Wert: variExpr,
				})
			}

			return getypisiertast.Variantaufruf{
				Variant:        variant,
				Konstruktor:    konstruktor,
				Argumenten:     argumenten,
				Strukturfelden: strukturfelden,
				Varianttyp:     varianttyp,

				LPos: lPos,
			}, nil
		} else if terminal.Leiste != nil {
			var (
				expressionen []getypisiertast.Expression

				lTyp getypisiertast.ITyp = getypisiertast.Nichtunifiziert{}
				lPos getypisiertast.Span = getypisiertast.NeuSpan(terminal.Pos, terminal.EndPos)
			)

			for _, it := range terminal.Leiste.Expressionen {
				v, feh := exprNamensauflösung(k, s, l, it)
				if feh != nil {
					return nil, feh
				}
				expressionen = append(expressionen, v)
			}

			return getypisiertast.Leiste{
				Werte: expressionen,
				LTyp:  lTyp,
				LPos:  lPos,
			}, nil
		} else if terminal.Funktionsaufruf != nil {
			var (
				funktion    getypisiertast.SymbolURL
				argumenten  []getypisiertast.Expression
				rückgabetyp getypisiertast.ITyp = getypisiertast.Nichtunifiziert{}
				lPos        getypisiertast.Span = getypisiertast.NeuSpan(terminal.Pos, terminal.EndPos)
			)

			_, funktion, feh = l.auflöseFunkSig(terminal.Funktionsaufruf.Name, getypisiertast.NeuSpan(terminal.Pos, terminal.EndPos))
			if feh != nil {
				return nil, feh
			}

			for _, arg := range terminal.Funktionsaufruf.Argumente {
				garg, feh := exprNamensauflösung(k, s, l, arg)
				if feh != nil {
					return nil, feh
				}
				argumenten = append(argumenten, garg)
			}

			return getypisiertast.Funktionsaufruf{
				Funktion:    funktion,
				Argumenten:  argumenten,
				Rückgabetyp: rückgabetyp,
				LPos:        lPos,
			}, nil
		} else if terminal.Variable != nil {
			_, gefunden := s.suche(*terminal.Variable)

			if !gefunden {
				return nil, fehlerberichtung.NeuFehler(getypisiertast.NeuSpan(terminal.Pos, terminal.EndPos), "variable »%s« nicht gefunden", *terminal.Variable)
			}

			return getypisiertast.Variable{
				Name: *terminal.Variable,
				ITyp: getypisiertast.Nichtunifiziert{},
				LPos: getypisiertast.NeuSpan(terminal.Pos, terminal.EndPos),
			}, nil
		} else if terminal.Zeichenkette != nil {
			return getypisiertast.Zeichenkette{
				Wert: *terminal.Zeichenkette,
				LPos: getypisiertast.NeuSpan(terminal.Pos, terminal.EndPos),
			}, nil
		} else if terminal.Strukturaktualisierung != nil {
			var (
				wert   getypisiertast.Expression
				felden []getypisiertast.Strukturaktualisierungsfeld
				lpos   getypisiertast.Span = getypisiertast.NeuSpan(terminal.Pos, terminal.EndPos)
			)

			wert, feh = exprNamensauflösung(k, s, l, terminal.Strukturaktualisierung.Struktur)
			if feh != nil {
				return nil, feh
			}

			for _, it := range terminal.Strukturaktualisierung.Felden {
				feldwert, feh := exprNamensauflösung(k, s, l, it.Wert)
				if feh != nil {
					return nil, feh
				}
				felden = append(felden, getypisiertast.Strukturaktualisierungsfeld{
					Name: it.Name,
					Wert: feldwert,
				})
			}

			return getypisiertast.Strukturaktualisierung{
				Wert:   wert,
				Felden: felden,
				LPos:   lpos,
			}, nil
		} else if terminal.Nativ != nil {
			n := getypisiertast.Nativ{
				Code: map[string]string{},
			}
			var feh error
			n.LTyp, feh = typ(l, terminal.Nativ.Typ, []string{})
			if feh != nil {
				return nil, feh
			}
			for _, it := range terminal.Nativ.Code {
				n.Code[it.Language] = it.Code[1 : len(it.Code)-1]
			}
			return n, nil
		}
	} else {
		links, feh := exprNamensauflösung(k, s, l, *astExpr.Links)
		if feh != nil {
			return nil, feh
		}
		rechts, feh := exprNamensauflösung(k, s, l, *astExpr.Rechts)
		if feh != nil {
			return nil, feh
		}
		switch *astExpr.Op {
		case ast.BinOpAdd:
			return getypisiertast.ValBinaryOperator{
				Links:  links,
				Rechts: rechts,
				Art:    getypisiertast.BinOpAdd,
				LTyp:   getypisiertast.Nichtunifiziert{},
				LPos:   getypisiertast.NeuSpan(astExpr.Pos, astExpr.EndPos),
			}, nil
		case ast.BinOpSub:
			return getypisiertast.ValBinaryOperator{
				Links:  links,
				Rechts: rechts,
				Art:    getypisiertast.BinOpSub,
				LTyp:   getypisiertast.Nichtunifiziert{},
				LPos:   getypisiertast.NeuSpan(astExpr.Pos, astExpr.EndPos),
			}, nil
		case ast.BinOpMul:
			return getypisiertast.ValBinaryOperator{
				Links:  links,
				Rechts: rechts,
				Art:    getypisiertast.BinOpMul,
				LTyp:   getypisiertast.Nichtunifiziert{},
				LPos:   getypisiertast.NeuSpan(astExpr.Pos, astExpr.EndPos),
			}, nil
		case ast.BinOpDiv:
			return getypisiertast.ValBinaryOperator{
				Links:  links,
				Rechts: rechts,
				Art:    getypisiertast.BinOpDiv,
				LTyp:   getypisiertast.Nichtunifiziert{},
				LPos:   getypisiertast.NeuSpan(astExpr.Pos, astExpr.EndPos),
			}, nil
		case ast.BinOpPow:
			return getypisiertast.ValBinaryOperator{
				Links:  links,
				Rechts: rechts,
				Art:    getypisiertast.BinOpPow,
				LTyp:   getypisiertast.Nichtunifiziert{},
				LPos:   getypisiertast.NeuSpan(astExpr.Pos, astExpr.EndPos),
			}, nil
		case ast.BinOpMod:
			return getypisiertast.ValBinaryOperator{
				Links:  links,
				Rechts: rechts,
				Art:    getypisiertast.BinOpMod,
				LTyp:   getypisiertast.Nichtunifiziert{},
				LPos:   getypisiertast.NeuSpan(astExpr.Pos, astExpr.EndPos),
			}, nil
		case ast.BinOpVerketten:
			return getypisiertast.ValBinaryOperator{
				Links:  links,
				Rechts: rechts,
				Art:    getypisiertast.BinOpVerketten,
				LTyp:   getypisiertast.Nichtunifiziert{},
				LPos:   getypisiertast.NeuSpan(astExpr.Pos, astExpr.EndPos),
			}, nil
		case ast.BinOpGleich:
			return getypisiertast.LogikBinaryOperator{
				Links:  links,
				Rechts: rechts,
				Art:    getypisiertast.BinOpGleich,
				LPos:   getypisiertast.NeuSpan(astExpr.Pos, astExpr.EndPos),
			}, nil
		case ast.BinOpNichtGleich:
			return getypisiertast.LogikBinaryOperator{
				Links:  links,
				Rechts: rechts,
				Art:    getypisiertast.BinOpNichtGleich,
				LPos:   getypisiertast.NeuSpan(astExpr.Pos, astExpr.EndPos),
			}, nil
		case ast.BinOpWeniger:
			return getypisiertast.LogikBinaryOperator{
				Links:  links,
				Rechts: rechts,
				Art:    getypisiertast.BinOpWeniger,
				LPos:   getypisiertast.NeuSpan(astExpr.Pos, astExpr.EndPos),
			}, nil
		case ast.BinOpWenigerGleich:
			return getypisiertast.LogikBinaryOperator{
				Links:  links,
				Rechts: rechts,
				Art:    getypisiertast.BinOpWenigerGleich,
				LPos:   getypisiertast.NeuSpan(astExpr.Pos, astExpr.EndPos),
			}, nil
		case ast.BinOpGrößer:
			return getypisiertast.LogikBinaryOperator{
				Links:  links,
				Rechts: rechts,
				Art:    getypisiertast.BinOpGrößer,
				LPos:   getypisiertast.NeuSpan(astExpr.Pos, astExpr.EndPos),
			}, nil
		case ast.BinOpGrößerGleich:
			return getypisiertast.LogikBinaryOperator{
				Links:  links,
				Rechts: rechts,
				Art:    getypisiertast.BinOpGrößerGleich,
				LPos:   getypisiertast.NeuSpan(astExpr.Pos, astExpr.EndPos),
			}, nil
		case ast.BinOpFeld:
			var v getypisiertast.Variable
			var ok bool
			if v, ok = rechts.(getypisiertast.Variable); !ok {
				return nil, fehlerberichtung.NeuFehler(rechts.Pos(), "ist kein feld")
			}
			return getypisiertast.Feldzugriff{
				Links: links,
				Feld:  v.Name,
				LTyp:  getypisiertast.Nichtunifiziert{},
				LPos:  rechts.Pos(),
			}, nil
		}
	}

	panic("unhandled case " + repr.String(astExpr))
}

func typNamensauflösung(k *Kontext, l *lokalekontext, t ast.Typdeklarationen) (getypisiertast.Typ, error) {
	g := getypisiertast.Typ{}
	g.SymbolURL = getypisiertast.SymbolURL{
		Paket: l.inModul,
		Name:  t.Name,
	}
	g.Generischeargumenten = t.Generischeargumenten

	for _, it := range t.Datenfelden {
		feld := getypisiertast.Datenfeld{
			Name: it.Name,
		}
		typ, err := typ(l, it.Typ, t.Generischeargumenten)
		if err != nil {
			return getypisiertast.Typ{}, err
		}
		feld.Typ = typ
		g.Datenfelden = append(g.Datenfelden, feld)
	}
	for _, it := range t.Varianten {
		vari := getypisiertast.Variant{
			Name: it.Name,
		}
		for _, feld := range it.Datenfelden {
			typ, err := typ(l, feld.Typ, t.Generischeargumenten)
			if err != nil {
				return getypisiertast.Typ{}, err
			}
			vari.Datenfelden = append(vari.Datenfelden, getypisiertast.Datenfeld{
				Name: feld.Name,
				Typ:  typ,
			})
		}
		g.Varianten = append(g.Varianten, vari)
	}

	return g, nil
}

func typ(l *lokalekontext, t ast.Typ, generischeargumenten []string) (getypisiertast.ITyp, error) {
	if t.Typkonstruktor != nil {
		_, haupt, err := l.auflöseTyp(t.Typkonstruktor.Name, getypisiertast.NeuSpan(t.Pos, t.EndPos))
		if err != nil {
			return nil, err
		}

		var ts []getypisiertast.ITyp
		for _, it := range t.Typkonstruktor.Generischeargumenten {
			t, err := typ(l, it, generischeargumenten)
			if err != nil {
				return nil, err
			}
			ts = append(ts, t)
		}

		return getypisiertast.Typnutzung{
			SymbolURL:            haupt,
			Generischeargumenten: ts,
		}, nil
	} else if t.Typvariable != nil {
		s := string(*t.Typvariable)
		for _, it := range generischeargumenten {
			if it == s {
				return getypisiertast.Typvariable{Name: string(*t.Typvariable)}, nil
			}
		}
		return nil, fehlerberichtung.NeuFehler(getypisiertast.NeuSpan(t.Pos, t.EndPos), "typvariable »%s« nicht gefunden", s)
	}
	panic("unreachable")
}

func funkZuSignatur(l *lokalekontext, t ast.Funktiondeklaration) (getypisiertast.Funktionssignatur, error) {
	r := getypisiertast.Funktionssignatur{}
	r.Generischeargumenten = t.Generischeargumenten

	if t.Rückgabetyp == nil {
		r.Rückgabetyp = getypisiertast.Typnutzung{
			SymbolURL: getypisiertast.SymbolURL{
				Paket: "Tawa/Eingebaut",
				Name:  "Einheit",
			},
		}
	} else {
		t, err := typ(l, *t.Rückgabetyp, t.Generischeargumenten)
		if err != nil {
			return getypisiertast.Funktionssignatur{}, err
		}
		r.Rückgabetyp = t
	}

	for _, it := range t.Formvariabeln {
		t, err := typ(l, it.Typ, t.Generischeargumenten)
		if err != nil {
			return getypisiertast.Funktionssignatur{}, nil
		}
		r.Formvariabeln = append(r.Formvariabeln, getypisiertast.Formvariable{
			Name: it.Name,
			Typ:  t,
		})
	}

	return r, nil
}

var defaultDependencies = []getypisiertast.Dependency{
	{
		Paket:      "Tawa/Eingebaut",
		ZeigeAlles: true,
	},
	{
		Paket:  "Tawa/Folge",
		Als:    "Folge",
		Zeigen: []string{"Folge"},
	},
	{
		Paket:  "Tawa/Vielleicht",
		Als:    "Vielleicht",
		Zeigen: []string{"Vielleicht"},
	},
	{
		Paket:  "Tawa/Leiste",
		Als:    "Leiste",
		Zeigen: []string{"Leiste"},
	},
}

func Lade(k *Kontext, paket string) error {
	if strings.HasPrefix(paket, "Tawa") {
		modul := ast.Modul{}
		builtin, err := standardbibliothek.StandardBibliothek.ReadFile(paket + ".tawa")
		if err != nil {
			return err
		}
		err = ast.Parser.ParseBytes(paket+".tawa", builtin, &modul)
		if err != nil {
			return err
		}
		g, err := zuGetypisierteAst(k, "Tawa", modul)
		if err != nil {
			return err
		}
		k.Module[paket] = g
		return nil
	}

	return errors.New("nicht implementiert")
}

func Auflösenamen(k *Kontext, m ast.Modul, modulePrefix string) (getypisiertast.Modul, error) {
	modul := getypisiertast.Modul{}
	l := &lokalekontext{
		k:                k,
		modul:            &modul,
		inModul:          modulePrefix + "/" + m.Package,
		lokaleFunktionen: map[string]getypisiertast.Funktionssignatur{},
		importieren:      defaultDependencies,
	}
	modul.Name = modulePrefix + "/" + m.Package
	if m.Nativauftakt != nil {
		modul.Nativcode = map[string]string{}
		for _, it := range m.Nativauftakt.Code {
			modul.Nativcode[it.Language] = it.Code[1 : len(it.Code)-1]
		}
	}
	if modul.Name == "Tawa/Eingebaut" {
		l.importieren = []getypisiertast.Dependency{}
	}
	for _, it := range m.Importierungen {
		dep := getypisiertast.Dependency{
			Paket: strings.Join(it.Import.Symbolen, "/"),
		}
		if it.Als != nil {
			dep.Als = strings.Join(it.Als.Symbolen, "/")
		}
		dep.Zeigen = it.Zeigen
		if it.ZeigenAlles != nil {
			dep.ZeigeAlles = true
		}
		if _, ok := k.Module[dep.Paket]; !ok {
			feh := Lade(k, dep.Paket)
			if feh != nil {
				return getypisiertast.Modul{}, feh
			}
		}
		l.importieren = append(l.importieren, dep)
		// TODO: geimportierte paket zum kontext hinzufuegen
	}
	modul.Zeigen = map[string]struct{}{}
	modul.Dependencies = append(modul.Dependencies, l.importieren...)
	if m.Zeigen.Alles != nil {
		modul.ZeigeAlles = true
	}
	if m.Zeigen.Symbolen != nil {
		for _, it := range *m.Zeigen.Symbolen {
			modul.Zeigen[it] = struct{}{}
		}
	}

	for _, it := range m.Deklarationen {
		if it.Typdeklarationen != nil {
			typ, feh := typNamensauflösung(k, l, *it.Typdeklarationen)
			if feh != nil {
				return getypisiertast.Modul{}, feh
			}
			modul.Typen = append(modul.Typen, typ)
		}
	}

	for _, it := range m.Deklarationen {
		if it.Funktiondeklaration != nil {
			sig, feh := funkZuSignatur(l, *it.Funktiondeklaration)
			if feh != nil {
				return getypisiertast.Modul{}, feh
			}
			l.lokaleFunktionen[it.Funktiondeklaration.Name] = sig
		}
	}

	for _, it := range m.Deklarationen {
		if it.Funktiondeklaration != nil {
			s := scopes{}
			s.neuScope()

			for _, it := range it.Funktiondeklaration.Formvariabeln {
				s.head().vars[it.Name] = getypisiertast.Nichtunifiziert{}
			}

			expr, feh := exprNamensauflösung(k, &s, l, it.Funktiondeklaration.Expression)
			if feh != nil {
				return getypisiertast.Modul{}, feh
			}

			modul.Funktionen = append(modul.Funktionen, getypisiertast.Funktion{
				SymbolURL: getypisiertast.SymbolURL{
					Paket: l.inModul,
					Name:  it.Funktiondeklaration.Name,
				},
				Funktionssignatur: l.lokaleFunktionen[it.Funktiondeklaration.Name],
				Expression:        expr,
			})
		}
	}

	return modul, nil
}
