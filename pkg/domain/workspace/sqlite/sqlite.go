package sqlite

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/odas0r/zet/pkg/domain/shared/sqlite"
	"github.com/odas0r/zet/pkg/domain/workspace"
)

type SQLiteWorkspaceRepository struct {
	db *sqlx.DB
}

type sqliteWorkspace struct {
	ID        uuid.UUID    `db:"id"`
	Path      string       `db:"path"`
	CreatedAt *sqlite.Time `db:"created_at"`
	UpdatedAt *sqlite.Time `db:"updated_at"`
}

type sqliteWorkspaceZettel struct {
	WorkspaceID uuid.UUID `db:"workspace_id"`
	ZettelID    uuid.UUID `db:"zettel_id"`
}

func NewFromWorkspace(w workspace.Workspace) sqliteWorkspace {
	return sqliteWorkspace{
		ID:        w.ID(),
		Path:      w.Path(),
		CreatedAt: &sqlite.Time{T: w.Timestamp().Created},
		UpdatedAt: &sqlite.Time{T: w.Timestamp().Updated},
	}
}

func (sw sqliteWorkspace) ToAggregate(zettelIDs []uuid.UUID) workspace.Workspace {
	w := workspace.Workspace{}
	w.SetID(sw.ID)
	w.SetPath(sw.Path)
	w.SetCreated(sw.CreatedAt.T)
	w.SetUpdated(sw.UpdatedAt.T)

	for _, zID := range zettelIDs {
		w.AddZettel(zID)
	}

	return w
}

func NewWorkspaceRepository(db *sqlx.DB) *SQLiteWorkspaceRepository {
	return &SQLiteWorkspaceRepository{db: db}
}

func (r *SQLiteWorkspaceRepository) FindWorkspaceByID(id uuid.UUID) (workspace.Workspace, error) {
	var sw sqliteWorkspace

	query := `
  select id, path, created_at, updated_at
  from workspace
  where id = $1
  `

	if err := r.db.Get(&sw, query, id); err != nil {
		if err == sql.ErrNoRows {
			return workspace.Workspace{}, workspace.ErrWorkspaceNotFound
		}
		return workspace.Workspace{}, err
	}

	// Fetch workspace zettels
	zettelsQuery := `
  select zettel_id
  from workspace_zettel
  where workspace_id = $1
  `
	var zettelIDs []uuid.UUID
	if err := r.db.Select(&zettelIDs, zettelsQuery, id); err != nil {
		return workspace.Workspace{}, err
	}

	return sw.ToAggregate(zettelIDs), nil
}

func (r *SQLiteWorkspaceRepository) FindAllWorkspaces() ([]workspace.Workspace, error) {
	query := `
  SELECT id, path, created_at, updated_at
  FROM workspace
  `

	var results []sqliteWorkspace
	if err := r.db.Select(&results, query); err != nil {
		return nil, err
	}

	workspaces := make([]workspace.Workspace, len(results))
	for i, row := range results {
		zettelIDs, err := r.findZettelIDsByWorkspaceID(row.ID)
		if err != nil {
			return nil, err
		}
		workspaces[i] = row.ToAggregate(zettelIDs)
	}

	return workspaces, nil
}

func (r *SQLiteWorkspaceRepository) findZettelIDsByWorkspaceID(workspaceID uuid.UUID) ([]uuid.UUID, error) {
	query := `
  select zettel_id
  from workspace_zettel
  where workspace_id = $1
  `
	var zettelIDs []uuid.UUID
	if err := r.db.Select(&zettelIDs, query, workspaceID); err != nil {
		return nil, err
	}
	return zettelIDs, nil
}

func (r *SQLiteWorkspaceRepository) Save(w workspace.Workspace) error {
	internal := NewFromWorkspace(w)

	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	query := `
  insert into workspace (id, path, created_at, updated_at)
	values (:id, :path, :created_at, :updated_at)
	on conflict (id) do
	update set path = excluded.path, updated_at = excluded.updated_at
  `

	_, err = tx.NamedExec(query, internal)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Save workspace zettels
	err = r.saveWorkspaceZettels(tx, internal.ID, w.ListZettelIDs())
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *SQLiteWorkspaceRepository) Update(w workspace.Workspace) error {
	internal := NewFromWorkspace(w)

	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	query := `
  update workspace
  set path = :path, updated_at = :updated_at
  where id = :id
  `

	_, err = tx.NamedExec(query, internal)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Save workspace zettels
	err = r.saveWorkspaceZettels(tx, internal.ID, w.ListZettelIDs())
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *SQLiteWorkspaceRepository) Delete(id uuid.UUID) error {
	query := `delete from workspace where id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *SQLiteWorkspaceRepository) saveWorkspaceZettels(tx *sqlx.Tx, workspaceID uuid.UUID, zettelIDs []uuid.UUID) error {
	// Delete existing workspace zettels
	delQuery := `delete from workspace_zettel where workspace_id = $1`
	_, err := tx.Exec(delQuery, workspaceID)
	if err != nil {
		return err
	}

	// Insert new workspace zettels
	insQuery := `
  insert into workspace_zettel (workspace_id, zettel_id)
  values ($1, $2)
  `
	for _, zID := range zettelIDs {
		_, err = tx.Exec(insQuery, workspaceID, zID)
		if err != nil {
			return err
		}
	}
	return nil
}
