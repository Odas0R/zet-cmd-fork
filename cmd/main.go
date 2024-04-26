package main

import (
	"errors"
	"log"
	"os"

	"github.com/urfave/cli/v2"

	// Ensure migrations are imported
	_ "github.com/odas0r/zet/migrations"
)

func main() {
	app := &cli.App{
		Name:                 "zet",
		Version:              "0.0.1",
		Usage:                "A simple way to manage your zettelkasten, via command line and web interface.",
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

					// TODO: initialize the database migrations

					return errors.New("TODO: implement")
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
