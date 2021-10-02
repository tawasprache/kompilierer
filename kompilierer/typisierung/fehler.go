package typisierung

import (
	"fmt"

	"github.com/alecthomas/participle/v2/lexer"
)

func neuFehler(l lexer.Position, format string, a ...interface{}) error {
	return fmt.Errorf("%s: %s", l, fmt.Sprintf(format, a...))
}
