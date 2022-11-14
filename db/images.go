package db

import (
	"github.com/WikimeCorp/WikimeBackend/types"
	"github.com/WikimeCorp/WikimeBackend/types/myerrors"
	"go.mongodb.org/mongo-driver/bson"
)

// AddImageToAnime adding image to anime
func AddImageToAnime(animeID types.AnimeID, filePath string) error {
	res, err := animeCollection.UpdateOne(ctx, bson.M{"_id": animeID}, bson.M{"$push": bson.M{"Images": filePath}})
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return &myerrors.ErrAnimeNotFound{AnimeID: animeID}
	}

	return nil
}

// SetPoster ...
func SetPoster(animeID types.AnimeID, filePath string) error {
	res, err := animeCollection.UpdateOne(ctx, bson.M{"_id": animeID}, bson.M{"$set": bson.M{"Poster": filePath}})
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return &myerrors.ErrAnimeNotFound{AnimeID: animeID}
	}

	return nil
}
