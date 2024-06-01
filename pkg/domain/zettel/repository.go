package zettel

import (
	"errors"

	"github.com/google/uuid"
)

var (
	// ErrZettelNotFound is returned when a zettel is not found.
	ErrZettelNotFound = errors.New("error: zettel not found")
)

// TODO: @Guilherme
type ZettelRepository interface {
	FindByID(id uuid.UUID) (Zettel, error)
	Save(zettel Zettel) error
	Update(z Zettel) error
	Delete(id uuid.UUID) error
	AddLink(from, to uuid.UUID) error
	RemoveLink(from, to uuid.UUID) error
}
