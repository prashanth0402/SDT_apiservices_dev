package common

import (
	"errors"
	"log"
	"strings"

	"github.com/BurntSushi/toml"
)

// --------------------------------------------------------------------
// function reads the constants from the config.toml file
// --------------------------------------------------------------------
func ReadTomlConfig(filename string) interface{} {
	var f interface{}
	if _, err := toml.DecodeFile(filename, &f); err != nil {
		log.Println(err)
	}
	return f
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

// ----------------------------------------------------------------
// Function to CapitalizeText capitalizes the first letter of each word in a string.
// ----------------------------------------------------------------
func CapitalizeText(input string) string {
	words := strings.Fields(input) // Split the input into words
	var capitalizedWords []string

	for _, word := range words {
		// Capitalize the first letter of the word
		capitalizedWord := strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		capitalizedWords = append(capitalizedWords, capitalizedWord)
	}

	// Join the capitalized words back into a string
	return strings.Join(capitalizedWords, " ")
}

//----------------------------------------------------------------
// Creating custom error
// ----------------------------------------------------------------

func CustomError(pErrorMsg string) error {
	err := errors.New(pErrorMsg)
	return err
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
