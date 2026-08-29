package apptest_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"gomander/internal/apptest"
	commandtest "gomander/internal/command/domain/test"
)

func TestClosingTheApp(t *testing.T) {
	t.Run("Should stop every running command and let the app close", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)
		project := givenAnOpenedProject(h)

		first := commandtest.NewCommandBuilder().WithProjectId(project.Id).WithPosition(0).Build()
		second := commandtest.NewCommandBuilder().WithProjectId(project.Id).WithPosition(1).Build()
		h.GivenCommands(first, second)

		assert.NoError(t, h.UseCases.RunCommand.Execute(first.Id))
		assert.NoError(t, h.UseCases.RunCommand.Execute(second.Id))

		// Act
		prevent := h.ClosingTheApp()

		// Assert
		assert.False(t, prevent)
		assert.ElementsMatch(t, []string{first.Id, second.Id}, h.StoppedProcessIds())
		assert.Empty(t, h.UseCases.GetRunningCommandIds.Execute())
	})
}

// The harness mirrors buildDeps in cmd/gomander/main.go by hand, so this is the
// drift the compiler cannot see: an operation added to the registry and left
// unwired here would make every test of it pass against a nil use case.
func TestTheHarnessWiresEveryOperationTheAppExposes(t *testing.T) {
	// Arrange
	h := apptest.New(t)
	registry := reflect.ValueOf(h.UseCases)

	// Act
	unwired := make([]string, 0)
	for i := range registry.NumField() {
		name := registry.Type().Field(i).Name
		if name == "GetTranslation" || name == "GetSupportedLanguages" {
			// Localization reads the locale files embedded in the desktop binary
			continue
		}
		if registry.Field(i).IsNil() {
			unwired = append(unwired, name)
		}
	}

	// Assert
	assert.Empty(t, unwired)
}
