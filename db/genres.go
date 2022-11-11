package db

import (
	"fmt"
	"log"

	"github.com/WikimeCorp/WikimeBackend/utils"
	"go.mongodb.org/mongo-driver/bson"
)

var genresList []string

func initGenres() ([]string, error) {
	_genreList := &genres{}
	err := genresCollection.FindOne(ctx, bson.M{"_id": "Genres"}).Decode(_genreList)
	return _genreList.Genres, err
}

func CheckGenres(genres []string) (res bool, badGenres []string) {
	res = true
	if len(genres) > len(genresList) {
		return false, []string{fmt.Sprintf("Too many genres, max %d", len(genresList))}
	}
	for _, genre := range genres {
		_, inGenres := utils.BinarySearch(genresList, genre)
		if !inGenres {
			res = false
			badGenres = append(badGenres, genre)
		}
	}

	return res, badGenres
}

func GetGenres() (ans []string) {
	copy(ans, genresList)
	return ans
}

func init() {
	tmpList, err := initGenres()
	if err != nil {
		log.Fatal(err)
	}
	for _, genre := range tmpList {
		genresList = utils.InsertInSorted(genresList, genre)
	}
}
