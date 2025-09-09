package test

import (
	"github.com/google/uuid"

	"gomander/internal/command/domain"
	commandgroupdomain "gomander/internal/commandgroup/domain"
)

type CommandGroupData struct {
	Id        string
	ProjectId string
	Name      string
	Position  int
	Commands  []domain.Command
}

type CommandGroupBuilder struct {
	data *CommandGroupData
}

func NewCommandGroupBuilder() *CommandGroupBuilder {
	return &CommandGroupBuilder{
		data: &CommandGroupData{
			Id:        uuid.New().String(),
			ProjectId: uuid.New().String(),
			Name:      "Test Command Group",
			Position:  0,
			Commands:  make([]domain.Command, 0),
		},
	}
}

func (b *CommandGroupBuilder) WithId(id string) *CommandGroupBuilder {
	b.data.Id = id
	return b
}

func (b *CommandGroupBuilder) WithProjectId(projectId string) *CommandGroupBuilder {
	b.data.ProjectId = projectId
	return b
}

func (b *CommandGroupBuilder) WithName(name string) *CommandGroupBuilder {
	b.data.Name = name
	return b
}

func (b *CommandGroupBuilder) WithPosition(position int) *CommandGroupBuilder {
	b.data.Position = position
	return b
}

func (b *CommandGroupBuilder) WithCommands(commands ...domain.Command) *CommandGroupBuilder {
	b.data.Commands = commands
	return b
}

func (b *CommandGroupBuilder) Build() commandgroupdomain.CommandGroup {
	return commandgroupdomain.CommandGroup{
		Id:        b.data.Id,
		ProjectId: b.data.ProjectId,
		Name:      b.data.Name,
		Commands:  b.data.Commands,
		Position:  b.data.Position,
	}
}

type CommandGroupCommandRelationData struct {
	CommandId      string
	CommandGroupId string
	Position       int
}

type CommandGroupCommandRelationBuilder struct {
	data *CommandGroupCommandRelationData
}
