package helper

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func init() {
	validate.RegisterTagNameFunc(func(f reflect.StructField) string {
		name := strings.SplitN(f.Tag.Get("json"), ",", 2)[0]
		if name == "" || name == "-" {
			return f.Name
		}
		return name
	})
}

func ValidateStruct(s interface{}) map[string][]string {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}
	ve, ok := err.(validator.ValidationErrors)
	if !ok {
		return map[string][]string{"request": {err.Error()}}
	}
	errs := map[string][]string{}
	for _, fe := range ve {
		errs[fe.Field()] = append(errs[fe.Field()], validationMessage(fe))
	}
	return errs
}

func validationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "wajib diisi"
	case "email":
		return "Format email tidak valid"
	case "min":
		return fmt.Sprintf("minimal %s karakter", fe.Param())
	case "len":
		return fmt.Sprintf("harus %s karakter", fe.Param())
	case "numeric":
		return "harus berupa angka"
	case "gte":
		return fmt.Sprintf("minimal %s", fe.Param())
	case "lte":
		return fmt.Sprintf("maksimal %s", fe.Param())
	default:
		return "tidak valid"
	}
}