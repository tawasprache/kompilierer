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
}

func Gleich(l, r Typ) bool {
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
	return g
}
