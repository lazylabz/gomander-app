// Package domain holds the shapes the Wails frontend sends and receives, and
// the mapping between them and the entities behind them. Every controller
// signature names a type from here, so what the UI sees is stated in one place
// and renaming a field of an entity cannot quietly reshape it.
//
// The package name is load-bearing rather than descriptive: Wails names the
// generated TypeScript namespace after the Go package of the types a bound
// method mentions, so calling this package anything else would move every
// model in models.ts. The frontend pays for that move once, in the change that
// takes the serialization tags off the entities.
package domain

import (
	"gomander/internal/helpers/array"
)

// Optional maps a value the frontend may not get at all — no Project open, a
// cancelled file dialog — so that absence stays absent instead of arriving as
// a zero value.
func Optional[T, R any](value *T, from func(T) R) *R {
	if value == nil {
		return nil
	}
	mapped := from(*value)
	return &mapped
}

// mapSlice keeps a nil slice nil, because the two differ on the wire: a nil
// slice serializes as null and an empty one as [].
func mapSlice[T, R any](values []T, from func(T) R) []R {
	if values == nil {
		return nil
	}
	return array.Map(values, from)
}
