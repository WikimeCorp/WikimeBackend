package db

import (
	"errors"
	"fmt"

	. "github.com/WikimeCorp/WikimeBackend/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AddAnime(anime *Anime) (ansAnimeId AnimeID, err error) {
	animeID, err := createAnimeDoc(anime.Title, anime.OriginTitle, anime.Author)
	if err != nil {
		return 0, err
	}

	err = EditAnime(anime)
	if err != nil {
		return
	}

	err = addToAdded(anime.Author, animeID)
	if err != nil {
		return
	}

	err = createRatingDoc(animeID)
	if err != nil {
		return
	}

	err = createCommentsDoc(animeID)
	if err != nil {
		return
	}

	return animeID, nil
}

// Rate is rate function on behalf of the user
//
// It must be guaranteed that the user with the `id` exists
func Rate(animeID AnimeID, userID UserID, rate AnimeRating) error {
	if rate < 1 || rate > 5 {
		return fmt.Errorf("rate must be 1<=rate<=5, not %d", rate)
	}

	oneRateProj := bson.M{"Rated": bson.M{"$elemMatch": bson.M{"AnimeId": animeID}}, "Watched": 0}

	user := &User{}

	err := usersCollection.FindOneAndUpdate(ctx,
		bson.M{"_id": userID, "Rated.AnimeId": animeID},
		bson.M{"$set": bson.M{"Rated.$.Rate": rate}},
		options.FindOneAndUpdate().SetProjection(oneRateProj),
	).Decode(user)

	if err == nil {
		err := ChangeRating(animeID, user.Rated[0].Rate, rate)
		return err
	} else if errors.Is(err, mongo.ErrNoDocuments) {
		err := addToRated(userID, animeID, rate)
		if err != nil {
			return err
		}
		err = addRate(animeID, rate)
		return err
	} else if err != nil {
		return err
	}
	return nil
}
