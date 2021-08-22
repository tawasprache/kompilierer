package interpreter

import (
	"fmt"
	"strconv"
	"strings"
)

type vollwert struct {
	links  lvalue
	rechts wert
}

func vrechts(v wert) *vollwert {
	return &vollwert{rechts: v}
}

type wert interface {
	isWert()
	AlsString() string
}

type IsWert struct{}

func (IsWert) isWert() {}

type logik struct {
	IsWert

	val bool
}

func (l logik) AlsString() string {
	if l.val {
		return "Wahr"
	}
	return "Falsch"
}

type ganz struct {
	IsWert

	val int64
}

func (l ganz) AlsString() string {
	return strconv.FormatInt(l.val, 10)
}

type ptr struct {
	IsWert

	w **wert
}

func (l ptr) AlsString() string {
	return fmt.Sprintf("zeiger auf %s", (*(*l.w)).AlsString())
}

type lvalue struct {
	IsWert

	w **wert
}

func (l lvalue) AlsString() string {
	return fmt.Sprintf("%s", (*(*l.w)).AlsString())
}

type struktur struct {
	IsWert

	fields map[string]*wert
}

func (l struktur) AlsString() string {
	var s strings.Builder
	s.WriteString("(\n")
	for k, v := range l.fields {
		fmt.Fprintf(&s, "%s ist %s\n", k, (*v).AlsString())
	}
	s.WriteString(")")
	return s.String()
}
