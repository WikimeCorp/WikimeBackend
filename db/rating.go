package db

import (
	"errors"
	"fmt"

	. "github.com/WikimeCorp/WikimeBackend/types"
	"github.com/WikimeCorp/WikimeBackend/types/dbtypes"
	inerr "github.com/WikimeCorp/WikimeBackend/types/myerrors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetRating gets a rating by anime id
func GetRating(id AnimeID) (*dbtypes.Rating, error) {
	ans := &dbtypes.Rating{}

	err := ratingCollection.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(ans)
	if err != nil {
		return nil, err
	}

	return ans, nil
}

func createRatingDoc(id AnimeID) error {
	_, err := ratingCollection.InsertOne(ctx, dbtypes.Rating{ID: id})
	return err
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

	rating, err := GetRating(id)
	if err != nil {
		return err
	}

	if from == to {
		return nil
	}

	rateName, ok := _rateNames[from]
	if !ok {
		return fmt.Errorf("invalid 'from' argument: %d. Must be 1, 2, 3, 4 or 5", from)
	}
	var fromRate uint32
	switch from {
	case 5:
		fromRate = rating.Five
	case 4:
		fromRate = rating.Four
	case 3:
		fromRate = rating.Three
	case 2:
		fromRate = rating.Two
	case 1:
		fromRate = rating.One
	}
	if fromRate <= 0 {
		return fmt.Errorf("the value of `from` is zero")
	}

	fromBson := bson.E{rateName, -1}

	rateName, ok = _rateNames[to]
	if !ok {
		return fmt.Errorf("invalid 'to' argument: %d. Must be 1, 2, 3, 4 or 5", from)
	}
	toBson := bson.E{rateName, 1}

	average := float64(rating.Five*5+rating.Four*4+rating.Three*3+rating.Two*2+rating.One-uint32(from)+uint32(to)) /
		float64(rating.Five+rating.Four+rating.Three+rating.Two+rating.One)

	_, err = ratingCollection.UpdateByID(ctx, id, bson.M{
		"$inc": bson.D{
			fromBson,
			toBson,
		},
		"$set": bson.D{
			{"Average", average},
		},
	})
	return err
}

// IncFavorite increases or decreases the `InFavorites` field for anime with `id`.
//
// To reduce it, you need to pass a negative number to `inc`.
func IncFavorite(id AnimeID, inc int) error {
	_, err := ratingCollection.UpdateByID(ctx, id, bson.M{
		"$inc": bson.D{
			{"InFavorites", inc},
		},
	})
	return err
}

func addRate(id AnimeID, rate AnimeRating) (err error) {
	rateName, ok := _rateNames[rate]
	if !ok {
		return fmt.Errorf("invalid 'rate' argument: %d. Must be 1, 2, 3, 4 or 5", rate)
	}

	rating := &dbtypes.Rating{}
	err = ratingCollection.FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{"$inc": bson.M{rateName: 1}}).Decode(rating)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return &inerr.ErrAnimeNotFound{id}
		}
		return err
	}

	average := float64(rating.Five*5+rating.Four*4+rating.Three*3+rating.Two*2+rating.One+uint32(rate)) /
		float64(rating.Five+rating.Four+rating.Three+rating.Two+rating.One+1)

	_, err = ratingCollection.UpdateByID(ctx, id, bson.M{
		"$set": bson.M{"Average": average},
	})

	return err
}

// NEED DELETE
func SetAverage(anime AnimeID, average float64) error {
	_, err := ratingCollection.UpdateByID(ctx, anime, bson.M{"$set": bson.M{"Average": average}})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return &inerr.ErrAnimeNotFound{anime}
		}
		return err
	}

	return err
}
