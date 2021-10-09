package typisierung

import (
	"Tawa/kompilierer/getypisiertast"
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2/lexer"
	"github.com/ztrue/tracerr"
)

func gleichErr(p lexer.Position, art string, a getypisiertast.ITyp, b getypisiertast.ITyp) error {
	return neuFehler(p, "%s nicht gleich: »%s« »%s«", art, a, b)
}

func neuFehler(l lexer.Position, format string, a ...interface{}) error {
	return tracerr.Wrap(fmt.Errorf("%s: %s", l, fmt.Sprintf(format, a...)))
}

type verketteterFehler struct {
	e []error
}

func (v verketteterFehler) Error() string {
	var s []string
	for _, it := range v.e {
		s = append(s, it.Error())
	}
	return strings.Join(s, "\n")
}

func fehlerVerketten(l error, r error) error {
	if l == nil && r == nil {
		return nil
	}
	var (
		ll []error
		rr []error
	)
	switch v := l.(type) {
	case verketteterFehler:
		ll = append(ll, v.e...)
	case nil:
		break
	default:
		ll = []error{v}
	}
	switch v := r.(type) {
	case verketteterFehler:
		rr = append(rr, v.e...)
	case nil:
		break
	default:
		rr = []error{v}
	}
	return verketteterFehler{
		e: append(ll, rr...),
	}
}
