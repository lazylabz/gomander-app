package commandinfrastructure

import (
	"github.com/google/uuid"
)

type CommandModelBuilder struct {
	command CommandModel
}

func NewCommandModelBuilder() *CommandModelBuilder {
	return &CommandModelBuilder{
		command: CommandModel{
			Id:               uuid.New().String(),
			ProjectId:        uuid.New().String(),
			Name:             "Default Command",
			Command:          "echo 'hello'",
			WorkingDirectory: "/app",
			Position:         0,
		},
	}
}

func (b *CommandModelBuilder) WithId(id string) *CommandModelBuilder {
	b.command.Id = id
	return b
}

func (b *CommandModelBuilder) WithProjectId(projectId string) *CommandModelBuilder {
	b.command.ProjectId = projectId
	return b
}

func (b *CommandModelBuilder) WithName(name string) *CommandModelBuilder {
	b.command.Name = name
	return b
}

func (b *CommandModelBuilder) WithCommand(command string) *CommandModelBuilder {
	b.command.Command = command
	return b
}

func (b *CommandModelBuilder) WithWorkingDirectory(workingDirectory string) *CommandModelBuilder {
	b.command.WorkingDirectory = workingDirectory
	return b
}

func (b *CommandModelBuilder) WithPosition(position int) *CommandModelBuilder {
	b.command.Position = position
	return b
}

func (b *CommandModelBuilder) Build() CommandModel {
	return b.command
}

func (b *CommandModelBuilder) BuildPtr() *CommandModel {
	return &b.command
}
