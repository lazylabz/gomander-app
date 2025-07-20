package main

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

func (l *Logger) info(message string) {
	runtime.LogInfo(l.ctx, message)
}

func (l *Logger) debug(message string) {
	runtime.LogDebug(l.ctx, message)
}

func (l *Logger) error(message string) {
	runtime.LogError(l.ctx, message)
}
