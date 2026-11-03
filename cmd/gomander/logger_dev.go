//go:build !production

package main

import "log/slog"

// Leaving the logger unset gives development builds Wails' own console logger.
func productionLogger() *slog.Logger {
	return nil
}
