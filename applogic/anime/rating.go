package anime

import (
	"github.com/WikimeCorp/WikimeBackend/db"
	"github.com/WikimeCorp/WikimeBackend/types"
)

func SetRating(userID types.UserID, animeID types.AnimeID, rating types.AnimeRating) error {
	err := db.Rate(animeID, userID, rating)
	return err
}
