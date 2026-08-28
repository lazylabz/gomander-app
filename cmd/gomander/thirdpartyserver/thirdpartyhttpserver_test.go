package thirdpartyserver_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"gomander/cmd/gomander/thirdpartyserver"
	"gomander/internal/apptest"
	commandtest "gomander/internal/command/domain/test"
	commandgrouptest "gomander/internal/commandgroup/domain/test"
	projecttest "gomander/internal/project/domain/test"
	"gomander/internal/usecases"
)

func TestNewThirdPartyIntegrationsServer_DiscoveryHandler(t *testing.T) {
	t.Run("GET /discovery should return discovery info", func(t *testing.T) {
		// Arrange
		testServer := serving(t, usecases.Registry{})
		defer testServer.Close()

		// Act
		resp, err := http.Get(testServer.URL + "/discovery")
		assert.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.JSONEq(t, `{"app": "Gomander"}`, string(body))
	})

	t.Run("POST /discovery should return 405 Method Not Allowed", func(t *testing.T) {
		// Arrange
		testServer := serving(t, usecases.Registry{})
		defer testServer.Close()

		// Act
		resp, err := http.Post(testServer.URL+"/discovery", "application/json", nil)
		assert.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	})
}

func TestNewThirdPartyIntegrationsServer_GetCommandsHandler(t *testing.T) {
	t.Run("GET /commands should return commands list with status", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)

		project := projecttest.NewProjectBuilder().Build()
		h.GivenProjects(project)
		h.GivenOpenedProject(project.Id)

		running := commandtest.NewCommandBuilder().
			WithId("cmd-1").
			WithProjectId(project.Id).
			WithName("Command 1").
			WithPosition(0).
			Build()
		stopped := commandtest.NewCommandBuilder().
			WithId("cmd-2").
			WithProjectId(project.Id).
			WithName("Command 2").
			WithPosition(1).
			Build()
		h.GivenCommands(running, stopped)

		assert.NoError(t, h.UseCases.RunCommand.Execute(running.Id))

		testServer := serving(t, h.UseCases)
		defer testServer.Close()

		// Act
		resp, err := http.Get(testServer.URL + "/commands")
		assert.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, []map[string]interface{}{
			{
				"id":     "cmd-1",
				"name":   "Command 1",
				"status": "running",
			},
			{
				"id":     "cmd-2",
				"name":   "Command 2",
				"status": "stopped",
			},
		}, decoded(t, resp))
	})

	t.Run("POST /commands should return 405 Method Not Allowed", func(t *testing.T) {
		// Arrange
		testServer := serving(t, usecases.Registry{})
		defer testServer.Close()

		// Act
		resp, err := http.Post(testServer.URL+"/commands", "application/json", nil)
		assert.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	})
}

func TestNewThirdPartyIntegrationsServer_RunCommandHandler(t *testing.T) {
	t.Run("POST /commands/{id}/run should run the command", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)

		project := projecttest.NewProjectBuilder().Build()
		h.GivenProjects(project)
		h.GivenOpenedProject(project.Id)

		command := commandtest.NewCommandBuilder().WithProjectId(project.Id).Build()
		h.GivenCommands(command)

		testServer := serving(t, h.UseCases)
		defer testServer.Close()

		// Act
		resp, err := http.Post(testServer.URL+"/commands/"+command.Id+"/run", "application/json", nil)
		assert.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, []string{command.Id}, h.UseCases.GetRunningCommandIds.Execute())
	})

	t.Run("GET /commands/{id}/run should return 405 Method Not Allowed", func(t *testing.T) {
		// Arrange
		testServer := serving(t, usecases.Registry{})
		defer testServer.Close()

		// Act
		resp, err := http.Get(testServer.URL + "/commands/cmd-1/run")
		assert.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	})

	t.Run("POST /commands/{id}/run should report a missing project", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)

		command := commandtest.NewCommandBuilder().WithProjectId("deleted-project").Build()
		h.GivenCommands(command)

		testServer := serving(t, h.UseCases)
		defer testServer.Close()

		// Act
		resp, err := http.Post(testServer.URL+"/commands/"+command.Id+"/run", "application/json", nil)
		assert.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestNewThirdPartyIntegrationsServer_StopCommandHandler(t *testing.T) {
	t.Run("POST /commands/{id}/stop should stop the command", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)

		project := projecttest.NewProjectBuilder().Build()
		h.GivenProjects(project)
		h.GivenOpenedProject(project.Id)

		command := commandtest.NewCommandBuilder().WithProjectId(project.Id).Build()
		h.GivenCommands(command)

		assert.NoError(t, h.UseCases.RunCommand.Execute(command.Id))

		testServer := serving(t, h.UseCases)
		defer testServer.Close()

		// Act
		resp, err := http.Post(testServer.URL+"/commands/"+command.Id+"/stop", "application/json", nil)
		assert.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, []string{command.Id}, h.StoppedProcessIds())
	})

	t.Run("GET /commands/{id}/stop should return 405 Method Not Allowed", func(t *testing.T) {
		// Arrange
		testServer := serving(t, usecases.Registry{})
		defer testServer.Close()

		// Act
		resp, err := http.Get(testServer.URL + "/commands/cmd-1/stop")
		assert.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	})

	t.Run("POST /commands/{id}/stop should report a command that is not there", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)

		testServer := serving(t, h.UseCases)
		defer testServer.Close()

		// Act
		resp, err := http.Post(testServer.URL+"/commands/deleted-command/stop", "application/json", nil)
		assert.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.Empty(t, h.StoppedProcessIds())
	})
}

func TestNewThirdPartyIntegrationsServer_RunCommandGroupHandler(t *testing.T) {
	t.Run("POST /command-groups/{id}/run should run the command group", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)

		project := projecttest.NewProjectBuilder().Build()
		h.GivenProjects(project)
		h.GivenOpenedProject(project.Id)

		first := commandtest.NewCommandBuilder().WithProjectId(project.Id).WithPosition(0).Build()
		second := commandtest.NewCommandBuilder().WithProjectId(project.Id).WithPosition(1).Build()
		h.GivenCommands(first, second)

		group := commandgrouptest.NewCommandGroupBuilder().
			WithProjectId(project.Id).
			WithCommands(first, second).
			Build()
		h.GivenCommandGroups(group)

		testServer := serving(t, h.UseCases)
		defer testServer.Close()

		// Act
		resp, err := http.Post(testServer.URL+"/command-groups/"+group.Id+"/run", "application/json", nil)
		assert.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, []string{first.Id, second.Id}, h.UseCases.GetRunningCommandIds.Execute())
	})

	t.Run("GET /command-groups/{id}/run should return 405 Method Not Allowed", func(t *testing.T) {
		// Arrange
		testServer := serving(t, usecases.Registry{})
		defer testServer.Close()

		// Act
		resp, err := http.Get(testServer.URL + "/command-groups/group-1/run")
		assert.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	})

	t.Run("POST /command-groups/{id}/run should report a missing project", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)

		command := commandtest.NewCommandBuilder().WithProjectId("deleted-project").Build()
		h.GivenCommands(command)

		group := commandgrouptest.NewCommandGroupBuilder().
			WithProjectId("deleted-project").
			WithCommands(command).
			Build()
		h.GivenCommandGroups(group)

		testServer := serving(t, h.UseCases)
		defer testServer.Close()

		// Act
		resp, err := http.Post(testServer.URL+"/command-groups/"+group.Id+"/run", "application/json", nil)
		assert.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("POST /command-groups/{id}/run should report a command group that is not there", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)

		testServer := serving(t, h.UseCases)
		defer testServer.Close()

		// Act
		resp, err := http.Post(testServer.URL+"/command-groups/deleted-command-group/run", "application/json", nil)
		assert.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.Empty(t, h.StartedProcesses())
	})
}

func TestNewThirdPartyIntegrationsServer_StopCommandGroupHandler(t *testing.T) {
	t.Run("POST /command-groups/{id}/stop should stop the command group", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)

		project := projecttest.NewProjectBuilder().Build()
		h.GivenProjects(project)
		h.GivenOpenedProject(project.Id)

		first := commandtest.NewCommandBuilder().WithProjectId(project.Id).WithPosition(0).Build()
		second := commandtest.NewCommandBuilder().WithProjectId(project.Id).WithPosition(1).Build()
		h.GivenCommands(first, second)

		group := commandgrouptest.NewCommandGroupBuilder().
			WithProjectId(project.Id).
			WithCommands(first, second).
			Build()
		h.GivenCommandGroups(group)

		assert.NoError(t, h.UseCases.RunCommandGroup.Execute(group.Id))

		testServer := serving(t, h.UseCases)
		defer testServer.Close()

		// Act
		resp, err := http.Post(testServer.URL+"/command-groups/"+group.Id+"/stop", "application/json", nil)
		assert.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, []string{first.Id, second.Id}, h.StoppedProcessIds())
	})

	t.Run("GET /command-groups/{id}/stop should return 405 Method Not Allowed", func(t *testing.T) {
		// Arrange
		testServer := serving(t, usecases.Registry{})
		defer testServer.Close()

		// Act
		resp, err := http.Get(testServer.URL + "/command-groups/group-1/stop")
		assert.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	})

	t.Run("POST /command-groups/{id}/stop should report a command group that is not there", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)

		testServer := serving(t, h.UseCases)
		defer testServer.Close()

		// Act
		resp, err := http.Post(testServer.URL+"/command-groups/deleted-command-group/stop", "application/json", nil)
		assert.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.Empty(t, h.StoppedProcessIds())
	})
}

func TestNewThirdPartyIntegrationsServer_GetCommandGroupsHandler(t *testing.T) {
	t.Run("GET /command-groups should return command groups list with running commands info", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)

		project := projecttest.NewProjectBuilder().Build()
		h.GivenProjects(project)
		h.GivenOpenedProject(project.Id)

		first := commandtest.NewCommandBuilder().WithProjectId(project.Id).WithPosition(0).Build()
		second := commandtest.NewCommandBuilder().WithProjectId(project.Id).WithPosition(1).Build()
		third := commandtest.NewCommandBuilder().WithProjectId(project.Id).WithPosition(2).Build()
		h.GivenCommands(first, second, third)

		firstGroup := commandgrouptest.NewCommandGroupBuilder().
			WithId("group-1").
			WithProjectId(project.Id).
			WithName("Group 1").
			WithPosition(0).
			WithCommands(first, second).
			Build()
		secondGroup := commandgrouptest.NewCommandGroupBuilder().
			WithId("group-2").
			WithProjectId(project.Id).
			WithName("Group 2").
			WithPosition(1).
			WithCommands(second, third).
			Build()
		h.GivenCommandGroups(firstGroup, secondGroup)

		assert.NoError(t, h.UseCases.RunCommand.Execute(first.Id))
		assert.NoError(t, h.UseCases.RunCommand.Execute(third.Id))

		testServer := serving(t, h.UseCases)
		defer testServer.Close()

		// Act
		resp, err := http.Get(testServer.URL + "/command-groups")
		assert.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, []map[string]interface{}{
			{
				"id":              "group-1",
				"name":            "Group 1",
				"commands":        float64(2),
				"runningCommands": float64(1),
			},
			{
				"id":              "group-2",
				"name":            "Group 2",
				"commands":        float64(2),
				"runningCommands": float64(1),
			},
		}, decoded(t, resp))
	})

	t.Run("POST /command-groups should return 405 Method Not Allowed", func(t *testing.T) {
		// Arrange
		testServer := serving(t, usecases.Registry{})
		defer testServer.Close()

		// Act
		resp, err := http.Post(testServer.URL+"/command-groups", "application/json", nil)
		assert.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	})
}

func TestThirdPartyIntegrationsServer_StartAndStop(t *testing.T) {
	t.Run("Should start and stop server without errors", func(t *testing.T) {
		// Arrange
		server := thirdpartyserver.NewThirdPartyIntegrationsServer(usecases.Registry{})
		err := server.RegisterHandlers()
		assert.NoError(t, err)

		// Act - Start the server
		server.Start()

		// Wait a bit for server to fully start
		time.Sleep(100 * time.Millisecond)

		// Assert - Check if server is running by making a request to the discovery endpoint
		var resp *http.Response
		var requestErr error

		assert.Eventually(t, func() bool {
			serverAddr := server.Server.Addr
			resp, requestErr = http.Get(fmt.Sprintf("http://%s/discovery", serverAddr))
			return requestErr == nil
		}, 500*time.Millisecond, 100*time.Millisecond, "Server should respond to /discovery requests")

		assert.NoError(t, requestErr)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Body.Close()

		// Act - Stop the server
		err = server.Stop()

		// Assert - Server stopped without error
		assert.NoError(t, err)
	})
}

func serving(t *testing.T, registry usecases.Registry) *httptest.Server {
	t.Helper()

	server := thirdpartyserver.NewThirdPartyIntegrationsServer(registry)
	if err := server.RegisterHandlers(); err != nil {
		t.Fatalf("failed to register the handlers: %v", err)
	}

	return httptest.NewServer(server.Server.Handler)
}

func decoded(t *testing.T, resp *http.Response) []map[string]interface{} {
	t.Helper()

	var decoded []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		t.Fatalf("failed to decode the response body: %v", err)
	}

	return decoded
}
