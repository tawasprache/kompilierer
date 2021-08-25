package codegenerierungaufgccjit

import (
	"Tawa/libgccjit"
	"fmt"
)

type kontext struct {
	funktionen map[string]*libgccjit.Function
	typen      map[string]libgccjit.IType
	namen      []map[string]*wert
	i          int
	ii         int
	löschen    *libgccjit.Function
	neu        *libgccjit.Function
}

func (c *kontext) pushScope() {
	c.i++
	c.namen = append(c.namen, make(map[string]*wert))
}

func (c *kontext) name() string {
	c.ii++

	return fmt.Sprintf("__tawa__struktur__%d", c.ii)
}

func (c *kontext) popScope() {
	c.i--
	c.namen = c.namen[:len(c.namen)-1]
}

func (c *kontext) namemiti(s string) string {
	return fmt.Sprintf("%s%d", s, c.i)
}

func (c *kontext) lookup(id string) *wert {
	for i := len(c.namen) - 1; i >= 0; i-- {
		val, ok := c.namen[i][id]
		if ok {
			return val
		}
	}

	panic("nicht gefunden: " + id)
}

func (c *kontext) assign(id string, v *wert) {
	for i := len(c.namen) - 1; i >= 0; i-- {
		_, ok := c.namen[i][id]
		if ok {
			c.namen[i][id] = v
			return
		}
	}

	panic("nicht gefunden: " + id)
}

func (c *kontext) top() map[string]*wert {
	return c.namen[len(c.namen)-1]
}

type wert struct {
	lvalue libgccjit.ILValue
	rvalue libgccjit.IRValue
	typ    libgccjit.IType
}

func nurrechts(i libgccjit.IRValue, t libgccjit.IType) *wert {
	return &wert{rvalue: i, typ: t}
}

func links(i libgccjit.ILValue, t libgccjit.IType) *wert {
	return &wert{lvalue: i, rvalue: i, typ: t}
}
