package typisierung

type Art struct {
}

func (a *Art) IstGleich(b *Art) bool {
	return true
}

type Kontext struct {
	Variabeln map[string]Art
}

type KontextStack []*Kontext

func (k *KontextStack) Push() *Kontext {
	es := &Kontext{}
	es.Variabeln = map[string]Art{}
	*k = append(*k, es)
	return es
}

func (k *KontextStack) Pop() {
	*k = (*k)[:len(*k)-1]
}

func (k *KontextStack) Top() *Kontext {
	return (*k)[len(*k)-1]
}

func (k *KontextStack) LookupVariable(s string) (Art, bool) {
	for _, it := range *k {
		if v, ok := it.Variabeln[s]; ok {
			return v, ok
		}
	}
	return Art{}, false
}

type VollKontext struct {
	KontextStack
}
