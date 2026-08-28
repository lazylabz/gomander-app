package infrastructure

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	facadetest "gomander/internal/facade/test"
	"gomander/internal/project/domain"
)

func TestPackageJSONFile_Read(t *testing.T) {
	t.Run("Should read the scripts as the commands of a project", func(t *testing.T) {
		// Arrange
		fs := new(facadetest.MockFsFacade)
		fs.On("ReadFile", "/home/user/app/package.json").Return([]byte(`{
			"name": "My NPM Project",
			"scripts": {"start": "node index.js"}
		}`), nil)

		// Act
		blueprint, err := NewPackageJSONFile(fs).Read("/home/user/app/package.json")

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "My NPM Project", blueprint.Name)
		assert.Equal(t, []domain.BlueprintCommand{
			{Id: "start", Name: "start", Command: "node index.js"},
		}, blueprint.Commands)
		assert.Empty(t, blueprint.CommandGroups)
	})

	// The folder used to be found with path.Dir, which only knows slashes: on
	// Windows it found no separator, answered ".", and every imported Command
	// ran wherever the app's process happened to be.
	t.Run("Should find the folder in a path spelled the way the operating system spells it", func(t *testing.T) {
		// Arrange
		filePath := filepath.FromSlash("/home/user/app/package.json")

		fs := new(facadetest.MockFsFacade)
		fs.On("ReadFile", filePath).Return([]byte(`{"name": "My NPM Project"}`), nil)

		// Act
		blueprint, err := NewPackageJSONFile(fs).Read(filePath)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, filepath.FromSlash("/home/user/app"), blueprint.WorkingDirectory)
	})

	t.Run("Should report a file that cannot be read", func(t *testing.T) {
		// Arrange
		fs := new(facadetest.MockFsFacade)
		fs.On("ReadFile", "/home/user/app/package.json").Return([]byte(nil), os.ErrNotExist)

		// Act
		blueprint, err := NewPackageJSONFile(fs).Read("/home/user/app/package.json")

		// Assert
		assert.ErrorIs(t, err, os.ErrNotExist)
		assert.Nil(t, blueprint)
	})

	t.Run("Should report a manifest that is not JSON", func(t *testing.T) {
		// Arrange
		fs := new(facadetest.MockFsFacade)
		fs.On("ReadFile", "/home/user/app/package.json").Return([]byte(`{"scripts": {`), nil)

		// Act
		blueprint, err := NewPackageJSONFile(fs).Read("/home/user/app/package.json")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, blueprint)
	})
}
