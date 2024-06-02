package sqlite

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/odas0r/zet/pkg/database"
	"github.com/odas0r/zet/pkg/domain/shared/sqlite"
	"github.com/odas0r/zet/pkg/domain/shared/timestamp"
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

	Links []sqliteLink `db:"-"`
}

type sqliteLink struct {
	From      uuid.UUID    `db:"zettel_id"`
	To        uuid.UUID    `db:"link_id"`
	CreatedAt *sqlite.Time `db:"created_at"`
	UpdatedAt *sqlite.Time `db:"updated_at"`
}

// NewFromZettel takes in an aggregate root and returns a struct that can be
// used to interact with the database
func NewFromZettel(z zettel.Zettel) sqliteZettel {
	var links []sqliteLink
	for _, link := range z.Links() {
		links = append(links, sqliteLink{
			From:      link.From,
			To:        link.To,
			CreatedAt: &sqlite.Time{T: link.Timestamp.Created},
			UpdatedAt: &sqlite.Time{T: link.Timestamp.Updated},
		})
	}

	return sqliteZettel{
		ID:      z.ID(),
		Title:   z.Title(),
		Content: z.Content(),
		Kind:    z.Kind(),
		Created: &sqlite.Time{T: z.Timestamp().Created},
		Updated: &sqlite.Time{T: z.Timestamp().Updated},
		Links:   links,
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

	var domainLinks []zettel.Link
	for _, sl := range sz.Links {
		domainLinks = append(domainLinks, zettel.Link{
			From: sl.From,
			To:   sl.To,
			Timestamp: timestamp.Timestamp{
				Created: sl.CreatedAt.T,
				Updated: sl.UpdatedAt.T,
			},
		})
	}
	z.SetLinks(domainLinks)

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

	// Fetch links
	linksQuery := `
  select zettel_id, link_id, created_at, updated_at
  from link
  where zettel_id = $1
  `
	var links []sqliteLink
	if err := r.db.Select(&links, linksQuery, id); err != nil {
		return zettel.Zettel{}, err
	}

	// Set links in the struct
	sz.Links = links

	return sz.ToAggregate(), nil
}

func (r *SQLiteRepository) Save(z zettel.Zettel) error {
	internal := NewFromZettel(z)

	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	query := `
  insert into zettel (id, title, content, kind, updated_at, created_at)
	values (:id, :title, :content, :kind, :updated_at, :created_at)
	on conflict (id) do
	update set title = excluded.title, content = excluded.content, kind = excluded.kind, updated_at = excluded.updated_at
  `

	_, err = tx.NamedExec(query, internal)
	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			return zettel.ErrZettelNotFound
		}
		return err
	}

	err = r.saveLinks(tx, internal.ID, internal.Links)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *SQLiteRepository) saveLinks(tx *sqlx.Tx, zettelID uuid.UUID, links []sqliteLink) error {
	delQuery := `delete from link where zettel_id = $1`
	_, err := tx.Exec(delQuery, zettelID)
	if err != nil {
		return err
	}

	insQuery := `
  insert into link (zettel_id, link_id, created_at, updated_at)
  values (:zettel_id, :link_id, :created_at, :updated_at)
  `
	for _, link := range links {
		_, err = tx.NamedExec(insQuery, link)
		if err != nil {
			return err
		}
	}
	return nil
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
