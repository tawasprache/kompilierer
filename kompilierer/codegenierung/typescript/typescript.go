package typescript

import (
	"Tawa/kompilierer/codegenierung"
	"Tawa/kompilierer/getypisiertast"
	_ "embed"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/alecthomas/repr"
)

func init() {
	codegenierung.UnterbauRegistrieren("typescript", typescriptUnterbau{})
}

type typescriptUnterbau struct{}

//go:embed tsconfig.json
var tsc []byte

//go:embed "Js Helpers.ts"
var jshelpers []byte

func (t typescriptUnterbau) Pregen(o codegenierung.Optionen) error {
	path_ := path.Join(o.Outpath, "tsconfig.json")
	feh := ioutil.WriteFile(path_, tsc, 0o666)
	if feh != nil {
		return feh
	}
	path_ = path.Join(o.Outpath, "Js Helpers.ts")
	feh = ioutil.WriteFile(path_, jshelpers, 0o666)
	if feh != nil {
		return feh
	}
	return nil
}

func zuIdent(s string) string {
	return strings.ReplaceAll(s, "/", "_")
}

func symZuIdent(e getypisiertast.SymbolURL, aktuellePaket string) string {
	if aktuellePaket == e.Paket {
		return e.Name
	}
	return zuIdent(e.Paket) + "." + e.Name
}

func isVoid(e getypisiertast.ITyp) bool {
	switch k := e.(type) {
	case getypisiertast.Typnutzung:
		if k.SymbolURL.Paket == "Tawa/Eingebaut" {
			switch k.SymbolURL.Name {
			case "Einheit":
				return true
			}
		}
	}
	return false
}

func istBool(url getypisiertast.SymbolURL) bool {
	return url.Paket == "Tawa/Eingebaut" && url.Name == "Logik"
}

func typZuIdent(e getypisiertast.ITyp, aktuellePaket string) string {
	switch k := e.(type) {
	case getypisiertast.Typnutzung:
		if k.SymbolURL.Paket == "Tawa/Eingebaut" {
			switch k.SymbolURL.Name {
			case "Einheit":
				return "void"
			case "Ganz":
				return "number"
			case "Logik":
				return "boolean"
			case "Zeichenkette":
				return "string"
			}
		}
		ident := symZuIdent(k.SymbolURL, aktuellePaket)
		if len(k.Generischeargumenten) > 0 {
			var gens []string
			for _, it := range k.Generischeargumenten {
				gens = append(gens, typZuIdent(it, aktuellePaket))
			}
			ident += fmt.Sprintf(`<%s>`, strings.Join(gens, `, `))
		}
		return ident
	case getypisiertast.Typvariable:
		return k.Name
	}
	panic("e")
}

func genExpr(f *codegenierung.Filebuilder, expr getypisiertast.Expression, aktuellePaket string) {
	switch e := expr.(type) {
	case getypisiertast.Ganzzahl:
		f.AddK(`%d`, e.Wert)
	case getypisiertast.Zeichenkette:
		f.AddK(`%s`, e.Wert)
	case getypisiertast.Variable:
		f.AddK(`%s`, e.Name)
	case getypisiertast.Variantaufruf:
		if istBool(e.Variant) {
			if e.Konstruktor == "Wahr" {
				f.AddK(`true`)
			} else {
				f.AddK(`false`)
			}
			return
		}
		if len(e.Argumenten) == 0 {
			f.AddK(`{__variant: {__tag: '%s'}}`, e.Konstruktor)
		} else {
			f.AddK(`{__variant: {__tag: '%s', __data: [`, e.Konstruktor)
			for idx, it := range e.Argumenten {
				genExpr(f, it, aktuellePaket)
				if idx != len(e.Argumenten)-1 {
					f.AddK(`, `)
				}
			}
			f.AddK(`] }}`)
		}
	case getypisiertast.Funktionsaufruf:
		f.AddK(`%s(`, symZuIdent(e.Funktion, aktuellePaket))
		for idx, it := range e.Argumenten {
			genExpr(f, it, aktuellePaket)
			if idx != len(e.Argumenten)-1 {
				f.AddK(`, `)
			}
		}
		f.AddK(`)`)
	case getypisiertast.Pattern:
		f.AddK(`((__pat) => {`)
		f.Einzug++
		f.AddNL()

		var boolean bool
		switch k := e.Wert.Typ().(type) {
		case getypisiertast.Typnutzung:
			boolean = istBool(k.SymbolURL)
		}

		if boolean {
			f.Add(`switch (__pat) {`)
		} else {
			f.Add(`switch (__pat.__variant.__tag) {`)
		}

		for _, it := range e.Mustern {
			if boolean {
				switch it.Konstruktor {
				case "Wahr":
					f.AddI(`case true:`)
				case "Falsch":
					f.AddI(`case false:`)
				}
			} else {
				f.AddI(`case "%s":`, it.Konstruktor)
			}
			for _, vari := range it.Variablen {
				f.Add(`let %s = __pat.__variant.__data[%d]`, vari.Name, vari.VonFeld)
			}
			f.AddE(`return `)
			genExpr(f, it.Expression, aktuellePaket)
			f.AddNL()
			f.Einzug--
		}
		f.Add(`}`)

		f.Einzug--
		f.AddE(`})(`)
		genExpr(f, e.Wert, aktuellePaket)
		f.AddK(`)`)
		f.AddNL()
	case getypisiertast.ValBinaryOperator:
		f.AddK(`(`)
		genExpr(f, e.Links, aktuellePaket)

		switch e.Art {
		case getypisiertast.BinOpAdd:
			f.AddK(` + `)
		case getypisiertast.BinOpSub:
			f.AddK(` - `)
		case getypisiertast.BinOpMul:
			f.AddK(` * `)
		case getypisiertast.BinOpDiv:
			f.AddK(` / `)
		case getypisiertast.BinOpPow:
			f.AddK(` ** `)
		case getypisiertast.BinOpMod:
			f.AddK(` % `)
		default:
			panic("e")
		}

		genExpr(f, e.Rechts, aktuellePaket)
		f.AddK(`)`)
	case getypisiertast.LogikBinaryOperator:
		if e.Art == getypisiertast.BinOpGleich || e.Art == getypisiertast.BinOpNichtGleich {
			if e.Art == getypisiertast.BinOpNichtGleich {
				f.AddK(`!`)
			}
			f.AddK(`$JsHelpers.eq(`)
			genExpr(f, e.Links, aktuellePaket)
			f.AddK(`,`)
			genExpr(f, e.Rechts, aktuellePaket)
			f.AddK(`)`)
			return
		}
		f.AddK(`(`)
		genExpr(f, e.Links, aktuellePaket)

		switch e.Art {
		case getypisiertast.BinOpWeniger:
			f.AddK(` < `)
		case getypisiertast.BinOpWenigerGleich:
			f.AddK(` <= `)
		case getypisiertast.BinOpGrößer:
			f.AddK(` > `)
		case getypisiertast.BinOpGrößerGleich:
			f.AddK(` >= `)
		default:
			panic("e")
		}

		genExpr(f, e.Rechts, aktuellePaket)
		f.AddK(`)`)
	default:
		panic("e " + repr.String(e))
	}
}

func generischeString(s []string) string {
	if len(s) > 0 {
		return fmt.Sprintf(`<%s>`, strings.Join(s, ", "))
	}
	return ""
}

func (t typescriptUnterbau) CodegenModul(o codegenierung.Optionen, m getypisiertast.Modul) error {
	f := codegenierung.Filebuilder{}

	f.Add(`import * as $JsHelpers from "Js Helpers"`)
	for _, it := range m.Dependencies {
		f.Add(`import * as %s from "%s"`, zuIdent(it), it)
	}

	for _, it := range m.Typen {
		f.AddI(`export interface %s%s {`, it.SymbolURL.Name, generischeString(it.Generischeargumenten))

		if len(it.Varianten) > 0 {
			f.AddI(`__variant:`)

			for _, vari := range it.Varianten {
				f.AddI(`| {`)
				f.Add(`__tag: '%s'`, vari.Name)
				if len(vari.Datenfelden) > 0 {
					f.AddE(`__data: [`)
					for idx, it := range vari.Datenfelden {
						f.AddK(`%s`, typZuIdent(it.Typ, m.Name))
						if idx != len(vari.Datenfelden)-1 {
							f.AddK(`, `)
						}
					}
					f.AddK(`]`)
					f.AddNL()
				}
				f.AddD(`}`)
			}

			f.AddD(``)
		}

		f.AddD(`}`)
	}
	for _, it := range m.Funktionen {
		var sig []string

		for _, arg := range it.Funktionssignatur.Formvariabeln {
			sig = append(sig, fmt.Sprintf(`%s: %s`, arg.Name, typZuIdent(arg.Typ, m.Name)))
		}

		f.Add(`export function %s%s(%s): %s`, it.SymbolURL.Name, generischeString(it.Funktionssignatur.Generischeargumenten), strings.Join(sig, ", "), typZuIdent(it.Funktionssignatur.Rückgabetyp, m.Name))
		f.AddI(`{`)
		if !isVoid(it.Funktionssignatur.Rückgabetyp) {
			f.AddE(`return `)
		} else {
			f.AddE(``)
		}
		genExpr(&f, it.Expression, m.Name)
		f.AddNL()
		f.AddD(`}`)

	}

	// AUSGABE

	target := path.Join(o.Outpath, m.Name+".ts")
	repr.Println(m)
	println(target)

	fehler := os.MkdirAll(path.Dir(target), 0o777)
	if fehler != nil {
		return fehler
	}

	fehler = ioutil.WriteFile(target, []byte(f.String()), 0o666)
	if fehler != nil {
		return fehler
	}

	return nil
}
