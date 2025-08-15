package infrastructure

import (
	"gomander/internal/command/domain"
)

func ToDomainCommand(commandModel CommandModel) domain.Command {
	return domain.Command{
		Id:               commandModel.Id,
		Name:             commandModel.Name,
		Command:          commandModel.Command,
		WorkingDirectory: commandModel.WorkingDirectory,
		Position:         commandModel.Position,
		ProjectId:        commandModel.ProjectId,
	}
}

func ToCommandModel(domainCommand *domain.Command) CommandModel {
	return CommandModel{
		Id:               domainCommand.Id,
		Name:             domainCommand.Name,
		Command:          domainCommand.Command,
		WorkingDirectory: domainCommand.WorkingDirectory,
		Position:         domainCommand.Position,
		ProjectId:        domainCommand.ProjectId,
	}
}
