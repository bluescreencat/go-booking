package validator

import (
	"regexp"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
)

type (
	ErrorFieldData struct {
		FailedField string      `json:"failedField"`
		Tag         string      `json:"tag"`
		Value       interface{} `json:"value"`
	}
)

const REGEXP_PASSWORD = `(?=.*\d)(?=.*[a-z])(?=.*[A-Z])(?=.*\W)`

// const REGEXP_PASSWORD = `(?=.*[a-z])(?=.*[A-Z])(?=.*\W)`

type CustomValidator struct {
	validate *validator.Validate
}

var customValidator *CustomValidator
var once sync.Once

func (customValidator *CustomValidator) Validate(data any) []ErrorFieldData {
	validationErrors := []ErrorFieldData{}
	errs := customValidator.validate.Struct(data)
	if errs == nil {
		return nil
	}
	for _, err := range errs.(validator.ValidationErrors) {
		var element ErrorFieldData
		element.FailedField = err.Field()
		element.Tag = err.Tag()
		validationErrors = append(validationErrors, element)
	}
	return validationErrors
}

func validatePasswordFormat(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	haveNumeric := regexp.MustCompile(`\d`).MatchString(password)
	haveUpperCase := regexp.MustCompile(`[a-z]`).MatchString(password)
	haveLowerCase := regexp.MustCompile(`[A-Z]`).MatchString(password)
	haveSpecialCharacters := strings.ContainsAny(password, "!@#$^&*()_+")
	isMatched := haveNumeric && haveUpperCase && haveLowerCase && haveSpecialCharacters
	return isMatched
}

func GetCustomValidatorInstance() *CustomValidator {
	if customValidator == nil {
		once.Do(func() {
			validate := validator.New()
			validate.RegisterValidation("password-format", validatePasswordFormat)
			customValidator = &CustomValidator{validate: validate}
		})
	}
	return customValidator
}
