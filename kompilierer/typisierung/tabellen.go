package typisierung

import (
	"Tawa/kompilierer/getypisiertast"
	"fmt"
	"strings"
)

type nametabelleEintrag interface {
	istNametabelleEintrag()
}

type istNametabelleEintragImpl struct{}

func (istNametabelleEintragImpl) istNametabelleEintrag() {}

type funktionEintrag struct {
	istNametabelleEintragImpl

	SymURL   getypisiertast.SymbolURL
	Sig      getypisiertast.Funktionssignatur
	Rumpf    getypisiertast.Funktion
	HatRumpf bool
}

type typEintrag struct {
	istNametabelleEintragImpl

	SymURL getypisiertast.SymbolURL
	Typ    getypisiertast.Typ
}

type variantEintrag struct {
	istNametabelleEintragImpl

	SymURL  getypisiertast.SymbolURL
	Typ     getypisiertast.Typ
	Variant getypisiertast.Variant
}

type konfliktEintrag struct {
	istNametabelleEintragImpl
	von []getypisiertast.SymbolURL
}

func (k konfliktEintrag) Error() string {
	return fmt.Sprintf("Ich weiß nicht, welches von %s Sie willen.", k.von)
}

type privatEintrag struct {
	istNametabelleEintragImpl
	name getypisiertast.SymbolURL
}

func (p privatEintrag) Error() string {
	return fmt.Sprintf("%s ist privat", p.name)
}

type tabelle struct {
	namen map[string]nametabelleEintrag
}

func neuTabelle() tabelle {
	return tabelle{
		namen: map[string]nametabelleEintrag{},
	}
}

func symURL(k nametabelleEintrag) getypisiertast.SymbolURL {
	switch v := k.(type) {
	case funktionEintrag:
		return v.SymURL
	case typEintrag:
		return v.SymURL
	case variantEintrag:
		return v.SymURL
	case privatEintrag:
		return v.name
	}
	return getypisiertast.SymbolURL{}
}

func (t *tabelle) hinzufügen(name string, neuEintrag nametabelleEintrag) {
	altWert, ok := t.namen[name]
	if ok {
		switch art := neuEintrag.(type) {
		case funktionEintrag:
			if art.HatRumpf {
				fnk, ok := altWert.(funktionEintrag)
				if ok {
					if !fnk.HatRumpf {
						t.namen[name] = art
						return
					}
				}
			}
		case privatEintrag:
			return
		}
		switch alteArt := altWert.(type) {
		case privatEintrag:
			t.namen[name] = neuEintrag
		case konfliktEintrag:
			alteArt.von = append(alteArt.von, symURL(neuEintrag))
			t.namen[name] = alteArt
		default:
			t.namen[name] = konfliktEintrag{
				von: []getypisiertast.SymbolURL{symURL(neuEintrag)},
			}
		}
		return
	}
	t.namen[name] = neuEintrag
}

func hat(s map[string]struct{}, ss string) bool {
	_, ok := s[ss]
	return ok
}

func (l *lokalekontext) tabelleErstellen() tabelle {
	tabelle := neuTabelle()

	for name, signatur := range l.lokaleFunktionen {
		tabelle.hinzufügen(name, funktionEintrag{
			SymURL: getypisiertast.SymbolURL{
				Paket: l.inModul,
				Name:  name,
			},
			Sig: signatur,
		})
	}

	hinzufügen := func(l *getypisiertast.Modul, symFunc func(getypisiertast.SymbolURL) string, varFunc func(string) string) {
		for _, funktion := range l.Funktionen {
			if l.ZeigeAlles || hat(l.Zeigen, funktion.SymbolURL.Name) {
				tabelle.hinzufügen(symFunc(funktion.SymbolURL), funktionEintrag{
					SymURL:   funktion.SymbolURL,
					Sig:      funktion.Funktionssignatur,
					Rumpf:    funktion,
					HatRumpf: true,
				})
			} else {
				tabelle.hinzufügen(symFunc(funktion.SymbolURL), privatEintrag{name: funktion.SymbolURL})
			}
		}
		for _, t := range l.Typen {
			if l.ZeigeAlles || hat(l.Zeigen, t.SymbolURL.Name) {
				tabelle.hinzufügen(symFunc(t.SymbolURL), typEintrag{
					SymURL: t.SymbolURL,
					Typ:    t,
				})
				for _, variant := range t.Varianten {
					tabelle.hinzufügen(varFunc(variant.Name), variantEintrag{
						SymURL:  t.SymbolURL,
						Typ:     t,
						Variant: variant,
					})
				}
			} else {
				tabelle.hinzufügen(symFunc(t.SymbolURL), privatEintrag{name: t.SymbolURL})
				for _, variant := range t.Varianten {
					tabelle.hinzufügen(varFunc(variant.Name), privatEintrag{name: t.SymbolURL})
				}
			}
		}
	}

	hinzufügen(l.modul, func(su getypisiertast.SymbolURL) string {
		return su.Name
	}, func(s string) string {
		return s
	})

	for _, it := range l.importieren {
		modul := l.k.Module[it.Paket]
		if it.ZeigeAlles {
			hinzufügen(&modul, func(su getypisiertast.SymbolURL) string {
				return su.Name
			}, func(s string) string {
				return s
			})
		} else if it.Als != "" {
			hinzufügen(&modul, func(su getypisiertast.SymbolURL) string {
				return it.Als + ":" + su.Name
			}, func(s string) string {
				return it.Als + ":" + s
			})
		} else {
			hinzufügen(&modul, func(su getypisiertast.SymbolURL) string {
				return strings.ReplaceAll(su.Paket, "/", ":") + ":" + su.Name
			}, func(s string) string {
				return strings.ReplaceAll(it.Paket, "/", ":") + ":" + s
			})
		}

		for _, zeige := range it.Zeigen {
			for _, funktion := range modul.Funktionen {
				if funktion.SymbolURL.Name == zeige {
					if modul.ZeigeAlles || hat(modul.Zeigen, zeige) {
						tabelle.hinzufügen(funktion.SymbolURL.Name, funktionEintrag{
							SymURL:   funktion.SymbolURL,
							Sig:      funktion.Funktionssignatur,
							Rumpf:    funktion,
							HatRumpf: true,
						})
					} else {
						tabelle.hinzufügen(funktion.SymbolURL.Name, privatEintrag{name: funktion.SymbolURL})
					}
				}
			}
			for _, t := range modul.Typen {
				if t.SymbolURL.Name == zeige {
					if modul.ZeigeAlles || hat(modul.Zeigen, zeige) {
						tabelle.hinzufügen(t.SymbolURL.Name, typEintrag{
							SymURL: t.SymbolURL,
							Typ:    t,
						})
						for _, variant := range t.Varianten {
							tabelle.hinzufügen(variant.Name, variantEintrag{
								SymURL:  t.SymbolURL,
								Typ:     t,
								Variant: variant,
							})
						}
					} else {
						tabelle.hinzufügen(t.SymbolURL.Name, privatEintrag{name: t.SymbolURL})
						for _, variant := range t.Varianten {
							tabelle.hinzufügen(variant.Name, privatEintrag{name: t.SymbolURL})
						}
					}
				}
			}
		}
	}

	return tabelle
}
