package typen

type Art interface {
	String() string
	IstGleich(Art) bool
}

type Primitiv struct {
	Name string
}

func (p Primitiv) IstGleich(a Art) bool {
	v, ok := a.(Primitiv)
	if !ok {
		return false
	}

	return p.Name == v.Name
}

func (p Primitiv) String() string {
	return p.Name
}
