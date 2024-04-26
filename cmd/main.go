package main

import (
	"errors"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:                 "zet",
		Version:              "0.0.1",
		Usage:                "A zettelkasten under a terminal approach",
		UsageText:            "A simple way to manage your zettelkasten",
		Flags:                []cli.Flag{},
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			{
				Name:  "serve",
				Usage: "Starts the web server",
				Action: func(c *cli.Context) error {
					if c.NArg() == 0 {
						return nil
					}

					return errors.New("TODO: implement")
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
