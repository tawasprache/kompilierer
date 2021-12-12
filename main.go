package main

import (
	"Tawa/dokumentation"
	"Tawa/kompilierer/ast"
	"Tawa/kompilierer/codegenierung"
	"Tawa/kompilierer/typisierung"
	"Tawa/langserv"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/alecthomas/repr"
	"github.com/urfave/cli/v2"

	_ "Tawa/kompilierer/codegenierung/typescript"
)

func main() {
	a := cli.App{
		Name:  "tawac",
		Usage: "the tawa parser",
		Commands: []*cli.Command{
			{
				Name:  "parse",
				Usage: "parse a file",
				Action: func(c *cli.Context) error {
					fi, err := os.Open(c.Args().First())
					if err != nil {
						return err
					}
					defer fi.Close()

					dat := ast.Modul{}
					err = ast.Parser.Parse(c.Args().First(), fi, &dat)
					if err != nil {
						return err
					}

					repr.Println(dat)

					return nil
				},
			},
			{
				Name:  "lsp",
				Usage: "lsp",
				Action: func(c *cli.Context) error {
					langserv.StartServer()
					return nil
				},
			},
			{
				Name: "grammar",
				Action: func(c *cli.Context) error {
					println(ast.Parser.String())
					println(ast.TerminalParser.String())
					return nil
				},
			},
			{
				Name:  "typecheck",
				Usage: "typecheck a file",
				Action: func(c *cli.Context) error {
					fi, err := os.Open(c.Args().First())
					if err != nil {
						return err
					}
					defer fi.Close()

					dat := ast.Modul{}
					err = ast.Parser.Parse(c.Args().First(), fi, &dat)
					if err != nil {
						return err
					}

					ktx := typisierung.NeuKontext()
					genannt, err := typisierung.Auflösenamen(ktx, dat, "User")
					if err != nil {
						return err
					}

					getypt, err := typisierung.Typiere(ktx, genannt, "User")
					if err != nil {
						return err
					}

					_ = getypt
					return nil
				},
			},
			{
				Name:  "doku",
				Usage: "doku",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name: "standardbibliothek",
					},
					&cli.StringFlag{
						Name:     "ausgabe-ordner",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					ktx := typisierung.NeuKontext()
					opath := c.String("ausgabe-ordner")

					if c.IsSet("standardbibliothek") {
						for m, v := range ktx.Module {
							feh := ioutil.WriteFile(path.Join(opath, strings.ReplaceAll(m, "/", ":")+".html"), []byte(dokumentation.Dokumentation(v)), 0o666)
							if feh != nil {
								return feh
							}
						}
					} else {
						fi, err := os.Open(c.Args().First())
						if err != nil {
							return err
						}
						defer fi.Close()

						dat := ast.Modul{}
						err = ast.Parser.Parse(c.Args().First(), fi, &dat)
						if err != nil {
							return err
						}

						genannt, err := typisierung.Auflösenamen(ktx, dat, "User")
						if err != nil {
							return err
						}

						getypt, err := typisierung.Typiere(ktx, genannt, "User")
						if err != nil {
							return err
						}

						feh := ioutil.WriteFile(path.Join(opath, strings.ReplaceAll(getypt.Name, "/", ":")+".html"), []byte(dokumentation.Dokumentation(getypt)), 0o666)
						if feh != nil {
							return feh
						}
					}

					return nil
				},
			},
			{
				Name:  "compile",
				Usage: "compile a file",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "backend",
						Required: true,
					},
					&cli.StringFlag{
						Name: "js-out",
					},
					&cli.StringFlag{
						Name: "html-out",
					},
				},
				Action: func(c *cli.Context) error {
					fi, err := os.Open(c.Args().First())
					if err != nil {
						return err
					}
					defer fi.Close()

					dat := ast.Modul{}
					err = ast.Parser.Parse(c.Args().First(), fi, &dat)
					if err != nil {
						return err
					}

					ktx := typisierung.NeuKontext()
					genannt, err := typisierung.Auflösenamen(ktx, dat, "User")
					if err != nil {
						return err
					}

					getypt, err := typisierung.Typiere(ktx, genannt, "User")
					if err != nil {
						return err
					}

					ktx.Module[getypt.Name] = getypt

					os.Mkdir(c.Args().Get(1), 0o777)

					unterbau := codegenierung.GetUnterbau(c.String("backend"))

					o := codegenierung.Optionen{
						Outpath:     c.Args().Get(1),
						JSOutfile:   c.String("js-out"),
						HTMLOutfile: c.String("html-out"),
						Entry:       getypt.Name,
					}

					feh := unterbau.Pregen(o)
					if feh != nil {
						return feh
					}

					var modulen = []string{}
					for _, it := range getypt.Dependencies {
						modulen = append(modulen, it.Paket)
					}

					for _, it := range append(modulen, getypt.Name) {
						feh := unterbau.CodegenModul(o, ktx.Module[it])
						if feh != nil {
							return feh
						}
					}

					feh = unterbau.Postgen(o)
					if feh != nil {
						return feh
					}

					return nil
				},
			},
		},
	}
	err := a.Run(os.Args)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
}
