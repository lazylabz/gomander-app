package domain

import (
	"errors"
	"fmt"
)

// DefaultLogsFontSize is xterm.js's own default, the size the logs had before
// it became configurable.
const (
	DefaultLogsFontSize = 15
	MinLogsFontSize     = 8
	MaxLogsFontSize     = 32
)

var ErrLogsFontSizeOutOfRange = errors.New("logs font size out of range")

func validateLogsFontSize(size int) error {
	if size < MinLogsFontSize || size > MaxLogsFontSize {
		return fmt.Errorf("%d is not between %d and %d: %w", size, MinLogsFontSize, MaxLogsFontSize, ErrLogsFontSizeOutOfRange)
	}

	return nil
}
