package dbrequests

import (
	"go.mongodb.org/mongo-driver/bson"
)

// GetAnimesSortedByRatingWithGenres generate mongodb request for *see func name*
func GetAnimesSortedByRatingWithGenres(genres []string) bson.A {

	genresArray := bson.A{}
	if len(genres) != 0 {
		for _, genre := range genres {
			genresArray = append(genresArray, bson.D{{"Genres", genre}})
		}
	} else {
		genresArray = append(genresArray, bson.D{})
	}

	ans := bson.A{
		bson.D{
			{Key: "$match",
				Value: bson.D{
					{Key: "$and",
						Value: genresArray,
					},
				},
			},
		},
		bson.D{
			{Key: "$sort", Value: bson.D{{Key: "Rating.Average", Value: -1}}},
		},
	}

	return ans
}
