package utility

import "strings"

//
// 1. Null / Empty Checker
//
func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

// IsNull returns true if error exists and prints error code + message
func IsNull(list []any) bool {
	return list == nil
}

// IsError returns true if error exists and prints error code + message
func IsError(err error) bool {
	return err == nil
}

//
// 2. Case-Sensitive Exact Match (Fast)
//
func MultipleStringComparer(value string, list []string) bool {
	return ExistsInSet(value, BuildStringSet(list))
}

func StringComparer(value string, value2 string) bool {
	return strings.Compare(ToLower(value), ToLower(value2)) == 0
}

func StringContains(value string, value2 string) bool {
	return strings.Contains(ToLower(value), ToLower(value2))
}

// 4. Fast Contains (substring check)
func MultipleContains(text string, words []string) bool {
	text = strings.ToLower(text)
	for _, word := range words {
		return strings.Contains(text, word) == true
	}
	return false
}

func BuildStringSet(list []string) map[string]struct{} {
	set := make(map[string]struct{}, len(list))
	for _, v := range list {
		set[strings.ToLower(v)] = struct{}{}
	}
	return set
}

func ExistsInSet(value string, set map[string]struct{}) bool {
	_, exists := set[strings.ToLower(value)]
	return exists
}

func HandleNull[T any](data []T) []T {
	if len(data) == 0 {
		// Create an empty slice of the same type
		return []T{}
	}
	return data
}

func HandleEmptyMap(Data map[string]string) map[string]string {
	if Data == nil {
		return make(map[string]string)
	}
	return Data
}

func Ternary[T any](condition bool, trueVal, falseVal T) T {
	if condition {
		return trueVal
	}
	return falseVal
}

func TernaryFunc[T any](condition bool, trueFunc, falseFunc func() T) T {
	if condition {
		return trueFunc()
	}
	return falseFunc()
}

func RemoveDuplicateStrings(arr []string) []string {
	uniqueMap := make(map[string]bool)
	result := []string{}

	for _, item := range arr {
		if !uniqueMap[item] {
			uniqueMap[item] = true
			result = append(result, item)
		}
	}

	return result
}
