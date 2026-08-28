package infrastructure

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	facadetest "gomander/internal/facade/test"
	"gomander/internal/project/domain"
)

// exportedByAnEarlierGomander is a file the app wrote before the format moved
// out of the domain, kept verbatim. Both directions are pinned against it: a
// file a user already has must still import, and a file the app writes today
// must still open in a Gomander that has not been updated.
const exportedByAnEarlierGomander = `{
  "version": 1,
  "name": "Gomander",
  "workingDirectory": "",
  "commands": [
    {
      "id": "cmd-1",
      "name": "Dev server",
      "command": "pnpm dev",
      "workingDirectory": "web"
    },
    {
      "id": "cmd-2",
      "name": "Tests",
      "command": "pnpm test",
      "workingDirectory": ""
    }
  ],
  "commandGroups": [
    {
      "id": "group-1",
      "name": "Everything",
      "commandIds": [
        "cmd-1",
        "cmd-2"
      ]
    }
  ]
}`

var exportedBlueprint = domain.Blueprint{
	Name: "Gomander",
	Commands: []domain.BlueprintCommand{
		{Id: "cmd-1", Name: "Dev server", Command: "pnpm dev", WorkingDirectory: "web"},
		{Id: "cmd-2", Name: "Tests", Command: "pnpm test"},
	},
	CommandGroups: []domain.BlueprintCommandGroup{
		{Id: "group-1", Name: "Everything", CommandIds: []string{"cmd-1", "cmd-2"}},
	},
}

func TestProjectFileV1_Read(t *testing.T) {
	t.Run("Should read a file exported by an earlier Gomander", func(t *testing.T) {
		// Arrange
		fs := new(facadetest.MockFsFacade)
		fs.On("ReadFile", "/home/user/gomander.json").Return([]byte(exportedByAnEarlierGomander), nil)

		// Act
		blueprint, err := NewProjectFileV1(fs).Read("/home/user/gomander.json")

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, exportedBlueprint, *blueprint)
	})

	t.Run("Should report a file that cannot be read", func(t *testing.T) {
		// Arrange
		fs := new(facadetest.MockFsFacade)
		fs.On("ReadFile", "/home/user/gone.json").Return([]byte(nil), os.ErrNotExist)

		// Act
		blueprint, err := NewProjectFileV1(fs).Read("/home/user/gone.json")

		// Assert
		assert.ErrorIs(t, err, os.ErrNotExist)
		assert.Nil(t, blueprint)
	})

	t.Run("Should report a file that is not the format", func(t *testing.T) {
		// Arrange
		fs := new(facadetest.MockFsFacade)
		fs.On("ReadFile", "/home/user/gomander.json").Return([]byte(`{"version": 1, "commands": [`), nil)

		// Act
		blueprint, err := NewProjectFileV1(fs).Read("/home/user/gomander.json")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, blueprint)
	})
}

func TestProjectFileV1_Write(t *testing.T) {
	t.Run("Should write the same bytes an earlier Gomander wrote", func(t *testing.T) {
		// Arrange
		fs := new(facadetest.MockFsFacade)

		var written []byte
		fs.On("WriteFile", "/home/user/gomander.json", mock.Anything, os.FileMode(0644)).
			Run(func(args mock.Arguments) { written = args.Get(1).([]byte) }).
			Return(nil)

		// Act
		err := NewProjectFileV1(fs).Write("/home/user/gomander.json", exportedBlueprint)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, exportedByAnEarlierGomander, string(written))
	})

	t.Run("Should report a destination that cannot be written to", func(t *testing.T) {
		// Arrange
		fs := new(facadetest.MockFsFacade)
		fs.On("WriteFile", mock.Anything, mock.Anything, mock.Anything).Return(os.ErrPermission)

		// Act
		err := NewProjectFileV1(fs).Write("/read/only/gomander.json", exportedBlueprint)

		// Assert
		assert.ErrorIs(t, err, os.ErrPermission)
	})
}
