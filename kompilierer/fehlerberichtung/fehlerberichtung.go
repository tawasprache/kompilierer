package fehlerberichtung

import (
	"Tawa/kompilierer/getypisiertast"
	"fmt"
	"strings"

	"github.com/ztrue/tracerr"
)

func GleichErr(span getypisiertast.Span, art string, a getypisiertast.ITyp, b getypisiertast.ITyp) error {
	return NeuFehler(span, "%s nicht gleich: »%s« »%s«", art, a, b)
}

func NeuFehler(span getypisiertast.Span, format string, a ...interface{}) error {
	return tracerr.Wrap(PositionError{
		Text: fmt.Sprintf(format, a...),
		Span: span,
	})
}

type PositionError struct {
	Text string
	Span getypisiertast.Span
}

func (p PositionError) Error() string {
	return fmt.Sprintf("%s: %s", p.Span.Von, p.Text)
}

type VerketteterFehler struct {
	Fehler []error
}

func (v VerketteterFehler) Error() string {
	var s []string
	for _, it := range v.Fehler {
		s = append(s, it.Error())
	}
	return strings.Join(s, "\n")
}

func FehlerVerketten(l error, r error) error {
	if l == nil && r == nil {
		return nil
	}
	var (
		ll []error
		rr []error
	)
	switch v := l.(type) {
	case VerketteterFehler:
		ll = append(ll, v.Fehler...)
	case nil:
		break
	default:
		ll = []error{v}
	}
	switch v := r.(type) {
	case VerketteterFehler:
		rr = append(rr, v.Fehler...)
	case nil:
		break
	default:
		rr = []error{v}
	}
	return VerketteterFehler{
		Fehler: append(ll, rr...),
	}
}
