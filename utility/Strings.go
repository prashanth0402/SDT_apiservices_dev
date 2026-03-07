package utility

import "strings"

func ToLower(s string) string {
	return strings.ToLower(s)
}

func ToUpper(s string) string {
	return strings.ToUpper(s)
}

// ----------------------------------------------------------------
// Function to CapitalizeText capitalizes the first letter of each word in a string.
// ----------------------------------------------------------------
func CapitalizeText(input string) string {
	words := strings.Fields(input) // Split the input into words
	var capitalizedWords []string

	for _, word := range words {
		// Capitalize the first letter of the word
		capitalizedWord := ToUpper(string(word[0])) + ToLower(word[1:])
		capitalizedWords = append(capitalizedWords, capitalizedWord)
	}

	// Join the capitalized words back into a string
	return strings.Join(capitalizedWords, " ")
}
