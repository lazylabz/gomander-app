package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddFailurePatterns, downAddFailurePatterns)
}

// upAddFailurePatterns adds the failure_patterns column to the project table.
func upAddFailurePatterns(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		ALTER TABLE project 
		ADD COLUMN failure_patterns TEXT
	`)
	return err
}

// downAddFailurePatterns removes the failure_patterns column from the project table.
func downAddFailurePatterns(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		ALTER TABLE project 
		DROP COLUMN failure_patterns
	`)
	return err
}
