package typisierung

import (
	"Tawa/parser"
)

func definiert(v *VollKontext, e *parser.Expression) error {
	if e.Bedingung != nil {
		if err := definiert(v, &e.Bedingung.Wenn); err != nil {
			return err
		}
		if err := definiert(v, &e.Bedingung.Werden); err != nil {
			return err
		}
		if e.Bedingung.Sonst != nil {
			if err := definiert(v, e.Bedingung.Sonst); err != nil {
				return err
			}
		}
	} else if e.Definierung != nil {
		if _, ok := v.LookupVariable(e.Definierung.Variable); ok {
			return NeuFehler(e.Pos, "redefinition von »%s«", e.Definierung.Variable)
		}
		if err := definiert(v, &e.Definierung.Wert); err != nil {
			return err
		}
		v.KontextStack.Top().Variabeln[e.Definierung.Variable] = Art{}
	} else if e.Zuweisung != nil {
		if _, ok := v.LookupVariable(e.Zuweisung.Variable); !ok {
			return NeuFehler(e.Pos, "»%s« nicht deklariert", e.Definierung.Variable)
		}
		if err := definiert(v, &e.Zuweisung.Wert); err != nil {
			return err
		}
	} else if e.Variable != nil {
		if _, ok := v.LookupVariable(*e.Variable); !ok {
			return NeuFehler(e.Pos, "»%s« nicht deklariert", *e.Variable)
		}
	} else if e.Block != nil {
		for _, it := range e.Block {
			if err := definiert(v, &it); err != nil {
				return err
			}
		}
	}
	return nil
}

func VariabelnDefinierung(v *VollKontext, d *parser.Datei) error {
	for _, fnk := range d.Funktionen {
		v.Push()
		if err := definiert(v, &fnk.Expression); err != nil {
			return err
		}
		v.Pop()
	}
	return nil
}
