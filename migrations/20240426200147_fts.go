package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upFts, downFts)
}

func upFts(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.Exec(`
create virtual table zettel_fts
  using fts5(id, title, content, path, tokenize = porter, content = 'zettel', content_rowid = 'id');

create trigger zettel_after_insert after insert on zettel begin
  insert into zettel_fts(rowid, id, title, content, path)
    values (new.id, new.id, new.title, new.content, new.path);
end;

create trigger zettel_after_update after update on zettel begin
  insert into zettel_fts(zettel_fts, rowid)
    values('delete', old.id);
  insert into zettel_fts(rowid, id, title, content, path)
    values (new.id, new.id, new.title, new.content, new.path);
end;

create trigger zettel_after_delete after delete on zettel begin
  insert into zettel_fts(zettel_fts, rowid)
    values('delete', old.id);
end;
`)
	if err != nil {
		return err
	}
	return nil
}

func downFts(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	_, err := tx.Exec(`
	drop table zettel_fts;
	drop trigger zettel_after_insert;
	drop trigger zettel_after_update;
	drop trigger zettel_after_delete;
	`)
	if err != nil {
		return err
	}
	return nil
}
