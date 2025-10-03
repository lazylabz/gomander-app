package migrations

import (
	"context"
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddLogLineLimitToUserConfig, downAddLogLineLimitToUserConfig)
}

func upAddLogLineLimitToUserConfig(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.ExecContext(ctx, `
		ALTER TABLE user_config ADD COLUMN log_line_limit INTEGER DEFAULT 100;
	`)
	return err
}

func downAddLogLineLimitToUserConfig(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	_, err := tx.ExecContext(ctx, `
		ALTER TABLE user_config DROP COLUMN log_line_limit;
	`)
	return err
}
