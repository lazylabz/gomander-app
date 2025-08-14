package infrastructure

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/pressly/goose/v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"gomander/internal/config/domain"
	_ "gomander/migrations"
)

func TestGormConfigRepository_GetOrCreate(t *testing.T) {
	t.Run("Should create config if not exists", func(t *testing.T) {
		repo, tmpDbFilePath := arrange(nil, nil)
		config, err := repo.GetOrCreate()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if config == nil {
			t.Fatal("Expected config, got nil")
		}
		if config.LastOpenedProjectId != "" {
			t.Errorf("Expected empty LastOpenedProjectId, got %s", config.LastOpenedProjectId)
		}
		if len(config.EnvironmentPaths) != 0 {
			t.Errorf("Expected 0 environment paths, got %d", len(config.EnvironmentPaths))
		}
		_ = os.Remove(tmpDbFilePath)
	})
	t.Run("Should return existing config with environment paths", func(t *testing.T) {
		preloadedConfig := &ConfigModel{Id: 1, LastOpenedProjectId: "proj-123"}
		preloadedPaths := []*EnvironmentPathModel{
			{Id: "path1", Path: "/usr/bin"},
			{Id: "path2", Path: "/usr/local/bin"},
		}
		repo, tmpDbFilePath := arrange(preloadedConfig, preloadedPaths)
		config, err := repo.GetOrCreate()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if config == nil {
			t.Fatal("Expected config, got nil")
		}
		if config.LastOpenedProjectId != "proj-123" {
			t.Errorf("Expected LastOpenedProjectId 'proj-123', got %s", config.LastOpenedProjectId)
		}
		if len(config.EnvironmentPaths) != 2 {
			t.Errorf("Expected 2 environment paths, got %d", len(config.EnvironmentPaths))
		}
		_ = os.Remove(tmpDbFilePath)
	})
}

func TestGormConfigRepository_Update(t *testing.T) {
	t.Run("Should save config and environment paths", func(t *testing.T) {
		repo, tmpDbFilePath := arrange(nil, nil)
		config := &domain.Config{
			LastOpenedProjectId: "proj-999",
			EnvironmentPaths: []domain.EnvironmentPath{
				{Id: "path1", Path: "/bin"},
				{Id: "path2", Path: "/sbin"},
			},
		}
		_, err := repo.GetOrCreate()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		err = repo.Update(config)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		got, err := repo.GetOrCreate()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if got.LastOpenedProjectId != "proj-999" {
			t.Errorf("Expected LastOpenedProjectId 'proj-999', got %s", got.LastOpenedProjectId)
		}
		if len(got.EnvironmentPaths) != 2 {
			t.Errorf("Expected 2 environment paths, got %d", len(got.EnvironmentPaths))
		}
		_ = os.Remove(tmpDbFilePath)
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
