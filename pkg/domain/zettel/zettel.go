package zettel

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidZettelKind = errors.New("invalid zettel kind")
	ErrMissingValues     = errors.New("missing values")
)

type ZettelKind string

const (
	Permanent ZettelKind = "permanent"
	Fleet     ZettelKind = "fleet"
)

// Zettel is an aggregate root that represents a zettel in the domain
type Zettel struct {
	id        uuid.UUID
	title     string
	content   string
	kind      ZettelKind
	createdAt time.Time
	updatedAt time.Time
}

func New(title, content string, kind ZettelKind) (Zettel, error) {
	if kind != Permanent && kind != Fleet {
		return Zettel{}, ErrInvalidZettelKind
	} else if title == "" || content == "" {
		return Zettel{}, ErrMissingValues
	}

	return Zettel{
		id:        uuid.New(),
		title:     title,
		content:   content,
		kind:      kind,
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}, nil
}

func (z Zettel) ID() uuid.UUID {
	return z.id
}

func (z Zettel) Title() string {
	return z.title
}

func (z Zettel) Content() string {
	return z.content
}

func (z Zettel) Kind() ZettelKind {
	return z.kind
}

func (z Zettel) CreatedAt() time.Time {
	return z.createdAt
}

func (z Zettel) UpdatedAt() time.Time {
	return z.updatedAt
}
