package main

import (
	"Tawa/codegenerierung"
	"Tawa/parser"
	"Tawa/typisierung"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	a := cli.App{
		Name:  "tawac",
		Usage: "the tawa compiler",
		Commands: []*cli.Command{
			{
				Name:  "compile",
				Usage: "compile a file",
				Action: func(c *cli.Context) error {
					fi, err := os.Open(c.Args().First())
					if err != nil {
						return err
					}
					defer fi.Close()

					dat := parser.Datei{}
					err = parser.Parser.Parse(c.Args().First(), fi, &dat)
					if err != nil {
						return err
					}

					typisierung.PrüfDatei(&dat)
					codegenerierung.CodegenZuDatei(&dat, c.Args().Get(1))

					return nil
				},
			},
		},
	}
	err := a.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
