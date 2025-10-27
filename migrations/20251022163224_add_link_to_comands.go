package migrations

import (
	"context"
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddLinkToComands, downAddLinkToComands)
}

func upAddLinkToComands(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.ExecContext(ctx, `
		ALTER TABLE command ADD COLUMN link TEXT;
	`)
	return err
}

func downAddLinkToComands(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	_, err := tx.ExecContext(ctx, `
		ALTER TABLE command DROP COLUMN link;
	`)
	return err
}
