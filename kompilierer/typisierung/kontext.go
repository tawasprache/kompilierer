package typisierung

import (
	"Tawa/kompilierer/ast"
	"Tawa/kompilierer/getypisiertast"
	"Tawa/standardbibliothek"
)

type Kontext struct {
	Module map[string]getypisiertast.Modul
}

func zuGetypisierteAst(k *Kontext, p string, m ast.Modul) (getypisiertast.Modul, error) {
	genannt, feh := Auflösenamen(k, m, p)
	if feh != nil {
		return getypisiertast.Modul{}, feh
	}
	getypt, feh := Typiere(k, genannt, p)
	if feh != nil {
		return getypisiertast.Modul{}, feh
	}
	return getypt, nil
}

func ladeEingebaute(k *Kontext, path string) {
	modul := ast.Modul{}
	builtin, _ := standardbibliothek.StandardBibliothek.ReadFile(path + ".tawa")
	err := ast.Parser.ParseBytes(path+".tawa", builtin, &modul)
	if err != nil {
		panic(err)
	}
	g, err := zuGetypisierteAst(k, "Tawa", modul)
	if err != nil {
		panic(err)
	}
	k.Module[path] = g
}

func NeuKontext() *Kontext {
	k := &Kontext{}
	k.Module = map[string]getypisiertast.Modul{}

	ladeEingebaute(k, "Tawa/Eingebaut")
	ladeEingebaute(k, "Tawa/Folge")
	ladeEingebaute(k, "Tawa/Vielleicht")

	return k
}
