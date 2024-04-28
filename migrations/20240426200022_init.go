package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upInit, downInit)
}

func upInit(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.Exec(`
create table zettel (
    id text not null primary key,
    title text not null,
    content text not null,
    kind text not null check (kind in ('fleet', 'permanent')),
    created_at text not null default (strftime('%Y-%m-%dT%H:%M:%fZ')),
    updated_at text not null default (strftime('%Y-%m-%dT%H:%M:%fZ')),
    check (created_at <= updated_at),
    check (kind in ('fleet', 'permanent'))
) strict;

create index zettel_created_idx on zettel (created_at);

create trigger zettel_updated_timestamp after update on zettel begin
  -- use ISO8601/RFC3339
  update zettel set updated_at = strftime('%Y-%m-%dT%H:%M:%fZ') where id = old.id;
end;

create table link (
    zettel_id text not null,
    link_id text not null,
    created_at text not null default (strftime('%Y-%m-%dT%H:%M:%fZ')),
    updated_at text not null default (strftime('%Y-%m-%dT%H:%M:%fZ')),

    primary key (zettel_id, link_id),

    foreign key (zettel_id) references zettel(id) on delete cascade,
    foreign key (link_id) references zettel(id) on delete cascade
) strict;

create index link_created_idx on link (created_at);

create trigger link_updated_timestamp after update on link begin
  -- use ISO8601/RFC3339
  update link set updated_at = strftime('%Y-%m-%dT%H:%M:%fZ') where id = old.id;
end;
	`)
	if err != nil {
		return err
	}
	return nil
}

func downInit(ctx context.Context, tx *sql.Tx) error {
	if _, err := tx.Exec(`
drop table link;
drop table zettel;
`); err != nil {
		return err
	}
	return nil
}
