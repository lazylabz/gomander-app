package array_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gomander/internal/helpers/array"
)

func TestFilter(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	result := array.Filter(input, func(n int) bool { return n%2 == 0 })
	assert.Equal(t, []int{2, 4}, result)
}

func TestMap(t *testing.T) {
	input := []int{1, 2, 3}
	result := array.Map(input, func(n int) string { return string(rune('a' + n - 1)) })
	assert.Equal(t, []string{"a", "b", "c"}, result)
}

func TestIndexOf(t *testing.T) {
	input := []string{"foo", "bar", "baz"}
	assert.Equal(t, 1, array.IndexOf(input, "bar"))
	assert.Equal(t, -1, array.IndexOf(input, "qux"))
}

func TestContains(t *testing.T) {
	input := []int{10, 20, 30}
	assert.True(t, array.Contains(input, 20))
	assert.False(t, array.Contains(input, 40))
}
