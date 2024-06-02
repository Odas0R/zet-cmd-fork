package workspace

import (
	"errors"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/odas0r/zet/pkg/domain/shared/timestamp"
)

var (
	ErrInvalidPath    = errors.New("invalid workspace path")
	ErrZettelNotFound = errors.New("zettel not found in workspace")
)

type Workspace struct {
	id        uuid.UUID
	path      string
	timestamp timestamp.Timestamp
	zettelIDs map[uuid.UUID]struct{}
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
		timestamp: timestamp.New(),
		zettelIDs: make(map[uuid.UUID]struct{}),
	}, nil
}

// Getters
func (w *Workspace) ID() uuid.UUID                  { return w.id }
func (w *Workspace) Path() string                   { return w.path }
func (w *Workspace) Timestamp() timestamp.Timestamp { return w.timestamp }

// Setters
func (w *Workspace) SetID(id uuid.UUID)           { w.id = id }
func (w *Workspace) SetPath(path string)          { w.path = path }
func (w *Workspace) SetCreated(created time.Time) { w.timestamp.Created = created }
func (w *Workspace) SetUpdated(updated time.Time) { w.timestamp.Updated = updated }

func (w *Workspace) AddZettel(zID uuid.UUID) {
	w.zettelIDs[zID] = struct{}{}
}

func (w *Workspace) RemoveZettel(id uuid.UUID) error {
	if _, exists := w.zettelIDs[id]; !exists {
		return ErrZettelNotFound
	}
	delete(w.zettelIDs, id)
	return nil
}

func (w *Workspace) ListZettelIDs() []uuid.UUID {
	zs := make([]uuid.UUID, 0, len(w.zettelIDs))
	for id := range w.zettelIDs {
		zs = append(zs, id)
	}
	return zs
}
