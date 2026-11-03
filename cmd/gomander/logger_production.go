//go:build production

package main

import (
	"log/slog"
	"os"
)

// Wails v3 discards every log line in production builds; v2 still wrote the
// errors, which are what explain a refused quit.
func productionLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
}
