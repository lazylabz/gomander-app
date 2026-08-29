package test

import (
	"github.com/google/uuid"

	"gomander/internal/project/domain"
)

type ProjectData struct {
	Id               string
	Name             string
	WorkingDirectory string
}

type ProjectBuilder struct {
	data *ProjectData
}

func NewProjectBuilder() *ProjectBuilder {
	return &ProjectBuilder{
		data: &ProjectData{
			Id:               uuid.New().String(),
			Name:             "Test Project",
			WorkingDirectory: "/app",
		},
	}
}

func (b *ProjectBuilder) WithId(id string) *ProjectBuilder {
	b.data.Id = id
	return b
}

func (b *ProjectBuilder) WithName(name string) *ProjectBuilder {
	b.data.Name = name
	return b
}

func (b *ProjectBuilder) WithWorkingDirectory(workingDirectory string) *ProjectBuilder {
	b.data.WorkingDirectory = workingDirectory
	return b
}

func (b *ProjectBuilder) Build() domain.Project {
	return domain.Project{
		Id:               b.data.Id,
		Name:             b.data.Name,
		WorkingDirectory: b.data.WorkingDirectory,
	}
}
