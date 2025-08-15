package testutils

import (
	"github.com/google/uuid"
)

type CommandGroupData struct {
	Id        string
	ProjectId string
	Name      string
	Position  int
	Commands  []CommandData
}

type CommandGroupBuilder struct {
	data *CommandGroupData
}

func NewCommandGroup() *CommandGroupBuilder {
	return &CommandGroupBuilder{
		data: &CommandGroupData{
			Id:        uuid.New().String(),
			ProjectId: uuid.New().String(),
			Name:      "Test Command Group",
			Position:  0,
			Commands:  make([]CommandData, 0),
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

func (b *CommandGroupBuilder) WithCommands(commands ...CommandData) *CommandGroupBuilder {
	b.data.Commands = commands
	return b
}

func (b *CommandGroupBuilder) Data() CommandGroupData {
	return *b.data
}

type CommandGroupCommandRelationData struct {
	CommandId      string
	CommandGroupId string
	Position       int
}

type CommandGroupCommandRelationBuilder struct {
	data *CommandGroupCommandRelationData
}
