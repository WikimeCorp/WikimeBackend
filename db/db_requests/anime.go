package dbrequests

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

// GetAnimesSortedByRatingWithGenres generate mongodb request for *see func name*
func GetAnimesSortedByRatingWithGenres(genres []string) bson.A {

	genresArray := bson.A{bson.D{}}
	if len(genres) != 0 {
		//genresD := make([]bson.D, len(genres))
		for _, genre := range genres {
			//genresD[idx] = bson.D{{"Genres", genre}}
			genresArray = append(genresArray, bson.D{{"Genres", genre}})
		}
		//genresArray = genresD
		fmt.Println(genresArray)
		fmt.Println(bson.A{
			bson.D{{"Genres", "Комедия"}},
			bson.D{{"Genres", "Детектив"}},
		})

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
