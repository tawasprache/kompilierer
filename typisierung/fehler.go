package typisierung

import (
	"fmt"

	"github.com/alecthomas/participle/v2/lexer"
)

type Fehler struct {
	lexer.Position
	Err error
}

func NeuFehler(p lexer.Position, format string, a ...interface{}) *Fehler {
	return &Fehler{
		Position: p,
		Err:      fmt.Errorf(format, a...),
	}
}

func (f *Fehler) Error() string {
	return fmt.Sprintf("%s: %s", f.Position, f.Err)
}
