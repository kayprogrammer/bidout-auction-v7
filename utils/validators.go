package utils

import (
	"fmt"
	"reflect"
	"strings"
    // "net/mail"

    "github.com/go-playground/locales/en"
    "github.com/go-playground/universal-translator"
    "github.com/go-playground/validator/v10"
)

var (
    customValidator *validator.Validate
    translator      ut.Translator
)

// CustomValidationError is a custom error type that implements the error interface
type CustomValidationError struct {
    Status   	string			`json:"status"`
    Message 	string			`json:"message"`
    Data 		interface{}		`json:"data"`

}

func (e *CustomValidationError) Error() string {
    return e.Message
}

// Email Validator
// func customEmailValidator(fl validator.FieldLevel) bool {
//     email := fl.Field().String()
//     _, err := mail.ParseAddress(email)
//     return err == nil
// }

// Initialize the custom validator and translator
func init() {
    customValidator = validator.New()
    en := en.New()
    uni := ut.New(en, en)
    translator, _ = uni.GetTranslator("en")

	customValidator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

    // Register Validators
    // customValidator.RegisterValidation("email", customEmailValidator)
}

// Register translations
func registerTranslations(param string) {
    // Register custom error messages for each validation tag
    registerTranslation := func (tag string, translation string, translator ut.Translator) {
        customValidator.RegisterTranslation(tag, translator, func(ut ut.Translator) error {
            return ut.Add(tag, translation, true)
        }, func(ut ut.Translator, fe validator.FieldError) string {
            t, _ := ut.T(tag, fe.Field())
            return t
        })
    }

    registerTranslation("required", "This field is required.", translator)
    minErrMsg := fmt.Sprintf("%s characters min", param)
    registerTranslation("min", minErrMsg, translator)
    maxErrMsg := fmt.Sprintf("%s characters max", param)
    registerTranslation("max", maxErrMsg, translator)
}

// CustomValidator is a custom validator that uses "github.com/go-playground/validator/v10"
type CustomValidator struct{}

// Validate performs the validation of the given struct
func (cv *CustomValidator) Validate(i interface{}) error {
    if err := customValidator.Struct(i); err != nil {
        err := err.(validator.ValidationErrors)
        return cv.translateValidationErrors(err)
    }
    return nil
}

// translateValidationErrors translates the validation errors to custom errors
func (cv *CustomValidator) translateValidationErrors(errs validator.ValidationErrors) error {
    errData := make(map[string]string)
	for _, err := range errs {
        registerTranslations(err.Param())
		errData[err.Field()] = err.Translate(translator)
    }
    return &CustomValidationError{
		Status:   "failure",
		Message: "Invalid Entry",
		Data: errData,
	}
}

// New creates a new instance of CustomValidator
func Validator() *CustomValidator {
    return &CustomValidator{}
}
