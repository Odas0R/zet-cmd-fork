package zettel

import (
	"context"
	"errors"

	"github.com/odas0r/zet/internal/model"
)

var (
	// ErrZettelNotFound is returned when a zettel is not found.
	ErrZettelNotFound = errors.New("error: zettel not found")
)

type ZettelRepository interface {
	Get(ctx context.Context, zettel *model.Zettel) error
	Save(ctx context.Context, zettel *model.Zettel) error
	SaveBulk(ctx context.Context, zettels ...*model.Zettel) error
	Link(ctx context.Context, zettel *model.Zettel, links []*model.Zettel) error
	LinkBulk(ctx context.Context, links ...*model.Link) error
	Unlink(ctx context.Context, zettel *model.Zettel, links []*model.Zettel) error
	Remove(ctx context.Context, zettel *model.Zettel) error
	RemoveBulk(ctx context.Context, zettels ...*model.Zettel) error
	LastOpened(ctx context.Context, zettel *model.Zettel) error
	InsertHistory(ctx context.Context, zettel *model.Zettel) error
	History(ctx context.Context) ([]*model.Zettel, error)
	ListFleet(ctx context.Context) ([]*model.Zettel, error)
	ListPermanent(ctx context.Context) ([]*model.Zettel, error)
	ListAll(ctx context.Context) ([]*model.Zettel, error)
	Backlinks(ctx context.Context, zet *model.Zettel) ([]*model.Zettel, error)
	Search(ctx context.Context, query string) ([]*model.Zettel, error)
	Reset(ctx context.Context) error
}
