package runner

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	commanddomain "gomander/internal/command/domain"
	"gomander/internal/event"
	"testing"
	"time"
)

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Info(message string) {
	m.Called(message)
	return
}

func (m *MockLogger) Debug(message string) {
	m.Called(message)
	return
}

func (m *MockLogger) Error(message string) {
	m.Called(message)
	return
}

type MockEventEmitter struct {
	mock.Mock
}

func (m *MockEventEmitter) EmitEvent(event event.Event, payload interface{}) {
	m.Called(event, payload)
	return
}

func TestDefaultRunner_RunCommand(t *testing.T) {
	t.Run("Should run command with success and emit events for each line", func(t *testing.T) {
		logger := new(MockLogger)
		emitter := new(MockEventEmitter)

		r := NewDefaulRunner(logger, emitter)

		emitter.On("EmitEvent", event.ProcessStarted, "1").Return()
		emitter.On("EmitEvent", event.ProcessFinished, "1").Return()
		emitter.On("EmitEvent", event.NewLogEntry, map[string]string{
			"id":   "1",
			"line": "'a'",
		}).Return()
		emitter.On("EmitEvent", event.NewLogEntry, map[string]string{
			"id":   "1",
			"line": "'b'",
		}).Return()
		emitter.On("EmitEvent", event.NewLogEntry, map[string]string{
			"id":   "1",
			"line": "'c'",
		}).Return()
		logger.On("Info", mock.Anything).Return()
		logger.On("Debug", mock.Anything).Return()

		err := r.RunCommand(commanddomain.Command{
			Id:               "1",
			ProjectId:        "1",
			Name:             "Test",
			Command:          "echo 'a'&& echo 'b'&& echo 'c'",
			WorkingDirectory: "/",
			Position:         0,
		}, []string{}, "")

		time.Sleep(100 * time.Millisecond)

		assert.Empty(t, r.runningCommands)

		assert.NoError(t, err)
	})
}
