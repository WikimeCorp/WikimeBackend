package db

import (
	"errors"
	"time"

	. "github.com/WikimeCorp/WikimeBackend/types"
	inerr "github.com/WikimeCorp/WikimeBackend/types/myerrors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func createAnimeDoc(title string, originTitle string, author UserID) (AnimeID, error) {
	user, err := CheckUser(author)
	if err != nil {
		return 0, err
	}

	if !user {
		return 0, &inerr.ErrUserNotFound{author}
	}

	animeID, err := getNextID[AnimeID]("AnimeId")
	if err != nil {
		return 0, err
	}

	anime := Anime{
		ID:          animeID,
		Title:       title,
		OriginTitle: originTitle,
		DateAdded:   time.Now(),
	}
	_, err = animeCollection.InsertOne(ctx, anime)

	return animeID, err
}

func CheckAnime(id AnimeID) (bool, error) {
	err := animeCollection.FindOne(ctx, bson.M{"_id": id}).Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func EditAnime(animeObjPtr *Anime) error {
	err := animeCollection.FindOneAndReplace(ctx, bson.M{"_id": animeObjPtr.ID}, &animeObjPtr).Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return &inerr.ErrAnimeNotFound{animeObjPtr.ID}
		}
		return err
	}

	return nil
}

func GetAnime(id AnimeID) (*Anime, error) {
	ans := &Anime{}
	err := animeCollection.FindOne(ctx, bson.M{"_id": id}).Decode(ans)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, &inerr.ErrAnimeNotFound{id}
		}
		return nil, err
	}

	return ans, nil
}
