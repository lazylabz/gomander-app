package logger

import (
	"context"
)

// LogSink is where the desktop shell's log lines end up.
type LogSink interface {
	LogInfo(ctx context.Context, message string)
	LogDebug(ctx context.Context, message string)
	LogError(ctx context.Context, message string)
}

type DefaultLogger struct {
	ctx  context.Context
	sink LogSink
}

func NewDefaultLogger(ctx context.Context, sink LogSink) *DefaultLogger {
	return &DefaultLogger{
		ctx:  ctx,
		sink: sink,
	}
}

func (l *DefaultLogger) Info(message string) {
	l.sink.LogInfo(l.ctx, message)
}

func (l *DefaultLogger) Debug(message string) {
	l.sink.LogDebug(l.ctx, message)
}

func (l *DefaultLogger) Error(message string) {
	l.sink.LogError(l.ctx, message)
}
