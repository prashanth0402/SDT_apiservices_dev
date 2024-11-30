package common

import (
	"database/sql"
	"errors"
	"log"
	"strings"
	"time"

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

// --------------------------------------------------------------------
// function convert the time and date format to customized format
// --------------------------------------------------------------------
func ChangeTimeFormat(pCustomizeLayout string, pInput string) (string, error) {
	log.Println("ChangeTimeFormat (+)")
	var lFormattedValue string

	Layout := ""
	length := len(pInput)
	if length == 19 {
		Layout = "02-01-2006 15:04:05"
	} else if length == 5 {
		Layout = "15:04"
	} else if length == 8 {
		Layout = "15:04:05"
	} else {
		Layout = "02-01-2006 15:04"
	}
	lTimevalue, lErr1 := time.Parse(Layout, pInput)
	if lErr1 != nil {
		log.Println("Error in Parse Timing:", lErr1)
		return lFormattedValue, lErr1
	} else {
		lFormattedValue = lTimevalue.Format(pCustomizeLayout)
	}

	log.Println("ChangeTimeFormat (-)")
	return lFormattedValue, nil
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

func GetLookUpDetails(pDB *sql.DB, pHeaderCode string) (map[string]string, error) {
	lId := ""
	lCode := ""
	lookupDetails := make(map[string]string)
	lsqlString := `SELECT coalesce (xld.id,'') , coalesce (xld.Code,'')
	FROM xx_lookup_details xld,xx_lookup_header xlh 
	where xlh.Code  = '` + pHeaderCode + `' and xlh.id = xld.headerid  ;`

	// log.Println("lsqlString", lsqlString)
	lRows, lErr := pDB.Query(lsqlString)
	if lErr != nil {
		log.Println("CGLD01", lErr)
		return lookupDetails, lErr
	} else {
		for lRows.Next() {
			lErr := lRows.Scan(&lId, &lCode)
			if lErr != nil {
				log.Println("CGLD02", lErr)
				return lookupDetails, lErr
			} else {
				lookupDetails[lId] = lCode
			}
		}
	}
	// log.Println("lookupDetails", lookupDetails)
	return lookupDetails, nil
}

func HandleNull[T any](data []T) []T {
	if len(data) == 0 {
		// Create an empty slice of the same type
		return []T{}
	}
	return data
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
