package typen

import (
	"Tawa/kompilierer/v2/ast"
	"fmt"
)

type namenaufslösungkontext struct {
	fehler               []error
	sichtbarkeitsbereich *Sichtbarkeitsbereich
}

func (n *namenaufslösungkontext) push() {
	n.sichtbarkeitsbereich = &Sichtbarkeitsbereich{
		übergeordneterSichtbarkeitsbereich: n.sichtbarkeitsbereich,
		namen:                              map[string]Objekt{},
	}
}

func (n *namenaufslösungkontext) pop() {
	n.sichtbarkeitsbereich = n.sichtbarkeitsbereich.übergeordneterSichtbarkeitsbereich
}

func neuNamenaufslösungkontext() *namenaufslösungkontext {
	k := &namenaufslösungkontext{}
	k.sichtbarkeitsbereich = Welt

	return k
}

func (k *namenaufslösungkontext) sucheTyp(a ast.Typnutzung) (*Typname, error) {
	switch t := a.(type) {
	case *ast.Typkonstruktor:
		_, o := k.sichtbarkeitsbereich.Suchen(t.Symbolkette.String())
		switch a := o.(type) {
		case *Typname:
			return a, nil
		case nil:
			return nil, fmt.Errorf("%s nicht gefunden", a)
		default:
			return nil, fmt.Errorf("%s ist kein Typ", a)
		}
	default:
		panic("e")
	}
}

func (k *namenaufslösungkontext) Visit(n ast.Node) ast.Visitor {
	switch x := n.(type) {
	case *ast.Funktiondeklaration:
		k.sichtbarkeitsbereich.Hinzufügen(
			&Funktion{
				objekt: objekt{
					sichtbarkeitsbereich: k.sichtbarkeitsbereich,
					name:                 x.Name.Name,
					paket:                "TODO",
					pos:                  x.Anfang(),
				},
				Argumenten: func() []*Typname {
					var r []*Typname

					for _, arg := range x.Argumenten.Argumenten {
						t, feh := k.sucheTyp(arg.Typ)
						if feh != nil {
							k.fehler = append(k.fehler, feh)
							r = append(r, nil /* TODO */)
						} else {
							r = append(r, t)
						}
					}

					return r
				}(),
				Rückgabetyp: func() *Typname {
					if x.Rückgabetyp == nil {
						return nil
					}

					t, feh := k.sucheTyp(x.Rückgabetyp)
					if feh != nil {
						k.fehler = append(k.fehler, feh)
						// TODO: fehlertyp
					}

					return t
				}(),
			},
		)
		k.push()
	case *ast.Typdeklaration:
		k.sichtbarkeitsbereich.Hinzufügen(
			&Typname{
				objekt: objekt{
					sichtbarkeitsbereich: k.sichtbarkeitsbereich,
					name:                 x.Name.Name,
					paket:                "TODO",
					pos:                  x.Anfang(),
				},
			},
		)
	case *ast.Sei:
		k.sichtbarkeitsbereich.Hinzufügen(
			&Variable{
				objekt: objekt{
					sichtbarkeitsbereich: k.sichtbarkeitsbereich,
					name:                 x.Variable.Name,
					paket:                "TODO",
					pos:                  x.Variable.Anfang(),
				},
			},
		)
	case *ast.Typkonstruktor:
		_, o := k.sichtbarkeitsbereich.Suchen(x.Symbolkette.String())
		switch a := o.(type) {
		case *Typname:
			return k
		case nil:
			k.fehler = append(k.fehler, fmt.Errorf("%s nicht gefunden", a))
		default:
			k.fehler = append(k.fehler, fmt.Errorf("%s ist kein Typ", a))
		}
	case *ast.IdentExpression:
		_, o := k.sichtbarkeitsbereich.Suchen(x.Ident.String())
		switch a := o.(type) {
		case *Variable:
			return k
		case nil:
			k.fehler = append(k.fehler, fmt.Errorf("%s nicht gefunden", a))
		default:
			k.fehler = append(k.fehler, fmt.Errorf("%s ist kein Variable", a))
		}
	}
	return k
}

func (k *namenaufslösungkontext) EndVisit(n ast.Node) {
	switch n.(type) {
	case *ast.Funktiondeklaration:
		k.pop()
	}
}

func Namenaufslösung(a *ast.Datei) []error {
	k := neuNamenaufslösungkontext()
	ast.Walk(k, a)
	return k.fehler
}
