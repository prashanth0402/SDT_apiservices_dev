package utility

import "strings"

// filtered := Filter(nums, func(n int) bool {
// 	return n > 3
// })
func Filter[T any](items []T, condition func(T) bool) []T {
	var result []T
	for _, item := range items {
		if condition(item) {
			result = append(result, item)
		}
	}
	return result
}

// value, found := Find(nums, func(n int) bool {
// 	return n == 4
// })

func Find[T any](items []T, condition func(T) bool) (T, bool) {
	for _, item := range items {
		if condition(item) {
			return item, true
		}
	}
	var zero T
	return zero, false
}

func Some[T any](items []T, condition func(T) bool) bool {
	for _, item := range items {
		if condition(item) {
			return true
		}
	}
	return false
}

// sum := Reduce(nums, 0, func(acc, n int) int {
// 	return acc + n
// })
func Reduce[T any, R any](items []T, initial R, reducer func(R, T) R) R {
	result := initial
	for _, item := range items {
		result = reducer(result, item)
	}
	return result
}

func FlatMap[T any, R any](items []T, mapper func(T) []R) []R {
	var result []R
	for _, item := range items {
		result = append(result, mapper(item)...)
	}
	return result
}

func ConvertArrayToString(array []string) string {
	var builder strings.Builder

	for i, str := range array {
		builder.WriteString("'")
		builder.WriteString(str)
		builder.WriteString("'")
		// Add comma if it's not the last element
		if i != len(array)-1 {
			builder.WriteString(",")
		}
	}
	if array == nil {
		return "''"
	}

	return builder.String()
}
