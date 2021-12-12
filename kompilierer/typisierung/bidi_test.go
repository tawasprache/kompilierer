package typisierung

import (
	"Tawa/kompilierer/ast"
	"Tawa/kompilierer/fehlerberichtung"
	"embed"
	"testing"

	"github.com/ztrue/tracerr"
)

//go:embed test/*
var testDateien embed.FS

func TestDateien(t *testing.T) {
	files, feh := testDateien.ReadDir("test/erwarte-gut")
	if feh != nil {
		panic("fehler: " + feh.Error())
	}
	for _, datei := range files {
		k := NeuKontext()

		data, feh := testDateien.ReadFile("test/erwarte-gut/" + datei.Name())
		t.Log(datei.Name())
		if feh != nil {
			panic("fehler: " + feh.Error())
		}

		modul := ast.Modul{}
		feh = ast.Parser.ParseBytes(datei.Name(), data, &modul)
		if feh != nil {
			panic(feh)
		}

		_, feh = zuGetypisierteAst(k, "Tawa", modul)
		if feh != nil {
			v, ok := feh.(fehlerberichtung.VerketteterFehler)
			if ok {
				for _, it := range v.Fehler {
					t.Logf("%s", tracerr.Sprint(it))
				}
			} else {
				t.Logf("%s", tracerr.Sprint(feh))
			}
			t.FailNow()
		}
	}
	files, feh = testDateien.ReadDir("test/erwarte-schlecht")
	if feh != nil {
		panic("fehler: " + feh.Error())
	}
	for _, datei := range files {
		k := NeuKontext()

		data, feh := testDateien.ReadFile("test/erwarte-schlecht/" + datei.Name())
		t.Log(datei.Name())
		if feh != nil {
			panic("fehler: " + feh.Error())
		}

		modul := ast.Modul{}
		feh = ast.Parser.ParseBytes(datei.Name(), data, &modul)
		if feh != nil {
			panic(feh)
		}

		_, feh = zuGetypisierteAst(k, "Tawa", modul)
		if feh == nil {
			t.Fatalf("fehler erwartet")
		}
		println(feh.Error())
	}
}
