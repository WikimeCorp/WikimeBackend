package user

import "github.com/WikimeCorp/WikimeBackend/types"

// UserModel is contract for getting user
type UserModel struct {
	UserID    types.UserID
	Nickname  string
	Role      string
	Favorites []types.AnimeID
	Watched   []types.AnimeID
	Added     []types.AnimeID
	Rated     []struct {
		ID   types.AnimeID
		Rate types.AnimeRating
	}
}
