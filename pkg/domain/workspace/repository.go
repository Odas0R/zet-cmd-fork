package workspace

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrWorkspaceNotFound = errors.New("workspace not found")
)

type Repository interface {
	FindWorkspaceByID(id uuid.UUID) (Workspace, error)
	FindAllWorkspaces() ([]Workspace, error)
	Save(workspace Workspace) error
	Update(w Workspace) error
	Delete(id uuid.UUID) error
}
