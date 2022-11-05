package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type ErrBaseEndpointError struct {
	Message   string `json:"error_message"`
	ErrorCode int    `json:"error"`
}

func (e *ErrBaseEndpointError) Error() string {
	return fmt.Sprintf("Error: %s. Error code: %d", e.Message, e.ErrorCode)
}

const (
	badJSONStruct = iota + 1
	badOuterToken
	badValidate
	userNotFound
	jwtTokenNotFound
	jwtTokenTimeout
	jwtTokenInvalidSignature
	notFound
	animeNotFound
	internalServerError
)

var ErrBadJSONStruct = ErrBaseEndpointError{Message: "Bad json", ErrorCode: badJSONStruct}
var ErrBadOuterToken = ErrBaseEndpointError{Message: "Token of outer service is wrong", ErrorCode: badOuterToken}
var ErrUserNotFound = ErrBaseEndpointError{Message: "User not found", ErrorCode: userNotFound}
var ErrJWTTokenNotFound = ErrBaseEndpointError{Message: "JWT token not found, check 'Authorization' header", ErrorCode: jwtTokenNotFound}
var ErrJWTTokenTimeout = ErrBaseEndpointError{Message: "JWT token timeout", ErrorCode: jwtTokenTimeout}
var ErrJWTTokenInvalidSignature = ErrBaseEndpointError{Message: "JWT token has invalid signature", ErrorCode: jwtTokenInvalidSignature}
var ErrNotFound = ErrBaseEndpointError{Message: "Page not found", ErrorCode: notFound}
var ErrAnimeNotFound = ErrBaseEndpointError{Message: "Anime not found", ErrorCode: animeNotFound}
var ErrInternalServerError = ErrBaseEndpointError{Message: "Internal server error", ErrorCode: internalServerError}

func SetErrorInResponce(err *ErrBaseEndpointError, w http.ResponseWriter, status int) {
	w.WriteHeader(status)
	jsonAns, _ := json.Marshal(err)

	w.Write(jsonAns)
}

func ErrValidate(errs []string) *ErrBaseEndpointError {
	err := ErrBaseEndpointError{ErrorCode: badValidate, Message: "Bad validation: " + strings.Join(errs, "; ")}
	return &err
}
