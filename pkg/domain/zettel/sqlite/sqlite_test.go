package sqlite_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/odas0r/zet/pkg/domain/zettel"
)

func TestSQLite_SaveZettel(t *testing.T) {
	z, err := zettel.New("title", "content", zettel.Fleet)
	if err != nil {
		t.Error(err)
	}
	if err := repo.Save(z); err != nil {
		t.Error(err)
	}
}

func TestSQLite_GetZettel(t *testing.T) {
	type testCase struct {
		name        string
		id          uuid.UUID
		expectedErr error
	}

	z, err := zettel.New("title", "content", zettel.Fleet)
	if err != nil {
		t.Error(err)
	}

	if err := repo.Save(z); err != nil {
		t.Error(err)
	}

	testCases := []testCase{
		{
			name:        "Found zettel with ID",
			id:          z.ID(),
			expectedErr: nil,
		},
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

func TestSQLite_UpdateZettel(t *testing.T) {
	type testCase struct {
		name        string
		zettel      func() zettel.Zettel
		expectedErr error
	}

	testCases := []testCase{
		{
			name: "Can update zettel",
			zettel: func() zettel.Zettel {
				zet := createZettel(t)
				zet.SetTitle("new title")
				zet.SetBody("new body")
				return zet
			},
			expectedErr: nil,
		},
		{
			name: "No zettel by ID",
			zettel: func() zettel.Zettel {
				zet := createZettel(t)
				zet.SetID(uuid.MustParse("f47ac10b-58cc-0372-8567-0e02b2c3d479"))
				return zet
			},
			expectedErr: zettel.ErrZettelNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			z := tc.zettel()
			if err := repo.Update(z); err != tc.expectedErr {
				t.Errorf("expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}
}

func TestSQLite_AddLink(t *testing.T) {
	type testCase struct {
		name        string
		setup      func() (zettel.Zettel, zettel.Zettel)
		expectedErr error
	}

	testCases := []testCase{
		{
			name: "Can add link",
			setup: func() (zettel.Zettel, zettel.Zettel) {
				z1 := createZettel(t)
				z2 := createZettel(t)

				z1.AddLink(z2)
				return z1, z2
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			z1, z2 := tc.setup()
			if err := repo.AddLink(z1, z2); err != tc.expectedErr {
				t.Errorf("expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}
}


func createZettel(t *testing.T) zettel.Zettel {
	z, err := zettel.New("title", "content", zettel.Fleet)
	if err != nil {
		t.Error(err)
	}
	if err := repo.Save(z); err != nil {
		t.Error(err)
	}
	return z
}
