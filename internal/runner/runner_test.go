package runner_test

import (
	"runtime"
	"strings"
	"testing"
	"time"

	commanddomain "gomander/internal/command/domain"
	"gomander/internal/event"
	"gomander/internal/logger/test"
	"gomander/internal/runner"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func validWorkingDirectory() string {
	if runtime.GOOS == "windows" {
		return "C:\\"
	}
	return "/"
}

func TestDefaultRunner_RunCommand(t *testing.T) {
	commandId := "1"

	t.Run("Should run command with success and emit events for each line", func(t *testing.T) {
		// Arrange
		logger := new(test.MockLogger)
		emitter := new(event.MockEventEmitter)

		r := runner.NewDefaultRunner(logger, emitter)

		emitter.On("EmitEvent", event.ProcessStarted, commandId).Return()
		emitter.On("EmitEvent", event.ProcessFinished, commandId).Return()
		mockEmitterLogEntry(emitter, commandId, "a")
		mockEmitterLogEntry(emitter, commandId, "b")
		mockEmitterLogEntry(emitter, commandId, "c")

		// Check first line
		emitter.On("EmitEvent", event.NewLogEntry, mock.MatchedBy(func(
			data map[string]string) bool {
			return strings.Contains(data["line"], "echo") && strings.Contains(data["line"], "a") && strings.Contains(data["line"], "b") && strings.Contains(data["line"], "c")
		})).Return()

		logger.On("Info", mock.Anything).Return()
		logger.On("Debug", mock.Anything).Return()

		// Act
		err := r.RunCommand(&commanddomain.Command{
			Id:               commandId,
			ProjectId:        commandId,
			Name:             "Test",
			Command:          "echo 'a'&& echo 'b'&& echo 'c'",
			WorkingDirectory: validWorkingDirectory(),
			Position:         0,
		}, []string{"/test"}, "/test")
		r.WaitForCommand(commandId)

		// Assert
		assert.NoError(t, err)
		assert.Empty(t, r.GetRunningCommands())
		mock.AssertExpectationsForObjects(t, emitter, logger)
	})

	t.Run("Should log error when executing an invalid command", func(t *testing.T) {
		// Arrange
		logger := new(test.MockLogger)
		emitter := new(event.MockEventEmitter)

		r := runner.NewDefaultRunner(logger, emitter)

		emitter.On("EmitEvent", event.ProcessStarted, commandId).Return()
		emitter.On("EmitEvent", event.ProcessFinished, commandId).Return()
		// Not an amazing matcher, but different OSes will have different error messages
		emitter.On("EmitEvent", event.NewLogEntry, mock.Anything).Return()

		logger.On("Info", mock.Anything).Return()
		logger.On("Error", mock.Anything).Return()
		logger.On("Debug", mock.Anything).Return()

		// Act
		err := r.RunCommand(&commanddomain.Command{
			Id:               commandId,
			ProjectId:        commandId,
			Name:             "Test",
			Command:          "definitely-not-a-real-command-12345",
			WorkingDirectory: validWorkingDirectory(),
			Position:         0,
		}, []string{}, "")
		r.WaitForCommand(commandId)

		// Assert
		assert.NoError(t, err)
		assert.Empty(t, r.GetRunningCommands())
		mock.AssertExpectationsForObjects(t, emitter, logger)
	})
}

func TestDefaultRunner_StopRunningCommand(t *testing.T) {
	t.Run("Should stop running command", func(t *testing.T) {
		// Arrange
		logger := new(test.MockLogger)
		emitter := new(event.MockEventEmitter)

		r := runner.NewDefaultRunner(logger, emitter)

		commandId := "1"

		emitter.On("EmitEvent", event.ProcessStarted, commandId).Return()

		// Sometimes, in CI, this event is not emitted fast enough, so we use Maybe()
		emitter.On("EmitEvent", event.ProcessFinished, commandId).Maybe().Return()
		emitter.On("EmitEvent", event.NewLogEntry, mock.Anything).Return()

		logger.On("Info", mock.Anything).Return()
		logger.On("Debug", mock.Anything).Return()
		// Depends on OS
		logger.On("Error", mock.Anything).Maybe().Return()

		// Act
		err := r.RunCommand(&commanddomain.Command{
			Id:               commandId,
			ProjectId:        commandId,
			Name:             "Test",
			Command:          infiniteCmd(),
			WorkingDirectory: validWorkingDirectory(),
			Position:         0,
		}, []string{}, "")
		assert.NoError(t, err)

		assert.Eventually(t, func() bool {
			return assert.NotEmpty(t, r.GetRunningCommands())
		}, 1*time.Second, 20*time.Millisecond)

		time.Sleep(500 * time.Millisecond) // Give some time for the command to start and some logs to be emitted

		err = r.StopRunningCommand(commandId)
		r.WaitForCommand(commandId)

		// Assert
		assert.NoError(t, err)
		assert.Empty(t, r.GetRunningCommands())
		mock.AssertExpectationsForObjects(t, emitter, logger)
	})

	t.Run("Should return error when stopping non-existing command", func(t *testing.T) {
		// Arrange
		logger := new(test.MockLogger)
		emitter := new(event.MockEventEmitter)

		r := runner.NewDefaultRunner(logger, emitter)

		// Act
		err := r.StopRunningCommand("non-existing-command")

		// Assert
		assert.Error(t, err)
		assert.Equal(t, "No running command with id: non-existing-command", err.Error())
	})
}

func TestDefaultRunner_StopAllRunningCommands(t *testing.T) {
	t.Run("Should stop all running commands", func(t *testing.T) {
		// Arrange
		logger := new(test.MockLogger)
		emitter := new(event.MockEventEmitter)

		r := runner.NewDefaultRunner(logger, emitter)

		cmd1Id := "1"
		cmd2Id := "2"

		emitter.On("EmitEvent", event.ProcessStarted, cmd1Id).Return()
		emitter.On("EmitEvent", event.ProcessStarted, cmd2Id).Return()
		emitter.On("EmitEvent", event.ProcessFinished, cmd1Id).Return()
		emitter.On("EmitEvent", event.ProcessFinished, cmd2Id).Return()

		emitter.On("EmitEvent", event.NewLogEntry, mock.Anything).Return()

		logger.On("Info", mock.Anything).Return()
		logger.On("Debug", mock.Anything).Return()
		// Depends on OS
		logger.On("Error", mock.Anything).Maybe().Return()

		// Act
		err := r.RunCommand(&commanddomain.Command{
			Id:               cmd1Id,
			ProjectId:        cmd1Id,
			Name:             "Test",
			Command:          infiniteCmd(),
			WorkingDirectory: validWorkingDirectory(),
			Position:         0,
		}, []string{}, "")
		assert.NoError(t, err)

		err = r.RunCommand(&commanddomain.Command{
			Id:               cmd2Id,
			ProjectId:        cmd1Id,
			Name:             "Test",
			Command:          infiniteCmd(),
			WorkingDirectory: validWorkingDirectory(),
			Position:         0,
		}, []string{}, "")
		assert.NoError(t, err)

		assert.Eventually(t, func() bool {
			return assert.NotEmpty(t, r.GetRunningCommands())
		}, 1*time.Second, 20*time.Millisecond)

		time.Sleep(500 * time.Millisecond) // Give some time for the command to start and some logs to be emitted

		errs := r.StopAllRunningCommands()

		r.WaitForCommand(cmd1Id)
		r.WaitForCommand(cmd2Id)

		// Assert
		assert.Empty(t, errs)
		assert.Empty(t, r.GetRunningCommands())
	})
}

func mockEmitterLogEntry(emitter *event.MockEventEmitter, id string, line string) {
	if runtime.GOOS == "windows" {
		emitter.On("EmitEvent", event.NewLogEntry, map[string]string{
			"id":   id,
			"line": "'" + line + "'",
		}).Return()
	} else {
		emitter.On("EmitEvent", event.NewLogEntry, map[string]string{
			"id":   id,
			"line": line,
		}).Return()
	}
}

func infiniteCmd() string {
	if runtime.GOOS == "windows" {
		return "ping -t 127.0.0.1"
	}
	return "ping 127.0.0.1"
}
