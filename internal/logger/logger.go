package logger

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Logger struct {
	ctx context.Context
}

func NewLogger(ctx context.Context) *Logger {
	return &Logger{
		ctx: ctx,
	}
}

func (l *Logger) Info(message string) {
	runtime.LogInfo(l.ctx, message)
}

func (l *Logger) Debug(message string) {
	runtime.LogDebug(l.ctx, message)
}

func (l *Logger) Error(message string) {
	runtime.LogError(l.ctx, message)
}
