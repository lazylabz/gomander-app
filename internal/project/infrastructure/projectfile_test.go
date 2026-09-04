package infrastructure

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	facadetest "gomander/internal/facade/test"
	"gomander/internal/project/domain"
)

func TestProjectFile_Read(t *testing.T) {
	t.Run("Should read a version 1 file as having no link and no error patterns", func(t *testing.T) {
		// Arrange
		fs := new(facadetest.MockFsFacade)
		fs.On("ReadFile", "/home/user/gomander.json").Return([]byte(exportedByAnEarlierGomander), nil)

		// Act
		blueprint, err := NewProjectFile(fs).Read("/home/user/gomander.json")

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, exportedBlueprint, *blueprint)
		for _, command := range blueprint.Commands {
			assert.Empty(t, command.Link)
			assert.Empty(t, command.ErrorPatterns)
		}
	})

	t.Run("Should read a version 2 file with each command's link and error patterns", func(t *testing.T) {
		// Arrange
		fs := new(facadetest.MockFsFacade)
		fs.On("ReadFile", "/home/user/gomander.json").Return([]byte(exportedV2), nil)

		// Act
		blueprint, err := NewProjectFile(fs).Read("/home/user/gomander.json")

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, exportedV2Blueprint, *blueprint)
	})

	t.Run("Should refuse a version the app does not know", func(t *testing.T) {
		// Arrange
		fs := new(facadetest.MockFsFacade)
		fs.On("ReadFile", "/home/user/gomander.json").Return([]byte(`{"version": 3, "name": "Gomander"}`), nil)

		// Act
		blueprint, err := NewProjectFile(fs).Read("/home/user/gomander.json")

		// Assert
		assert.EqualError(t, err, "unsupported project file version 3")
		assert.Nil(t, blueprint)
	})
}

func TestProjectFile_Write(t *testing.T) {
	t.Run("Should write version 2 so a later import restores link and error patterns", func(t *testing.T) {
		// Arrange
		fs := new(facadetest.MockFsFacade)

		var written []byte
		fs.On("WriteFile", "/home/user/gomander.json", mock.Anything, os.FileMode(0644)).
			Run(func(args mock.Arguments) { written = args.Get(1).([]byte) }).
			Return(nil)

		// Act
		err := NewProjectFile(fs).Write("/home/user/gomander.json", exportedV2Blueprint)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, exportedV2, string(written))

		roundTripFs := new(facadetest.MockFsFacade)
		roundTripFs.On("ReadFile", "/home/user/gomander.json").Return(written, nil)

		got, err := NewProjectFile(roundTripFs).Read("/home/user/gomander.json")
		require.NoError(t, err)
		assert.Equal(t, exportedV2Blueprint, *got)
	})
}

func TestProjectFile_WriteThenRead(t *testing.T) {
	t.Run("Should restore a command that had no link and no error patterns", func(t *testing.T) {
		// Arrange
		blueprint := domain.Blueprint{
			Name: "Gomander",
			Commands: []domain.BlueprintCommand{
				{Id: "cmd-1", Name: "Dev server", Command: "pnpm dev", WorkingDirectory: "web"},
			},
		}

		fs := new(facadetest.MockFsFacade)
		var written []byte
		fs.On("WriteFile", "/tmp/out.json", mock.Anything, os.FileMode(0644)).
			Run(func(args mock.Arguments) { written = args.Get(1).([]byte) }).
			Return(nil)

		require.NoError(t, NewProjectFile(fs).Write("/tmp/out.json", blueprint))

		readFs := new(facadetest.MockFsFacade)
		readFs.On("ReadFile", "/tmp/out.json").Return(written, nil)

		// Act
		got, err := NewProjectFile(readFs).Read("/tmp/out.json")

		// Assert
		require.NoError(t, err)
		assert.Equal(t, "", got.Commands[0].Link)
		assert.Equal(t, []string{}, got.Commands[0].ErrorPatterns)
	})
}
