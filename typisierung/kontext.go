package typisierung

import "Tawa/typen"

type Kontext struct {
	Arten     map[string]typen.Art
	Variabeln map[string]typen.Art
}

type KontextStack []*Kontext

func (k *KontextStack) Push() *Kontext {
	es := &Kontext{}
	es.Variabeln = map[string]typen.Art{}
	es.Arten = map[string]typen.Art{}
	*k = append(*k, es)
	return es
}

func (k *KontextStack) Pop() {
	*k = (*k)[:len(*k)-1]
}

func (k *KontextStack) Top() *Kontext {
	return (*k)[len(*k)-1]
}

func (k *KontextStack) LookupVariable(s string) (typen.Art, bool) {
	for _, it := range *k {
		if v, ok := it.Variabeln[s]; ok {
			return v, ok
		}
	}
	return nil, false
}

func (k *KontextStack) LookupArt(s string) (typen.Art, bool) {
	for _, it := range *k {
		if v, ok := it.Arten[s]; ok {
			return v, ok
		}
	}
	return nil, false
}

type VollKontext struct {
	Funktionen map[string]typen.Funktion
	KontextStack
}

var (
	logikArt     = typen.Logik{}
	ganzArt      = typen.Primitiv{Name: "ganz"}
	g8Art        = typen.Primitiv{Name: "g8"}
	g16Art       = typen.Primitiv{Name: "g16"}
	g32Art       = typen.Primitiv{Name: "g32"}
	g64Art       = typen.Primitiv{Name: "g64"}
	vzlosganzArt = typen.Primitiv{Name: "vzlosganz"}
	vzlosg8Art   = typen.Primitiv{Name: "vzlosg8"}
	vzlosg16Art  = typen.Primitiv{Name: "vzlosg16"}
	vzlosg32Art  = typen.Primitiv{Name: "vzlosg32"}
	vzlosg64Art  = typen.Primitiv{Name: "vzlosg64"}
	nichtsArt    = typen.Nichts{}
)

func NeuVollKontext() *VollKontext {
	return &VollKontext{
		Funktionen: map[string]typen.Funktion{},
		KontextStack: KontextStack{
			&Kontext{
				Arten: map[string]typen.Art{
					"logik": logikArt,

					// ganzzahl
					"ganz": ganzArt,
					"g8":   g8Art,
					"g16":  g16Art,
					"g32":  g32Art,
					"g64":  g64Art,

					// vorzeichenlose ganzzahl
					"vzlosganz": vzlosganzArt,
					"vzlosg8":   vzlosg8Art,
					"vzlosg16":  vzlosg16Art,
					"vzlosg32":  vzlosg32Art,
					"vzlosg64":  vzlosg64Art,

					"nichts": nichtsArt,
				},
				Variabeln: map[string]typen.Art{},
			},
		},
	}
}
