// Package ordering owns where a thing sits in a list: a Command among its
// Project's Commands, a Command Group among its Project's Groups. Positions run
// 0, 1, 2 … with no gaps, and this is the only place that decides them - an
// operation says it wants to append, to rearrange or to close the gap a removal
// left, and the sequence comes out dense either way.
package ordering

import (
	"cmp"
	"slices"
)

// List is one kind of ordered thing - a Command among its Project's Commands, a
// Command Group among its Project's Groups - described by how to read its id and
// how to read and write its Position. Each domain declares one, and every
// operation on that domain's order goes through it.
type List[T any] struct {
	idOf       func(T) string
	positionOf func(T) int
	placeAt    func(*T, int)
}

func NewList[T any](idOf func(T) string, positionOf func(T) int, placeAt func(*T, int)) List[T] {
	return List[T]{
		idOf:       idOf,
		positionOf: positionOf,
		placeAt:    placeAt,
	}
}

// End answers the Position an item appended to items takes: behind every one
// already there, so a list that reached storage with a gap in it still cannot
// hand out the same Position twice.
func (l List[T]) End(items []T) int {
	end := 0

	for _, item := range items {
		if position := l.positionOf(item); position >= end {
			end = position + 1
		}
	}

	return end
}

// Rearrange renumbers items so they follow orderedIds, saving the ones that
// moved. Items orderedIds does not name keep their order behind the ones it
// does, so a partial order still leaves the list dense.
func (l List[T]) Rearrange(items []T, orderedIds []string, save func(*T) error) error {
	ordered := l.inPositionOrder(items)

	slices.SortStableFunc(ordered, func(a, b T) int {
		return cmp.Compare(l.rankIn(orderedIds, a), l.rankIn(orderedIds, b))
	})

	return l.renumber(ordered, save)
}

// CloseGaps renumbers items into a dense sequence without changing their order,
// saving the ones that moved. It is what a removal leaves to do.
func (l List[T]) CloseGaps(items []T, save func(*T) error) error {
	return l.renumber(l.inPositionOrder(items), save)
}

func (l List[T]) rankIn(orderedIds []string, item T) int {
	if rank := slices.Index(orderedIds, l.idOf(item)); rank >= 0 {
		return rank
	}

	return len(orderedIds)
}

func (l List[T]) inPositionOrder(items []T) []T {
	ordered := slices.Clone(items)

	slices.SortStableFunc(ordered, func(a, b T) int {
		return cmp.Compare(l.positionOf(a), l.positionOf(b))
	})

	return ordered
}

func (l List[T]) renumber(ordered []T, save func(*T) error) error {
	for position := range ordered {
		if l.positionOf(ordered[position]) == position {
			continue
		}

		l.placeAt(&ordered[position], position)

		if err := save(&ordered[position]); err != nil {
			return err
		}
	}

	return nil
}
