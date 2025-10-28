package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddErrorPatternToCommands, downAddErrorPatternToCommands)
}

func upAddErrorPatternToCommands(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.ExecContext(ctx, `
		ALTER TABLE command ADD COLUMN error_patterns TEXT;
	`)
	return err
}

func downAddErrorPatternToCommands(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	_, err := tx.ExecContext(ctx, `
		ALTER TABLE command DROP COLUMN error_patterns;
	`)
	return err
}
