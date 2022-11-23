package user

import "github.com/WikimeCorp/WikimeBackend/types"

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
