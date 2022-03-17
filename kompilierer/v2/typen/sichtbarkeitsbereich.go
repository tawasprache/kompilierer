package typen

type Sichtbarkeitsbereich struct {
	übergeordneterSichtbarkeitsbereich *Sichtbarkeitsbereich

	namen map[string]Objekt
}

var Wahr Objekt = &Strukturfall{
	objekt: objekt{
		sichtbarkeitsbereich: nil,
		name:                 "Wahr",
		paket:                "Eingebaut",
	},

	Fallname: "Wahr",
}

var Falsch Objekt = &Strukturfall{
	objekt: objekt{
		sichtbarkeitsbereich: nil,
		name:                 "Falsch",
		paket:                "Eingebaut",
	},

	Fallname: "Falsch",
}

var Wahrheitswert = &Strukturtyp{
	Fälle: []*Strukturfall{
		Wahr.(*Strukturfall),
		Falsch.(*Strukturfall),
	},
}

func init() {
	Wahr.(*Strukturfall).ÜbergeordneterStrukturtyp = Wahrheitswert
	Falsch.(*Strukturfall).ÜbergeordneterStrukturtyp = Wahrheitswert
}

var Welt *Sichtbarkeitsbereich = &Sichtbarkeitsbereich{
	übergeordneterSichtbarkeitsbereich: nil,
	namen: map[string]Objekt{
		"Ganzzahl": &Strukturtyp{
			objekt: objekt{
				sichtbarkeitsbereich: nil,
				name:                 "Ganzzahl",
				paket:                "Eingebaut",
			},
		},
		"Zeichenkette": &Strukturtyp{
			objekt: objekt{
				sichtbarkeitsbereich: nil,
				name:                 "Zeichenkette",
				paket:                "Eingebaut",
			},
		},
		"Wahrheitswert": Wahrheitswert,
		"Wahr":          Wahr,
		"Falsch":        Falsch,
	},
}

func (s *Sichtbarkeitsbereich) Hinzufügen(o Objekt) Objekt {
	s.namen[o.Name()] = o
	return o
}

func (s *Sichtbarkeitsbereich) Suchen(n string) (*Sichtbarkeitsbereich, Objekt) {
	if v, ok := s.namen[n]; ok {
		return s, v
	}

	if s.übergeordneterSichtbarkeitsbereich == nil {
		return nil, nil
	}

	return s.übergeordneterSichtbarkeitsbereich.Suchen(n)
}

func (s *Sichtbarkeitsbereich) Ersetzen(n string, o Objekt) {
	if _, ok := s.namen[n]; ok {
		s.namen[n] = o
	}

	if s.übergeordneterSichtbarkeitsbereich == nil {
		return
	}

	s.übergeordneterSichtbarkeitsbereich.Ersetzen(n, o)
}
