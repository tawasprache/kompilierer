package codegenerierung

import (
	"Tawa/typen"
	"go/types"

	"github.com/llir/llvm/ir/value"
)

type namedThing interface{ isNamedThing() }
type NamedThingImpl struct{}

func (n NamedThingImpl) isNamedThing() {}

type LLVMMutableValue struct {
	NamedThingImpl
	value.Value
}
type LLVMValue struct {
	NamedThingImpl
	value.Value
}
type LLVMType struct {
	NamedThingImpl
	types.Type
	fields map[string]int
}
type TypenTyp struct {
	NamedThingImpl
	Art typen.Art
}

type ctx struct {
	names           []map[string]namedThing
	stringConstants map[string]value.Value
	i               int
}

func (c *ctx) pushScope() {
	c.names = append(c.names, make(map[string]namedThing))
}

func (c *ctx) popScope() {
	c.names = c.names[:len(c.names)-1]
}

func (c *ctx) lookup(id string) namedThing {
	for i := len(c.names) - 1; i >= 0; i-- {
		val, ok := c.names[i][id]
		if ok {
			return val
		}
	}

	panic("could not lookup " + id)
}

func (c *ctx) assign(id string, v namedThing) {
	for i := len(c.names) - 1; i >= 0; i-- {
		_, ok := c.names[i][id]
		if ok {
			c.names[i][id] = v
			return
		}
	}

	panic("could not find " + id)
}

func (c *ctx) top() map[string]namedThing {
	return c.names[len(c.names)-1]
}
