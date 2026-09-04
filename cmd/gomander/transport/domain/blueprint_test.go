package domain_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	transport "gomander/cmd/gomander/transport/domain"
	projectdomain "gomander/internal/project/domain"
)

var blueprint = projectdomain.Blueprint{
	Name:             "Gomander",
	WorkingDirectory: "/home/user/app",
	Commands: []projectdomain.BlueprintCommand{
		{
			Id:               "cmd-1",
			Name:             "Dev server",
			Command:          "pnpm dev",
			WorkingDirectory: "web",
			Link:             "http://localhost:3000",
			ErrorPatterns:    []string{"ERROR", "FATAL"},
		},
	},
	CommandGroups: []projectdomain.BlueprintCommandGroup{
		{Id: "group-1", Name: "Everything", CommandIds: []string{"cmd-1"}},
	},
}

func TestProjectBlueprint(t *testing.T) {
	t.Run("Should reach the frontend under the names it reads", func(t *testing.T) {
		// Act
		sent, err := json.Marshal(transport.FromBlueprint(blueprint))

		// Assert
		require.NoError(t, err)
		assert.JSONEq(t, `{
			"name": "Gomander",
			"workingDirectory": "/home/user/app",
			"commands": [
				{
					"id": "cmd-1",
					"name": "Dev server",
					"command": "pnpm dev",
					"workingDirectory": "web",
					"link": "http://localhost:3000",
					"errorPatterns": ["ERROR", "FATAL"]
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
		returned := transport.FromBlueprint(blueprint).ToDomain()

		// Assert
		assert.Equal(t, blueprint, returned)
	})

	t.Run("Should give the frontend empty lists for a Blueprint that carries none", func(t *testing.T) {
		// Act
		sent, err := json.Marshal(transport.FromBlueprint(projectdomain.Blueprint{}))

		// Assert
		require.NoError(t, err)
		assert.JSONEq(t, `{
			"name": "",
			"workingDirectory": "",
			"commands": [],
			"commandGroups": []
		}`, string(sent))
	})

	t.Run("Should give the frontend an empty error pattern list when a command carries none", func(t *testing.T) {
		// Act
		sent := transport.FromBlueprint(projectdomain.Blueprint{
			Commands: []projectdomain.BlueprintCommand{
				{Id: "cmd-1", Name: "Tests", Command: "pnpm test"},
			},
		})

		// Assert
		assert.Equal(t, "", sent.Commands[0].Link)
		assert.Equal(t, []string{}, sent.Commands[0].ErrorPatterns)
		assert.Equal(t, []string{}, sent.ToDomain().Commands[0].ErrorPatterns)
	})
}
