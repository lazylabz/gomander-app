package test

import (
	"github.com/google/uuid"

	"gomander/internal/command/domain"
)

type CommandData struct {
	Id               string
	ProjectId        string
	Name             string
	Command          string
	WorkingDirectory string
	Position         int
	Link             string
}

type CommandBuilder struct {
	data *CommandData
}

func NewCommandBuilder() *CommandBuilder {
	return &CommandBuilder{
		data: &CommandData{
			Id:               uuid.New().String(),
			ProjectId:        uuid.New().String(),
			Name:             "Test Command",
			Command:          "echo 'hello'",
			WorkingDirectory: "/app",
			Position:         0,
			Link:             "",
		},
	}
}

func (b *CommandBuilder) WithId(id string) *CommandBuilder {
	b.data.Id = id
	return b
}

func (b *CommandBuilder) WithProjectId(projectId string) *CommandBuilder {
	b.data.ProjectId = projectId
	return b
}

func (b *CommandBuilder) WithName(name string) *CommandBuilder {
	b.data.Name = name
	return b
}

func (b *CommandBuilder) WithCommand(command string) *CommandBuilder {
	b.data.Command = command
	return b
}

func (b *CommandBuilder) WithWorkingDirectory(workingDirectory string) *CommandBuilder {
	b.data.WorkingDirectory = workingDirectory
	return b
}

func (b *CommandBuilder) WithPosition(position int) *CommandBuilder {
	b.data.Position = position
	return b
}

func (b *CommandBuilder) WithLink(link string) *CommandBuilder {
	b.data.Link = link
	return b
}

func (b *CommandBuilder) Build() domain.Command {
	return domain.Command{
		Id:               b.data.Id,
		ProjectId:        b.data.ProjectId,
		Name:             b.data.Name,
		Command:          b.data.Command,
		WorkingDirectory: b.data.WorkingDirectory,
		Position:         b.data.Position,
		Link:             b.data.Link,
	}
}
