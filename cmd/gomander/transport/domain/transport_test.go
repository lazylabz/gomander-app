package domain_test

import (
	"encoding/json"
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
}

func TestCommandGroup(t *testing.T) {
	commandGroup := commandgroupdomain.CommandGroup{
		Id:        "group-id",
		ProjectId: "project-id",
		Name:      "group-name",
		Commands:  []commanddomain.Command{command},
		Position:  2,
	}

	t.Run("Should reach the frontend under the field names it already receives", func(t *testing.T) {
		// Act
		encoded, err := json.Marshal(transport.FromCommandGroup(commandGroup))

		// Assert
		require.NoError(t, err)
		assert.JSONEq(t, `{
			"id": "group-id",
			"projectId": "project-id",
			"name": "group-name",
			"commands": [`+commandJSON+`],
			"position": 2
		}`, string(encoded))
	})

	t.Run("Should carry every field back to the entity", func(t *testing.T) {
		// Act
		roundTripped := transport.FromCommandGroup(commandGroup).ToDomain()

		// Assert
		assert.Equal(t, commandGroup, roundTripped)
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
}

func TestProjectExport(t *testing.T) {
	projectExport := projectdomain.ProjectExportJSONv1{
		Version:          1,
		Name:             "exported-name",
		WorkingDirectory: "/home/dev/project",
		Commands: []projectdomain.CommandJSONv1{{
			Id:               "command-id",
			Name:             "command-name",
			Command:          "echo hi",
			WorkingDirectory: "apps/api",
		}},
		CommandGroups: []projectdomain.CommandGroupJSONv1{{
			Id:         "group-id",
			Name:       "group-name",
			CommandIds: []string{"command-id"},
		}},
	}

	t.Run("Should reach the frontend under the field names it already receives", func(t *testing.T) {
		// Act
		encoded, err := json.Marshal(transport.FromProjectExport(projectExport))

		// Assert
		require.NoError(t, err)
		assert.JSONEq(t, `{
			"version": 1,
			"name": "exported-name",
			"workingDirectory": "/home/dev/project",
			"commands": [{
				"id": "command-id",
				"name": "command-name",
				"command": "echo hi",
				"workingDirectory": "apps/api"
			}],
			"commandGroups": [{
				"id": "group-id",
				"name": "group-name",
				"commandIds": ["command-id"]
			}]
		}`, string(encoded))
	})

	t.Run("Should carry every field back to the export format", func(t *testing.T) {
		// Act
		roundTripped := transport.FromProjectExport(projectExport).ToDomain()

		// Assert
		assert.Equal(t, projectExport, roundTripped)
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
