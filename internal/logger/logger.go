package logger

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Logger interface {
	Info(message string)
	Debug(message string)
	Error(message string)
}

type DefaultLogger struct {
	ctx context.Context
}

func NewDefaultLogger(ctx context.Context) *DefaultLogger {
	return &DefaultLogger{
		ctx: ctx,
	}
}

func (l *DefaultLogger) Info(message string) {
	runtime.LogInfo(l.ctx, message)
}

func (l *DefaultLogger) Debug(message string) {
	runtime.LogDebug(l.ctx, message)
}

func (l *DefaultLogger) Error(message string) {
	runtime.LogError(l.ctx, message)
}
