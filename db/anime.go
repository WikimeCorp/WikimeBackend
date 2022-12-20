package db

import (
	"errors"
	"time"

	dbrequests "github.com/WikimeCorp/WikimeBackend/db/db_requests"
	. "github.com/WikimeCorp/WikimeBackend/types"
	"github.com/WikimeCorp/WikimeBackend/types/dbtypes"
	inerr "github.com/WikimeCorp/WikimeBackend/types/myerrors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func createAnimeDoc(title string, originTitle string, author UserID, poster string) (*dbtypes.Anime, error) {
	user, err := CheckUser(author)
	if err != nil {
		return nil, err
	}

	if !user {
		return nil, &inerr.ErrUserNotFound{author}
	}

	animeID, err := getNextID[AnimeID]("AnimeID")
	if err != nil {
		return nil, err
	}

	anime := dbtypes.Anime{
		ID:          animeID,
		Title:       title,
		OriginTitle: originTitle,
		DateAdded:   time.Now(),
		Rating:      &dbtypes.Rating{},
		Poster:      &poster,
	}
	_, err = animeCollection.InsertOne(ctx, anime)

	return &anime, err
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

func ReplaceAnime(animeObjPtr *dbtypes.Anime) error {
	err := animeCollection.FindOneAndReplace(ctx, bson.M{"_id": animeObjPtr.ID}, &animeObjPtr).Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return &inerr.ErrAnimeNotFound{animeObjPtr.ID}
		}
		return err
	}

	return nil
}

func EditTextFieldsAnime(animeObj *dbtypes.Anime) error {
	animeID := animeObj.ID

	req := bson.M{}
	if animeObj.Title != "" {
		req["Title"] = animeObj.Title
	}
	if animeObj.OriginTitle != "" {
		req["OriginTitle"] = animeObj.OriginTitle
	}
	if animeObj.Genres != nil && len(animeObj.Genres) != 0 {
		req["Genres"] = animeObj.Genres
	}
	if animeObj.Description != "" {
		req["Description"] = animeObj.Description
	}
	if animeObj.Director != "" {
		req["Director"] = animeObj.Director
	}
	if true { // Need add valid check
		req["ReleaseDate"] = animeObj.ReleaseDate
	}

	ans, err := animeCollection.UpdateOne(ctx, bson.M{"_id": animeID}, bson.M{"$set": req})
	if err != nil {
		return err
	}

	if ans.MatchedCount == 0 {
		return &inerr.ErrAnimeNotFound{AnimeID: animeID}
	}

	return nil
}

func GetAnimeByID(id AnimeID) (*dbtypes.Anime, error) {
	ans := &dbtypes.Anime{}
	err := animeCollection.FindOne(ctx, bson.M{"_id": id}).Decode(ans)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, &inerr.ErrAnimeNotFound{id}
		}
		return nil, err
	}

	return ans, nil
}

func GetAnimes(genres []string, sortBy string, order int) (ansList []*dbtypes.Anime, err error) {
	opts := options.Find()
	opts.SetCursorType(options.NonTailable)
	opts.SetSort(bson.D{{sortBy, order}})

	cur, err := animeCollection.Find(ctx, bson.M{"Genres": genres}, opts)
	if err != nil {
		return ansList, err
	}

	for cur.Next(ctx) {
		curAnine := &dbtypes.Anime{}
		err = cur.Decode(curAnine)
		if err != nil {
			return ansList, err
		}
		ansList = append(ansList, curAnine)
	}

	return ansList, err
}

func DecodeAnimesToIDListFromCursor(cursor *mongo.Cursor) ([]AnimeID, error) {
	results := make([]AnimeID, 0)

	for cursor.Next(ctx) {

		elem := struct {
			ID AnimeID `bson:"_id"`
		}{}

		err := cursor.Decode(&elem)
		if err != nil {
			return nil, err
		}
		results = append(results, elem.ID)
	}

	return results, nil
}

func DecodeAnimesFromCursor(cursor *mongo.Cursor) ([]*dbtypes.Anime, error) {
	results := make([]*dbtypes.Anime, 0)

	for cursor.Next(ctx) {

		elem := dbtypes.Anime{}

		err := cursor.Decode(&elem)
		if err != nil {
			return nil, err
		}
		results = append(results, &elem)
	}

	return results, nil
}

func GetAnimeIDsSortedByRating(genres []string, order int8, skip, limit int) ([]AnimeID, error) {
	cursor, err := animeCollection.Aggregate(ctx, dbrequests.GetAnimesSortedByRatingWithGenres(genres, order, limit, skip))
	if err != nil {
		return nil, err
	}
	return DecodeAnimesToIDListFromCursor(cursor)
}

func GetAnimeIDsSortedByFavorites(genres []string, order int8, skip, limit int) ([]AnimeID, error) {
	cursor, err := animeCollection.Aggregate(ctx, dbrequests.GetAnimeSortedByFavoritesWithGenres(genres, order, limit, skip))
	if err != nil {
		return nil, err
	}

	return DecodeAnimesToIDListFromCursor(cursor)
}

func GetAnimeIDsSortedByAddingDate(genres []string, order int8, skip, limit int) ([]AnimeID, error) {
	cursor, err := animeCollection.Aggregate(ctx, dbrequests.GetAnimeSortedByAddingDateWithGenres(genres, order, limit, skip))
	if err != nil {
		return nil, err
	}

	return DecodeAnimesToIDListFromCursor(cursor)
}

func GetAnimeIDsSortedByReleaseDate(genres []string, order int8, skip, limit int) ([]AnimeID, error) {
	cursor, err := animeCollection.Aggregate(ctx, dbrequests.GetAnimeSortedByReleaseDateWithGenres(genres, order, limit, skip))
	if err != nil {
		return nil, err
	}

	return DecodeAnimesToIDListFromCursor(cursor)
}

func GetAnimeAddedByUser(userID UserID) ([]AnimeID, error) {
	cursor, err := animeCollection.Find(ctx, bson.M{"Author": userID})
	if err != nil {
		return nil, err
	}
	return DecodeAnimesToIDListFromCursor(cursor)
}

func SearchAnime(text string) ([]AnimeID, error) {
	opt := options.Find().SetSort(bson.D{{"Rating.Average", -1}})
	opt = opt.SetLimit(15)
	cur, err := animeCollection.Find(ctx, bson.M{"$text": bson.M{"$search": text}}, opt)
	if err != nil {
		return nil, err
	}
	return DecodeAnimesToIDListFromCursor(cur)
}
