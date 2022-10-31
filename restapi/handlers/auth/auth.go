package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/WikimeCorp/WikimeBackend/applogic/auth"
	"github.com/WikimeCorp/WikimeBackend/dependencies"
	apiErrors "github.com/WikimeCorp/WikimeBackend/restapi/errors"
	"github.com/WikimeCorp/WikimeBackend/types/myerrors"
	"github.com/go-playground/validator/v10"
)

func oAuthVk(w http.ResponseWriter, req *http.Request) {
	form := AuthRequest{}
	decodeErr := json.NewDecoder(req.Body).Decode(&form)
	if decodeErr != nil {
		fmt.Println(decodeErr)
		apiErrors.SetErrorInResponce(&apiErrors.ErrBadJSONStruct, w, http.StatusBadRequest)
		return
	}

	err := dependencies.Validate.Struct(form)

	if err != nil {
		tmpErrors := make([]string, 0)
		for _, err := range err.(validator.ValidationErrors) {
			tmpErrors = append(tmpErrors, err.Error())
		}

		err := apiErrors.ErrValidate(tmpErrors)
		apiErrors.SetErrorInResponce(err, w, http.StatusBadRequest)
		return
	}

	tokenStr, err := auth.VkAuth(form.AuthToken)

	if err != nil {
		if errors.Is(err, myerrors.ErrOuterTokenIsWrong) {
			apiErrors.SetErrorInResponce(&apiErrors.ErrBadOuterToken, w, http.StatusBadRequest)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	ans := AuthResponse{AuthToken: tokenStr}
	ansJSON, _ := json.Marshal(ans)
	w.Write(ansJSON)

}

// OAuthVkHandler return OAuthVk handler
func OAuthVkHandler() func(w http.ResponseWriter, req *http.Request) {
	return oAuthVk
}
