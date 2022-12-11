package user

import (
	"github.com/WikimeCorp/WikimeBackend/types"
	"github.com/WikimeCorp/WikimeBackend/types/dbtypes"
)

// UserModel is contract for getting user
type UserModel struct {
	UserID    types.UserID    `json:"userId"`
	Nickname  string          `json:"nickname"`
	Role      types.Role      `json:"role"`
	Favorites []types.AnimeID `json:"favorites"`
	Watched   []types.AnimeID `json:"watched"`
	Rated     []RatedAnime    `json:"rated"`
	Avatar    string          `json:"avatar"`
}

type RatedAnime struct {
	ID   types.AnimeID     `json:"id"`
	Rate types.AnimeRating `json:"Rate"`
}

func UserModelFromDBUser(user *dbtypes.User) *UserModel {
	ans := UserModel{
		UserID:    user.ID,
		Nickname:  user.Nickname,
		Role:      types.Role(user.Role),
		Favorites: user.Favorites,
		Watched:   user.Watched,
		Rated:     []RatedAnime{},
		Avatar:    user.Avatar,
	}
	for _, el := range user.Rated {
		ans.Rated = append(ans.Rated, RatedAnime{ID: el.ID, Rate: el.Rate})
	}

	return &ans
}
