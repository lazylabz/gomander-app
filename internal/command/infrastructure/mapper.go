package infrastructure

import (
	"strings"

	"gomander/internal/command/domain"
	"gomander/internal/helpers/array"
)

func ToDomainCommand(commandModel CommandModel) domain.Command {
	return domain.Command{
		Id:               commandModel.Id,
		Name:             commandModel.Name,
		Command:          commandModel.Command,
		WorkingDirectory: commandModel.WorkingDirectory,
		Position:         commandModel.Position,
		ProjectId:        commandModel.ProjectId,
		Link:             commandModel.Link,
		ErrorPatterns:    array.Filter(strings.Split(commandModel.ErrorPatterns, "\n"), func(s string) bool { return s != "" }),
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
		Link:             domainCommand.Link,
		ErrorPatterns:    strings.Join(domainCommand.ErrorPatterns, "\n"),
	}
}
