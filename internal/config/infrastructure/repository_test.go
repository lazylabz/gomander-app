package infrastructure

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"gomander/internal/config/domain"
	"gomander/internal/testdb"
)

type testHelper struct {
	t    *testing.T
	repo *GormConfigRepository
}

func newTestHelper(t *testing.T,
	preloadedConfig *ConfigModel, preloadedPaths []*EnvironmentPathModel) *testHelper {
	t.Helper() // IMPORTANT: This marks the function as a helper, so error traces will point to the test instead of here

	repo := arrange(
		t,
		preloadedConfig,
		preloadedPaths,
	)

	helper := &testHelper{
		t:    t,
		repo: repo,
	}

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
	})
	t.Run("Should return existing config with environment paths", func(t *testing.T) {
		// Arrange
		preloadedConfig := &ConfigModel{Id: 1, LastOpenedProjectId: "proj-123"}
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
		}, config)
	})
}

func TestGormConfigRepository_Update(t *testing.T) {
	t.Run("Should save config and environment paths", func(t *testing.T) {
		// Arrange
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
		}, got)
	})
}

func arrange(t *testing.T, preloadedConfig *ConfigModel, preloadedPaths []*EnvironmentPathModel) (repo *GormConfigRepository) {
	t.Helper()

	ctx := context.Background()
	gormDb := testdb.New(t)

	if preloadedConfig != nil {
		err := gorm.G[ConfigModel](gormDb).Create(ctx, preloadedConfig)
		if err != nil {
			t.Fatalf("failed to preload the config: %v", err)
		}
	}
	for _, m := range preloadedPaths {
		err := gorm.G[EnvironmentPathModel](gormDb).Create(ctx, m)
		if err != nil {
			t.Fatalf("failed to preload an environment path: %v", err)
		}
	}
	repo = NewGormConfigRepository(gormDb, ctx)
	return
}
