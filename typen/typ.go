package typen

type Typ interface{ istTyp() }

type istTypImpl struct{}

func (i istTypImpl) istTyp() {}

type Funktion struct {
	Eingabe []Typ
	Ausgabe Typ

	istTypImpl
}

type Integer struct {
	istTypImpl
}

type Logik struct {
	istTypImpl
}

type Nichts struct {
	istTypImpl
}
