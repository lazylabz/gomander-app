package runner_test

import (
	"runtime"
	"testing"
	"time"

	commanddomain "gomander/internal/command/domain"
	"gomander/internal/event"
	"gomander/internal/runner"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Info(message string) {
	m.Called(message)
}

func (m *MockLogger) Debug(message string) {
	m.Called(message)
}

func (m *MockLogger) Error(message string) {
	m.Called(message)
}

type MockEventEmitter struct {
	mock.Mock
}

func (m *MockEventEmitter) EmitEvent(event event.Event, payload interface{}) {
	m.Called(event, payload)
}

const (
	DefaultTimeout         = 5 * time.Second       // seconds
	DefaultPollingInterval = 20 * time.Millisecond // milliseconds
)

func TestDefaultRunner_RunCommand(t *testing.T) {
	t.Run("Should run command with success and emit events for each line", func(t *testing.T) {
		logger := new(MockLogger)
		emitter := new(MockEventEmitter)

		r := runner.NewDefaultRunner(logger, emitter)

		emitter.On("EmitEvent", event.ProcessStarted, "1").Return()
		emitter.On("EmitEvent", event.ProcessFinished, "1").Return()
		mockEmitterLogEntry(emitter, "1", "a")
		mockEmitterLogEntry(emitter, "1", "b")
		mockEmitterLogEntry(emitter, "1", "c")
		logger.On("Info", mock.Anything).Return()
		logger.On("Debug", mock.Anything).Return()

		err := r.RunCommand(&commanddomain.Command{
			Id:               "1",
			ProjectId:        "1",
			Name:             "Test",
			Command:          "echo 'a'&& echo 'b'&& echo 'c'",
			WorkingDirectory: "/",
			Position:         0,
		}, []string{}, "")
		assert.NoError(t, err)

		assert.Eventually(t, func() bool {
			return assert.Empty(t, r.GetRunningCommands()) &&
				mock.AssertExpectationsForObjects(t, emitter, logger)
		}, DefaultTimeout, DefaultPollingInterval)

	})
	t.Run("Should log error when executing an invalid command", func(t *testing.T) {
		logger := new(MockLogger)
		emitter := new(MockEventEmitter)

		r := runner.NewDefaultRunner(logger, emitter)

		emitter.On("EmitEvent", event.ProcessStarted, "1").Return()
		emitter.On("EmitEvent", event.ProcessFinished, "1").Return()
		// Not an amazing matcher, but different OSes will have different error messages
		emitter.On("EmitEvent", event.NewLogEntry, mock.Anything).Return()

		logger.On("Info", mock.Anything).Return()
		logger.On("Error", mock.Anything).Return()
		logger.On("Debug", mock.Anything).Return()

		err := r.RunCommand(&commanddomain.Command{
			Id:               "1",
			ProjectId:        "1",
			Name:             "Test",
			Command:          "definitely-not-a-real-command-12345",
			WorkingDirectory: "/",
			Position:         0,
		}, []string{}, "")
		assert.NoError(t, err)

		assert.Eventually(t, func() bool {
			return assert.Empty(t, r.GetRunningCommands()) &&
				mock.AssertExpectationsForObjects(t, emitter, logger)
		}, DefaultTimeout, DefaultPollingInterval)
	})
}

func TestDefaultRunner_StopRunningCommand(t *testing.T) {
	t.Run("Should stop running command", func(t *testing.T) {
		logger := new(MockLogger)
		emitter := new(MockEventEmitter)

		r := runner.NewDefaultRunner(logger, emitter)

		emitter.On("EmitEvent", event.ProcessStarted, "1").Return()
		emitter.On("EmitEvent", event.ProcessFinished, "1").Return()
		emitter.On("EmitEvent", event.NewLogEntry, mock.Anything).Return()
		logger.On("Info", mock.Anything).Return()
		logger.On("Debug", mock.Anything).Return()
		// Depends on OS
		logger.On("Error", mock.Anything).Maybe().Return()

		err := r.RunCommand(&commanddomain.Command{
			Id:               "1",
			ProjectId:        "1",
			Name:             "Test",
			Command:          infiniteCmd(),
			WorkingDirectory: "/",
			Position:         0,
		}, []string{}, "")
		assert.NoError(t, err)

		assert.Eventually(t, func() bool {
			return assert.NotEmpty(t, r.GetRunningCommands())
		}, DefaultTimeout, DefaultPollingInterval)

		err = r.StopRunningCommand("1")
		assert.NoError(t, err)

		assert.Eventually(t, func() bool {
			return assert.Empty(t, r.GetRunningCommands()) &&
				mock.AssertExpectationsForObjects(t, emitter, logger)
		}, DefaultTimeout, DefaultPollingInterval)
	})
	t.Run("Should return error when stopping non-existing command", func(t *testing.T) {
		logger := new(MockLogger)
		emitter := new(MockEventEmitter)

		r := runner.NewDefaultRunner(logger, emitter)

		err := r.StopRunningCommand("non-existing-command")
		assert.Error(t, err)
		assert.Equal(t, "No running command with id: non-existing-command", err.Error())
	})
}

func TestDefaultRunner_StopAllRunningCommands(t *testing.T) {
	t.Run("Should stop all running commands", func(t *testing.T) {
		logger := new(MockLogger)
		emitter := new(MockEventEmitter)

		r := runner.NewDefaultRunner(logger, emitter)

		emitter.On("EmitEvent", event.ProcessStarted, "1").Return()
		emitter.On("EmitEvent", event.ProcessStarted, "2").Return()
		emitter.On("EmitEvent", event.ProcessFinished, "1").Return()
		emitter.On("EmitEvent", event.ProcessFinished, "2").Return()
		emitter.On("EmitEvent", event.NewLogEntry, mock.Anything).Return()
		logger.On("Info", mock.Anything).Return()
		logger.On("Debug", mock.Anything).Return()
		// Depends on OS
		logger.On("Error", mock.Anything).Maybe().Return()

		err := r.RunCommand(&commanddomain.Command{
			Id:               "1",
			ProjectId:        "1",
			Name:             "Test",
			Command:          infiniteCmd(),
			WorkingDirectory: "/",
			Position:         0,
		}, []string{}, "")
		assert.NoError(t, err)

		err = r.RunCommand(&commanddomain.Command{
			Id:               "2",
			ProjectId:        "1",
			Name:             "Test",
			Command:          infiniteCmd(),
			WorkingDirectory: "/",
			Position:         0,
		}, []string{}, "")
		assert.NoError(t, err)

		assert.Eventually(t, func() bool {
			return assert.NotEmpty(t, r.GetRunningCommands())
		}, DefaultTimeout, DefaultPollingInterval)

		errs := r.StopAllRunningCommands()

		assert.Eventually(t, func() bool {
			return assert.Empty(t, errs) &&
				assert.Empty(t, r.GetRunningCommands())
		}, DefaultTimeout, DefaultPollingInterval)
	})
}

func mockEmitterLogEntry(emitter *MockEventEmitter, id string, line string) {
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
