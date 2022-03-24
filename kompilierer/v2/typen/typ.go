package typen

type Typ interface {
	Basis() Typ
}

type Signature struct {
	Argumenten  []Typ
	Rückgabetyp Typ
}

func (s *Signature) Basis() Typ {
	return s
}

func Gleich(l, r Typ) bool {
	l = l.Basis()
	r = r.Basis()

	if l == r {
		return true
	}

	switch l := l.(type) {
	case *Strukturtyp:
		if r, ok := r.(*Strukturtyp); ok {
			return l.Name() == r.Name() && l.Paket() == r.Paket()
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

			if !Gleich(l.Rückgabetyp, r.Rückgabetyp) {
				return false
			}
			if len(l.Argumenten) != len(r.Argumenten) {
				return false
			}
			for idx := range l.Argumenten {
				if !Gleich(l.Argumenten[idx], r.Argumenten[idx]) {
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

func (s *Strukturtyp) Typ() Typ {
	return s
}

func (s *Strukturtyp) Feld(n string) (f *Strukturfeld) {
	for _, es := range s.Felden {
		if es.Name == n {
			return es
		}
	}
	return nil
}

func (s *Strukturtyp) Fall(n string) (f *Strukturfall) {
	for _, es := range s.Fälle {
		if es.name == n {
			return es
		}
	}
	return nil
}

type Fehlertyp struct {
}

func (f *Fehlertyp) Basis() Typ {
	return f
}

func (s *Strukturtyp) Basis() Typ {
	return s
}
