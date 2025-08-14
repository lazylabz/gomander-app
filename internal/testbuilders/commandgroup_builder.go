package testbuilders

import (
	commanddomain "gomander/internal/command/domain"
	groupdomain "gomander/internal/commandgroup/domain"

	"github.com/google/uuid"
)

type CommandGroupBuilder struct {
	group groupdomain.CommandGroup
}

func NewCommandGroupBuilder() *CommandGroupBuilder {
	return &CommandGroupBuilder{
		group: groupdomain.CommandGroup{
			Id:        uuid.New().String(),
			ProjectId: uuid.New().String(),
			Name:      "Default Command Group",
			Commands:  []commanddomain.Command{},
			Position:  0,
		},
	}
}

func (b *CommandGroupBuilder) WithId(id string) *CommandGroupBuilder {
	b.group.Id = id
	return b
}

func (b *CommandGroupBuilder) WithProjectId(projectId string) *CommandGroupBuilder {
	b.group.ProjectId = projectId
	return b
}

func (b *CommandGroupBuilder) WithName(name string) *CommandGroupBuilder {
	b.group.Name = name
	return b
}

func (b *CommandGroupBuilder) WithCommands(commands []commanddomain.Command) *CommandGroupBuilder {
	b.group.Commands = commands
	return b
}

func (b *CommandGroupBuilder) WithPosition(position int) *CommandGroupBuilder {
	b.group.Position = position
	return b
}

func (b *CommandGroupBuilder) Build() groupdomain.CommandGroup {
	return b.group
}

func (b *CommandGroupBuilder) BuildPtr() *groupdomain.CommandGroup {
	return &b.group
}

// Convenience methods

func (b *CommandGroupBuilder) AddCommand(command commanddomain.Command) *CommandGroupBuilder {
	b.group.Commands = append(b.group.Commands, command)
	return b
}

func (b *CommandGroupBuilder) AddCommandBuilder(commandBuilder *CommandBuilder) *CommandGroupBuilder {
	// Ensure command and group have the same ProjectId
	command := commandBuilder.WithProjectId(b.group.ProjectId).Build()
	return b.AddCommand(command)
}

func (b *CommandGroupBuilder) AddCommands(commands ...commanddomain.Command) *CommandGroupBuilder {
	b.group.Commands = append(b.group.Commands, commands...)
	return b
}
