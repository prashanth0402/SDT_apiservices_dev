package utility

import "strings"

//
// 1. Null / Empty Checker
//
func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

func IsNullorEmptylist(list []any) bool {
	return list == nil
}

//
// 2. Case-Sensitive Exact Match (Fast)
//
func UseStringComparer(value string, list []string) bool {
	return ExistsInSet(value, BuildStringSet(list))
}

// 4. Fast Contains (substring check)
func ContainsWord(text string, words []string) bool {
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
