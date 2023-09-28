package validator

import (
	"github.com/go-playground/validator/v10"
	"strings"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

type File struct {
	Name string `validate:"required,max=255,fsname"`
	Size int64  `validate:"lte=5000000"`
	Type string `validate:"filetype"`
}

// InitializeValidator инициализирует валидатор и регистрирует пользовательские функции валидации.
func InitializeValidator() *CustomValidator {
	val := validator.New()
	val.RegisterValidation("fsname", ValidateFSName)
	return &CustomValidator{validator: val}
}

// ValidateFSName проверяет, что имя файла или папки не содержит недопустимых символов.
func ValidateFSName(fl validator.FieldLevel) bool {
	name := fl.Field().String()
	invalidChars := []string{"/", "\\", ":", "*", "?", "<", ">", "|", "\""}
	for _, ch := range invalidChars {
		if strings.Contains(name, ch) {
			return false
		}
	}
	return true
}
