package zettel

import (
	"errors"
)

var (
	// ErrZettelNotFound is returned when a zettel is not found.
	ErrZettelNotFound = errors.New("error: zettel not found")
)

// TODO: @Guilherme
type Repository interface {
	Get(zettel Zettel) error
	Save(zettel Zettel) error
	SaveBulk(zettels ...Zettel) error
	RemoveBulk(zettels ...Zettel) error
	LastOpened(zettel Zettel) error
	InsertHistory(zettel Zettel) error
	History() ([]Zettel, error)
	List() ([]Zettel, error)
	ListFleet() ([]Zettel, error)
	ListPermanent() ([]Zettel, error)
}
