package infrastructure

import (
	"gomander/internal/command/infrastructure"
	"gomander/internal/commandgroup/domain"
	"gomander/internal/helpers/array"
)

func ToDomainCommandGroup(commandGroupModel CommandGroupModel) *domain.CommandGroup {
	return &domain.CommandGroup{
		Id:        commandGroupModel.Id,
		Name:      commandGroupModel.Name,
		ProjectId: commandGroupModel.ProjectId,
		Position:  commandGroupModel.Position,
		Commands:  array.Map(commandGroupModel.Commands, infrastructure.ToDomainCommand),
	}
}

func ToCommandGroupModel(domainCommandGroup *domain.CommandGroup) CommandGroupModel {
	return CommandGroupModel{
		Id:        domainCommandGroup.Id,
		Name:      domainCommandGroup.Name,
		ProjectId: domainCommandGroup.ProjectId,
		Position:  domainCommandGroup.Position,
	}
}
