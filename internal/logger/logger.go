package logger

// LogSink is where the desktop shell's log lines end up.
type LogSink interface {
	LogInfo(message string)
	LogDebug(message string)
	LogError(message string)
}

type DefaultLogger struct {
	sink LogSink
}

func NewDefaultLogger(sink LogSink) *DefaultLogger {
	return &DefaultLogger{
		sink: sink,
	}
}

func (l *DefaultLogger) Info(message string) {
	l.sink.LogInfo(message)
}

func (l *DefaultLogger) Debug(message string) {
	l.sink.LogDebug(message)
}

func (l *DefaultLogger) Error(message string) {
	l.sink.LogError(message)
}
