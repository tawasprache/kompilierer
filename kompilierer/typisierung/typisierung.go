package typisierung

import (
	"Tawa/kompilierer/ast"
	"Tawa/kompilierer/getypisiertast"
)

func ZuGetypisierteAst(k *Kontext, modulePrefix string, d ast.Modul) (getypisiertast.Modul, error) {
	modul := getypisiertast.Modul{}
	l := &lokalekontext{
		k:                k,
		modul:            &modul,
		inModul:          modulePrefix + "/" + d.Package,
		lokaleFunktionen: map[string]getypisiertast.Funktionssignatur{},
		importieren:      []string{"Tawa/Eingebaut"},
	}
	modul.Name = modulePrefix + "/" + d.Package
	if modul.Name == "Tawa/Eingebaut" {
		l.importieren = []string{}
	}
	for _, it := range d.Importierungen {
		l.importieren = append(l.importieren, it.Import)
		// TODO: geimportierte paket zum kontext hinzufuegen
	}
	modul.Dependencies = append(modul.Dependencies, l.importieren...)

	for _, it := range d.Deklarationen {
		if it.Funktiondeklaration != nil {
			sig, err := funkZuSignatur(l, *it.Funktiondeklaration)
			if err != nil {
				return getypisiertast.Modul{}, err
			}
			l.lokaleFunktionen[it.Funktiondeklaration.Name] = sig
		}
	}

	for _, it := range d.Deklarationen {
		if it.Typdeklarationen != nil {
			it, err := typDeklZu(l, *it.Typdeklarationen)
			if err != nil {
				return getypisiertast.Modul{}, err
			}
			modul.Typen = append(modul.Typen, it)
		} else if it.Funktiondeklaration != nil {
			it, err := funkZu(l, *it.Funktiondeklaration)
			if err != nil {
				return getypisiertast.Modul{}, err
			}
			modul.Funktionen = append(modul.Funktionen, it)
		}
	}

	k.Module[modul.Name] = modul
	return modul, nil
}
