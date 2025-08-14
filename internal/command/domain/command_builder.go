package domain

import (
	"github.com/google/uuid"
)

type CommandBuilder struct {
	command Command
}

func NewCommandBuilder() *CommandBuilder {
	return &CommandBuilder{
		command: Command{
			Id:               uuid.New().String(),
			ProjectId:        uuid.New().String(),
			Name:             "Default Command",
			Command:          "echo 'hello'",
			WorkingDirectory: "/app",
			Position:         0,
		},
	}
}

func (b *CommandBuilder) WithId(id string) *CommandBuilder {
	b.command.Id = id
	return b
}

func (b *CommandBuilder) WithProjectId(projectId string) *CommandBuilder {
	b.command.ProjectId = projectId
	return b
}

func (b *CommandBuilder) WithName(name string) *CommandBuilder {
	b.command.Name = name
	return b
}

func (b *CommandBuilder) WithCommand(command string) *CommandBuilder {
	b.command.Command = command
	return b
}

func (b *CommandBuilder) WithWorkingDirectory(workingDirectory string) *CommandBuilder {
	b.command.WorkingDirectory = workingDirectory
	return b
}

func (b *CommandBuilder) WithPosition(position int) *CommandBuilder {
	b.command.Position = position
	return b
}

func (b *CommandBuilder) Build() Command {
	return b.command
}

func (b *CommandBuilder) BuildPtr() *Command {
	return &b.command
}
