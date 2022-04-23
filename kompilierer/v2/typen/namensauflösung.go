package typen

import (
	"Tawa/kompilierer/v2/ast"
	"Tawa/kompilierer/v2/fehlerberichtung"
)

type namenaufslösungkontext struct {
	fehler               []error
	sichtbarkeitsbereich *Sichtbarkeitsbereich
	ktx                  *Kontext
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

func neuNamenaufslösungkontext(ktx *Kontext) *namenaufslösungkontext {
	k := &namenaufslösungkontext{}
	k.sichtbarkeitsbereich = Welt
	k.ktx = ktx

	return k
}

func (k *Sichtbarkeitsbereich) sucheObjekt(a ast.Typnutzung) (Objekt, error) {
	switch t := a.(type) {
	case *ast.IdentExpression:
		_, o := k.Suchen(t.Ident.String())
		if o == nil {
			return nil, fehlerberichtung.Neu(fehlerberichtung.TypNichtGefunden, a)
		}
		return o, nil
	case *ast.SelektorExpression:
		übergeordneterObjekt, feh := k.sucheObjekt(t.Objekt)
		if feh != nil {
			println("feld fail")
			return nil, feh
		}
		for _, objekt := range übergeordneterObjekt.Kindobjekte() {
			if objekt.Name() == t.Feld.Name {
				return objekt, nil
			}
		}
		return nil, fehlerberichtung.Neu(fehlerberichtung.FeldNichtGefunden, a)
	default:
		panic("e")
	}
}

func (k *Sichtbarkeitsbereich) sucheTyp(a ast.Typnutzung) (Typ, error) {
	v, feh := k.sucheObjekt(a)
	if feh != nil {
		return nil, feh
	}

	switch t := v.(type) {
	case Typ:
		return t, nil
	default:
		return nil, fehlerberichtung.Neu(fehlerberichtung.IstKeinTyp, a)
	}
}

func (k *namenaufslösungkontext) Visit(n ast.Node) ast.Visitor {
	switch x := n.(type) {
	case *ast.Funktiondeklaration:
		k.ktx.Defs[x.Name] = k.sichtbarkeitsbereich.Hinzufügen(
			&Funktion{
				objekt: objekt{
					sichtbarkeitsbereich: k.sichtbarkeitsbereich,
					name:                 x.Name.Name,
					paket:                "TODO",
					pos:                  x.Anfang(),
					typ: &Signature{
						Argumenten: func() []Typ {
							var r []Typ

							for _, arg := range x.Argumenten.Argumente {
								t, feh := k.sichtbarkeitsbereich.sucheTyp(arg.Typ)
								if feh != nil {
									k.fehler = append(k.fehler, feh)
									r = append(r, nil /* TODO */)
								} else {
									r = append(r, t)
								}
							}

							return r
						}(),
						Rückgabetyp: func() Typ {
							if x.Rückgabetyp == nil {
								return nil
							}

							t, feh := k.sichtbarkeitsbereich.sucheTyp(x.Rückgabetyp)
							if feh != nil {
								k.fehler = append(k.fehler, feh)
								// TODO: fehlertyp
							}

							return t
						}(),
					},
				},
			},
		)
		k.push()
	case *ast.Muster:
		k.push()
		for _, variable := range x.Pattern.Variabeln {
			k.ktx.Defs[variable] = k.sichtbarkeitsbereich.Hinzufügen(
				&Variable{
					objekt: objekt{
						sichtbarkeitsbereich: k.sichtbarkeitsbereich,
						name:                 variable.Name,
						paket:                "TODO",
						pos:                  variable.Anfang(),
					},
				},
			)
		}
	case *ast.Typdeklaration:
		strukturTyp := &Strukturtyp{
			objekt: objekt{
				sichtbarkeitsbereich: k.sichtbarkeitsbereich,
				name:                 x.Name.Name,
				paket:                "TODO",
				pos:                  x.Anfang(),
			},
		}
		strukturTyp.Felden = func() (r []*Strukturfeld) {
			for _, feld := range x.Felden {
				t, feh := k.sichtbarkeitsbereich.sucheTyp(feld.Typ)
				if feh != nil {
					panic(feh)
				}
				r = append(r, &Strukturfeld{
					Name:                      feld.Name.Name,
					Typ:                       t,
					ÜbergeordneterStrukturtyp: strukturTyp,
				})
			}
			return r
		}()
		strukturTyp.Fälle = func() (r []*Strukturfall) {
			for _, fall := range x.Fälle {
				s := &Strukturfall{
					objekt: objekt{
						sichtbarkeitsbereich: k.sichtbarkeitsbereich,
						name:                 fall.Name.Name,
						paket:                "TODO",
						pos:                  x.Anfang(),
					},
				}
				s.Fallname = fall.Name.Name
				for _, feld := range fall.Felden {
					t, feh := k.sichtbarkeitsbereich.sucheTyp(feld.Typ)
					if feh != nil {
						panic(feh)
					}
					s.Felden = append(s.Felden, Strukturfeld{
						Name:                      feld.Name.Name,
						Typ:                       t,
						ÜbergeordneterStrukturtyp: strukturTyp,
					})
				}
				r = append(r, s)
			}
			return r
		}()
		k.ktx.Defs[x.Name] = strukturTyp.sichtbarkeitsbereich.Hinzufügen(strukturTyp)
	case *ast.Argument:
		// fehler handled in funktiondeklaration
		t, _ := k.sichtbarkeitsbereich.sucheTyp(x.Typ)
		k.ktx.Defs[x.Name] = k.sichtbarkeitsbereich.Hinzufügen(
			&Variable{
				objekt: objekt{
					sichtbarkeitsbereich: k.sichtbarkeitsbereich,
					name:                 x.Name.Name,
					paket:                "TODO",
					pos:                  x.Anfang(),
					typ:                  t,
				},
			},
		)
	case *ast.Sei:
		k.ktx.Defs[x.Variable] = k.sichtbarkeitsbereich.Hinzufügen(
			&Variable{
				objekt: objekt{
					sichtbarkeitsbereich: k.sichtbarkeitsbereich,
					name:                 x.Variable.Name,
					paket:                "TODO",
					pos:                  x.Variable.Anfang(),
				},
			},
		)
	// case *ast.Typkonstruktor:
	// 	_, o := k.sichtbarkeitsbereich.Suchen(x.Ident.String())
	// 	switch o.(type) {
	// 	case *Strukturtyp:
	// 		k.ktx.Benutzern[x.Ident] = o
	// 		return k
	// 	case nil:
	// 		k.fehler = append(k.fehler, fehlerberichtung.Neu(fehlerberichtung.TypNichtGefunden, n))
	// 	default:
	// 		k.fehler = append(k.fehler, fehlerberichtung.Neu(fehlerberichtung.IstKeinTyp, n))
	// 	}
	case *ast.IdentExpression:
		_, o := k.sichtbarkeitsbereich.Suchen(x.Ident.String())
		switch o.(type) {
		// case *Variable:
		// 	k.ktx.Benutzern[x.Ident] = o
		// 	return k
		case nil:
			k.fehler = append(k.fehler, fehlerberichtung.Neu(fehlerberichtung.VarNichtGefunden, n))
		default:
			k.ktx.Benutzern[x.Ident] = o
			return k
			// k.fehler = append(k.fehler, fehlerberichtung.Neu(fehlerberichtung.IstKeinVariable, n))
		}
	case *ast.StrukturwertExpression:
		o, _ := k.sichtbarkeitsbereich.sucheObjekt(x.Typ)
		var strukturtyp *Strukturtyp
		switch objekt := o.(type) {
		case *Strukturfall:
			k.ktx.Benutzern[x.Typ] = objekt
			strukturtyp = objekt.ÜbergeordneterStrukturtyp
		case *Strukturtyp:
			k.ktx.Benutzern[x.Typ] = objekt
			strukturtyp = objekt
		case nil:
			k.fehler = append(k.fehler, fehlerberichtung.Neu(fehlerberichtung.TypNichtGefunden, n))
			return k
		}
	aus:
		for _, feld := range x.Felden {
			v := strukturtyp.Feld(feld.Name.String())
			if v != nil {
				continue aus
			}
			k.fehler = append(k.fehler, fehlerberichtung.Neu(fehlerberichtung.FeldNichtGefunden, n))
		}
	}
	return k
}

func (k *namenaufslösungkontext) EndVisit(n ast.Node) {
	switch n.(type) {
	case *ast.Funktiondeklaration, *ast.Muster:
		k.ktx.Sichtbarkeitsbereichen[n] = k.sichtbarkeitsbereich
		k.pop()
	}
}

type Kontext struct {
	Defs                   map[ast.Node]Objekt
	Benutzern              map[ast.Node]Objekt
	Sichtbarkeitsbereichen map[ast.Node]*Sichtbarkeitsbereich
}

func NeuKontext() *Kontext {
	return &Kontext{
		Defs:                   map[ast.Node]Objekt{},
		Benutzern:              map[ast.Node]Objekt{},
		Sichtbarkeitsbereichen: map[ast.Node]*Sichtbarkeitsbereich{},
	}
}

func Namenaufslösung(a *ast.Datei, ktx *Kontext) []error {
	k := neuNamenaufslösungkontext(ktx)
	ast.Walk(k, a)
	return k.fehler
}
