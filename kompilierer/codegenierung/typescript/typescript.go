package typescript

import (
	"Tawa/kompilierer/codegenierung"
	"Tawa/kompilierer/getypisiertast"
	"Tawa/kompilierer/typisierung"
	"embed"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/alecthomas/repr"
	"github.com/evanw/esbuild/pkg/api"
)

func init() {
	codegenierung.UnterbauRegistrieren("typescript", typescriptUnterbau{})
}

type typescriptUnterbau struct{}

//go:embed tsconfig.json
var tsc []byte

//go:embed "Js Helpers.ts"
var jshelpers []byte

//go:embed index.html.tmpl
var html string

//go:embed million/src
var million embed.FS

var tmpl = template.Must(template.New("s").Parse(html))

type tmplData struct {
	Script template.JS
}

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

func (t typescriptUnterbau) Postgen(o codegenierung.Optionen) error {
	var files = []string{
		"drivers/children.ts",
		"drivers/props.ts",
		"createElement.ts",
		"index.ts",
		"jsx.ts",
		"m.ts",
		"patch.ts",
		"schedule.ts",
		"types.ts",
	}
	feh := os.MkdirAll(path.Join(o.Outpath, "million", "drivers"), 0o777)
	if feh != nil {
		return feh
	}
	for _, it := range files {
		data, feh := million.ReadFile(path.Join("million", "src", it))
		if feh != nil {
			return feh
		}
		feh = ioutil.WriteFile(path.Join(o.Outpath, "million", it), data, 0o666)
		if feh != nil {
			return feh
		}
	}
	opts := api.BuildOptions{
		EntryPoints:       []string{path.Join(o.Outpath, o.Entry+".ts")},
		Outfile:           o.JSOutfile,
		Tsconfig:          path.Join(o.Outpath, "tsconfig.json"),
		TreeShaking:       api.TreeShakingTrue,
		Bundle:            true,
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		Write:             true,
	}
	if o.JSOutfile != "" {
		result := api.Build(opts)
		if len(result.Errors) > 0 {
			return fmt.Errorf("%s", repr.String(result))
		}
	}
	if o.HTMLOutfile != "" {
		opts.Write = false
		opts.Outfile = "/out.js"
		result := api.Build(opts)
		if len(result.Errors) > 0 {
			return fmt.Errorf("%s", repr.String(result))
		}
		var s strings.Builder
		fehler := tmpl.Execute(&s, tmplData{
			Script: template.JS(result.OutputFiles[0].Contents),
		})
		if fehler != nil {
			return fehler
		}
		fehler = ioutil.WriteFile(o.HTMLOutfile, []byte(s.String()), 0o666)
		if fehler != nil {
			return fehler
		}
	}

	return nil
}

func typEinfach(b getypisiertast.ITyp) bool {
	if typisierung.TypGleich(b, getypisiertast.TypGanz) {
		return true
	}
	if typisierung.TypGleich(b, getypisiertast.TypZeichenkette) {
		return true
	}
	if typisierung.TypGleich(b, getypisiertast.TypLogik) {
		return true
	}
	return false
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
		f.AddI(`{`)
		if e.Konstruktor != "" {
			f.Add(`__variant: {__tag: '%s'},`, e.Konstruktor)
			if len(e.Argumenten) > 0 {
				f.AddI(`__data: [`)
				for _, it := range e.Argumenten {
					genExpr(f, it, aktuellePaket)
					f.AddK(`, `)
				}
				f.AddD(`]`)
			}
		}
		for _, it := range e.Strukturfelden {
			f.AddE(`%s:`, it.Name)
			genExpr(f, it.Wert, aktuellePaket)
			f.AddK(`,`)
			f.AddNL()
		}
		f.AddD(`}`)
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
			if typEinfach(e.Links.Typ()) {
				f.AddK(`(`)
				genExpr(f, e.Links, aktuellePaket)
				if e.Art == getypisiertast.BinOpNichtGleich {
					f.AddK(` != `)
				} else {
					f.AddK(` == `)
				}
				genExpr(f, e.Rechts, aktuellePaket)
				f.AddK(`)`)
				return
			}
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
	case getypisiertast.Strukturaktualisierung:
		f.AddK(`$JsHelpers.update(`)
		genExpr(f, e.Wert, aktuellePaket)
		f.AddK(`, `)
		f.AddK(`{`)
		for idx, feld := range e.Felden {
			f.AddK(`%s: `, feld.Name)
			genExpr(f, feld.Wert, aktuellePaket)
			if idx != len(e.Felden)-1 {
				f.AddK(`,`)
			}
		}
		f.AddK(`}`)
		f.AddK(`)`)
	case getypisiertast.Feldzugriff:
		f.AddK(`(`)
		genExpr(f, e.Links, aktuellePaket)
		f.AddK(`).%s`, e.Feld)
	case getypisiertast.Nativ:
		var (
			v  string
			ok bool
		)
		if v, ok = e.Code["typescript"]; !ok {
			panic("kein nativ fur typescript")
		}
		f.AddK(v)
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
		f.Add(`import * as %s from "%s"`, zuIdent(it.Paket), it.Paket)
	}

	if m.Nativcode != nil {
		if v, ok := m.Nativcode["typescript"]; ok {
			f.Add("%s", v)
		}
	}

	for _, it := range m.Typen {
		f.AddI(`export interface %s%s {`, it.SymbolURL.Name, generischeString(it.Generischeargumenten))

		for _, it := range it.Datenfelden {
			f.Add(`%s: %s`, it.Name, typZuIdent(it.Typ, m.Name))
		}

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

	if m.Name == "User/Haupt" {
		f.Add(`Haupt()`)
	}

	// AUSGABE

	target := path.Join(o.Outpath, m.Name+".ts")
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
