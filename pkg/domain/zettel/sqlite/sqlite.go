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
	ID        uuid.UUID   `db:"id"`
	Title     string      `db:"title"`
	Content   string      `db:"content"`
	Kind      string      `db:"kind"`
	CreatedAt sqlite.Time `db:"created_at"`
	UpdatedAt sqlite.Time `db:"updated_at"`
}

// NewFromZettel takes in an aggregate root and returns a struct that can be
// used to interact with the database
func NewFromZettel(z zettel.Zettel) sqliteZettel {
	return sqliteZettel{
		ID:      z.ID(),
		Title:   z.Title(),
		Content: z.Content(),
		Kind:    string(z.Kind()),
	}
}

// ToAggregate converts the struct to the aggregate root
func (sz sqliteZettel) ToAggregate() zettel.Zettel {
	z := zettel.Zettel{}

	z.SetID(sz.ID)
	z.SetTitle(sz.Title)
	z.SetBody(sz.Content)
	// cast the string to the Kind type
	z.SetKind(zettel.Kind(sz.Kind))
	// z.SetCreatedAt(sz.createdAt.T)
	// z.SetUpdatedAt(sz.updatedAt.T)

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

func (r *SQLiteRepository) Save(z zettel.Zettel) error {
	internal := NewFromZettel(z)

	query := `
  insert into zettel (id, title, content, kind)
	values (:id, :title, :content, :kind)
	on conflict (id) do
	update set title = excluded.title, content = excluded.content, kind = excluded.kind
  `

	_, err := r.db.NamedExec(query, internal)
	return err
}

func (r *SQLiteRepository) Update(z zettel.Zettel) error {
	sz := NewFromZettel(z)

	query := `
  update zettel
  set title = :title, content = :content, kind = :kind
  where id = :id
  `
	_, err := r.db.NamedExec(query, sz)
	return err
}

func (r *SQLiteRepository) Delete(id uuid.UUID) error {
	query := `
  delete from zettel
  where id = $1
  `

	_, err := r.db.Exec(query, id)
	return err
}
