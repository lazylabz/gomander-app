package domain_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	transport "gomander/cmd/gomander/transport/domain"
	commanddomain "gomander/internal/command/domain"
	commandgroupdomain "gomander/internal/commandgroup/domain"
	configdomain "gomander/internal/config/domain"
	projectdomain "gomander/internal/project/domain"
)

var command = commanddomain.Command{
	Id:               "command-id",
	ProjectId:        "project-id",
	Name:             "command-name",
	Command:          "echo hi",
	WorkingDirectory: "apps/api",
	Position:         3,
	Link:             "http://localhost:3000",
	ErrorPatterns:    []string{"ERR", "FATAL"},
}

const commandJSON = `{
	"id": "command-id",
	"projectId": "project-id",
	"name": "command-name",
	"command": "echo hi",
	"workingDirectory": "apps/api",
	"position": 3,
	"link": "http://localhost:3000",
	"errorPatterns": ["ERR", "FATAL"]
}`

func TestCommand(t *testing.T) {
	t.Run("Should reach the frontend under the field names it already receives", func(t *testing.T) {
		// Act
		encoded, err := json.Marshal(transport.FromCommand(command))

		// Assert
		require.NoError(t, err)
		assert.JSONEq(t, commandJSON, string(encoded))
	})

	t.Run("Should carry every field back to the entity", func(t *testing.T) {
		// Act
		roundTripped := transport.FromCommand(command).ToDomain()

		// Assert
		assert.Equal(t, command, roundTripped)
	})

	t.Run("Should say everything the entity says", func(t *testing.T) {
		assertMirrors(t, command, transport.FromCommand(command))
	})
}

func TestCommandGroup(t *testing.T) {
	commandGroup := commandgroupdomain.CommandGroup{
		Id:        "group-id",
		ProjectId: "project-id",
		Name:      "group-name",
		Commands:  []commanddomain.Command{command},
		Position:  2,
	}

	// The entity the DTO stands for: what the read path still carries, once
	// narrowed to the Commands it names.
	named := commandgroupdomain.CommandGroupWithCommandIds{
		Id:         "group-id",
		ProjectId:  "project-id",
		Name:       "group-name",
		CommandIds: []string{command.Id},
		Position:   2,
	}

	t.Run("Should reach the frontend naming the commands the group holds", func(t *testing.T) {
		// Act
		encoded, err := json.Marshal(transport.FromCommandGroup(commandGroup))

		// Assert
		require.NoError(t, err)
		assert.JSONEq(t, `{
			"id": "group-id",
			"projectId": "project-id",
			"name": "group-name",
			"commandIds": ["command-id"],
			"position": 2
		}`, string(encoded))
	})

	t.Run("Should carry every field back to the entity", func(t *testing.T) {
		// Act
		roundTripped := transport.FromCommandGroup(commandGroup).ToDomain()

		// Assert
		assert.Equal(t, named, roundTripped)
	})

	t.Run("Should say everything the entity says", func(t *testing.T) {
		assertMirrors(t, named, transport.FromCommandGroup(commandGroup))
	})
}

func TestProject(t *testing.T) {
	project := projectdomain.Project{
		Id:               "project-id",
		Name:             "project-name",
		WorkingDirectory: "/home/dev/project",
	}

	t.Run("Should reach the frontend under the field names it already receives", func(t *testing.T) {
		// Act
		encoded, err := json.Marshal(transport.FromProject(project))

		// Assert
		require.NoError(t, err)
		assert.JSONEq(t, `{
			"id": "project-id",
			"name": "project-name",
			"workingDirectory": "/home/dev/project"
		}`, string(encoded))
	})

	t.Run("Should carry every field back to the entity", func(t *testing.T) {
		// Act
		roundTripped := transport.FromProject(project).ToDomain()

		// Assert
		assert.Equal(t, project, roundTripped)
	})

	t.Run("Should say everything the entity says", func(t *testing.T) {
		assertMirrors(t, project, transport.FromProject(project))
	})
}

func TestConfig(t *testing.T) {
	config := configdomain.Config{
		LastOpenedProjectId: "project-id",
		EnvironmentPaths:    []configdomain.EnvironmentPath{{Id: "path-id", Path: "/usr/local/bin"}},
		Locale:              "es",
	}

	t.Run("Should reach the frontend under the field names it already receives", func(t *testing.T) {
		// Act
		encoded, err := json.Marshal(transport.FromConfig(config))

		// Assert
		require.NoError(t, err)
		assert.JSONEq(t, `{
			"lastOpenedProjectId": "project-id",
			"environmentPaths": [{"id": "path-id", "path": "/usr/local/bin"}],
			"locale": "es"
		}`, string(encoded))
	})

	t.Run("Should carry every field back to the entity", func(t *testing.T) {
		// Act
		roundTripped := transport.FromConfig(config).ToDomain()

		// Assert
		assert.Equal(t, config, roundTripped)
	})

	t.Run("Should say everything the entity says", func(t *testing.T) {
		assertMirrors(t, config, transport.FromConfig(config))
	})
}

func TestAbsence(t *testing.T) {
	t.Run("Should leave a value the frontend does not get absent", func(t *testing.T) {
		// Act
		mapped := transport.Optional(nil, transport.FromProject)

		// Assert
		assert.Nil(t, mapped)
	})

	t.Run("Should tell a list that is absent apart from one that is empty", func(t *testing.T) {
		// Act
		absent, err := json.Marshal(transport.FromCommands(nil))
		empty, emptyErr := json.Marshal(transport.FromCommands([]commanddomain.Command{}))

		// Assert
		require.NoError(t, err)
		require.NoError(t, emptyErr)
		assert.Equal(t, "null", string(absent))
		assert.Equal(t, "[]", string(empty))
	})
}

// assertMirrors catches a field the entity carries and the DTO does not, which
// the compiler cannot: a keyed struct literal stays valid when a field is added
// to the type it builds, so the mapping would drop the new one silently. It
// compares field names rather than the two serialized forms, because only the
// DTO carries serialization tags now.
func assertMirrors(t *testing.T, entity, dto any) {
	t.Helper()

	assert.Equal(t, fieldNames(entity), fieldNames(dto))
}

func fieldNames(value any) []string {
	structType := reflect.TypeOf(value)
	names := make([]string, 0, structType.NumField())
	for i := range structType.NumField() {
		names = append(names, structType.Field(i).Name)
	}
	return names
}
