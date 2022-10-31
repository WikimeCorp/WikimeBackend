package errors

import (
	"encoding/json"
	"net/http"
	"strings"
)

type ErrBaseEndpointError struct {
	Message   string `json:"error_message"`
	ErrorCode int    `json:"error"`
}

const (
	badJSONStruct = iota + 1
	badOuterToken
	badValidate
	userNotFound
)

var ErrBadJSONStruct = ErrBaseEndpointError{Message: "Bad json", ErrorCode: badJSONStruct}
var ErrBadOuterToken = ErrBaseEndpointError{Message: "Token of outer service is wrong", ErrorCode: badOuterToken}
var ErrUserNotFound = ErrBaseEndpointError{Message: "Token of outer service is wrong", ErrorCode: userNotFound}

func SetErrorInResponce(err *ErrBaseEndpointError, w http.ResponseWriter, status int) {
	w.WriteHeader(status)
	jsonAns, _ := json.Marshal(err)

	w.Write(jsonAns)
}

func ErrValidate(errs []string) *ErrBaseEndpointError {
	err := ErrBaseEndpointError{ErrorCode: badValidate, Message: "Bad validation: " + strings.Join(errs, "; ")}
	return &err
}
