package dokumentation

import (
	"Tawa/kompilierer/getypisiertast"
	"embed"
	"fmt"
	"html/template"
	"strings"
)

//go:embed tmpl
var tmpls_ embed.FS

func zuLink(s getypisiertast.SymbolURL, art string) string {
	ersetzte := strings.ReplaceAll(s.Paket, "/", ":")
	return fmt.Sprintf(`<a href="%s.html#%s-%s"><small>%s</small>:%s</a>`, "./"+ersetzte, art, s.Name, ersetzte, s.Name)
}

func zuLinkT(i getypisiertast.ITyp) string {
	switch t := i.(type) {
	case getypisiertast.Typnutzung:
		var s []string
		for _, it := range t.Generischeargumenten {
			s = append(s, zuLinkT(it))
		}
		if len(t.Generischeargumenten) > 0 {
			return fmt.Sprintf("%s[%s]", zuLink(t.SymbolURL, "type"), strings.Join(s, ","))
		} else {
			return zuLink(t.SymbolURL, "type")
		}
	default:
		return t.String()
	}
}

var fns = map[string]interface{}{
	"typZuHTML": func(a getypisiertast.Typ) template.HTML {
		var s strings.Builder

		s.WriteString("typ ")
		s.WriteString(a.SymbolURL.Name)
		if len(a.Generischeargumenten) > 0 {
			s.WriteString(fmt.Sprintf("[%s]", strings.Join(a.Generischeargumenten, ", ")))
		}
		s.WriteString(" ist\n")

		for _, it := range a.Datenfelden {
			s.WriteString(fmt.Sprintf("\t%s: %s\n", it.Name, zuLinkT(it.Typ)))
		}

		if len(a.Datenfelden) > 0 && len(a.Varianten) > 0 {
			s.WriteString("\n")
		}

		for _, it := range a.Varianten {
			s.WriteString(fmt.Sprintf("\t| %s", it.Name))
			if len(it.Datenfelden) == 0 {
				s.WriteString("\n")
				continue
			}
			s.WriteString("(")
			for idx, typ := range it.Datenfelden {
				s.WriteString(fmt.Sprintf("%s: %s", typ.Name, zuLinkT(typ.Typ)))
				if idx != len(it.Datenfelden)-1 {
					s.WriteString(", ")
				}
			}
			s.WriteString(")\n")
		}

		s.WriteString("beende")

		return template.HTML(s.String())
	},
	"funkZuHTML": func(a getypisiertast.Funktion) template.HTML {
		var s strings.Builder

		s.WriteString("funk ")
		s.WriteString(a.SymbolURL.Name)
		if len(a.Funktionssignatur.Generischeargumenten) > 0 {
			s.WriteString(fmt.Sprintf("[%s]", strings.Join(a.Funktionssignatur.Generischeargumenten, ", ")))
		}
		s.WriteString("(")
		for idx, it := range a.Funktionssignatur.Formvariabeln {
			s.WriteString(fmt.Sprintf("%s: %s", it.Name, zuLinkT(it.Typ)))
			if idx != len(a.Funktionssignatur.Generischeargumenten)-1 {
				s.WriteString(", ")
			}
		}
		s.WriteString("): ")
		s.WriteString(zuLinkT(a.Funktionssignatur.Rückgabetyp))

		return template.HTML(s.String())
	},
}

var tmpls = template.Must(template.New("doc.tmpl").Funcs(fns).ParseFS(tmpls_, "tmpl/*.tmpl"))

type data struct {
	Name string
	T    []getypisiertast.Typ
	F    []getypisiertast.Funktion
}

func Dokumentation(g getypisiertast.Modul) string {
	var ku strings.Builder
	feh := tmpls.Execute(&ku, data{
		Name: g.Name,
		T:    g.Typen,
		F:    g.Funktionen,
	})
	if feh != nil {
		panic(feh)
	}
	return ku.String()
}
