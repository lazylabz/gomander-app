package infrastructure

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"gomander/internal/domainerrors"
	"gomander/internal/project/domain"
	"gomander/internal/testdb"
)

type testHelper struct {
	t    *testing.T
	repo *GormProjectRepository
}

func newTestHelper(t *testing.T, preloadedProjects []*ProjectModel) *testHelper {
	t.Helper()

	repo := arrange(t, preloadedProjects)

	helper := &testHelper{
		t:    t,
		repo: repo,
	}

	return helper
}

func TestGormProjectRepository_GetAll(t *testing.T) {
	t.Run("Should return all projects", func(t *testing.T) {
		// Arrange
		preloadedProjects := []*ProjectModel{
			{Id: "p1", Name: "Project 1", WorkingDirectory: "/tmp/1"},
			{Id: "p2", Name: "Project 2", WorkingDirectory: "/tmp/2"},
		}
		expectedProjects := []domain.Project{
			{Id: "p1", Name: "Project 1", WorkingDirectory: "/tmp/1"},
			{Id: "p2", Name: "Project 2", WorkingDirectory: "/tmp/2"},
		}
		h := newTestHelper(t, preloadedProjects)

		// Act
		projects, err := h.repo.GetAll()

		// Assert
		assert.NoError(t, err)
		assert.Len(t, projects, 2)
		assert.Equal(t, expectedProjects, projects)
	})
}

func TestGormProjectRepository_Get(t *testing.T) {
	t.Run("Should return project when it exists", func(t *testing.T) {
		// Arrange
		preloadedProjects := []*ProjectModel{
			{Id: "p1", Name: "Project 1", WorkingDirectory: "/tmp/1"},
		}
		expectedProject := domain.Project{Id: "p1", Name: "Project 1", WorkingDirectory: "/tmp/1"}
		h := newTestHelper(t, preloadedProjects)

		// Act
		project, err := h.repo.Get("p1")

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedProject, project)
	})
	t.Run("Should report a project that does not exist as not found", func(t *testing.T) {
		// Arrange
		h := newTestHelper(t, nil)

		// Act
		_, err := h.repo.Get("nonexistent")

		// Assert
		assert.ErrorIs(t, err, domainerrors.ErrNotFound)
	})
}

func TestGormProjectRepository_Find(t *testing.T) {
	t.Run("Should return the project when it exists", func(t *testing.T) {
		// Arrange
		preloadedProjects := []*ProjectModel{
			{Id: "p1", Name: "Project 1", WorkingDirectory: "/tmp/1"},
		}
		expectedProject := domain.Project{Id: "p1", Name: "Project 1", WorkingDirectory: "/tmp/1"}
		h := newTestHelper(t, preloadedProjects)

		// Act
		project, found, err := h.repo.Find("p1")

		// Assert
		assert.NoError(t, err)
		assert.True(t, found)
		assert.Equal(t, expectedProject, project)
	})
	t.Run("Should report absence without an error when the project does not exist", func(t *testing.T) {
		// Arrange
		h := newTestHelper(t, nil)

		// Act
		_, found, err := h.repo.Find("nonexistent")

		// Assert
		assert.NoError(t, err)
		assert.False(t, found)
	})
}

func TestGormProjectRepository_Create(t *testing.T) {
	t.Run("Should create a new project", func(t *testing.T) {
		// Arrange
		h := newTestHelper(t, nil)
		newProject := domain.Project{Id: "p3", Name: "Project 3", WorkingDirectory: "/tmp/3"}

		// Act
		err := h.repo.Create(newProject)

		// Assert
		assert.NoError(t, err)

		// Verify the project was created
		project, err := h.repo.Get("p3")
		assert.NoError(t, err)
		assert.Equal(t, "p3", project.Id)
	})
}

func TestGormProjectRepository_Update(t *testing.T) {
	t.Run("Should update an existing project", func(t *testing.T) {
		// Arrange
		preloadedProjects := []*ProjectModel{
			{Id: "p1", Name: "Old Name", WorkingDirectory: "/tmp/old"},
		}
		h := newTestHelper(t, preloadedProjects)
		updated := domain.Project{Id: "p1", Name: "New Name", WorkingDirectory: "/tmp/new"}

		// Act
		err := h.repo.Update(updated)

		// Assert
		assert.NoError(t, err)

		// Verify the project was updated
		project, err := h.repo.Get("p1")
		assert.NoError(t, err)
		assert.Equal(t, "New Name", project.Name)
		assert.Equal(t, "/tmp/new", project.WorkingDirectory)
	})
}

func TestGormProjectRepository_Delete(t *testing.T) {
	t.Run("Should delete an existing project", func(t *testing.T) {
		// Arrange
		preloadedProjects := []*ProjectModel{
			{Id: "p1", Name: "To Delete", WorkingDirectory: "/tmp/del"},
		}
		h := newTestHelper(t, preloadedProjects)

		// Act
		err := h.repo.Delete("p1")

		// Assert
		assert.NoError(t, err)

		// Verify the project was deleted
		_, found, err := h.repo.Find("p1")
		assert.NoError(t, err)
		assert.False(t, found)
	})
}

func arrange(t *testing.T, preloadedProjects []*ProjectModel) (repo *GormProjectRepository) {
	t.Helper()

	ctx := context.Background()
	gormDb := testdb.New(t)

	for _, m := range preloadedProjects {
		err := gorm.G[ProjectModel](gormDb).Create(ctx, m)
		if err != nil {
			t.Fatalf("failed to preload project: %v", err)
		}
	}
	repo = NewGormProjectRepository(gormDb, ctx)
	return
}
