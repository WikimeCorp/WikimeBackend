package other

import (
	"net/http"

	"encoding/json"

	"github.com/WikimeCorp/WikimeBackend/dependencies"
	apiErrors "github.com/WikimeCorp/WikimeBackend/restapi/errors"
	"github.com/go-playground/validator/v10"
)

func NotFoundEndpoint(w http.ResponseWriter, req *http.Request) {
	apiErrors.SetErrorInResponce(&apiErrors.ErrNotFound, w, http.StatusNotFound)
}

func CheckRequestJSONData[T any](w http.ResponseWriter, req *http.Request, data *T) error {
	err := json.NewDecoder(req.Body).Decode(data)

	if err != nil {
		apiErrors.SetErrorInResponce(&apiErrors.ErrBadJSONStruct, w, http.StatusBadRequest)
		return err
	}

	err = dependencies.Validate.Struct(data)

	if err != nil {
		tmpErrors := make([]string, 0)
		for _, err := range err.(validator.ValidationErrors) {
			tmpErrors = append(tmpErrors, err.Error())
		}

		err := apiErrors.ErrValidate(tmpErrors)
		apiErrors.SetErrorInResponce(err, w, http.StatusBadRequest)
		return err
	}

	return nil
}
