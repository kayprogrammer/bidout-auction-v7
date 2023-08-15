package utils

import (
	r "crypto/rand"
	"encoding/base64"
	"math/rand"
	"time"
	"reflect"
	"github.com/shopspring/decimal"
	"log"
)

func GetRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	randomStr := make([]byte, length)
	for i := range randomStr {
		randomStr[i] = charset[rand.Intn(len(charset))]
	}
	return string(randomStr)
}

// Generates a random integer with a specified number of digits
func GetRandomInt(size int) int {
	if size <= 0 {
		return 0
	}

	// Calculate the min and max possible values for the specified size
	min := intPow(10, size-1)
	max := intPow(10, size) - 1

	// Initialize the random number generator
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate a random integer within the range [min, max]
	return rand.Intn(max-min+1) + min
}

// intPow calculates the power of base^exponent for integers
func intPow(base, exponent int) int {
	result := 1
	for i := 0; i < exponent; i++ {
		result *= base
	}
	return result
}

// generateRandomPassword generates a random password for the test database
func GenerateRandomPassword() string {
	const passwordLength = 16 // You can adjust the password length as needed
	rb := make([]byte, passwordLength)
	r.Read(rb)
	return base64.URLEncoding.EncodeToString(rb)
}

// Check if keys exist in map
func KeysExistInMap(keys []string, myMap map[string]interface{}) bool {
    for _, key := range keys {
        if _, ok := myMap[key]; !ok {
            return false
        }
    }
    return true
}

func AssignFields(src interface{}, dest interface{}) {
	srcValue := reflect.ValueOf(src)
	destValue := reflect.ValueOf(dest).Elem()
	for i := 0; i < srcValue.NumField(); i++ {
		srcField := srcValue.Field(i)

		if !srcField.IsNil() && srcField.Kind() == reflect.Ptr {
			destField := destValue.FieldByName(srcValue.Type().Field(i).Name)
			log.Println(destField.Type())
			log.Println(srcField.Elem())
			
			if destField.IsValid() {
				if srcField.Elem().Type().ConvertibleTo(destField.Type()) {
					destField.Set(srcField.Elem().Convert(destField.Type()))
				} else if destField.Type() == reflect.TypeOf(decimal.Decimal{}) && srcField.Elem().Kind() == reflect.Float64 {
					decimalValue := decimal.NewFromFloat(srcField.Elem().Float())
					destField.Set(reflect.ValueOf(decimalValue))
				} else if destField.Type() == reflect.TypeOf(time.Time{}) && srcField.Elem().Kind() == reflect.String {
					layout := "2006-01-02T15:04:05.000Z" // Change this to match your input date format
					dateString := srcField.Elem().String()
					parsedTime, _ := time.Parse(layout, dateString)
					destField.Set(reflect.ValueOf(parsedTime))
				} else {
					destField.Set(srcField.Elem())
				}
			}
		}
	}
}

func TimeParser(timeStr string) time.Time {
	parsedTime, err := time.Parse("2006-01-02T15:04:05.999Z07:00", timeStr)
	log.Println(err)
	// if err != nil {
	// 	parsedTime, err := time.Parse("2006-01-02T15:04:05.000+01:00", timeStr)
	// }
	return parsedTime.UTC()
}

func DecimalParser(num float64) decimal.Decimal {
	parsedDecimal := decimal.NewFromFloat(num)
	return parsedDecimal.Round(2)
}