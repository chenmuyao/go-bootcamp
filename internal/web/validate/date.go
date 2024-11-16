package validate

import (
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var SupportedValidators = map[string]validator.Func{
	"date": DateValidate,
}

func UseValidators(validators ...string) {
	if e, ok := binding.Validator.Engine().(*validator.Validate); ok {
		for _, v := range validators {
			if f, ok := SupportedValidators[v]; ok {
				e.RegisterValidation(v, f)
			}
		}
	}
}

func DateValidate(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()
	if len(dateStr) == 0 {
		return true
	}
	_, err := time.Parse("2006-01-02", dateStr)
	return err == nil
}
