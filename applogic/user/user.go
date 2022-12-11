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
	ans := UserModelFromDBUser(user)
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

func DeleteFromFavorites(userID types.UserID, animeID types.AnimeID) error {
	err := db.DeleteAnimeFromFavorites(animeID, userID)
	return err
}

func DeleteFromWatched(userID types.UserID, animeID types.AnimeID) error {
	err := db.DeleteAnimeFromWatched(animeID, userID)
	return err
}

func ChangeRole(userID types.UserID, role types.Role) error {
	return db.ChangeRole(userID, role)
}

func GetUsersByRole(role types.Role) ([]UserModel, error) {
	ans, err := db.GetUsersByRole(role)
	if err != nil {
		return nil, err
	}

	res := make([]UserModel, 0)
	for _, el := range ans {
		user := UserModelFromDBUser(el)
		res = append(res, *user)
	}

	return res, nil
}
