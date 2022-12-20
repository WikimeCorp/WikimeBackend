package user

import (
	"github.com/WikimeCorp/WikimeBackend/applogic/user"
	"github.com/WikimeCorp/WikimeBackend/types"
)

type ChangeNicknameRequest struct {
	Nickname string `json:"nickname" validate:"required"`
}

type AddToFavoritesRequest struct {
	AnimeID *types.AnimeID `json:"animeId" validate:"required"`
}

type AddToWatchedRequest struct {
	AddToFavoritesRequest
}

type DeleteFromWatchedRequest struct {
	AddToFavoritesRequest
}

type DeleteFromFavoritesRequest struct {
	AddToFavoritesRequest
}

type ShortUser struct {
	ID       types.UserID `json:"id"`
	Avatar   string       `json:"avatar"`
	Nickname string       `json:"nickname"`
}

type User struct {
	user.UserModel
	Added []types.AnimeID `json:"added"`
}
