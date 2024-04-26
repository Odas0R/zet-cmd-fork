package model

type Zettel struct {
	ID        string `db:"id"`
	Title     string `db:"title"`
	Content   string `db:"content"`
	Path      string `db:"path"`
	Type      string `db:"type"`
	CreatedAt Time   `db:"created_at"`
	UpdatedAt Time   `db:"updated_at"`

	// Auxiliary fields (not stored in the database)
	Lines []string
	Links []*Zettel
}
