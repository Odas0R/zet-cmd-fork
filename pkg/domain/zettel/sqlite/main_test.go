package sqlite_test

import (
	"log"
	"testing"

	"github.com/odas0r/zet/pkg/database"
	"github.com/odas0r/zet/pkg/domain/zettel/sqlite"
)

var repo *sqlite.SQLiteRepository

func TestMain(m *testing.M) {
	var err error
	repo, err = sqlite.New(
		database.New(database.Options{
			URL:                "../../../../zettel.db",
			MaxOpenConnections: 1,
			MaxIdleConnections: 1,
			LogQueries:         true,
		}),
	)
	if err != nil {
		log.Fatalf("Failed to set up test repository: %v", err)
	}

	// Run the tests
	m.Run()
}
