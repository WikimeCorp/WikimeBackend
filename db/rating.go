package db

import (
	"fmt"
	"log"

	dbrequests "github.com/WikimeCorp/WikimeBackend/db/db_requests"
	. "github.com/WikimeCorp/WikimeBackend/types"
	"github.com/WikimeCorp/WikimeBackend/types/dbtypes"
	"github.com/WikimeCorp/WikimeBackend/types/myerrors"
	inerr "github.com/WikimeCorp/WikimeBackend/types/myerrors"
	"go.mongodb.org/mongo-driver/bson"
)

// GetRating gets a rating by anime id
func GetRating(id AnimeID) (*dbtypes.Rating, error) {
	ans := &dbtypes.Anime{}

	err := animeCollection.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(ans)
	if err != nil {
		return nil, err
	}

	return ans.Rating, nil
}

var _rateNames = map[AnimeRating]string{
	5: "Five",
	4: "Four",
	3: "Three",
	2: "Two",
	1: "One",
}

// ChangeRating reduces the `from` field and increases the `to` field and recalculates the average
func ChangeRating(id AnimeID, from AnimeRating, to AnimeRating) error {
	fromRateName, ok := _rateNames[from]
	if !ok {
		return fmt.Errorf("invalid 'rate' argument: %d. Must be 1, 2, 3, 4 or 5", from)
	}
	toRateName, ok := _rateNames[to]
	if !ok {
		return fmt.Errorf("invalid 'rate' argument: %d. Must be 1, 2, 3, 4 or 5", to)
	}
	if from == to {
		return nil
	}

	ans, err := animeCollection.UpdateByID(ctx, id, dbrequests.ChangeRating(fromRateName, toRateName))
	log.Println("ChangeRating err: ", err)
	if err != nil {
		return err
	}
	if ans.MatchedCount == 0 {
		return &myerrors.ErrAnimeNotFound{id}
	}

	return nil
}

// IncFavorite increases or decreases the `InFavorites` field for anime with `id`.
//
// To reduce it, you need to pass a negative number to `inc`.
func IncFavorite(id AnimeID, inc int) error {
	ans, err := animeCollection.UpdateByID(ctx, id, bson.M{
		"$inc": bson.D{
			{"Rating.InFavorites", inc},
		},
	})
	if err != nil {
		return err
	}
	if ans.MatchedCount == 0 {
		return &myerrors.ErrAnimeNotFound{id}
	}
	return nil
}

// IncWatched increases or decreases the `Watched` field for anime with `id`.
//
// To reduce it, you need to pass a negative number to `inc`.
func IncWatched(id AnimeID, inc int) error {
	ans, err := animeCollection.UpdateByID(ctx, id, bson.M{
		"$inc": bson.D{
			{"Rating.Watched", inc},
		},
	})
	if err != nil {
		return err
	}
	if ans.MatchedCount == 0 {
		return &myerrors.ErrAnimeNotFound{id}
	}
	return nil
}

func addRate(id AnimeID, rate AnimeRating) (err error) {
	rateName, ok := _rateNames[rate]
	if !ok {
		return fmt.Errorf("invalid 'rate' argument: %d. Must be 1, 2, 3, 4 or 5", rate)
	}

	ans, err := animeCollection.UpdateByID(ctx, id, dbrequests.AddRate(rateName))
	log.Println("addRate err: ", err)
	if err != nil {
		return err
	}
	if ans.MatchedCount == 0 {
		return &myerrors.ErrAnimeNotFound{id}
	}

	return nil
}

// NEED DELETE
func SetAverage(anime AnimeID, average float64) error {
	ans, err := animeCollection.UpdateByID(ctx, anime, bson.M{"$set": bson.M{"Rating.Average": average}})

	if err != nil {
		return err
	}
	if ans.MatchedCount == 0 {
		return &inerr.ErrAnimeNotFound{anime}
	}

	return err
}
