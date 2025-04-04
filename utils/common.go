package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// WriteJson writes JSON data to the response
func WriteJson(w http.ResponseWriter, statusCode int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}

// parseTime converts a time string to a time.Time object
func ParseTime(timeStr string) (time.Time, error) {
	formats := []string{
		time.RFC3339,                 // "2006-01-02T15:04:05Z07:00"
		"2006-01-02 15:04:05",        // Standard SQL format
		"2006-01-02T15:04:05",        // ISO without timezone
		"2006-01-02 15:04:05.000000", // SQLite format with microseconds
	}

	var parsedTime time.Time
	var parserErr error
	for _, format := range formats {
		parsedTime, parserErr = time.Parse(format, timeStr)
		if parserErr == nil {
			break
		}
	}

	if parserErr != nil {
		log.Printf("Error parsing time: %v", parserErr)
		return time.Time{}, parserErr
	}

	return parsedTime, nil

}
