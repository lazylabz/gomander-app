package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddNotificationsEnabledToUserConfig, downAddNotificationsEnabledToUserConfig)
}

func upAddNotificationsEnabledToUserConfig(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `ALTER TABLE user_config ADD COLUMN notifications_enabled BOOLEAN NOT NULL DEFAULT 1;`)

	return err
}

func downAddNotificationsEnabledToUserConfig(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `ALTER TABLE user_config DROP COLUMN notifications_enabled;`)

	return err
}
