package logger

import (
	"context"

	"gomander/internal/facade"
)

type Logger interface {
	Info(message string)
	Debug(message string)
	Error(message string)
}

type DefaultLogger struct {
	ctx     context.Context
	runtime facade.RuntimeFacade
}

func NewDefaultLogger(ctx context.Context, runtime facade.RuntimeFacade) *DefaultLogger {
	return &DefaultLogger{
		ctx:     ctx,
		runtime: runtime,
	}
}

func (l *DefaultLogger) Info(message string) {
	l.runtime.LogInfo(l.ctx, message)
}

func (l *DefaultLogger) Debug(message string) {
	l.runtime.LogDebug(l.ctx, message)
}

func (l *DefaultLogger) Error(message string) {
	l.runtime.LogError(l.ctx, message)
}
