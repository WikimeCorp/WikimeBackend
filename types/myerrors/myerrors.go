package myerrors

import (
	"errors"
	"fmt"

	. "github.com/WikimeCorp/WikimeBackend/types"
)

type ErrUserNotFound struct {
	UserID UserID
}

func (e *ErrUserNotFound) Error() string {
	return fmt.Sprintf("user with id %d not found", e.UserID)
}

type ErrAnimeNotFound struct {
	AnimeID AnimeID
}

func (e *ErrAnimeNotFound) Error() string {
	return fmt.Sprintf("anime with id %d not found", e.AnimeID)
}

var ErrOuterTokenIsWrong = errors.New("token of outer service is wrong")

var ErrNoDocumentsFromDB = errors.New("db returned 0 documents")

var ErrIncorrectPartsNumberOfJWTToken = errors.New("The number of parts of the jwt token is incorrect. There should be 3 (there should be two dots '.' in the jwt token)")

var ErrJWTTokenTimeout = errors.New("JWT token timeout")
var ErrJWTTokenInvalidSignature = errors.New("JWT token has invalid signature")
