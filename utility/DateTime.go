package utility

import (
	"errors"
	"time"
)

// Supported input date formats
var supportedLayouts = []string{
	time.RFC3339,          // 2006-01-02T15:04:05Z07:00
	"2006/01/02",          // 2026/03/09
	"02 Jan 2006",         // 09 Mar 2026
	"02 January 2006",     // 09 March 2026
	"02-01-2006 15:04:05", // 09-03-2026 10:20:30
	"02/01/2006 15:04:05", // 09/03/2026 10:20:30
}

var TimeFormats = map[string]string{
	"RFC3339":             time.RFC3339,
	"mm/dd/yyyy":          "01/02/2006",
	"dd/mm/yyyy":          "02/01/2006",
	"yyyy-mm-dd":          "2006-01-02",
	"dd-mm-yyyy":          "02-01-2006",
	"datetime":            "2006-01-02 15:04:05",
	"yyyy-mm-dd HH:MM:SS": "2006-01-02 15:04:05",
	"yyyy-mm-dd HH:MM":    "2006-01-02 15:04",
	"yyyy-mm-dd HH":       "2006-01-02 15",
}

// ConvertDateFormat converts any supported format into desired format
func ConvertDateFormat(inputDate string, outputLayout string) (string, error) {

	var parsedTime time.Time
	var err error

	for _, layout := range supportedLayouts {
		parsedTime, err = time.Parse(layout, inputDate)
		if err == nil {
			return parsedTime.Format(outputLayout), nil
		}
	}

	return "", errors.New("unsupported date format")
}
