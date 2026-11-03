package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddLogsFontSizeAndZoomLevelToUserConfig, downAddLogsFontSizeAndZoomLevelToUserConfig)
}

func upAddLogsFontSizeAndZoomLevelToUserConfig(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		ALTER TABLE user_config ADD COLUMN logs_font_size INTEGER NOT NULL DEFAULT 15;
		ALTER TABLE user_config ADD COLUMN zoom_level REAL NOT NULL DEFAULT 1;
	`)

	return err
}

func downAddLogsFontSizeAndZoomLevelToUserConfig(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		ALTER TABLE user_config DROP COLUMN zoom_level;
		ALTER TABLE user_config DROP COLUMN logs_font_size;
	`)

	return err
}
