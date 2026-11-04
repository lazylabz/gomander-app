package apptest_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"gomander/internal/apptest"
	commandtest "gomander/internal/command/domain/test"
	projecttest "gomander/internal/project/domain/test"
)

func runningCommandWatchingFor(t *testing.T, h *apptest.Harness, name string, errorPatterns ...string) string {
	t.Helper()

	project := projecttest.NewProjectBuilder().Build()
	h.GivenProjects(project)
	h.GivenOpenedProject(project.Id)

	command := commandtest.NewCommandBuilder().
		WithProjectId(project.Id).
		WithName(name).
		WithErrorPatterns(errorPatterns).
		Build()
	h.GivenCommands(command)
	assert.NoError(t, h.UseCases.RunCommand.Execute(command.Id))

	return command.Id
}

func TestNotifyingAnErrorPattern(t *testing.T) {
	t.Run("Should notify the user when a command's output matches an error pattern while the window is in the background", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)
		commandId := runningCommandWatchingFor(t, h, "api", "ERROR")
		h.GivenWindowInBackground()

		// Act
		h.ProcessPrints(commandId, "ERROR: connection refused")

		// Assert
		assert.Equal(t, []apptest.SentNotification{{
			Title: "api reported an error",
			Body:  "ERROR: connection refused",
		}}, h.SentNotifications())
	})

	t.Run("Should notify once per run, however many lines match", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)
		commandId := runningCommandWatchingFor(t, h, "api", "ERROR")
		h.GivenWindowInBackground()

		// Act
		h.ProcessPrints(commandId, "ERROR: first")
		h.ProcessPrints(commandId, "ERROR: second")

		// Assert
		assert.Equal(t, []apptest.SentNotification{{
			Title: "api reported an error",
			Body:  "ERROR: first",
		}}, h.SentNotifications())
	})

	t.Run("Should notify again when the command runs again", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)
		commandId := runningCommandWatchingFor(t, h, "api", "ERROR")
		h.GivenWindowInBackground()
		h.ProcessPrints(commandId, "ERROR: first run")
		h.ProcessExits(commandId)

		// Act
		assert.NoError(t, h.UseCases.RunCommand.Execute(commandId))
		h.ProcessPrints(commandId, "ERROR: second run")

		// Assert
		assert.Len(t, h.SentNotifications(), 2)
		assert.Equal(t, "ERROR: second run", h.SentNotifications()[1].Body)
	})

	t.Run("Should word the notification in the user's language", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)
		h.GivenLocale("es")
		commandId := runningCommandWatchingFor(t, h, "api", "ERROR")
		h.GivenWindowInBackground()

		// Act
		h.ProcessPrints(commandId, "ERROR: boom")

		// Assert
		assert.Equal(t, "api ha informado de un error", h.SentNotifications()[0].Title)
	})

	t.Run("Should show the matching line without colours and cut short", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)
		commandId := runningCommandWatchingFor(t, h, "api", "ERROR")
		h.GivenWindowInBackground()

		// Act
		h.ProcessPrints(commandId, "\x1b[31mERROR\x1b[0m: "+strings.Repeat("x", 300))

		// Assert
		body := []rune(h.SentNotifications()[0].Body)
		assert.Len(t, body, 200)
		assert.Equal(t, "ERROR: xxx", string(body[:10]))
		assert.Equal(t, '…', body[199])
	})

	t.Run("Should not notify while the user is looking at the window", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)
		commandId := runningCommandWatchingFor(t, h, "api", "ERROR")

		// Act
		h.ProcessPrints(commandId, "ERROR: boom")

		// Assert
		assert.Empty(t, h.SentNotifications())
	})

	t.Run("Should not notify when the user turned notifications off", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)
		commandId := runningCommandWatchingFor(t, h, "api", "ERROR")
		h.GivenWindowInBackground()
		config := userConfigOf(t, h)
		config.NotificationsEnabled = false
		assert.NoError(t, h.UseCases.SaveUserConfig.Execute(config))

		// Act
		h.ProcessPrints(commandId, "ERROR: boom")

		// Assert
		assert.Empty(t, h.SentNotifications())
	})

	t.Run("Should not notify a command that fails without matching an error pattern", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)
		commandId := runningCommandWatchingFor(t, h, "api", "ERROR")
		h.GivenWindowInBackground()

		// Act
		h.ProcessPrints(commandId, "exit status 1")
		h.ProcessExits(commandId)

		// Assert
		assert.Empty(t, h.SentNotifications())
	})

	t.Run("Should keep the app going when the operating system refuses the notification", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)
		commandId := runningCommandWatchingFor(t, h, "api", "ERROR")
		h.GivenWindowInBackground()
		h.GivenNotificationsThatFail(errors.New("no notification daemon"))

		// Act
		h.ProcessPrints(commandId, "ERROR: boom")

		// Assert
		assert.Empty(t, h.SentNotifications())
		assert.Equal(t, []string{commandId}, h.UseCases.GetRunningCommandIds.Execute())
		assert.NoError(t, h.UseCases.StopCommand.Execute(commandId))
	})
}
