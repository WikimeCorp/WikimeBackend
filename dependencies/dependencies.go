package dependencies

import "github.com/go-playground/validator/v10"

// Validate is validator for request struct
var Validate *validator.Validate

var CtxUserID = "user_id"

func init() {
	Validate = validator.New()
}
