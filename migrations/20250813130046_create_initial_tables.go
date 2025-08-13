package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateInitialTables, downCreateInitialTables)
}

func upCreateInitialTables(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.ExecContext(ctx, `
		CREATE TABLE command (
			id TEXT PRIMARY KEY,
			project_id TEXT,
			name TEXT,
			command TEXT,
			working_directory TEXT,
			position INTEGER
		);
        CREATE TABLE project (
            id TEXT PRIMARY KEY,
            name TEXT,
            working_directory TEXT
        );
		CREATE TABLE command_group (
            id TEXT PRIMARY KEY,
            project_id TEXT,
            name TEXT,
            position INTEGER
        );
		CREATE TABLE command_group_command (
            command_group_id TEXT NOT NULL,
            command_id TEXT NOT NULL,
            position INTEGER,
            PRIMARY KEY (command_id, command_group_id)
        );
		CREATE TABLE global_config (
            id TEXT PRIMARY KEY CHECK (id = 1),
            last_opened_project_id TEXT
        );
		CREATE TABLE global_config_extra_path (
            id TEXT,
            path TEXT
        );
	`)

	return err
}

func downCreateInitialTables(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	_, err := tx.ExecContext(ctx, `
		DROP TABLE command;
		DROP TABLE project;
		DROP TABLE command_group;
		DROP TABLE command_group_command;
		DROP TABLE global_config;
		DROP TABLE global_config_extra_path;
	`)
	return err
}
