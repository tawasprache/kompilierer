package typisierung

import (
	"Tawa/kompilierer/ast"
	"Tawa/kompilierer/getypisiertast"
	"Tawa/standardbibliothek"
)

type Kontext struct {
	Module map[string]getypisiertast.Modul
}

func ladeEingebaute(k *Kontext, path string) {
	modul := ast.Modul{}
	builtin, _ := standardbibliothek.StandardBibliothek.ReadFile(path + ".tawa")
	err := ast.Parser.ParseBytes(path+".tawa", builtin, &modul)
	if err != nil {
		panic(err)
	}
	g, err := ZuGetypisierteAst(k, "Tawa", modul)
	if err != nil {
		panic(err)
	}
	k.Module[path] = g
}

func NeuKontext() *Kontext {
	k := &Kontext{}
	k.Module = map[string]getypisiertast.Modul{}

	ladeEingebaute(k, "Tawa/Eingebaut")

	return k
}
