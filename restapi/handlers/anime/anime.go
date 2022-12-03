package anime

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/WikimeCorp/WikimeBackend/applogic/anime"
	"github.com/WikimeCorp/WikimeBackend/dependencies"
	"github.com/WikimeCorp/WikimeBackend/types"
	"github.com/WikimeCorp/WikimeBackend/types/myerrors"
	"github.com/go-playground/validator/v10"

	apiErrors "github.com/WikimeCorp/WikimeBackend/restapi/errors"
	"github.com/WikimeCorp/WikimeBackend/restapi/handlers/other"
	"github.com/gorilla/mux"
)

// GetAnimeByIDHandler ...
func GetAnimeByIDHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		animeID, _ := strconv.Atoi(mux.Vars(req)["anime_id"])
		anime, err := anime.GetAnimeByID(types.AnimeID(animeID))
		var errAnimeNotFound *myerrors.ErrAnimeNotFound
		if err != nil {
			if errors.As(err, &errAnimeNotFound) {
				apiErrors.SetErrorInResponce(&apiErrors.ErrAnimeNotFound, w, http.StatusNotFound)
				return
			}
			apiErrors.SetErrorInResponce(&apiErrors.ErrInternalServerError, w, http.StatusInternalServerError)
			return
		}

		jsonAns, _ := json.Marshal(anime)

		w.Write(jsonAns)
	}
}

// CreateAnimeHandler ...
func CreateAnimeHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		animeReq := &AnimeCreateRequest{}
		err := json.NewDecoder(req.Body).Decode(animeReq)

		if err != nil {
			fmt.Println(err)
			apiErrors.SetErrorInResponce(&apiErrors.ErrBadJSONStruct, w, http.StatusBadRequest)
			return
		}

		err = dependencies.Validate.Struct(animeReq)

		if err != nil {
			tmpErrors := make([]string, 0)
			for _, err := range err.(validator.ValidationErrors) {
				tmpErrors = append(tmpErrors, err.Error())
			}

			err := apiErrors.ErrValidate(tmpErrors)
			apiErrors.SetErrorInResponce(err, w, http.StatusBadRequest)
			return
		}

		animeID, err := anime.CreateAnime(animeReq.NewAnimeModel())
		if err != nil {
			log.Println("createAnimeEndpoint CreateAnime error", err)
			return
		}

		ans, err := json.Marshal(AnimeResponce{AnimeID: animeID})

		w.Write(ans)

	}
}

// GetAnimeByListIDHandler ...
func GetAnimeByListIDHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		animeListReq := AnimeByListIDRequest{}
		err := other.CheckRequestJSONData(w, req, &animeListReq)
		if err != nil {
			return
		}

		animeList, err := anime.GetAnimesByListID(animeListReq.IDs)

		var errAnimeNotFound *myerrors.ErrAnimeNotFound

		if err != nil {
			if errors.As(err, &errAnimeNotFound) {
				err := err.(*myerrors.ErrAnimeNotFound)
				apiErrors.SetErrorInResponce(
					apiErrors.ErrAnimeNotFound.SetNewMessage(
						fmt.Sprintf("Anime with anime id %d not found", err.AnimeID)),
					w, http.StatusNotFound)
				return
			}
			log.Fatal(err)
			return

		}

		animeListRes := AnimeByListIDResponce{Animes: animeList}
		ans, _ := json.Marshal(animeListRes)
		w.Write(ans)
	}
}

func GetAnimesHangler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		sortBy := req.URL.Query().Get("sortBy")
		if sortBy == "" {
			apiErrors.SetErrorInResponce(
				apiErrors.ErrBadRequest.SetNewMessage("request must contain sortBy parameter"),
				w,
				http.StatusBadRequest,
			)
			return
		}

		orderStr := req.URL.Query().Get("order")
		order := int8(-1)
		var err error
		if orderStr != "" {
			var orderInt int
			orderInt, err = strconv.Atoi(orderStr)
			order = int8(orderInt)
		}

		if err != nil || (order != -1 && order != 1) {
			apiErrors.SetErrorInResponce(
				apiErrors.ErrBadRequest.SetNewMessage("Invalid order, must be 1 or -1"),
				w,
				http.StatusBadRequest,
			)
			return
		}

		genresIn := req.URL.Query().Get("genres")
		genres := make([]string, 0)
		if genresIn != "" {
			genres = strings.Split(genresIn, ",")
		}

		animesIDs := make([]types.AnimeID, 0)

		if sortBy == "rating" {
			animesIDs, err = anime.GetAnimeSortedByRating(genres, order)
		} else if sortBy == "releaseDate" {
			animesIDs, err = anime.GetAnimeSortedByReleaseDate(genres, order)
		} else if sortBy == "dateAdded" {
			animesIDs, err = anime.GetAnimeSortedByAddingDate(genres, order)
		} else if sortBy == "favorites" {
			animesIDs, err = anime.GetAnimeSortedByFavorites(genres, order)
		} else {
			apiErrors.SetErrorInResponce(
				apiErrors.ErrBadRequest.SetNewMessage("Cannot sort by"+sortBy),
				w,
				http.StatusBadRequest,
			)
			return
		}

		if err != nil {
			switch cErr := err.(type) {
			case *myerrors.ErrWrongGenres:
				apiErrors.SetErrorInResponce(
					apiErrors.ErrBadRequest.SetNewMessage("Bad genres: "+strings.Join(cErr.Genres, ", ")),
					w,
					http.StatusBadRequest,
				)
				return
			}
		}

		ansStruct := struct {
			AnimeIDs []types.AnimeID `json:"animeIDs"`
		}{animesIDs}
		ans, _ := json.Marshal(ansStruct)
		w.Write(ans)

	}
}

func SetAverageEndpoint(w http.ResponseWriter, req *http.Request) {
	animeID, _ := strconv.Atoi(mux.Vars(req)["anime_id"])
	reqJson := struct {
		Average *float64 `validate:"required" json:"average"`
	}{}

	err := other.CheckRequestJSONData(w, req, &reqJson)
	if err != nil {
		return
	}

	err = anime.SetAverage(types.AnimeID(animeID), *reqJson.Average)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

}

func MostPopularHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		countStr := req.URL.Query().Get("count")
		count, err := strconv.Atoi(countStr)
		if err != nil || count < 0 {
			apiErrors.SetErrorInResponce(
				apiErrors.ErrBadRequest.SetNewMessage("Invalid 'count' type, must be uint"),
				w,
				http.StatusBadRequest)
		}

		animes, err := anime.GetMostPopular(count)
		if err != nil {
			switch err.(type) {
			default:
			}
			return
		}

		ans, _ := json.Marshal(animes)
		w.Write(ans)
	}
}
