package sqltest

import (
	"testing"

	"github.com/odas0r/zet/pkg/database"
	"github.com/pressly/goose"

	// Ensure migrations are imported
	_ "github.com/odas0r/zet/migrations"
)

// CreateDatabase for testing.
func CreateDatabase(t *testing.T) *database.Database {
	t.Helper()

	db := database.NewDatabase(database.DatabaseOptions{
		InMemory:           true,
		MaxOpenConnections: 1,
		MaxIdleConnections: 1,
	})

	if err := db.Connect(); err != nil {
		t.Fatal(err)
	}

	if err := goose.SetDialect("sqlite3"); err != nil {
		t.Fatal(err)
	}

	if err := goose.Up(db.DB.DB, "migrations"); err != nil {
		t.Fatal(err)
	}

	return db
}
