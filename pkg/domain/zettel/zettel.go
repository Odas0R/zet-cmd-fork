package zettel

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/odas0r/zet/pkg/domain/shared/timestamp"
)

var (
	ErrInvalidZettelKind = errors.New("invalid zettel kind")
	ErrMissingValues     = errors.New("missing values")
)

// Zettel is an aggregate root that represents a zettel in the domain
type Zettel struct {
	id        uuid.UUID
	content   *Content
	kind      Kind
	timestamp timestamp.Timestamp
}

func New(title, content string, kind Kind) (Zettel, error) {
	if kind != Permanent && kind != Fleet {
		return Zettel{}, ErrInvalidZettelKind
	} else if title == "" || content == "" {
		return Zettel{}, ErrMissingValues
	}

	return Zettel{
		id: uuid.New(),
		content: &Content{
			Body:  content,
			Title: title,
		},
		kind:      kind,
		timestamp: timestamp.New(),
	}, nil
}

func (z *Zettel) ID() uuid.UUID                  { return z.id }
func (z *Zettel) Title() string                  { return z.content.Title }
func (z *Zettel) Content() string                { return z.content.Body }
func (z *Zettel) Kind() Kind                     { return z.kind }
func (z *Zettel) Timestamp() timestamp.Timestamp { return z.timestamp }

func (z *Zettel) SetID(id uuid.UUID)           { z.id = id }
func (z *Zettel) SetKind(kind Kind)            { z.kind = kind }
func (z *Zettel) SetCreated(created time.Time) { z.timestamp.Created = created }
func (z *Zettel) SetUpdated(updated time.Time) { z.timestamp.Updated = updated }
func (z *Zettel) SetTitle(title string) {
	if z.content == nil {
		z.content = &Content{}
	}
	z.content.Title = title
}
func (z *Zettel) SetBody(body string) {
	if z.content == nil {
		z.content = &Content{}
	}
	z.content.Body = body
}
