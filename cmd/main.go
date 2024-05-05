package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/odas0r/zet/pkg/controllers"
	"github.com/odas0r/zet/pkg/database"
	"github.com/odas0r/zet/pkg/domain/workspace"
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
				Action: func(c *cli.Context) error {
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

					mux := http.NewServeMux()
					mux.HandleFunc("/", controller.HandleHome)
					mux.HandleFunc("/create", controller.HandleCreateForm)
					mux.HandleFunc("/create", controller.HandleCreate)
					mux.HandleFunc("/archive/{id}", controller.HandleArchive)
					mux.HandleFunc("/initialize", controller.HandleInitializeForm)
					mux.HandleFunc("/initialize", controller.HandleInitialize)

					log.Println("Listening on :3000")
					return http.ListenAndServe(":3000", mux)
				},
			},
			{
				Name:  "migrate",
				Usage: "Migrate the database to the latest version",
				Subcommands: []*cli.Command{
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
								log.Printf("%-3s %-2v %v\n", s.Source.Type, s.Source.Version, s.State)
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
