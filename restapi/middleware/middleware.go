package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/WikimeCorp/WikimeBackend/applogic/auth"
	"github.com/WikimeCorp/WikimeBackend/dependencies"
	apiErrors "github.com/WikimeCorp/WikimeBackend/restapi/errors"
	"github.com/WikimeCorp/WikimeBackend/types/myerrors"
)

// SetJSONHeader is middleware, that add "Context-Type" header as "application/json"
func SetJSONHeader(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h.ServeHTTP(w, r)
	})
}

func PrintRequestURL(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL, r.Body)
		h.ServeHTTP(w, r)
	})
}

// NeedAuthentication check authorization header and check JWT token valid
func NeedAuthentication(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			apiErrors.SetErrorInResponce(&apiErrors.ErrJWTTokenNotFound, w, http.StatusUnauthorized)
			return
		}

		payload, err := auth.CheckJWTValid(authHeader)

		if err != nil {
			if errors.Is(err, myerrors.ErrJWTTokenTimeout) {
				apiErrors.SetErrorInResponce(&apiErrors.ErrJWTTokenTimeout, w, http.StatusUnauthorized)
				return
			}

			if errors.Is(err, myerrors.ErrJWTTokenInvalidSignature) {
				apiErrors.SetErrorInResponce(&apiErrors.ErrJWTTokenInvalidSignature, w, http.StatusUnauthorized)
				return
			}
		}

		ctx := context.WithValue(r.Context(), dependencies.CtxUserID, payload.UserID)
		r = r.WithContext(ctx)

		h.ServeHTTP(w, r)
	})
}
