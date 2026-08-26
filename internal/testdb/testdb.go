// Package testdb builds the database the tests run against: an in-memory
// SQLite instance with the project's own migrations applied.
package testdb

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	_ "gomander/migrations"
)

var (
	counter atomic.Uint64

	// Goose keeps both the dialect and the registered migrations in package
	// level state, and migrating mutates the latter, so parallel tests take
	// turns at it.
	migrating sync.Mutex
)

// New opens a database private to the calling test and migrates it. The
// database only exists while the connection is open, so it is closed on
// cleanup.
func New(t *testing.T) *gorm.DB {
	t.Helper()

	// A named in-memory database with a shared cache: every connection in the
	// pool reaches the same data, and a unique name keeps tests out of each
	// other's.
	dsn := fmt.Sprintf("file:testdb%d?mode=memory&cache=shared", counter.Add(1))

	gormDb, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Error),
	})
	if err != nil {
		t.Fatalf("failed to open the test database: %v", err)
	}

	db, err := gormDb.DB()
	if err != nil {
		t.Fatalf("failed to reach the test database connection: %v", err)
	}

	// One connection, as the app does: SQLite takes one writer at a time, and
	// holding the connection open is what keeps the database alive.
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close the test database: %v", err)
		}
	})

	migrating.Lock()
	defer migrating.Unlock()

	if err := goose.SetDialect("sqlite3"); err != nil {
		t.Fatalf("failed to set the migration dialect: %v", err)
	}

	if err := goose.UpContext(context.Background(), db, "."); err != nil {
		t.Fatalf("failed to migrate the test database: %v", err)
	}

	return gormDb
}
