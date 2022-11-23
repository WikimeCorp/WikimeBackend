package user

import (
	"github.com/WikimeCorp/WikimeBackend/db"
	"github.com/WikimeCorp/WikimeBackend/types"
)

// GetUser return user from db
func GetUser(userID types.UserID) (*UserModel, error) {
	user, err := db.GetUser(userID)
	if err != nil {
		return nil, err
	}
	ans := &UserModel{
		UserID:    user.ID,
		Nickname:  user.Nickname,
		Role:      user.Role,
		Favorites: user.Favorites,
		Watched:   user.Watched,
		Rated: []struct {
			ID   types.AnimeID
			Rate types.AnimeRating
		}(user.Rated),
	}
	return ans, err
}

func SetNickname(userID types.UserID, newNickname string) error {
	err := db.EditNickname(userID, newNickname)
	return err
}

func AddToFavorites(userID types.UserID, animeID types.AnimeID) error {
	err := db.AddAnimeToFavorites(animeID, userID)
	return err
}

func AddToWatched(userID types.UserID, animeID types.AnimeID) error {
	err := db.AddAnimeToWatched(animeID, userID)
	return err
}
