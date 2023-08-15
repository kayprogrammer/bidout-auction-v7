package utils

import (
	"log"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

// Validates if a date has a correct format (ISO8601)
func DateValidator(fl validator.FieldLevel) bool {
	inputTimeString := fl.Field().String()
	iso8601Pattern := `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d{3})?(Z|[\+\-]\d{2}:\d{2})$`

	match, _ := regexp.MatchString(iso8601Pattern, inputTimeString)
	return match
}

// Validates if a closing date is greater than current date
func ClosingDateValidator(fl validator.FieldLevel) bool {
	inputTimeString := fl.Field().Interface().(string)
	// Parse the field value as time.Time
	parsedTime := TimeParser(inputTimeString)
	log.Println(parsedTime)

	// Get the current time
	currentTime := time.Now().UTC()

	// Compare the input time with the current time
	return parsedTime.After(currentTime)
}

// Validates if a file type is accepted
func FileTypeValidator(fl validator.FieldLevel) bool {
	fileType := fl.Field().Interface().(string)
	fileTypeFound := false
	for key := range ImageExtensions {
		if key == fileType {
			fileTypeFound = true
			break
		}
	}
	return fileTypeFound
}
