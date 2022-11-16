package user

import "github.com/WikimeCorp/WikimeBackend/types"

// UserModel is contract for getting user
type UserModel struct {
	UserID    types.UserID    `json:"userId"`
	Nickname  string          `json:"nickname"`
	Role      string          `json:"role"`
	Favorites []types.AnimeID `json:"favorites"`
	Watched   []types.AnimeID `json:"watched"`
	Rated     []struct {
		ID   types.AnimeID
		Rate types.AnimeRating
	} `json:"rated"`
}
