package utils

import (
	"time"
	"github.com/go-playground/validator/v10"
)

// Validates if a closing date is greater than current date
func ClosingDateValidator(fl validator.FieldLevel) bool {
	// Parse the field value as time.Time
	inputTime := fl.Field().Interface().(time.Time)

	// Get the current time
	currentTime := time.Now().UTC()

	// Compare the input time with the current time
	return inputTime.After(currentTime)
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