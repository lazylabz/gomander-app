package runner_test

import (
	"runtime"
	"strings"
	"testing"
	"time"

	commanddomain "gomander/internal/command/domain"
	"gomander/internal/event"
	test2 "gomander/internal/event/test"
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
		emitter := new(test2.MockEventEmitter)

		r := runner.NewDefaultRunner(logger, emitter)

		emitter.On("EmitEvent", event.ProcessStarted, commandId).Return()
		emitter.On("EmitEvent", event.ProcessFinished, commandId).Return()
		mockEmitterLogEntry(emitter, commandId, "a")
		mockEmitterLogEntry(emitter, commandId, "b")
		mockEmitterLogEntry(emitter, commandId, "c")

		// Check first line
		emitter.On("EmitEvent", event.NewLogEntry, mock.MatchedBy(func(
			data map[string]string) bool {
			return strings.Contains(data["line"], "echo")
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
		emitter := new(test2.MockEventEmitter)

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
		emitter := new(test2.MockEventEmitter)

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

	t.Run("Should not throw if trying to run an already running command", func(t *testing.T) {
		// Arrange
		logger := new(test.MockLogger)
		emitter := new(test2.MockEventEmitter)

		r := runner.NewDefaultRunner(logger, emitter)
		commandId := "1"

		emitter.On("EmitEvent", event.ProcessStarted, commandId).Return()
		emitter.On("EmitEvent", event.ProcessFinished, commandId).Maybe().Return()
		emitter.On("EmitEvent", event.NewLogEntry, mock.Anything).Return()

		logger.On("Info", mock.Anything).Return()
		logger.On("Debug", mock.Anything).Return()
		logger.On("Error", mock.Anything).Return()

		// Act
		command := commanddomain.Command{
			Id:               commandId,
			ProjectId:        commandId,
			Name:             "Test",
			Command:          infiniteCmd(),
			WorkingDirectory: validWorkingDirectory(),
			Position:         0,
		}
		err := r.RunCommand(&command, []string{}, "")
		assert.NoError(t, err)

		assert.Eventually(t, func() bool {
			return assert.NotEmpty(t, r.GetRunningCommands())
		}, 1*time.Second, 20*time.Millisecond)

		// Try to run the same command again
		err = r.RunCommand(&command, []string{}, "")

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 1, len(r.GetRunningCommands()))

		// Cleanup
		err = r.StopRunningCommand(commandId)
		assert.NoError(t, err)
		r.WaitForCommand(commandId)
	})
}

func TestDefaultRunner_StopAllRunningCommands(t *testing.T) {
	t.Run("Should stop all running commands", func(t *testing.T) {
		// Arrange
		logger := new(test.MockLogger)
		emitter := new(test2.MockEventEmitter)

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

func TestDefaultRunner_RunCommands(t *testing.T) {
	t.Run("Should run multiple commands successfully", func(t *testing.T) {
		// Arrange
		logger := new(test.MockLogger)
		emitter := new(test2.MockEventEmitter)

		r := runner.NewDefaultRunner(logger, emitter)

		cmd1Id := "1"
		cmd2Id := "2"

		emitter.On("EmitEvent", event.ProcessStarted, cmd1Id).Return()
		emitter.On("EmitEvent", event.ProcessFinished, cmd1Id).Maybe().Return()
		emitter.On("EmitEvent", event.ProcessStarted, cmd2Id).Return()
		emitter.On("EmitEvent", event.ProcessFinished, cmd2Id).Maybe().Return()
		emitter.On("EmitEvent", event.NewLogEntry, mock.MatchedBy(func(
			data map[string]string) bool {
			return strings.Contains(data["line"], "echo")
		})).Return()

		// Mock log entries for both commands
		mockEmitterLogEntry(emitter, cmd1Id, "command1 output")
		mockEmitterLogEntry(emitter, cmd2Id, "command2 output")

		logger.On("Info", mock.Anything).Return()
		logger.On("Debug", mock.Anything).Return()

		commands := []commanddomain.Command{
			{
				Id:               cmd1Id,
				ProjectId:        "project1",
				Name:             "Test Command 1",
				Command:          "echo 'command1 output'",
				WorkingDirectory: validWorkingDirectory(),
				Position:         0,
			},
			{
				Id:               cmd2Id,
				ProjectId:        "project1",
				Name:             "Test Command 2",
				Command:          "echo 'command2 output'",
				WorkingDirectory: validWorkingDirectory(),
				Position:         1,
			},
		}

		// Act
		err := r.RunCommands(commands, []string{"/test"}, "/test")

		// Wait for both commands to complete
		r.WaitForCommand(cmd1Id)
		r.WaitForCommand(cmd2Id)

		// Assert
		assert.NoError(t, err)
		assert.Empty(t, r.GetRunningCommands())
		mock.AssertExpectationsForObjects(t, emitter, logger)
	})

	t.Run("Should return error if any command fails to execute", func(t *testing.T) {
		// Arrange
		logger := new(test.MockLogger)
		emitter := new(test2.MockEventEmitter)

		r := runner.NewDefaultRunner(logger, emitter)

		cmd1Id := "1"
		cmd2Id := "2"

		// Mock for the first command to succeed
		emitter.On("EmitEvent", event.ProcessStarted, cmd1Id).Return()
		emitter.On("EmitEvent", event.ProcessFinished, cmd1Id).Maybe().Return()
		emitter.On("EmitEvent", event.ProcessFinished, cmd2Id).Maybe().Return()
		mockEmitterLogEntry(emitter, cmd1Id, "command1 output")

		// For the second command, we won't set expectations because it should
		// use a non-existent working directory, causing an error
		emitter.On("EmitEvent", event.NewLogEntry, mock.Anything).Return()

		logger.On("Info", mock.Anything).Return()
		logger.On("Debug", mock.Anything).Return()

		invalidWorkingDir := "/definitely/not/a/real/directory/12345"

		commands := []commanddomain.Command{
			{
				Id:               cmd1Id,
				ProjectId:        "project1",
				Name:             "Test Command 1",
				Command:          "echo 'command1 output'",
				WorkingDirectory: validWorkingDirectory(),
				Position:         0,
			},
			{
				Id:               cmd2Id,
				ProjectId:        "project1",
				Name:             "Test Command 2",
				Command:          "echo 'command2 output'",
				WorkingDirectory: invalidWorkingDir,
				Position:         1,
			},
		}

		// Act
		err := r.RunCommands(commands, []string{}, "")

		// Wait for the first command to complete
		r.WaitForCommand(cmd1Id)

		// Assert
		assert.Error(t, err)

		// Clean up any running commands
		r.StopAllRunningCommands()

		mock.AssertExpectationsForObjects(t, emitter, logger)
	})
}

func TestDefaultRunner_StopRunningCommands(t *testing.T) {
	t.Run("Should stop multiple running commands", func(t *testing.T) {
		// Arrange
		logger := new(test.MockLogger)
		emitter := new(test2.MockEventEmitter)

		r := runner.NewDefaultRunner(logger, emitter)

		cmd1Id := "1"
		cmd2Id := "2"

		emitter.On("EmitEvent", event.ProcessStarted, cmd1Id).Return()
		emitter.On("EmitEvent", event.ProcessStarted, cmd2Id).Return()
		emitter.On("EmitEvent", event.ProcessFinished, cmd1Id).Maybe().Return()
		emitter.On("EmitEvent", event.ProcessFinished, cmd2Id).Maybe().Return()
		emitter.On("EmitEvent", event.NewLogEntry, mock.Anything).Return()

		logger.On("Info", mock.Anything).Return()
		logger.On("Debug", mock.Anything).Return()
		// Depends on OS
		logger.On("Error", mock.Anything).Maybe().Return()

		// Start two long-running commands
		cmd1 := commanddomain.Command{
			Id:               cmd1Id,
			ProjectId:        "project1",
			Name:             "Test 1",
			Command:          infiniteCmd(),
			WorkingDirectory: validWorkingDirectory(),
			Position:         0,
		}
		cmd2 := commanddomain.Command{
			Id:               cmd2Id,
			ProjectId:        "project1",
			Name:             "Test 2",
			Command:          infiniteCmd(),
			WorkingDirectory: validWorkingDirectory(),
			Position:         1,
		}

		err := r.RunCommand(&cmd1, []string{}, "")
		assert.NoError(t, err)

		err = r.RunCommand(&cmd2, []string{}, "")
		assert.NoError(t, err)

		assert.Eventually(t, func() bool {
			return len(r.GetRunningCommands()) == 2
		}, 1*time.Second, 20*time.Millisecond)

		time.Sleep(500 * time.Millisecond) // Give some time for the commands to start

		// Act
		commands := []commanddomain.Command{cmd1, cmd2}
		err = r.StopRunningCommands(commands)

		// Wait for commands to stop
		r.WaitForCommand(cmd1Id)
		r.WaitForCommand(cmd2Id)

		// Assert
		assert.NoError(t, err)
		assert.Empty(t, r.GetRunningCommands())

		mock.AssertExpectationsForObjects(t, emitter, logger)
	})

	t.Run("Should not error when stopping non-running commands", func(t *testing.T) {
		// Arrange
		logger := new(test.MockLogger)
		emitter := new(test2.MockEventEmitter)

		r := runner.NewDefaultRunner(logger, emitter)

		// Define commands that aren't running
		cmd1 := commanddomain.Command{
			Id:               "1",
			ProjectId:        "project1",
			Name:             "Test 1",
			Command:          "echo test",
			WorkingDirectory: validWorkingDirectory(),
			Position:         0,
		}
		cmd2 := commanddomain.Command{
			Id:               "2",
			ProjectId:        "project1",
			Name:             "Test 2",
			Command:          "echo test",
			WorkingDirectory: validWorkingDirectory(),
			Position:         1,
		}

		// Act - try to stop commands that aren't running
		commands := []commanddomain.Command{cmd1, cmd2}
		err := r.StopRunningCommands(commands)

		// Assert
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, emitter, logger)
	})
}

func TestDefaultRunner_GetRunningCommandIds(t *testing.T) {
	t.Run("Should return empty list when no commands are running", func(t *testing.T) {
		// Arrange
		logger := new(test.MockLogger)
		logger.On("Info", mock.Anything).Return()
		logger.On("Debug", mock.Anything).Return()
		logger.On("Error", mock.Anything).Return()

		emitter := new(test2.MockEventEmitter)
		emitter.On("EmitEvent", mock.Anything, mock.Anything).Return()

		sut := runner.NewDefaultRunner(logger, emitter)

		// Act
		result := sut.GetRunningCommandIds()

		// Assert
		assert.Empty(t, result)
	})

	t.Run("Should return list of running command ids", func(t *testing.T) {
		// Arrange
		logger := new(test.MockLogger)
		logger.On("Info", mock.Anything).Return()
		logger.On("Debug", mock.Anything).Return()
		logger.On("Error", mock.Anything).Return()

		emitter := new(test2.MockEventEmitter)
		emitter.On("EmitEvent", mock.Anything, mock.Anything).Return()

		sut := runner.NewDefaultRunner(logger, emitter)

		// Create a few commands that will run for a short time
		command1 := &commanddomain.Command{
			Id:      "cmd-1",
			Command: infiniteCmd(),
		}
		command2 := &commanddomain.Command{
			Id:      "cmd-2",
			Command: infiniteCmd(),
		}

		// Start the commands
		_ = sut.RunCommand(command1, []string{}, validWorkingDirectory())
		_ = sut.RunCommand(command2, []string{}, validWorkingDirectory())

		// Give them a moment to start
		time.Sleep(10 * time.Millisecond)

		// Act
		result := sut.GetRunningCommandIds()

		// Assert
		assert.Len(t, result, 2)
		assert.Contains(t, result, "cmd-1")
		assert.Contains(t, result, "cmd-2")

		// Wait for the commands to finish so we don't affect other tests
		time.Sleep(200 * time.Millisecond)
	})
}

func mockEmitterLogEntry(emitter *test2.MockEventEmitter, id string, line string) {
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

func TestDefaultRunner_ErrorPatternDetection(t *testing.T) {
	t.Run("Should emit CommandErrorDetected event when error pattern is matched", func(t *testing.T) {
		// Arrange
		logger := new(test.MockLogger)
		emitter := new(test2.MockEventEmitter)

		r := runner.NewDefaultRunner(logger, emitter)

		commandId := "error-pattern-test"

		// Mock the standard events
		emitter.On("EmitEvent", event.ProcessStarted, commandId).Return()
		emitter.On("EmitEvent", event.ProcessFinished, commandId).Return()

		// Mock the error detection event - this is what we're testing
		emitter.On("EmitEvent", event.CommandErrorDetected, commandId).Return()

		// Mock log entries for the command output
		emitter.On("EmitEvent", event.NewLogEntry, mock.Anything).Return()

		logger.On("Info", mock.Anything).Return()
		logger.On("Debug", mock.Anything).Return()

		// Act - create a command with error patterns and output that matches them
		err := r.RunCommand(&commanddomain.Command{
			Id:               commandId,
			ProjectId:        commandId,
			Name:             "Error Pattern Test",
			Command:          "echo 'Starting...' && echo 'ERROR: Something went wrong' && echo 'Done'",
			WorkingDirectory: validWorkingDirectory(),
			Position:         0,
			ErrorPatterns: []string{
				"ERROR:",
				"FATAL:",
			},
		}, []string{}, "")

		r.WaitForCommand(commandId)

		// Assert
		assert.NoError(t, err)
		assert.Empty(t, r.GetRunningCommands())

		// Verify that CommandErrorDetected was called
		emitter.AssertCalled(t, "EmitEvent", event.CommandErrorDetected, commandId)
		mock.AssertExpectationsForObjects(t, emitter, logger)
	})

	t.Run("Should not emit CommandErrorDetected when no error pattern matches", func(t *testing.T) {
		// Arrange
		logger := new(test.MockLogger)
		emitter := new(test2.MockEventEmitter)

		r := runner.NewDefaultRunner(logger, emitter)

		commandId := "no-error-pattern-test"

		// Mock the standard events
		emitter.On("EmitEvent", event.ProcessStarted, commandId).Return()
		emitter.On("EmitEvent", event.ProcessFinished, commandId).Return()

		// Mock log entries for the command output
		emitter.On("EmitEvent", event.NewLogEntry, mock.Anything).Return()

		logger.On("Info", mock.Anything).Return()
		logger.On("Debug", mock.Anything).Return()

		// Act - create a command with error patterns but output that doesn't match
		err := r.RunCommand(&commanddomain.Command{
			Id:               commandId,
			ProjectId:        commandId,
			Name:             "No Error Pattern Test",
			Command:          "echo 'Starting...' && echo 'All good!' && echo 'Done'",
			WorkingDirectory: validWorkingDirectory(),
			Position:         0,
			ErrorPatterns: []string{
				"ERROR:",
				"FATAL:",
			},
		}, []string{}, "")

		r.WaitForCommand(commandId)

		// Assert
		assert.NoError(t, err)
		assert.Empty(t, r.GetRunningCommands())

		// Verify that CommandErrorDetected was NOT called
		emitter.AssertNotCalled(t, "EmitEvent", event.CommandErrorDetected, mock.Anything)
		mock.AssertExpectationsForObjects(t, emitter, logger)
	})

	t.Run("Should emit CommandErrorDetected for each matching line", func(t *testing.T) {
		// Arrange
		logger := new(test.MockLogger)
		emitter := new(test2.MockEventEmitter)

		r := runner.NewDefaultRunner(logger, emitter)

		commandId := "multiple-errors-test"

		// Mock the standard events
		emitter.On("EmitEvent", event.ProcessStarted, commandId).Return()
		emitter.On("EmitEvent", event.ProcessFinished, commandId).Return()

		// Mock the error detection event - will be called once per matching line (3 times in this case)
		emitter.On("EmitEvent", event.CommandErrorDetected, commandId).Return()

		// Mock log entries for the command output
		emitter.On("EmitEvent", event.NewLogEntry, mock.Anything).Return()

		logger.On("Info", mock.Anything).Return()
		logger.On("Debug", mock.Anything).Return()

		// Act - create a command with error patterns and output with multiple matching lines
		err := r.RunCommand(&commanddomain.Command{
			Id:               commandId,
			ProjectId:        commandId,
			Name:             "Multiple Errors Test",
			Command:          "echo 'ERROR: First error' && echo 'ERROR: Second error' && echo 'FATAL: Third error'",
			WorkingDirectory: validWorkingDirectory(),
			Position:         0,
			ErrorPatterns: []string{
				"ERROR:",
				"FATAL:",
			},
		}, []string{}, "")

		r.WaitForCommand(commandId)

		// Assert
		assert.NoError(t, err)
		assert.Empty(t, r.GetRunningCommands())

		// Verify that CommandErrorDetected was called (multiple times is OK - once per matching line)
		emitter.AssertCalled(t, "EmitEvent", event.CommandErrorDetected, commandId)
		mock.AssertExpectationsForObjects(t, emitter, logger)
	})
}

func infiniteCmd() string {
	if runtime.GOOS == "windows" {
		return "ping -t 127.0.0.1"
	}
	return "ping 127.0.0.1"
}
