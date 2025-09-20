package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddLocaleToUserConfig, downAddLocaleToUserConfig)
}

func upAddLocaleToUserConfig(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.ExecContext(ctx, `
		ALTER TABLE user_config ADD COLUMN locale TEXT DEFAULT 'en';
	`)

	return err
}

func downAddLocaleToUserConfig(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	_, err := tx.ExecContext(ctx, `
		ALTER TABLE user_config DROP COLUMN locale;
	`)
	return err
}