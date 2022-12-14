package db

import (
	"errors"
	"fmt"
	"log"

	. "github.com/WikimeCorp/WikimeBackend/types"
	"github.com/WikimeCorp/WikimeBackend/types/dbtypes"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AddAnime creates anime correctly
func AddAnime(anime *dbtypes.Anime) (ansAnimeID AnimeID, err error) {
	newAnime, err := createAnimeDoc(anime.Title, anime.OriginTitle, anime.Author, *anime.Poster)
	if err != nil {
		return 0, err
	}

	anime.ID = newAnime.ID
	anime.DateAdded = newAnime.DateAdded
	anime.Rating = newAnime.Rating

	err = ReplaceAnime(anime)
	if err != nil {
		return
	}

	// err = createRatingDoc(newAnime.ID)
	// if err != nil {
	// 	return
	// }

	err = createCommentsDoc(newAnime.ID)
	if err != nil {
		return
	}

	return newAnime.ID, nil
}

// Rate is rate function on behalf of the user
//
// It must be guaranteed that the user with the `id` exists
func Rate(animeID AnimeID, userID UserID, rate AnimeRating) error {
	if rate < 1 || rate > 5 {
		return fmt.Errorf("rate must be 1<=rate<=5, not %d", rate)
	}

	oneRateProj := bson.M{"Rated": bson.M{"$elemMatch": bson.M{"AnimeId": animeID}}, "Watched": 0}

	user := &dbtypes.User{}

	err := usersCollection.FindOneAndUpdate(ctx,
		bson.M{"_id": userID, "Rated.AnimeId": animeID},
		bson.M{"$set": bson.M{"Rated.$.Rate": rate}},
		options.FindOneAndUpdate().SetProjection(oneRateProj),
	).Decode(user)

	log.Println("SAFE METHOD ERR: ", err)

	if err == nil {
		err := ChangeRating(animeID, user.Rated[0].Rate, rate)
		log.Println("3: ", err)
		return err
	} else if errors.Is(err, mongo.ErrNoDocuments) {
		err := addToRated(userID, animeID, rate)
		log.Println("1: ", err)
		if err != nil {
			return err
		}
		err = addRate(animeID, rate)
		log.Println("2: ", err)
		return err
	} else if err != nil {
		return err
	}
	return nil
}

// AddAnimeToFavorites adding anime to the list of 'favorites' of a user with the ID 'User ID`
//
// It must be guaranteed that the user with the `id` exists
func AddAnimeToFavorites(animeID AnimeID, userID UserID) error {
	wasModified, err := addToFavorites(userID, animeID)

	if err != nil {
		return err
	}

	if wasModified == true {
		err = IncFavorite(animeID, 1)
		return err
	}
	return nil
}

// AddAnimeToWatched adding anime to the list of 'Watched' of a user with the ID 'User ID`
//
// It must be guaranteed that the user with the `id` exists
func AddAnimeToWatched(animeID AnimeID, userID UserID) error {
	wasModified, err := addToWatched(userID, animeID)
	if err != nil {
		return err
	}

	if wasModified == true {
		err = IncWatched(animeID, 1)
		return err
	}

	return nil
}

// DeleteAnimeFromFavorites removes anime from the list of `favorites` of a user with the id `userID`
//
// It must be guaranteed that the user with the `id` exists
func DeleteAnimeFromFavorites(animeID AnimeID, userID UserID) error {
	wasModified, err := deleteFromFavorites(userID, animeID)
	if err != nil {
		return err
	}

	if wasModified == true {
		err = IncFavorite(animeID, -1)
		return err
	}

	return nil
}

// DeleteAnimeFromWatched removes anime from the list of `Watched` of a user with the id `userID`
//
// It must be guaranteed that the user with the `id` exists
func DeleteAnimeFromWatched(animeID AnimeID, userID UserID) error {
	wasModified, err := deleteFromWatched(userID, animeID)
	if err != nil {
		return err
	}

	if wasModified == true {
		err = IncWatched(animeID, -1)
		return err
	}

	return nil
}
