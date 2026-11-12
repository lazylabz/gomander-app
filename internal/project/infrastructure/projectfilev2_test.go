package infrastructure

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gomander/internal/project/domain"
)

const exportedWithLinksAndErrorPatterns = `{
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
        "panic:"
      ]
    },
    {
      "id": "cmd-2",
      "name": "Tests",
      "command": "pnpm test",
      "workingDirectory": "",
      "link": "",
      "errorPatterns": null
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

var exportedBlueprintWithLinksAndErrorPatterns = domain.Blueprint{
	Name: "Gomander",
	Commands: []domain.BlueprintCommand{
		{
			Id:               "cmd-1",
			Name:             "Dev server",
			Command:          "pnpm dev",
			WorkingDirectory: "web",
			Link:             "http://localhost:3000",
			ErrorPatterns:    []string{"ERROR", "panic:"},
		},
		{Id: "cmd-2", Name: "Tests", Command: "pnpm test"},
	},
	CommandGroups: []domain.BlueprintCommandGroup{
		{Id: "group-1", Name: "Everything", CommandIds: []string{"cmd-1", "cmd-2"}},
	},
}

func TestDecodeProjectV2(t *testing.T) {
	t.Run("Should read the link and error patterns of each command", func(t *testing.T) {
		// Act
		blueprint, err := decodeProjectV2([]byte(exportedWithLinksAndErrorPatterns))

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, exportedBlueprintWithLinksAndErrorPatterns, *blueprint)
	})
}

func TestEncodeProjectV2(t *testing.T) {
	t.Run("Should write the bytes this version of the format is pinned to", func(t *testing.T) {
		// Act
		written, err := encodeProjectV2(exportedBlueprintWithLinksAndErrorPatterns)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, exportedWithLinksAndErrorPatterns, string(written))
	})
}
