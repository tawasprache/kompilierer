package erstellungsprozess

import (
	"testing"
)

func TestGraf(t *testing.T) {
	a := Knote("a.tawa")
	b := Knote("b.tawa")
	c := Knote("c.tawa")
	d := Knote("d.tawa")
	k := Knote("k.tawa")

	g := NeuGraf()

	g.KnoteHinzufügen(a)
	g.KnoteHinzufügen(b)
	g.KnoteHinzufügen(c)
	g.KnoteHinzufügen(d)
	g.KnoteHinzufügen(k)

	g.KanteHinzufügen(k, a)

	g.KanteHinzufügen(a, c)
	g.KanteHinzufügen(b, c)

	g.KanteHinzufügen(c, d)

	g.KnoteAktivieren(a)
	g.KnoteAktivieren(k)
}
