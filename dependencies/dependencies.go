package dependencies

import "github.com/go-playground/validator/v10"

// Validate is validator for request struct
var Validate *validator.Validate

func init() {
	Validate = validator.New()
}
