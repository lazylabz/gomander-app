package infrastructure

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"gomander/internal/config/domain"
	_ "gomander/migrations"
)

type testHelper struct {
	t      *testing.T
	repo   *GormConfigRepository
	dbPath string
}

func newTestHelper(t *testing.T,
	preloadedConfig *ConfigModel, preloadedPaths []*EnvironmentPathModel) *testHelper {
	t.Helper() // IMPORTANT: This marks the function as a helper, so error traces will point to the test instead of here

	repo, dbPath := arrange(
		preloadedConfig,
		preloadedPaths,
	)

	helper := &testHelper{
		t:      t,
		repo:   repo,
		dbPath: dbPath,
	}

	// Automatic cleanup when test finishes
	t.Cleanup(func() {
		assert.NoError(t, os.Remove(helper.dbPath), "Failed to cleanup test database")
	})

	return helper
}

func TestGormConfigRepository_GetOrCreate(t *testing.T) {
	t.Run("Should create config if not exists", func(t *testing.T) {
		helper := newTestHelper(t, nil, nil)
		config, err := helper.repo.GetOrCreate()
		assert.NoError(t, err)
		assert.NotNil(t, config)
		assert.Equal(t, "", config.LastOpenedProjectId)
		assert.Empty(t, config.EnvironmentPaths)
	})
	t.Run("Should return existing config with environment paths", func(t *testing.T) {
		preloadedConfig := &ConfigModel{Id: 1, LastOpenedProjectId: "proj-123"}
		preloadedPaths := []*EnvironmentPathModel{
			{Id: "path1", Path: "/usr/bin"},
			{Id: "path2", Path: "/usr/local/bin"},
		}
		helper := newTestHelper(t, preloadedConfig, preloadedPaths)
		config, err := helper.repo.GetOrCreate()
		assert.NoError(t, err)
		assert.Equal(t, &domain.Config{
			LastOpenedProjectId: "proj-123",
			EnvironmentPaths: []domain.EnvironmentPath{
				{Id: "path1", Path: "/usr/bin"},
				{Id: "path2", Path: "/usr/local/bin"},
			}}, config)
	})
}

func TestGormConfigRepository_Update(t *testing.T) {
	t.Run("Should save config and environment paths", func(t *testing.T) {
		preloadedConfig := &ConfigModel{Id: 1, LastOpenedProjectId: "proj-123"}
		preloadedPaths := []*EnvironmentPathModel{
			{Id: "path1", Path: "/usr/bin"},
			{Id: "path2", Path: "/usr/local/bin"},
		}

		helper := newTestHelper(t, preloadedConfig, preloadedPaths)

		newConfig := &domain.Config{
			LastOpenedProjectId: "proj-999",
			EnvironmentPaths: []domain.EnvironmentPath{
				{Id: "path1", Path: "/bin2"},
				{Id: "path2", Path: "/usr/local/bin2"},
			},
		}

		err := helper.repo.Update(newConfig)
		assert.NoError(t, err)

		got, err := helper.repo.GetOrCreate()
		assert.NoError(t, err)
		assert.Equal(t, &domain.Config{
			LastOpenedProjectId: "proj-999",
			EnvironmentPaths: []domain.EnvironmentPath{
				{Id: "path1", Path: "/bin2"},
				{Id: "path2", Path: "/usr/local/bin2"},
			}}, got)
	})
}

func arrange(preloadedConfig *ConfigModel, preloadedPaths []*EnvironmentPathModel) (repo *GormConfigRepository, tmpDbFilePath string) {
	ctx := context.Background()
	tmp := os.TempDir()
	id := uuid.New().String()
	tmpDbFilePath = filepath.Join(tmp, id+".db")
	gormDb, err := gorm.Open(sqlite.Open(tmpDbFilePath), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db, err := gormDb.DB()
	if err != nil {
		panic(err)
	}
	err = goose.SetDialect("sqlite3")
	if err != nil {
		panic(err)
	}
	err = goose.UpContext(ctx, db, ".")
	if err != nil {
		panic(err)
	}
	if preloadedConfig != nil {
		err = gorm.G[ConfigModel](gormDb).Create(ctx, preloadedConfig)
		if err != nil {
			panic(err)
		}
	}
	for _, m := range preloadedPaths {
		err = gorm.G[EnvironmentPathModel](gormDb).Create(ctx, m)
		if err != nil {
			panic(err)
		}
	}
	repo = &GormConfigRepository{db: gormDb, ctx: ctx}
	return
}
