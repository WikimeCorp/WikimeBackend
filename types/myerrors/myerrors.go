package myerrors

import (
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
