package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upWorkspace, downWorkspace)
}

func upWorkspace(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`
-- create the workspace table
create table if not exists workspace (
    id text primary key,
    path text not null,
    created_at text not null,
    updated_at text not null
);

-- create the workspace_zettel table
create table if not exists workspace_zettel (
    workspace_id text not null,
    zettel_id text not null,
    primary key (workspace_id, zettel_id),
    foreign key (workspace_id) references workspace(id) on delete cascade,
    foreign key (zettel_id) references zettel(id) on delete cascade
);

-- create indexes to improve query performance
create index if not exists idx_workspace_zettel_workspace_id on workspace_zettel(workspace_id);
create index if not exists idx_workspace_zettel_zettel_id on workspace_zettel(zettel_id);
	`)
	if err != nil {
		return err
	}

	return nil
}

func downWorkspace(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`
-- drop the workspace tables
drop table workspace_zettel;
drop table workspace;
`)
	if err != nil {
		return err
	}
	return nil
}
