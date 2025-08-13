package infrastructure

import (
	commandDomain "gomander/internal/command/domain"
	"gomander/internal/commandgroup/domain"
)

func ToDomainCommandGroup(commandGroupModel CommandGroupModel) *domain.CommandGroup {
	return &domain.CommandGroup{
		Id:        commandGroupModel.Id,
		Name:      commandGroupModel.Name,
		ProjectId: commandGroupModel.ProjectId,
		Position:  commandGroupModel.Position,
		Commands:  make([]*commandDomain.Command, 0),
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
