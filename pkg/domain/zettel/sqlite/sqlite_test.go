package sqlite_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/odas0r/zet/pkg/domain/zettel"
	"github.com/odas0r/zet/pkg/domain/zettel/sqlite"
	"github.com/odas0r/zet/pkg/test/sqltest"
)

func TestSQLite_Get(t *testing.T) {
	type testCase struct {
		name        string
		id          uuid.UUID
		expectedErr error
	}

	repo := sqlite.New(
		sqltest.CreateDatabase(t),
	)

	// Create a new zettel
	z, err := zettel.New("title", "content", zettel.Fleet)
	if err != nil {
		t.Fatal(err)
	}

	if err := repo.Save(z); err != nil {
		t.Fatal(err)
	}

	testCases := []testCase{
		{
			name:        "No zettel by ID",
			id:          uuid.MustParse("f47ac10b-58cc-0372-8567-0e02b2c3d479"),
			expectedErr: zettel.ErrZettelNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := repo.FindByID(tc.id)
			if err != tc.expectedErr {
				t.Errorf("expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}
}
