package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/odas0r/zet/pkg/controllers"
	"github.com/odas0r/zet/pkg/database"
	"github.com/odas0r/zet/pkg/domain/workspace"
	"github.com/odas0r/zet/pkg/router"
	"github.com/odas0r/zet/pkg/router/middleware"
	"github.com/pressly/goose/v3"
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
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "port",
						Aliases: []string{"p"},
						Value:   "3000",
						Usage:   "Port to listen on",
					},
					&cli.BoolFlag{
						Name:  "dev",
						Value: false,
						Usage: "Enable development mode",
					},
					&cli.StringFlag{
						Name:  "address",
						Value: "localhost",
						Usage: "Address to listen on",
					},
				},
				Action: func(c *cli.Context) error {
					port := c.String("port")
					dev := c.Bool("dev")
					// repo, err := sqlite.New(
					// 	database.New(database.Options{
					// 		URL:                "../../../../zettel.db",
					// 		MaxOpenConnections: 1,
					// 		MaxIdleConnections: 1,
					// 		LogQueries:         true,
					// 	}),
					// )
					// if err != nil {
					// 	return err
					// }
					controller := controllers.NewController(workspace.Workspace{})

					r := router.New()
					r.Use(middleware.WithLogger)

					r.HandleFunc("GET /", controller.HandleHome)

					r.Handle("GET /public/",
						http.StripPrefix("/public/", http.FileServer(http.Dir("public"))),
						middleware.WithDisableCache(dev),
					)

					log.Printf("Listening on :%s\n", port)
					return http.ListenAndServe(fmt.Sprintf(":%s", port), r)
				},
			},
			{
				Name:  "migrate",
				Usage: "Migrate the database to the latest version",
				Subcommands: []*cli.Command{
					{
						Name:  "create",
						Usage: "Create a new migration file",
						Action: func(c *cli.Context) error {
							if len(c.Args().Slice()) == 0 {
								return fmt.Errorf("missing migration name")
							}

							db := database.New(database.Options{
								URL: "./zettel.db",
							})

							if err := db.Connect(); err != nil {
								return err
							}

							name := c.Args().First()
							if err := goose.Create(db.DB.DB, "migrations", name, "go"); err != nil {
								return err
							}

							return nil
						},
					},
					{
						Name:  "up",
						Usage: "Migrate the database to the latest version",
						Action: func(c *cli.Context) error {
							db := database.New(database.Options{
								URL: "./zettel.db",
							})

							if err := db.Connect(); err != nil {
								return err
							}

							provider, err := goose.NewProvider(
								goose.DialectSQLite3,
								db.DB.DB,
								nil, // we're using golang migrations
							)
							if err != nil {
								return err
							}

							results, err := provider.Up(context.Background())
							if err != nil {
								return err
							}

							for _, r := range results {
								fmt.Println(r.String())
							}

							return nil
						},
					},
					{
						Name:  "down",
						Usage: "Downgrade the database by one version",
						Action: func(c *cli.Context) error {
							db := database.New(database.Options{
								URL: "./zettel.db",
							})

							if err := db.Connect(); err != nil {
								return err
							}

							provider, err := goose.NewProvider(
								goose.DialectSQLite3,
								db.DB.DB,
								nil, // we're using golang migrations
							)
							if err != nil {
								return err
							}

							result, err := provider.Down(context.Background())
							if err != nil {
								return err
							}

							fmt.Println(result.String())

							return nil
						},
					},
					{
						Name:  "status",
						Usage: "Print the status of the database",
						Action: func(c *cli.Context) error {
							db := database.New(database.Options{
								URL: "./zettel.db",
							})

							if err := db.Connect(); err != nil {
								return err
							}

							provider, err := goose.NewProvider(
								goose.DialectSQLite3,
								db.DB.DB,
								nil, // we're using golang migrations
							)
							if err != nil {
								return err
							}

							stats, err := provider.Status(context.Background())
							if err != nil {
								return err
							}

							fmt.Println("=== migration status ===")
							for _, s := range stats {
								// get the basename of a path
								file := filepath.Base(s.Source.Path)
								log.Printf("%s %v\n", file, s.State)
							}

							return nil
						},
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
