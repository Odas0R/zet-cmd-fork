package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upHistory, downHistory)
}

func upHistory(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.Exec(`
	create table history (
    zettel_id text not null primary key,
    created_at text not null default (strftime('%Y-%m-%dT%H:%M:%fZ')), -- use ISO8601/RFC3339

    foreign key (zettel_id) references zettel(id) on delete cascade
) strict;
`)
	if err != nil {
		return err
	}

	return nil
}

func downHistory(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`
drop table history;
`)
	if err != nil {
		return err
	}
	return nil
}
