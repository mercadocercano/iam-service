package validator

import (
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var slugRegexp = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-]*$`)

// RegisterCustomValidators registra las validaciones custom usadas por los
// request DTOs del servicio (por ejemplo `slug`). Es idempotente: registrar
// dos veces el mismo tag en el mismo engine de gin no genera error.
func RegisterCustomValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("slug", func(fl validator.FieldLevel) bool {
			return slugRegexp.MatchString(fl.Field().String())
		})
	}
}
