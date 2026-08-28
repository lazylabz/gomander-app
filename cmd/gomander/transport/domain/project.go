package domain

import (
	projectdomain "gomander/internal/project/domain"
)

type Project struct {
	Id               string `json:"id"`
	Name             string `json:"name"`
	WorkingDirectory string `json:"workingDirectory"`
}

func FromProject(project projectdomain.Project) Project {
	return Project{
		Id:               project.Id,
		Name:             project.Name,
		WorkingDirectory: project.WorkingDirectory,
	}
}

func FromProjects(projects []projectdomain.Project) []Project {
	return mapSlice(projects, FromProject)
}

func (p Project) ToDomain() projectdomain.Project {
	return projectdomain.Project{
		Id:               p.Id,
		Name:             p.Name,
		WorkingDirectory: p.WorkingDirectory,
	}
}
