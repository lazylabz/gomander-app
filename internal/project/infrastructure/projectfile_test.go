package infrastructure

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	facadetest "gomander/internal/facade/test"
)

func TestProjectFile_Read(t *testing.T) {
	t.Run("Should read a version 1 file as the commands it carried, without link or error patterns", func(t *testing.T) {
		// Arrange
		fs := new(facadetest.MockFsFacade)
		fs.On("ReadFile", "/home/user/gomander.json").Return([]byte(exportedByAnEarlierGomander), nil)

		// Act
		blueprint, err := NewProjectFile(fs).Read("/home/user/gomander.json")

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, exportedBlueprint, *blueprint)
	})

	t.Run("Should read a version 2 file with the link and error patterns it carries", func(t *testing.T) {
		// Arrange
		fs := new(facadetest.MockFsFacade)
		fs.On("ReadFile", "/home/user/gomander.json").Return([]byte(exportedWithLinksAndErrorPatterns), nil)

		// Act
		blueprint, err := NewProjectFile(fs).Read("/home/user/gomander.json")

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, exportedBlueprintWithLinksAndErrorPatterns, *blueprint)
	})

	t.Run("Should refuse a version it has no codec for", func(t *testing.T) {
		for name, contents := range map[string]string{
			"newer than any it knows": `{"version": 3, "name": "Gomander"}`,
			"missing":                 `{"name": "Gomander"}`,
			"zero":                    `{"version": 0, "name": "Gomander"}`,
			"negative":                `{"version": -1, "name": "Gomander"}`,
		} {
			t.Run(name, func(t *testing.T) {
				// Arrange
				fs := new(facadetest.MockFsFacade)
				fs.On("ReadFile", "/home/user/gomander.json").Return([]byte(contents), nil)

				// Act
				blueprint, err := NewProjectFile(fs).Read("/home/user/gomander.json")

				// Assert
				assert.ErrorContains(t, err, "unsupported project file version")
				assert.Nil(t, blueprint)
			})
		}
	})

	t.Run("Should report a file that is not JSON", func(t *testing.T) {
		// Arrange
		fs := new(facadetest.MockFsFacade)
		fs.On("ReadFile", "/home/user/gomander.json").Return([]byte("not json"), nil)

		// Act
		blueprint, err := NewProjectFile(fs).Read("/home/user/gomander.json")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, blueprint)
	})

	t.Run("Should report a file it cannot read", func(t *testing.T) {
		// Arrange
		fs := new(facadetest.MockFsFacade)
		fs.On("ReadFile", "/home/user/gomander.json").Return([]byte(nil), assert.AnError)

		// Act
		blueprint, err := NewProjectFile(fs).Read("/home/user/gomander.json")

		// Assert
		assert.ErrorIs(t, err, assert.AnError)
		assert.Nil(t, blueprint)
	})
}

func TestProjectFile_Write(t *testing.T) {
	t.Run("Should write the newest version of the format", func(t *testing.T) {
		// Arrange
		fs := new(facadetest.MockFsFacade)

		var written []byte
		fs.On("WriteFile", "/home/user/gomander.json", mock.Anything, os.FileMode(0644)).
			Run(func(args mock.Arguments) { written = args.Get(1).([]byte) }).
			Return(nil)

		// Act
		err := NewProjectFile(fs).Write("/home/user/gomander.json", exportedBlueprintWithLinksAndErrorPatterns)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, exportedWithLinksAndErrorPatterns, string(written))
	})
}
