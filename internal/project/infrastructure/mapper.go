package infrastructure

import "gomander/internal/project/domain"

func ToDomainProject(model ProjectModel) *domain.Project {
	return &domain.Project{
		Id:               model.Id,
		Name:             model.Name,
		WorkingDirectory: model.WorkingDirectory,
	}
}

func ToProjectModel(domainProject domain.Project) ProjectModel {
	return ProjectModel{
		Id:               domainProject.Id,
		Name:             domainProject.Name,
		WorkingDirectory: domainProject.WorkingDirectory,
	}
}
