package array

func Filter[T any](array []T, predicate func(T) bool) []T {
	result := make([]T, 0, len(array))
	for _, item := range array {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}
