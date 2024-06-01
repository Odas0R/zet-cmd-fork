package sqlite

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/odas0r/zet/pkg/database"
	"github.com/odas0r/zet/pkg/domain/shared/sqlite"
	"github.com/odas0r/zet/pkg/domain/zettel"
)

type SQLiteRepository struct {
	db *sqlx.DB
}

type sqliteZettel struct {
	ID      uuid.UUID    `db:"id"`
	Title   string       `db:"title"`
	Content string       `db:"content"`
	Kind    zettel.Kind  `db:"kind"`
	Created *sqlite.Time `db:"created_at"`
	Updated *sqlite.Time `db:"updated_at"`
}

// NewFromZettel takes in an aggregate root and returns a struct that can be
// used to interact with the database
func NewFromZettel(z zettel.Zettel) sqliteZettel {
	return sqliteZettel{
		ID:      z.ID(),
		Title:   z.Title(),
		Content: z.Content(),
		Kind:    z.Kind(),
		Created: &sqlite.Time{T: z.Timestamp().Created},
		Updated: &sqlite.Time{T: z.Timestamp().Updated},
	}
}

// ToAggregate converts the struct to the aggregate root
func (sz sqliteZettel) ToAggregate() zettel.Zettel {
	z := zettel.Zettel{}

	z.SetID(sz.ID)
	z.SetTitle(sz.Title)
	z.SetBody(sz.Content)
	z.SetKind(sz.Kind)
	z.SetCreated(sz.Created.T)
	z.SetUpdated(sz.Updated.T)

	return z
}

func New(database *database.Database) (*SQLiteRepository, error) {
	if err := database.Connect(); err != nil {
		return nil, err
	}

	return &SQLiteRepository{
		db: database.DB,
	}, nil
}

func (r *SQLiteRepository) FindByID(id uuid.UUID) (zettel.Zettel, error) {
	var sz sqliteZettel

	query := `
  select id, title, content, kind, created_at, updated_at
  from zettel
  where id = $1
  `

	if err := r.db.Get(&sz, query, id); err != nil {
		if err == sql.ErrNoRows {
			return zettel.Zettel{}, zettel.ErrZettelNotFound
		}
		return zettel.Zettel{}, err
	}

	return sz.ToAggregate(), nil
}

// Save takes in an aggregate root and saves it to the database as an upsert
func (r *SQLiteRepository) Save(z zettel.Zettel) error {
	internal := NewFromZettel(z)

	query := `
  insert into zettel (id, title, content, kind, updated_at, created_at)
	values (:id, :title, :content, :kind, :updated_at, :created_at)
	on conflict (id) do
	update set title = excluded.title, content = excluded.content, kind = excluded.kind, updated_at = excluded.updated_at
  `

	_, err := r.db.NamedExec(query, internal)
	return err
}

func (r *SQLiteRepository) Update(z zettel.Zettel) error {
	internal := NewFromZettel(z)

	query := `
	update zettel
	set title = :title, content = :content, kind = :kind, updated_at = :updated_at
	where id = :id
	`

	result, err := r.db.NamedExec(query, internal)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return zettel.ErrZettelNotFound
	}

	return nil
}

func (r *SQLiteRepository) Delete(id uuid.UUID) error {
	query := `
  delete from zettel
  where id = $1
  `

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return zettel.ErrZettelNotFound
	}

	return nil
}

func (r *SQLiteRepository) AddLink(from, to uuid.UUID) error {
	query := `
        insert into link (zettel_id, link_id, created_at, updated_at)
        values ($1, $2, strftime('%y-%m-%dt%h:%m:%fz'), strftime('%y-%m-%dt%h:%m:%fz'))
        on conflict (zettel_id, link_id) do
				update set updated_at = excluded.updated_at;
    `
	_, err := r.db.Exec(query, from, to)
	return err
}

func (r *SQLiteRepository) RemoveLink(from, to uuid.UUID) error {
	query := `
        DELETE FROM link
        WHERE zettel_id = $1 AND link_id = $2;
    `
	_, err := r.db.Exec(query, from, to)
	return err
}
