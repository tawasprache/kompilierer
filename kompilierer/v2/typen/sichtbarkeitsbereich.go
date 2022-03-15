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
	},
}

func (s *Sichtbarkeitsbereich) Hinzufügen(o Objekt) {
	s.namen[o.Name()] = o
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
