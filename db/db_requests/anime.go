package dbrequests

import (
	"go.mongodb.org/mongo-driver/bson"
)

func genresToBsonArray(genres []string) bson.A {
	genresArray := bson.A{}
	if len(genres) != 0 {
		for _, genre := range genres {
			genresArray = append(genresArray, bson.D{{Key: "Genres", Value: genre}})
		}
	} else {
		genresArray = append(genresArray, bson.D{})
	}
	return genresArray
}

func generateSortingWithGenres(field string, genres []string, order int8) bson.A {
	genresArray := genresToBsonArray(genres)

	return bson.A{
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
			{Key: "$sort", Value: bson.D{{Key: field, Value: order}}},
		},
	}
}

// GetAnimesSortedByRatingWithGenres generate mongodb request for *see func name*
func GetAnimesSortedByRatingWithGenres(genres []string, order int8) bson.A {
	return generateSortingWithGenres("Rating.Average", genres, order)
}

// GetAnimeSortedByAddingDateWithGenres generate mongodb request for *see func name*
func GetAnimeSortedByAddingDateWithGenres(genres []string, order int8) bson.A {
	return generateSortingWithGenres("DateAdded", genres, order)
}

// GetAnimeSortedByReleaseDateWithGenres generate mongodb request for *see func name*
func GetAnimeSortedByReleaseDateWithGenres(genres []string, order int8) bson.A {
	return generateSortingWithGenres("ReleaseDate", genres, order)
}

// GetAnimeSortedByFavoritesWithGenres generate mongodb request for *see func name*
func GetAnimeSortedByFavoritesWithGenres(genres []string, order int8) bson.A {
	return generateSortingWithGenres("Rating.InFavorites", genres, order)
}
