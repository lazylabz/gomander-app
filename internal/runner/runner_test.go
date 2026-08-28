package runner_test

import (
	"os/exec"
	"runtime"
	"strings"
	"testing"
	"time"

	commanddomain "gomander/internal/command/domain"
	"gomander/internal/event"
	test2 "gomander/internal/event/test"
	"gomander/internal/execution"
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

		// The opening line is the Command's own text, with no presentation on it
		emitter.On("EmitEvent", event.NewLogEntry, map[string]string{
			"id":   commandId,
			"line": "echo 'a'&& echo 'b'&& echo 'c'",
			"kind": "command",
		}).Return()

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
		}, execution.Environment{Paths: []string{"/test"}, BaseWorkingDirectory: "/test"})
		r.WaitForCommand(commandId)

		// Assert
		assert.NoError(t, err)
		assert.Empty(t, r.GetRunningCommandIds())
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
		}, execution.Environment{})
		r.WaitForCommand(commandId)

		// Assert
		assert.NoError(t, err)
		assert.Empty(t, r.GetRunningCommandIds())
		mock.AssertExpectationsForObjects(t, emitter, logger)
	})

	t.Run("Should stream the output of several commands running at once", func(t *testing.T) {
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

		mockEmitterLogEntry(emitter, cmd1Id, "command1 output")
		mockEmitterLogEntry(emitter, cmd2Id, "command2 output")

		logger.On("Info", mock.Anything).Return()
		logger.On("Debug", mock.Anything).Return()

		// Act
		err := r.RunCommand(&commanddomain.Command{
			Id:               cmd1Id,
			ProjectId:        "project1",
			Name:             "Test Command 1",
			Command:          "echo 'command1 output'",
			WorkingDirectory: validWorkingDirectory(),
			Position:         0,
		}, execution.Environment{Paths: []string{"/test"}, BaseWorkingDirectory: "/test"})
		assert.NoError(t, err)

		err = r.RunCommand(&commanddomain.Command{
			Id:               cmd2Id,
			ProjectId:        "project1",
			Name:             "Test Command 2",
			Command:          "echo 'command2 output'",
			WorkingDirectory: validWorkingDirectory(),
			Position:         1,
		}, execution.Environment{Paths: []string{"/test"}, BaseWorkingDirectory: "/test"})
		assert.NoError(t, err)

		r.WaitForCommand(cmd1Id)
		r.WaitForCommand(cmd2Id)

		// Assert
		assert.Empty(t, r.GetRunningCommandIds())
		mock.AssertExpectationsForObjects(t, emitter, logger)
	})

	t.Run("Should return an error when the working directory does not exist", func(t *testing.T) {
		// Arrange
		logger := new(test.MockLogger)
		emitter := new(test2.MockEventEmitter)

		r := runner.NewDefaultRunner(logger, emitter)

		emitter.On("EmitEvent", event.NewLogEntry, mock.Anything).Return()

		// Act
		err := r.RunCommand(&commanddomain.Command{
			Id:               commandId,
			ProjectId:        "project1",
			Name:             "Test",
			Command:          "echo 'never runs'",
			WorkingDirectory: "/definitely/not/a/real/directory/12345",
			Position:         0,
		}, execution.Environment{})

		// Assert
		assert.Error(t, err)
		assert.Empty(t, r.GetRunningCommandIds())
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
		}, execution.Environment{})
		assert.NoError(t, err)

		assert.Eventually(t, func() bool {
			return assert.NotEmpty(t, r.GetRunningCommandIds())
		}, 1*time.Second, 20*time.Millisecond)

		time.Sleep(500 * time.Millisecond) // Give some time for the command to start and some logs to be emitted

		err = r.StopRunningCommand(commandId)
		r.WaitForCommand(commandId)

		// Assert
		assert.NoError(t, err)
		assert.Empty(t, r.GetRunningCommandIds())
		mock.AssertExpectationsForObjects(t, emitter, logger)
	})

	t.Run("Should not write a termination error to the output of a stopped command", func(t *testing.T) {
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
		logger.On("Error", mock.Anything).Maybe().Return()

		// Act
		err := r.RunCommand(&commanddomain.Command{
			Id:               commandId,
			ProjectId:        commandId,
			Name:             "Test",
			Command:          infiniteCmd(),
			WorkingDirectory: validWorkingDirectory(),
			Position:         0,
		}, execution.Environment{})
		assert.NoError(t, err)

		assert.Eventually(t, func() bool {
			return assert.NotEmpty(t, r.GetRunningCommandIds())
		}, 1*time.Second, 20*time.Millisecond)

		assert.NoError(t, r.StopRunningCommand(commandId))
		r.WaitForCommand(commandId)

		// Assert
		for _, line := range streamedLines(emitter) {
			assert.NotContains(t, runner.ExpectedTerminationLogs, line,
				"stopping a Command is not an error the user should read in its output")
		}
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
		err := r.RunCommand(&command, execution.Environment{})
		assert.NoError(t, err)

		assert.Eventually(t, func() bool {
			return assert.NotEmpty(t, r.GetRunningCommandIds())
		}, 1*time.Second, 20*time.Millisecond)

		// Try to run the same command again
		err = r.RunCommand(&command, execution.Environment{})

		// Assert
		assert.NoError(t, err)
		assert.Len(t, r.GetRunningCommandIds(), 1)

		// Cleanup
		err = r.StopRunningCommand(commandId)
		assert.NoError(t, err)
		r.WaitForCommand(commandId)
	})

	t.Run("Should stop each of several commands running at once", func(t *testing.T) {
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

		assert.NoError(t, r.RunCommand(&cmd1, execution.Environment{}))
		assert.NoError(t, r.RunCommand(&cmd2, execution.Environment{}))

		assert.Eventually(t, func() bool {
			return len(r.GetRunningCommandIds()) == 2
		}, 1*time.Second, 20*time.Millisecond)

		time.Sleep(500 * time.Millisecond) // Give some time for the commands to start

		// Act
		err := r.StopRunningCommand(cmd1Id)
		assert.NoError(t, err)

		err = r.StopRunningCommand(cmd2Id)
		assert.NoError(t, err)

		r.WaitForCommand(cmd1Id)
		r.WaitForCommand(cmd2Id)

		// Assert
		assert.Empty(t, r.GetRunningCommandIds())
		mock.AssertExpectationsForObjects(t, emitter, logger)
	})

	t.Run("Should finish tearing down when a descendant outlives the process group", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("the process group is a unix concept; on windows every grandchild already outlives taskkill /PID")
		}
		if _, err := exec.LookPath("python3"); err != nil {
			t.Skip("needs python3 to leave the process group")
		}

		// Arrange
		logger := new(test.MockLogger)
		emitter := new(test2.MockEventEmitter)

		r := runner.NewDefaultRunner(logger, emitter)

		commandId := "1"

		logger.On("Info", mock.Anything).Return()
		logger.On("Debug", mock.Anything).Return()
		logger.On("Error", mock.Anything).Maybe().Return()
		emitter.On("EmitEvent", mock.Anything, mock.Anything).Return()

		// The python child calls setsid, so it leaves the group the runner signals
		// while still holding the stdout it inherited.
		escapes := "python3 -c 'import os,time; os.setsid(); time.sleep(30)' & " + infiniteCmd()

		// Act
		assert.NoError(t, r.RunCommand(&commanddomain.Command{
			Id:               commandId,
			ProjectId:        commandId,
			Name:             "Test",
			Command:          escapes,
			WorkingDirectory: validWorkingDirectory(),
			Position:         0,
		}, execution.Environment{}))

		assert.Eventually(t, func() bool {
			return assert.NotEmpty(t, r.GetRunningCommandIds())
		}, 1*time.Second, 20*time.Millisecond)

		time.Sleep(500 * time.Millisecond) // Give the descendant time to be spawned

		assert.NoError(t, r.StopRunningCommand(commandId))

		tornDown := make(chan struct{})
		go func() {
			r.WaitForCommand(commandId)
			close(tornDown)
		}()

		// Assert
		select {
		case <-tornDown:
		case <-time.After(15 * time.Second):
			t.Fatal("the survivor held the output pipe open, so the Command never finished tearing down")
		}
		assert.Empty(t, r.GetRunningCommandIds())
	})

	t.Run("Should not error when stopping a command that is not running", func(t *testing.T) {
		// Arrange
		logger := new(test.MockLogger)
		emitter := new(test2.MockEventEmitter)

		r := runner.NewDefaultRunner(logger, emitter)

		// Act
		err := r.StopRunningCommand("not-running")

		// Assert
		assert.NoError(t, err)
		mock.AssertExpectationsForObjects(t, emitter, logger)
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
		}, execution.Environment{})
		assert.NoError(t, err)

		err = r.RunCommand(&commanddomain.Command{
			Id:               cmd2Id,
			ProjectId:        cmd1Id,
			Name:             "Test",
			Command:          infiniteCmd(),
			WorkingDirectory: validWorkingDirectory(),
			Position:         0,
		}, execution.Environment{})
		assert.NoError(t, err)

		assert.Eventually(t, func() bool {
			return assert.NotEmpty(t, r.GetRunningCommandIds())
		}, 1*time.Second, 20*time.Millisecond)

		time.Sleep(500 * time.Millisecond) // Give some time for the command to start and some logs to be emitted

		errs := r.StopAllRunningCommands()

		r.WaitForCommand(cmd1Id)
		r.WaitForCommand(cmd2Id)

		// Assert
		assert.Empty(t, errs)
		assert.Empty(t, r.GetRunningCommandIds())
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
		_ = sut.RunCommand(command1, execution.Environment{BaseWorkingDirectory: validWorkingDirectory()})
		_ = sut.RunCommand(command2, execution.Environment{BaseWorkingDirectory: validWorkingDirectory()})

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
			"kind": "output",
		}).Return()
	} else {
		emitter.On("EmitEvent", event.NewLogEntry, map[string]string{
			"id":   id,
			"line": line,
			"kind": "output",
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
		}, execution.Environment{})

		r.WaitForCommand(commandId)

		// Assert
		assert.NoError(t, err)
		assert.Empty(t, r.GetRunningCommandIds())

		// Verify that CommandErrorDetected was called
		mock.AssertExpectationsForObjects(t, emitter, logger)
	})
}

// streamedLines is what the runner wrote to the Command's terminal. Safe to
// read once WaitForCommand has returned: nothing is emitting any more.
func streamedLines(emitter *test2.MockEventEmitter) []string {
	lines := make([]string, 0)

	for _, call := range emitter.Calls {
		if call.Method != "EmitEvent" || call.Arguments.Get(0) != event.NewLogEntry {
			continue
		}
		lines = append(lines, call.Arguments.Get(1).(map[string]string)["line"])
	}

	return lines
}

func infiniteCmd() string {
	if runtime.GOOS == "windows" {
		return "ping -t 127.0.0.1"
	}
	return "ping 127.0.0.1"
}
