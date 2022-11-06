package dependencies

import (
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
)

// Validate is validator for request struct
var Validate *validator.Validate

// Decoder is validator for request struct
var Decoder *schema.Decoder

var CtxUserID = "user_id"

func init() {
	Validate = validator.New()
	Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		return fld.Tag.Get("validatorName")
	})

	Decoder = schema.NewDecoder()
	Decoder.SetAliasTag("json")
}
