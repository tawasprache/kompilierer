package typisierung

import (
	"fmt"

	"github.com/alecthomas/participle/v2/lexer"
	"github.com/ztrue/tracerr"
)

func neuFehler(l lexer.Position, format string, a ...interface{}) error {
	return tracerr.Wrap(fmt.Errorf("%s: %s", l, fmt.Sprintf(format, a...)))
}
