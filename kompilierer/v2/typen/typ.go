package typen

type Typ interface {
	Basis() Typ
}

type Signature struct {
	Argumenten  []*Typname
	Rückgabetyp *Typname
}

func (s *Signature) Basis() Typ {
	return s
}

type Genanntetyp struct {
	Name  string
	Paket string

	basis Typ
}

func Gleich(l, r Typ) bool {
	l = l.Basis()
	r = r.Basis()

	if l == r {
		return true
	}

	switch l := l.(type) {
	case *Genanntetyp:
		if r, ok := r.(*Genanntetyp); ok {
			return l.Name == r.Name && l.Paket == r.Paket
		}
	case *Signature:
		if r, ok := r.(*Signature); ok {
			if l.Rückgabetyp == nil && r.Rückgabetyp == nil {
				return false
			} else if l.Rückgabetyp == nil {
				return false
			} else if r.Rückgabetyp == nil {
				return false
			}

			if !Gleich(l.Rückgabetyp.Typ(), r.Rückgabetyp.Typ()) {
				return false
			}
			if len(l.Argumenten) != len(r.Argumenten) {
				return false
			}
			for idx := range l.Argumenten {
				if !Gleich(l.Argumenten[idx].Typ(), r.Argumenten[idx].Typ()) {
					return false
				}
			}
			return true
		}
	default:
		panic("e")
	}

	return false
}

func (g *Genanntetyp) Basis() Typ {
	if g.basis != nil {
		return g.basis
	}
	return g
}

type Strukturfeld struct {
	Name string
	Typ  Typ

	ÜbergeordneterStrukturtyp *Strukturtyp
}

type Strukturfall struct {
	objekt

	Fallname string
	Felden   []Strukturfeld

	ÜbergeordneterStrukturtyp *Strukturtyp
}

type Strukturtyp struct {
	objekt

	Felden []*Strukturfeld
	Fälle  []*Strukturfall
}

func (s *Strukturtyp) Feld(n string) (f *Strukturfeld) {
	for _, es := range s.Felden {
		if es.Name == n {
			return es
		}
	}
	return nil
}

func (s *Strukturtyp) Basis() Typ {
	return s
}
