package ordering_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"gomander/internal/ordering"
)

type item struct {
	Id       string
	Position int
}

var itemOrder = ordering.NewList(
	func(t item) string { return t.Id },
	func(t item) int { return t.Position },
	func(t *item, position int) { t.Position = position },
)

func TestEnd(t *testing.T) {
	t.Run("Should place the first item at zero", func(t *testing.T) {
		// Arrange
		var empty []item

		// Act & Assert
		assert.Equal(t, 0, itemOrder.End(empty))
	})

	t.Run("Should place a new item behind the ones already there", func(t *testing.T) {
		// Arrange
		existing := []item{{Id: "a", Position: 0}, {Id: "b", Position: 1}}

		// Act & Assert
		assert.Equal(t, 2, itemOrder.End(existing))
	})

	t.Run("Should not hand out a position a gapped list already took", func(t *testing.T) {
		// Arrange
		gapped := []item{{Id: "a", Position: 0}, {Id: "b", Position: 2}}

		// Act & Assert
		assert.Equal(t, 3, itemOrder.End(gapped))
	})
}

func TestRearrange(t *testing.T) {
	t.Run("Should renumber the list to follow the given order", func(t *testing.T) {
		// Arrange
		items := []item{{Id: "a", Position: 0}, {Id: "b", Position: 1}, {Id: "c", Position: 2}}
		saved := recorder{}

		// Act
		err := itemOrder.Rearrange(items, []string{"c", "a", "b"}, saved.save)

		// Assert
		assert.NoError(t, err)
		assertDense(t, saved.merge(items))
		assert.Equal(t, []string{"c", "a", "b"}, ids(inOrder(saved.merge(items))))
	})

	t.Run("Should leave alone the items already in place", func(t *testing.T) {
		// Arrange
		items := []item{{Id: "a", Position: 0}, {Id: "b", Position: 1}, {Id: "c", Position: 2}}
		saved := recorder{}

		// Act
		err := itemOrder.Rearrange(items, []string{"a", "c", "b"}, saved.save)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, []string{"c", "b"}, ids(saved.items))
	})

	t.Run("Should keep the items the given order omits behind the ones it names", func(t *testing.T) {
		// Arrange
		items := []item{{Id: "a", Position: 0}, {Id: "b", Position: 1}, {Id: "c", Position: 2}}
		saved := recorder{}

		// Act
		err := itemOrder.Rearrange(items, []string{"c"}, saved.save)

		// Assert
		assert.NoError(t, err)
		merged := saved.merge(items)
		assertDense(t, merged)
		assert.Equal(t, []string{"c", "a", "b"}, ids(inOrder(merged)))
	})

	t.Run("Should renumber a list that arrives out of position order", func(t *testing.T) {
		// Arrange
		items := []item{{Id: "c", Position: 2}, {Id: "a", Position: 0}, {Id: "b", Position: 1}}
		saved := recorder{}

		// Act
		err := itemOrder.Rearrange(items, []string{"b", "a", "c"}, saved.save)

		// Assert
		assert.NoError(t, err)
		merged := saved.merge(items)
		assertDense(t, merged)
		assert.Equal(t, []string{"b", "a", "c"}, ids(inOrder(merged)))
	})

	t.Run("Should stop at the first item it cannot save", func(t *testing.T) {
		// Arrange
		items := []item{{Id: "a", Position: 0}, {Id: "b", Position: 1}}
		expectedErr := errors.New("storage is gone")

		// Act
		err := itemOrder.Rearrange(items, []string{"b", "a"}, func(*item) error { return expectedErr })

		// Assert
		assert.ErrorIs(t, err, expectedErr)
	})
}

func TestCloseGaps(t *testing.T) {
	t.Run("Should close the gap a removed item left behind", func(t *testing.T) {
		// Arrange
		remaining := []item{{Id: "a", Position: 0}, {Id: "c", Position: 2}}
		saved := recorder{}

		// Act
		err := itemOrder.CloseGaps(remaining, saved.save)

		// Assert
		assert.NoError(t, err)
		merged := saved.merge(remaining)
		assertDense(t, merged)
		assert.Equal(t, []string{"a", "c"}, ids(inOrder(merged)))
	})

	t.Run("Should renumber a list that never started at zero", func(t *testing.T) {
		// Arrange
		remaining := []item{{Id: "a", Position: 3}, {Id: "b", Position: 7}}
		saved := recorder{}

		// Act
		err := itemOrder.CloseGaps(remaining, saved.save)

		// Assert
		assert.NoError(t, err)
		assertDense(t, saved.merge(remaining))
	})

	t.Run("Should touch nothing when the list is already dense", func(t *testing.T) {
		// Arrange
		remaining := []item{{Id: "a", Position: 0}, {Id: "b", Position: 1}}
		saved := recorder{}

		// Act
		err := itemOrder.CloseGaps(remaining, saved.save)

		// Assert
		assert.NoError(t, err)
		assert.Empty(t, saved.items)
	})

	t.Run("Should stop at the first item it cannot save", func(t *testing.T) {
		// Arrange
		remaining := []item{{Id: "a", Position: 1}}
		expectedErr := errors.New("storage is gone")

		// Act
		err := itemOrder.CloseGaps(remaining, func(*item) error { return expectedErr })

		// Assert
		assert.ErrorIs(t, err, expectedErr)
	})
}

type recorder struct {
	items []item
}

func (r *recorder) save(item *item) error {
	r.items = append(r.items, *item)
	return nil
}

// merge answers the list as storage holds it once the saved items landed: the
// ones save was never called for keep the Position they arrived with.
func (r *recorder) merge(original []item) []item {
	merged := make([]item, 0, len(original))
	for _, item := range original {
		for _, saved := range r.items {
			if saved.Id == item.Id {
				item = saved
			}
		}
		merged = append(merged, item)
	}
	return merged
}

func inOrder(items []item) []item {
	ordered := make([]item, len(items))
	for _, item := range items {
		ordered[item.Position] = item
	}
	return ordered
}

func ids(items []item) []string {
	result := make([]string, 0, len(items))
	for _, item := range items {
		result = append(result, item.Id)
	}
	return result
}

// assertDense states the invariant the whole module exists for: Positions run
// 0, 1, 2 … with no gaps and no repeats.
func assertDense(t *testing.T, items []item) {
	t.Helper()

	taken := make(map[int]string, len(items))
	for _, item := range items {
		if item.Position < 0 || item.Position >= len(items) {
			t.Fatalf("position %d of %s is outside a list of %d", item.Position, item.Id, len(items))
		}
		if other, repeated := taken[item.Position]; repeated {
			t.Fatalf("%s and %s both sit at position %d", other, item.Id, item.Position)
		}
		taken[item.Position] = item.Id
	}
}
