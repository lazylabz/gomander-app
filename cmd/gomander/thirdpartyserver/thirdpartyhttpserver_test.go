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
	"github.com/stretchr/testify/mock"

	"gomander/cmd/gomander/thirdpartyserver"
	"gomander/internal/app"
	commandusecasestest "gomander/internal/command/application/usecases/test"
	commanddomain "gomander/internal/command/domain"
	commandgroupusecasestest "gomander/internal/commandgroup/application/usecases/test"
	commandgroupdomain "gomander/internal/commandgroup/domain"
)

func TestNewThirdPartyIntegrationsServer_DiscoveryHandler(t *testing.T) {
	t.Run("GET /discovery should return discovery info", func(t *testing.T) {
		// Arrange
		server := thirdpartyserver.NewThirdPartyIntegrationsServer(app.UseCases{})
		err := server.RegisterHandlers()
		assert.NoError(t, err)

		testServer := httptest.NewServer(server.Server.Handler)
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
		server := thirdpartyserver.NewThirdPartyIntegrationsServer(app.UseCases{})
		err := server.RegisterHandlers()
		assert.NoError(t, err)

		testServer := httptest.NewServer(server.Server.Handler)
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
		mockGetCommands := new(commandusecasestest.MockGetCommands)
		mockGetRunningCommandIds := new(commandusecasestest.MockGetRunningCommandIds)

		commands := []commanddomain.Command{
			{
				Id:   "cmd-1",
				Name: "Command 1",
			},
			{
				Id:   "cmd-2",
				Name: "Command 2",
			},
		}

		runningCommandIds := []string{"cmd-1"}

		mockGetCommands.On("Execute").Return(commands, nil)
		mockGetRunningCommandIds.On("Execute").Return(runningCommandIds)

		useCases := app.UseCases{
			GetCommands:          mockGetCommands,
			GetRunningCommandIds: mockGetRunningCommandIds,
		}

		server := thirdpartyserver.NewThirdPartyIntegrationsServer(useCases)
		err := server.RegisterHandlers()
		assert.NoError(t, err)

		testServer := httptest.NewServer(server.Server.Handler)
		defer testServer.Close()

		// Act
		resp, err := http.Get(testServer.URL + "/commands")
		assert.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result []map[string]interface{}
		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		err = json.Unmarshal(body, &result)
		assert.NoError(t, err)

		assert.Equal(t, result, []map[string]interface{}{
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
		})

		mock.AssertExpectationsForObjects(t, mockGetRunningCommandIds, mockGetCommands)
	})

	t.Run("POST /commands should return 405 Method Not Allowed", func(t *testing.T) {
		// Arrange
		server := thirdpartyserver.NewThirdPartyIntegrationsServer(app.UseCases{})
		err := server.RegisterHandlers()
		assert.NoError(t, err)

		testServer := httptest.NewServer(server.Server.Handler)
		defer testServer.Close()

		// Act
		resp, err := http.Post(testServer.URL+"/commands", "application/json", nil)
		assert.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	})

	t.Run("Should return 500 when GetCommands returns error", func(t *testing.T) {
		// Arrange
		mockGetCommands := new(commandusecasestest.MockGetCommands)
		expectedError := fmt.Errorf("database error")

		mockGetCommands.On("Execute").Return([]commanddomain.Command{}, expectedError)

		useCases := app.UseCases{
			GetCommands: mockGetCommands,
		}

		server := thirdpartyserver.NewThirdPartyIntegrationsServer(useCases)
		err := server.RegisterHandlers()
		assert.NoError(t, err)

		testServer := httptest.NewServer(server.Server.Handler)
		defer testServer.Close()

		// Act
		resp, err := http.Get(testServer.URL + "/commands")
		assert.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockGetCommands.AssertExpectations(t)
	})
}

// Test Run Command Handler
func TestNewThirdPartyIntegrationsServer_RunCommandHandler(t *testing.T) {
	t.Run("POST /commands/{id}/run should run the command", func(t *testing.T) {
		// Arrange
		mockRunCommand := new(commandusecasestest.MockRunCommands)
		commandId := "cmd-1"

		mockRunCommand.On("Execute", commandId).Return(nil)

		useCases := app.UseCases{
			RunCommand: mockRunCommand,
		}

		server := thirdpartyserver.NewThirdPartyIntegrationsServer(useCases)
		err := server.RegisterHandlers()
		assert.NoError(t, err)

		testServer := httptest.NewServer(server.Server.Handler)
		defer testServer.Close()

		// Act
		resp, err := http.Post(testServer.URL+"/commands/"+commandId+"/run", "application/json", nil)
		assert.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		mock.AssertExpectationsForObjects(t, mockRunCommand)
	})

	t.Run("GET /commands/{id}/run should return 405 Method Not Allowed", func(t *testing.T) {
		// Arrange
		server := thirdpartyserver.NewThirdPartyIntegrationsServer(app.UseCases{})
		err := server.RegisterHandlers()
		assert.NoError(t, err)

		testServer := httptest.NewServer(server.Server.Handler)
		defer testServer.Close()

		// Act
		resp, err := http.Get(testServer.URL + "/commands/cmd-1/run")
		assert.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	})

	t.Run("Should return 500 when RunCommand returns error", func(t *testing.T) {
		// Arrange
		mockRunCommand := new(commandusecasestest.MockRunCommands)
		commandId := "cmd-1"
		expectedError := fmt.Errorf("failed to run command")

		mockRunCommand.On("Execute", commandId).Return(expectedError)

		useCases := app.UseCases{
			RunCommand: mockRunCommand,
		}

		server := thirdpartyserver.NewThirdPartyIntegrationsServer(useCases)
		err := server.RegisterHandlers()
		assert.NoError(t, err)

		testServer := httptest.NewServer(server.Server.Handler)
		defer testServer.Close()

		// Act
		resp, err := http.Post(testServer.URL+"/commands/"+commandId+"/run", "application/json", nil)
		assert.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockRunCommand.AssertExpectations(t)
	})
}

// Test Stop Command Handler
func TestNewThirdPartyIntegrationsServer_StopCommandHandler(t *testing.T) {
	t.Run("POST /commands/{id}/stop should stop the command", func(t *testing.T) {
		// Arrange
		mockStopCommand := new(commandusecasestest.MockStopCommand)
		commandId := "cmd-1"

		mockStopCommand.On("Execute", commandId).Return(nil)

		useCases := app.UseCases{
			StopCommand: mockStopCommand,
		}

		server := thirdpartyserver.NewThirdPartyIntegrationsServer(useCases)
		err := server.RegisterHandlers()
		assert.NoError(t, err)

		testServer := httptest.NewServer(server.Server.Handler)
		defer testServer.Close()

		// Act
		resp, err := http.Post(testServer.URL+"/commands/"+commandId+"/stop", "application/json", nil)
		assert.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		mock.AssertExpectationsForObjects(t, mockStopCommand)
	})

	t.Run("GET /commands/{id}/stop should return 405 Method Not Allowed", func(t *testing.T) {
		// Arrange
		server := thirdpartyserver.NewThirdPartyIntegrationsServer(app.UseCases{})
		err := server.RegisterHandlers()
		assert.NoError(t, err)

		testServer := httptest.NewServer(server.Server.Handler)
		defer testServer.Close()

		// Act
		resp, err := http.Get(testServer.URL + "/commands/cmd-1/stop")
		assert.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	})

	t.Run("Should return 500 when StopCommand returns error", func(t *testing.T) {
		// Arrange
		mockStopCommand := new(commandusecasestest.MockStopCommand)
		commandId := "cmd-1"
		expectedError := fmt.Errorf("failed to stop command")

		mockStopCommand.On("Execute", commandId).Return(expectedError)

		useCases := app.UseCases{
			StopCommand: mockStopCommand,
		}

		server := thirdpartyserver.NewThirdPartyIntegrationsServer(useCases)
		err := server.RegisterHandlers()
		assert.NoError(t, err)

		testServer := httptest.NewServer(server.Server.Handler)
		defer testServer.Close()

		// Act
		resp, err := http.Post(testServer.URL+"/commands/"+commandId+"/stop", "application/json", nil)
		assert.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockStopCommand.AssertExpectations(t)
	})
}

// Test Run Command Group Handler
func TestNewThirdPartyIntegrationsServer_RunCommandGroupHandler(t *testing.T) {
	t.Run("POST /command-groups/{id}/run should run the command group", func(t *testing.T) {
		// Arrange
		mockRunCommandGroup := new(commandgroupusecasestest.MockRunCommandGroup)
		groupId := "group-1"

		mockRunCommandGroup.On("Execute", groupId).Return(nil)

		useCases := app.UseCases{
			RunCommandGroup: mockRunCommandGroup,
		}

		server := thirdpartyserver.NewThirdPartyIntegrationsServer(useCases)
		err := server.RegisterHandlers()
		assert.NoError(t, err)

		testServer := httptest.NewServer(server.Server.Handler)
		defer testServer.Close()

		// Act
		resp, err := http.Post(testServer.URL+"/command-groups/"+groupId+"/run", "application/json", nil)
		assert.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		mock.AssertExpectationsForObjects(t, mockRunCommandGroup)
	})

	t.Run("GET /command-groups/{id}/run should return 405 Method Not Allowed", func(t *testing.T) {
		// Arrange
		server := thirdpartyserver.NewThirdPartyIntegrationsServer(app.UseCases{})
		err := server.RegisterHandlers()
		assert.NoError(t, err)

		testServer := httptest.NewServer(server.Server.Handler)
		defer testServer.Close()

		// Act
		resp, err := http.Get(testServer.URL + "/command-groups/group-1/run")
		assert.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	})

	t.Run("Should return 500 when RunCommandGroup returns error", func(t *testing.T) {
		// Arrange
		mockRunCommandGroup := new(commandgroupusecasestest.MockRunCommandGroup)
		groupId := "group-1"
		expectedError := fmt.Errorf("failed to run command group")

		mockRunCommandGroup.On("Execute", groupId).Return(expectedError)

		useCases := app.UseCases{
			RunCommandGroup: mockRunCommandGroup,
		}

		server := thirdpartyserver.NewThirdPartyIntegrationsServer(useCases)
		err := server.RegisterHandlers()
		assert.NoError(t, err)

		testServer := httptest.NewServer(server.Server.Handler)
		defer testServer.Close()

		// Act
		resp, err := http.Post(testServer.URL+"/command-groups/"+groupId+"/run", "application/json", nil)
		assert.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockRunCommandGroup.AssertExpectations(t)
	})
}

// Test Stop Command Group Handler
func TestNewThirdPartyIntegrationsServer_StopCommandGroupHandler(t *testing.T) {
	t.Run("POST /command-groups/{id}/stop should stop the command group", func(t *testing.T) {
		// Arrange
		mockStopCommandGroup := new(commandgroupusecasestest.MockStopCommandGroup)
		groupId := "group-1"

		mockStopCommandGroup.On("Execute", groupId).Return(nil)

		useCases := app.UseCases{
			StopCommandGroup: mockStopCommandGroup,
		}

		server := thirdpartyserver.NewThirdPartyIntegrationsServer(useCases)
		err := server.RegisterHandlers()
		assert.NoError(t, err)

		testServer := httptest.NewServer(server.Server.Handler)
		defer testServer.Close()

		// Act
		resp, err := http.Post(testServer.URL+"/command-groups/"+groupId+"/stop", "application/json", nil)
		assert.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		mock.AssertExpectationsForObjects(t, mockStopCommandGroup)
	})

	t.Run("GET /command-groups/{id}/stop should return 405 Method Not Allowed", func(t *testing.T) {
		// Arrange
		server := thirdpartyserver.NewThirdPartyIntegrationsServer(app.UseCases{})
		err := server.RegisterHandlers()
		assert.NoError(t, err)

		testServer := httptest.NewServer(server.Server.Handler)
		defer testServer.Close()

		// Act
		resp, err := http.Get(testServer.URL + "/command-groups/group-1/stop")
		assert.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	})

	t.Run("Should return 500 when StopCommandGroup returns error", func(t *testing.T) {
		// Arrange
		mockStopCommandGroup := new(commandgroupusecasestest.MockStopCommandGroup)
		groupId := "group-1"
		expectedError := fmt.Errorf("failed to stop command group")

		mockStopCommandGroup.On("Execute", groupId).Return(expectedError)

		useCases := app.UseCases{
			StopCommandGroup: mockStopCommandGroup,
		}

		server := thirdpartyserver.NewThirdPartyIntegrationsServer(useCases)
		err := server.RegisterHandlers()
		assert.NoError(t, err)

		testServer := httptest.NewServer(server.Server.Handler)
		defer testServer.Close()

		// Act
		resp, err := http.Post(testServer.URL+"/command-groups/"+groupId+"/stop", "application/json", nil)
		assert.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockStopCommandGroup.AssertExpectations(t)
	})
}

// Test Get Command Groups Handler
func TestNewThirdPartyIntegrationsServer_GetCommandGroupsHandler(t *testing.T) {
	t.Run("GET /command-groups should return command groups list with running commands info", func(t *testing.T) {
		// Arrange
		mockGetCommandGroups := new(commandgroupusecasestest.MockGetCommandGroups)
		mockGetRunningCommandIds := new(commandusecasestest.MockGetRunningCommandIds)

		cmd1 := commanddomain.Command{Id: "cmd-1", Name: "Command 1", Command: "echo 1"}
		cmd2 := commanddomain.Command{Id: "cmd-2", Name: "Command 2", Command: "echo 2"}
		cmd3 := commanddomain.Command{Id: "cmd-3", Name: "Command 3", Command: "echo 3"}

		groups := []commandgroupdomain.CommandGroup{
			{
				Id:       "group-1",
				Name:     "Group 1",
				Commands: []commanddomain.Command{cmd1, cmd2},
			},
			{
				Id:       "group-2",
				Name:     "Group 2",
				Commands: []commanddomain.Command{cmd2, cmd3},
			},
		}

		runningCommandIds := []string{"cmd-1", "cmd-3"}

		mockGetCommandGroups.On("Execute").Return(groups, nil)
		mockGetRunningCommandIds.On("Execute").Return(runningCommandIds)

		useCases := app.UseCases{
			GetCommandGroups:     mockGetCommandGroups,
			GetRunningCommandIds: mockGetRunningCommandIds,
		}

		server := thirdpartyserver.NewThirdPartyIntegrationsServer(useCases)
		err := server.RegisterHandlers()
		assert.NoError(t, err)

		testServer := httptest.NewServer(server.Server.Handler)
		defer testServer.Close()

		// Act
		resp, err := http.Get(testServer.URL + "/command-groups")
		assert.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result []map[string]interface{}
		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		err = json.Unmarshal(body, &result)
		assert.NoError(t, err)

		// First group should have 2 commands with 1 running
		assert.Equal(t, result, []map[string]interface{}{
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
		})

		mock.AssertExpectationsForObjects(t, mockGetCommandGroups, mockGetRunningCommandIds)
	})

	t.Run("POST /command-groups should return 405 Method Not Allowed", func(t *testing.T) {
		// Arrange
		server := thirdpartyserver.NewThirdPartyIntegrationsServer(app.UseCases{})
		err := server.RegisterHandlers()
		assert.NoError(t, err)

		testServer := httptest.NewServer(server.Server.Handler)
		defer testServer.Close()

		// Act
		resp, err := http.Post(testServer.URL+"/command-groups", "application/json", nil)
		assert.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	})

	t.Run("Should return 500 when GetCommandGroups returns error", func(t *testing.T) {
		// Arrange
		mockGetCommandGroups := new(commandgroupusecasestest.MockGetCommandGroups)
		expectedError := fmt.Errorf("failed to get command groups")

		mockGetCommandGroups.On("Execute").Return([]commandgroupdomain.CommandGroup{}, expectedError)

		useCases := app.UseCases{
			GetCommandGroups: mockGetCommandGroups,
		}

		server := thirdpartyserver.NewThirdPartyIntegrationsServer(useCases)
		err := server.RegisterHandlers()
		assert.NoError(t, err)

		testServer := httptest.NewServer(server.Server.Handler)
		defer testServer.Close()

		// Act
		resp, err := http.Get(testServer.URL + "/command-groups")
		assert.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockGetCommandGroups.AssertExpectations(t)
	})
}

func TestThirdPartyIntegrationsServer_StartAndStop(t *testing.T) {
	t.Run("Should start and stop server without errors", func(t *testing.T) {
		// Arrange
		server := thirdpartyserver.NewThirdPartyIntegrationsServer(app.UseCases{})
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
