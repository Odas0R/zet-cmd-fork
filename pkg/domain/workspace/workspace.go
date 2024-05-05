package workspace

import (
	"errors"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/odas0r/zet/pkg/domain/shared/timestamp"
	"github.com/odas0r/zet/pkg/domain/zettel"
)

var (
	ErrInvalidPath = errors.New("invalid workspace path")
)

type Workspace struct {
	id        uuid.UUID
	path      string
	zettels   []zettel.Zettel
	timestamp timestamp.Timestamp
}

func New(path string) (Workspace, error) {
	if path == "" {
		return Workspace{}, ErrInvalidPath
	} else if _, err := os.Stat(path); err != nil {
		return Workspace{}, ErrInvalidPath
	}

	return Workspace{
		id:        uuid.New(),
		path:      path,
		zettels:   []zettel.Zettel{},
		timestamp: timestamp.New(),
	}, nil
}

func (w *Workspace) ID() uuid.UUID                  { return w.id }
func (w *Workspace) Path() string                   { return w.path }
func (w *Workspace) Zettels() []zettel.Zettel       { return w.zettels }
func (w *Workspace) Timestamp() timestamp.Timestamp { return w.timestamp }

func (w *Workspace) SetID(id uuid.UUID)           { w.id = id }
func (w *Workspace) SetPath(path string)          { w.path = path }
func (w *Workspace) SetCreated(created time.Time) { w.timestamp.Created = created }
func (w *Workspace) SetUpdated(updated time.Time) { w.timestamp.Updated = updated }
func (w *Workspace) SetZettels(z []zettel.Zettel) { w.zettels = z }
