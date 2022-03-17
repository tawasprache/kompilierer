package typen

type Sichtbarkeitsbereich struct {
	übergeordneterSichtbarkeitsbereich *Sichtbarkeitsbereich

	namen map[string]Objekt
}

var Welt = &Sichtbarkeitsbereich{
	übergeordneterSichtbarkeitsbereich: nil,
	namen: map[string]Objekt{
		"Ganzzahl": &Typname{
			objekt: objekt{
				sichtbarkeitsbereich: nil,
				name:                 "Ganzzahl",
				paket:                "Eingebaut",
			},
		},
		"Zeichenkette": &Typname{
			objekt: objekt{
				sichtbarkeitsbereich: nil,
				name:                 "Zeichenkette",
				paket:                "Eingebaut",
			},
		},
		"Wahrheitswert": &Typname{
			objekt: objekt{
				sichtbarkeitsbereich: nil,
				name:                 "Wahrheitswert",
				paket:                "Eingebaut",
			},
		},
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
