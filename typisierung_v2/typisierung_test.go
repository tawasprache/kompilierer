package typisierungv2

import (
	"Tawa/parser"
	"testing"

	"github.com/alecthomas/repr"
)

func checkExpr(ktx *kontext, expr *parser.Expression, gegenArt typ, gut bool, t *testing.T) {
	err := checkExpression(ktx, expr, gegenArt)
	if gut {
		if err != nil {
			t.Fatalf("expected no errors, got one: %s", err)
		}
	} else {
		if err == nil {
			t.Fatalf("expected an error, didn't get one: %s", err)
		}
	}
}

func TestInteger(t *testing.T) {
	expr := &parser.Expression{
		Integer: &parser.Integer{
			Value: 2,
		},
	}
	expr2 := &parser.Expression{
		Logik: &parser.Logik{
			Wert: "Wahr",
		},
	}
	kind := integer{}
	ctx := neuKontext()

	checkExpr(ctx, expr, kind, true, t)
	checkExpr(ctx, expr2, kind, false, t)
}

func TestFunktion(t *testing.T) {
	expr1 := &parser.Expression{
		Funktionsaufruf: &parser.Funktionsaufruf{
			Name: "ident",
			Argumente: []parser.Expression{
				{
					Integer: &parser.Integer{
						Value: 2,
					},
				},
			},
		},
	}
	expr2 := &parser.Expression{
		Funktionsaufruf: &parser.Funktionsaufruf{
			Name: "ident",
			Argumente: []parser.Expression{
				{
					Logik: &parser.Logik{
						Wert: "Wahr",
					},
				},
			},
		},
	}
	ekind := integer{}
	ctx := neuKontext()
	ctx.fns["ident"] = funktion{
		eingabe: []typ{
			kvar{n: "a"},
		},
		ausgabe: kvar{n: "a"},
	}

	checkExpr(ctx, expr1, ekind, true, t)
	checkExpr(ctx, expr2, ekind, false, t)
}

func TestFunktionZwei(t *testing.T) {
	// schlecht
	expr1 := &parser.Expression{
		Funktionsaufruf: &parser.Funktionsaufruf{
			Name: "ident",
			Argumente: []parser.Expression{
				{
					Integer: &parser.Integer{
						Value: 2,
					},
				},
				{
					Logik: &parser.Logik{
						Wert: "Wahr",
					},
				},
			},
		},
	}
	// gut
	expr2 := &parser.Expression{
		Funktionsaufruf: &parser.Funktionsaufruf{
			Name: "ident",
			Argumente: []parser.Expression{
				{
					Integer: &parser.Integer{
						Value: 2,
					},
				},
				{
					Integer: &parser.Integer{
						Value: 2,
					},
				},
			},
		},
	}
	ekind := integer{}
	ctx := neuKontext()
	ctx.fns["ident"] = funktion{
		eingabe: []typ{
			kvar{n: "a"},
			kvar{n: "a"},
		},
		ausgabe: kvar{n: "a"},
	}

	checkExpr(ctx, expr1, ekind, false, t)
	checkExpr(ctx, expr2, ekind, true, t)
}

func TestVoll(t *testing.T) {
	maybeT := generischerTyp{
		von: entweder{
			fallen: map[string]typ{
				"nur":    kvar{n: "a"},
				"nichts": nichts{},
			},
		},
		argumenten: []string{"a"},
	}
	maybeInt := maybeT.voll(map[string]typ{
		"a": integer{},
	})

	if !gleich(maybeInt.(entweder).fallen["nur"], integer{}) {
		t.Fatalf("a isn't an integer\n\tgot:%s\n\twant:%s", repr.String(maybeInt.(entweder).fallen["nur"]), repr.String(integer{}))
	}

	maybeLogik := maybeT.voll(map[string]typ{
		"a": logik{},
	})

	if !gleich(maybeLogik.(entweder).fallen["nur"], logik{}) {
		t.Fatalf("a isn't an logik\n\tgot:%s\n\twant:%s", repr.String(maybeInt.(entweder).fallen["nur"]), repr.String(logik{}))
	}

	maybeA := maybeT.voll(map[string]typ{
		"a": kvar{n: "a"},
	})

	if !gleich(maybeA.(entweder).fallen["nur"], kvar{n: "a"}) {
		t.Fatalf("a isn't an kvar\n\tgot:%s\n\twant:%s", repr.String(maybeInt.(entweder).fallen["nur"]), repr.String(kvar{n: "a"}))
	}
}

func TestFunktionDrei(t *testing.T) {
	ctx := neuKontext()

	maybeT := generischerTyp{
		von: entweder{
			fallen: map[string]typ{
				"nur":    kvar{n: "a"},
				"nichts": nichts{},
			},
		},
		argumenten: []string{"a"},
	}
	maybeInt := maybeT.voll(map[string]typ{
		"a": integer{},
	})

	ctx.fns["anMaybe"] = funktion{
		eingabe: []typ{
			kvar{n: "a"},
		},
		ausgabe: maybeT.voll(map[string]typ{
			"a": kvar{n: "a"},
		}),
	}
	ctx.fns["abMaybe"] = funktion{
		eingabe: []typ{
			maybeT.voll(map[string]typ{
				"a": kvar{n: "a"},
			}),
		},
		ausgabe: kvar{n: "a"},
	}

	expr1 := &parser.Expression{
		Funktionsaufruf: &parser.Funktionsaufruf{
			Name: "anMaybe",
			Argumente: []parser.Expression{
				{
					Integer: &parser.Integer{
						Value: 1,
					},
				},
			},
		},
	}
	expr2 := &parser.Expression{
		Funktionsaufruf: &parser.Funktionsaufruf{
			Name: "abMaybe",
			Argumente: []parser.Expression{
				*expr1,
			},
		},
	}

	checkExpr(ctx, expr1, maybeInt, true, t)
	checkExpr(ctx, expr2, integer{}, true, t)
}

func TestFunktionVier(t *testing.T) {
	ctx := neuKontext()

	maybeT := generischerTyp{
		von: entweder{
			fallen: map[string]typ{
				"nur":    kvar{n: "a"},
				"nichts": nichts{},
			},
		},
		argumenten: []string{"a"},
	}
	maybeInt := maybeT.voll(map[string]typ{
		"a": integer{},
	})

	ctx.fns["anMaybe"] = funktion{
		eingabe: []typ{
			kvar{n: "a"},
		},
		ausgabe: maybeT.voll(map[string]typ{
			"a": kvar{n: "a"},
		}),
	}
	ctx.fns["withDefault"] = funktion{
		eingabe: []typ{
			maybeT.voll(map[string]typ{
				"a": kvar{n: "a"},
			}),
			kvar{n: "a"},
		},
		ausgabe: kvar{n: "a"},
	}

	expr1 := &parser.Expression{
		Funktionsaufruf: &parser.Funktionsaufruf{
			Name: "anMaybe",
			Argumente: []parser.Expression{
				{
					Integer: &parser.Integer{
						Value: 1,
					},
				},
			},
		},
	}
	expr2 := &parser.Expression{
		Funktionsaufruf: &parser.Funktionsaufruf{
			Name: "withDefault",
			Argumente: []parser.Expression{
				*expr1,
				{
					Integer: &parser.Integer{
						Value: 2,
					},
				},
			},
		},
	}
	expr3 := &parser.Expression{
		Funktionsaufruf: &parser.Funktionsaufruf{
			Name: "withDefault",
			Argumente: []parser.Expression{
				*expr1,
				{
					Logik: &parser.Logik{
						Wert: "Wahr",
					},
				},
			},
		},
	}

	checkExpr(ctx, expr1, maybeInt, true, t)
	checkExpr(ctx, expr2, integer{}, true, t)
	checkExpr(ctx, expr3, integer{}, false, t)
}
