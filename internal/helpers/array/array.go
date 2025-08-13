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

func Map[T, R any](array []T, mapper func(T) R) []R {
	result := make([]R, 0, len(array))
	for _, item := range array {
		result = append(result, mapper(item))
	}
	return result
}
