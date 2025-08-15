package infrastructure

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"gomander/internal/project/domain"
	_ "gomander/migrations"
)

type testHelper struct {
	t      *testing.T
	repo   *GormProjectRepository
	dbPath string
}

func newTestHelper(t *testing.T, preloadedProjects []*ProjectModel) *testHelper {
	t.Helper()

	repo, dbPath := arrange(preloadedProjects)

	helper := &testHelper{
		t:      t,
		repo:   repo,
		dbPath: dbPath,
	}

	t.Cleanup(func() {
		assert.NoError(t, os.Remove(helper.dbPath), "Failed to cleanup test database")
	})

	return helper
}

func TestGormProjectRepository_GetAll(t *testing.T) {
	t.Run("Should return all projects", func(t *testing.T) {
		preloadedProjects := []*ProjectModel{
			{Id: "p1", Name: "Project 1", WorkingDirectory: "/tmp/1"},
			{Id: "p2", Name: "Project 2", WorkingDirectory: "/tmp/2"},
		}
		expectedProjects := []domain.Project{
			{Id: "p1", Name: "Project 1", WorkingDirectory: "/tmp/1"},
			{Id: "p2", Name: "Project 2", WorkingDirectory: "/tmp/2"},
		}
		h := newTestHelper(t, preloadedProjects)
		projects, err := h.repo.GetAll()
		assert.NoError(t, err)
		assert.Len(t, projects, 2)
		for i, project := range projects {
			assert.True(t, project.Equals(&expectedProjects[i]), "Expected project %v, got %v", expectedProjects[i], project)
		}
	})
}

func TestGormProjectRepository_Get(t *testing.T) {
	t.Run("Should return project when it exists", func(t *testing.T) {
		preloadedProjects := []*ProjectModel{
			{Id: "p1", Name: "Project 1", WorkingDirectory: "/tmp/1"},
		}
		expectedProject := domain.Project{Id: "p1", Name: "Project 1", WorkingDirectory: "/tmp/1"}
		h := newTestHelper(t, preloadedProjects)
		project, err := h.repo.Get("p1")
		assert.NoError(t, err)
		assert.NotNil(t, project)
		assert.Equal(t, "p1", project.Id)
		assert.True(t, project.Equals(&expectedProject))
	})
	t.Run("Should return nil when project does not exist", func(t *testing.T) {
		h := newTestHelper(t, nil)
		project, err := h.repo.Get("nonexistent")
		assert.NoError(t, err)
		assert.Nil(t, project)
	})
}

func TestGormProjectRepository_Create(t *testing.T) {
	t.Run("Should create a new project", func(t *testing.T) {
		h := newTestHelper(t, nil)
		newProject := domain.Project{Id: "p3", Name: "Project 3", WorkingDirectory: "/tmp/3"}
		err := h.repo.Create(newProject)
		assert.NoError(t, err)
		project, err := h.repo.Get("p3")
		assert.NoError(t, err)
		assert.NotNil(t, project)
		assert.Equal(t, "p3", project.Id)
	})
}

func TestGormProjectRepository_Update(t *testing.T) {
	t.Run("Should update an existing project", func(t *testing.T) {
		preloadedProjects := []*ProjectModel{
			{Id: "p1", Name: "Old Name", WorkingDirectory: "/tmp/old"},
		}
		h := newTestHelper(t, preloadedProjects)
		updated := domain.Project{Id: "p1", Name: "New Name", WorkingDirectory: "/tmp/new"}
		err := h.repo.Update(updated)
		assert.NoError(t, err)
		project, err := h.repo.Get("p1")
		assert.NoError(t, err)
		assert.NotNil(t, project)
		assert.Equal(t, "New Name", project.Name)
		assert.Equal(t, "/tmp/new", project.WorkingDirectory)
	})
}

func TestGormProjectRepository_Delete(t *testing.T) {
	t.Run("Should delete an existing project", func(t *testing.T) {
		preloadedProjects := []*ProjectModel{
			{Id: "p1", Name: "To Delete", WorkingDirectory: "/tmp/del"},
		}
		h := newTestHelper(t, preloadedProjects)
		err := h.repo.Delete("p1")
		assert.NoError(t, err)
		project, err := h.repo.Get("p1")
		assert.NoError(t, err)
		assert.Nil(t, project)
	})
}

func arrange(preloadedProjects []*ProjectModel) (repo *GormProjectRepository, tmpDbFilePath string) {
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
	for _, m := range preloadedProjects {
		err = gorm.G[ProjectModel](gormDb).Create(ctx, m)
		if err != nil {
			panic(err)
		}
	}
	repo = &GormProjectRepository{db: gormDb, ctx: ctx}
	return
}
