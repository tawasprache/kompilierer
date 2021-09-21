package codegenerierung

import (
	"fmt"

	"github.com/pontaoski/gccjit"
)

type kontext struct {
	funktionen map[string]*gccjit.Function
	typen      map[string]gccjit.IType
	namen      []map[string]*wert
	i          int
	ii         int
	löschen    *gccjit.Function
	neu        *gccjit.Function
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
	lvalue gccjit.ILValue
	rvalue gccjit.IRValue
	typ    gccjit.IType
}

func nurrechts(i gccjit.IRValue, t gccjit.IType) *wert {
	return &wert{rvalue: i, typ: t}
}

func links(i gccjit.ILValue, t gccjit.IType) *wert {
	return &wert{lvalue: i, rvalue: i, typ: t}
}
