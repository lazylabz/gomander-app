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

	"gomander/internal/project/domain"
	_ "gomander/migrations"
)

func TestGormProjectRepository_GetAllProjects(t *testing.T) {
	t.Run("Should return all projects", func(t *testing.T) {
		preloadedProjects := []*ProjectModel{
			{Id: "p1", Name: "Project 1", WorkingDirectory: "/tmp/1"},
			{Id: "p2", Name: "Project 2", WorkingDirectory: "/tmp/2"},
		}
		expectedProjects := []domain.Project{
			{Id: "p1", Name: "Project 1", WorkingDirectory: "/tmp/1"},
			{Id: "p2", Name: "Project 2", WorkingDirectory: "/tmp/2"},
		}
		repo, tmpDbFilePath := arrange(preloadedProjects)
		projects, err := repo.GetAllProjects()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(projects) != 2 {
			t.Errorf("Expected 2 projects, got %d", len(projects))
		}
		for i, project := range projects {
			if !project.Equals(&expectedProjects[i]) {
				t.Errorf("Expected project %v, got %v", expectedProjects[i], project)
			}
		}
		_ = os.Remove(tmpDbFilePath)
	})
}

func TestGormProjectRepository_GetProjectById(t *testing.T) {
	t.Run("Should return project when it exists", func(t *testing.T) {
		preloadedProjects := []*ProjectModel{
			{Id: "p1", Name: "Project 1", WorkingDirectory: "/tmp/1"},
		}
		expectedProject := domain.Project{Id: "p1", Name: "Project 1", WorkingDirectory: "/tmp/1"}
		repo, tmpDbFilePath := arrange(preloadedProjects)
		project, err := repo.GetProjectById("p1")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if project == nil || project.Id != "p1" {
			t.Errorf("Expected project with id p1, got %v", project)
		}
		if !project.Equals(&expectedProject) {
			t.Errorf("Expected project %v, got %v", expectedProject, project)
		}
		_ = os.Remove(tmpDbFilePath)
	})
	t.Run("Should return nil when project does not exist", func(t *testing.T) {
		repo, tmpDbFilePath := arrange(nil)
		project, err := repo.GetProjectById("nonexistent")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if project != nil {
			t.Errorf("Expected nil, got %v", project)
		}
		_ = os.Remove(tmpDbFilePath)
	})
}

func TestGormProjectRepository_CreateProject(t *testing.T) {
	t.Run("Should create a new project", func(t *testing.T) {
		repo, tmpDbFilePath := arrange(nil)
		newProject := domain.Project{Id: "p3", Name: "Project 3", WorkingDirectory: "/tmp/3"}
		err := repo.CreateProject(newProject)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		project, err := repo.GetProjectById("p3")
		if err != nil || project == nil || project.Id != "p3" {
			t.Errorf("Expected project with id p3, got %v (err: %v)", project, err)
		}
		_ = os.Remove(tmpDbFilePath)
	})
}

func TestGormProjectRepository_UpdateProject(t *testing.T) {
	t.Run("Should update an existing project", func(t *testing.T) {
		preloadedProjects := []*ProjectModel{
			{Id: "p1", Name: "Old Name", WorkingDirectory: "/tmp/old"},
		}
		repo, tmpDbFilePath := arrange(preloadedProjects)
		updated := domain.Project{Id: "p1", Name: "New Name", WorkingDirectory: "/tmp/new"}
		err := repo.UpdateProject(updated)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		project, _ := repo.GetProjectById("p1")
		if project == nil || project.Name != "New Name" || project.WorkingDirectory != "/tmp/new" {
			t.Errorf("Expected updated project, got %v", project)
		}
		_ = os.Remove(tmpDbFilePath)
	})
}

func TestGormProjectRepository_DeleteProject(t *testing.T) {
	t.Run("Should delete an existing project", func(t *testing.T) {
		preloadedProjects := []*ProjectModel{
			{Id: "p1", Name: "To Delete", WorkingDirectory: "/tmp/del"},
		}
		repo, tmpDbFilePath := arrange(preloadedProjects)
		err := repo.DeleteProject("p1")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		project, _ := repo.GetProjectById("p1")
		if project != nil {
			t.Errorf("Expected nil after delete, got %v", project)
		}
		_ = os.Remove(tmpDbFilePath)
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
