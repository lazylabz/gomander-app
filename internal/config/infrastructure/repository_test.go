package infrastructure

import (
	"context"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"gomander/internal/config/domain"
	_ "gomander/migrations"
)

type testHelper struct {
	t    *testing.T
	repo *GormConfigRepository
}

func newTestHelper(t *testing.T,
	preloadedConfig *ConfigModel, preloadedPaths []*EnvironmentPathModel) *testHelper {
	t.Helper() // IMPORTANT: This marks the function as a helper, so error traces will point to the test instead of here

	repo := arrange(
		preloadedConfig,
		preloadedPaths,
	)

	helper := &testHelper{
		t:    t,
		repo: repo,
	}

	t.Cleanup(func() {
		assert.NoError(t, repo.db.Exec("DELETE FROM user_config").Error, "Failed to cleanup test database")
	})

	return helper
}

func TestGormConfigRepository_GetOrCreate(t *testing.T) {
	t.Run("Should create config if not exists", func(t *testing.T) {
		// Arrange
		helper := newTestHelper(t, nil, nil)

		// Act
		config, err := helper.repo.GetOrCreate()

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, config)
		assert.Equal(t, "", config.LastOpenedProjectId)
		assert.Empty(t, config.EnvironmentPaths)
		assert.Equal(t, 100, config.LogLineLimit)
	})
	t.Run("Should return existing config with environment paths", func(t *testing.T) {
		// Arrange
		preloadedConfig := &ConfigModel{Id: 1, LastOpenedProjectId: "proj-123", LogLineLimit: 200}
		preloadedPaths := []*EnvironmentPathModel{
			{Id: "path1", Path: "/usr/bin"},
			{Id: "path2", Path: "/usr/local/bin"},
		}
		helper := newTestHelper(t, preloadedConfig, preloadedPaths)

		// Act
		config, err := helper.repo.GetOrCreate()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, &domain.Config{
			LastOpenedProjectId: "proj-123",
			EnvironmentPaths: []domain.EnvironmentPath{
				{Id: "path1", Path: "/usr/bin"},
				{Id: "path2", Path: "/usr/local/bin"},
			},
			LogLineLimit: 200,
		}, config)
	})
}

func TestGormConfigRepository_Update(t *testing.T) {
	t.Run("Should save config and environment paths", func(t *testing.T) {
		// Arrange
		preloadedConfig := &ConfigModel{Id: 1, LastOpenedProjectId: "proj-123", LogLineLimit: 100}
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
			LogLineLimit: 500,
		}

		// Act
		err := helper.repo.Update(newConfig)

		// Assert
		assert.NoError(t, err)

		// Verify the config was updated correctly
		got, err := helper.repo.GetOrCreate()
		assert.NoError(t, err)
		assert.Equal(t, &domain.Config{
			LastOpenedProjectId: "proj-999",
			EnvironmentPaths: []domain.EnvironmentPath{
				{Id: "path1", Path: "/bin2"},
				{Id: "path2", Path: "/usr/local/bin2"},
			},
			LogLineLimit: 500,
		}, got)
	})
}

func arrange(preloadedConfig *ConfigModel, preloadedPaths []*EnvironmentPathModel) (repo *GormConfigRepository) {
	ctx := context.Background()

	gormDb, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
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
	repo = NewGormConfigRepository(gormDb, ctx)
	return
}
