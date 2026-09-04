package infrastructure

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	facadetest "gomander/internal/facade/test"
	"gomander/internal/project/domain"
)

const exportedV2 = `{
  "version": 2,
  "name": "Gomander",
  "workingDirectory": "",
  "commands": [
    {
      "id": "cmd-1",
      "name": "Dev server",
      "command": "pnpm dev",
      "workingDirectory": "web",
      "link": "http://localhost:3000",
      "errorPatterns": [
        "ERROR",
        "FATAL"
      ]
    },
    {
      "id": "cmd-2",
      "name": "Tests",
      "command": "pnpm test",
      "workingDirectory": "",
      "link": "",
      "errorPatterns": []
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

var exportedV2Blueprint = domain.Blueprint{
	Name: "Gomander",
	Commands: []domain.BlueprintCommand{
		{
			Id:               "cmd-1",
			Name:             "Dev server",
			Command:          "pnpm dev",
			WorkingDirectory: "web",
			Link:             "http://localhost:3000",
			ErrorPatterns:    []string{"ERROR", "FATAL"},
		},
		{
			Id:               "cmd-2",
			Name:             "Tests",
			Command:          "pnpm test",
			WorkingDirectory: "",
			Link:             "",
			ErrorPatterns:    []string{},
		},
	},
	CommandGroups: []domain.BlueprintCommandGroup{
		{Id: "group-1", Name: "Everything", CommandIds: []string{"cmd-1", "cmd-2"}},
	},
}

func TestProjectFileV2_Read(t *testing.T) {
	t.Run("Should read a file that carries each command's link and error patterns", func(t *testing.T) {
		// Arrange
		fs := new(facadetest.MockFsFacade)
		fs.On("ReadFile", "/home/user/gomander.json").Return([]byte(exportedV2), nil)

		// Act
		blueprint, err := NewProjectFileV2(fs).Read("/home/user/gomander.json")

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, exportedV2Blueprint, *blueprint)
	})
}

func TestProjectFileV2_Write(t *testing.T) {
	t.Run("Should write each command's link and error patterns", func(t *testing.T) {
		// Arrange
		fs := new(facadetest.MockFsFacade)

		var written []byte
		fs.On("WriteFile", "/home/user/gomander.json", mock.Anything, os.FileMode(0644)).
			Run(func(args mock.Arguments) { written = args.Get(1).([]byte) }).
			Return(nil)

		// Act
		err := NewProjectFileV2(fs).Write("/home/user/gomander.json", exportedV2Blueprint)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, exportedV2, string(written))
	})
}
