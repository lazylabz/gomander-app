package apptest_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gomander/internal/apptest"
	commanddomain "gomander/internal/command/domain"
	commandgroupdomain "gomander/internal/commandgroup/domain"
	"gomander/internal/helpers/array"
	projectdomain "gomander/internal/project/domain"
	projecttest "gomander/internal/project/domain/test"
)

func givenAnOpenedProject(h *apptest.Harness) projectdomain.Project {
	project := projecttest.NewProjectBuilder().Build()
	h.GivenProjects(project)
	h.GivenOpenedProject(project.Id)

	return project
}

func commandsOf(t *testing.T, h *apptest.Harness) []commanddomain.Command {
	t.Helper()

	commands, err := h.UseCases.GetCommands.Execute()
	assert.NoError(t, err)

	return commands
}

func commandGroupsOf(t *testing.T, h *apptest.Harness) []commandgroupdomain.CommandGroup {
	t.Helper()

	commandGroups, err := h.UseCases.GetCommandGroups.Execute()
	assert.NoError(t, err)

	return commandGroups
}

func commandId(command commanddomain.Command) string { return command.Id }

func commandName(command commanddomain.Command) string { return command.Name }

func commandPosition(command commanddomain.Command) int { return command.Position }

// commandNamesOf resolves the Commands a Command Group names against the
// Project's own, so an assertion can read in names rather than in ids.
func commandNamesOf(commands []commanddomain.Command, commandIds []string) []string {
	namesById := make(map[string]string, len(commands))
	for _, command := range commands {
		namesById[command.Id] = command.Name
	}

	return array.Map(commandIds, func(commandId string) string { return namesById[commandId] })
}

func commandGroupId(commandGroup commandgroupdomain.CommandGroup) string { return commandGroup.Id }

func commandGroupPosition(commandGroup commandgroupdomain.CommandGroup) int {
	return commandGroup.Position
}
