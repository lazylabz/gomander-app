package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	projectdomain "gomander/internal/project/domain"
)

var blueprint = projectdomain.Blueprint{
	Name:             "Gomander",
	WorkingDirectory: "/home/user/app",
	Commands: []projectdomain.BlueprintCommand{
		{Id: "cmd-1", Name: "Dev server", Command: "pnpm dev", WorkingDirectory: "web"},
	},
	CommandGroups: []projectdomain.BlueprintCommandGroup{
		{Id: "group-1", Name: "Everything", CommandIds: []string{"cmd-1"}},
	},
}

func TestProjectBlueprint(t *testing.T) {
	t.Run("Should reach the frontend under the names it reads", func(t *testing.T) {
		// Act
		sent, err := json.Marshal(newProjectBlueprint(blueprint))

		// Assert
		assert.NoError(t, err)
		assert.JSONEq(t, `{
			"name": "Gomander",
			"workingDirectory": "/home/user/app",
			"commands": [
				{
					"id": "cmd-1",
					"name": "Dev server",
					"command": "pnpm dev",
					"workingDirectory": "web"
				}
			],
			"commandGroups": [
				{
					"id": "group-1",
					"name": "Everything",
					"commandIds": ["cmd-1"]
				}
			]
		}`, string(sent))
	})

	t.Run("Should come back from the frontend as the Blueprint it left as", func(t *testing.T) {
		// Act
		returned := newProjectBlueprint(blueprint).toDomain()

		// Assert
		assert.Equal(t, blueprint, returned)
	})
}
