package validatorx

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	validate.RegisterValidation("strong_password", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()
		return len(password) >= 6 && strings.ContainsAny(password, "0123456789")
	})
}

// devuelve UN SOLO mensaje
func ValidateStruct(s interface{}) (string, bool) {
	err := validate.Struct(s)
	if err == nil {
		return "", false
	}

	for _, err := range err.(validator.ValidationErrors) {
		return buildMessage(err), true
	}

	return "invalid request", true
}

// mensajes personalizados
func buildMessage(fe validator.FieldError) string {
	field := toSnakeCase(fe.Field())

	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s es requerido", field)
	case "email":
		return "email inválido"
	case "min":
		return fmt.Sprintf("%s debe tener al menos %s caracteres", field, fe.Param())
	case "max":
		return fmt.Sprintf("%s debe tener máximo %s caracteres", field, fe.Param())
	case "strong_password":
		return "la contraseña debe tener mínimo 6 caracteres y al menos un número"
	case "gt":
		return fmt.Sprintf("%s debe ser mayor que %s", field, fe.Param())
	case "gte":
		return fmt.Sprintf("%s debe ser mayor o igual a %s", field, fe.Param())
	case "oneof":
		return fmt.Sprintf("%s es inválido.", field)
	default:
		return fmt.Sprintf("%s es inválido", field)
	}
}

// helper
func toSnakeCase(str string) string {
	var result []rune
	for i, r := range str {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}
