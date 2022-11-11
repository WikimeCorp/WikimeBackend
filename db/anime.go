package db

import (
	"errors"
	"fmt"
	"time"

	. "github.com/WikimeCorp/WikimeBackend/types"
	"github.com/WikimeCorp/WikimeBackend/types/dbtypes"
	inerr "github.com/WikimeCorp/WikimeBackend/types/myerrors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func createAnimeDoc(title string, originTitle string, author UserID) (*dbtypes.Anime, error) {
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

func EditAnime(animeObjPtr *dbtypes.Anime) error {
	err := animeCollection.FindOneAndReplace(ctx, bson.M{"_id": animeObjPtr.ID}, &animeObjPtr).Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return &inerr.ErrAnimeNotFound{animeObjPtr.ID}
		}
		return err
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

func GetAnimeIDsSortedByRating(genres []string) ([]AnimeID, error) {
	matchAgg := bson.D{
		{Key: "$match",
			Value: bson.D{},
		},
	}
	fmt.Println("Genres: ", genres)
	if len(genres) != 0 {
		genresAsDObjects := make([]bson.D, len(genres))

		for idx, genre := range genres {
			genresAsDObjects[idx] = bson.D{{"Genres", genre}}
		}

		genresAObject := genresAsDObjects

		matchAgg = bson.D{
			{Key: "$match",
				Value: bson.D{
					{Key: "$and",
						Value: genresAObject,
					},
				},
			},
		}
	}
	fmt.Println("aggr ", matchAgg)
	pipeline := []bson.D{
		matchAgg,
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "Rating"},
					{"localField", "_id"},
					{"foreignField", "_id"},
					{"as", "RateTmp"},
				},
			},
		},
		bson.D{
			{"$set",
				bson.D{
					{"Average",
						bson.D{
							{"$getField",
								bson.D{
									{"field", "Average"},
									{"input", bson.D{{"$first", "$RateTmp"}}},
								},
							},
						},
					},
				},
			},
		},
		bson.D{{"$unset", "RateTmp"}},
		bson.D{{"$sort", bson.D{{"Average", -1}}}},
	}

	cursor, err := animeCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

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
