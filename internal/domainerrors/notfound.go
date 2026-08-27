// Package domainerrors holds the errors every entity domain reports, so a
// caller can tell them apart without knowing which repository raised them.
package domainerrors

import (
	"errors"
	"fmt"
)

// ErrNotFound reports that a lookup which required a row did not find one:
// absence is this error, and anything else is a storage failure.
var ErrNotFound = errors.New("not found")

func NotFound(entity string, id string) error {
	return fmt.Errorf("%s %q: %w", entity, id, ErrNotFound)
}
