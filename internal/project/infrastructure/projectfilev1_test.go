package infrastructure

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gomander/internal/project/domain"
)

// exportedByAnEarlierGomander is a file the app wrote before commands carried a
// link and error patterns, kept verbatim: a file a user already has must still
// import.
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

func TestDecodeProjectV1(t *testing.T) {
	t.Run("Should read a file exported by an earlier Gomander", func(t *testing.T) {
		// Act
		blueprint, err := decodeProjectV1([]byte(exportedByAnEarlierGomander))

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, exportedBlueprint, *blueprint)
	})
}
