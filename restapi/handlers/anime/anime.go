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
	"github.com/WikimeCorp/WikimeBackend/applogic/authentication"
	"github.com/WikimeCorp/WikimeBackend/applogic/user"
	"github.com/WikimeCorp/WikimeBackend/dependencies"
	"github.com/WikimeCorp/WikimeBackend/types"
	"github.com/WikimeCorp/WikimeBackend/types/myerrors"

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
		userID := req.Context().Value(dependencies.CtxUserID).(types.UserID)
		user, err := user.GetUser(userID)
		if err != nil {
			switch err.(type) {
			default:
				apiErrors.SetErrorInResponce(&apiErrors.ErrInternalServerError, w, http.StatusInternalServerError)
			}
			return
		}

		authAns := authentication.AnyOne(
			authentication.CheckAdmin(user),
			authentication.CheckModeratorRole(user),
		)

		if authAns.Bool() == false {
			apiErrors.SetErrorInResponce(apiErrors.ErrForbidden.SetNewMessage(authAns.MessageIfFalse()), w, http.StatusForbidden)
			return
		}

		animeReq := &AnimeCreateRequest{}
		err = other.CheckRequestJSONData(w, req, animeReq)
		if err != nil {
			return
		}

		animeObj := animeReq.NewAnimeModel()
		animeObj.Author = userID

		animeID, err := anime.CreateAnime(animeObj)
		if err != nil {
			log.Println("createAnimeEndpoint CreateAnime error", err)
			return
		}

		ans, err := json.Marshal(CreateAnimeResponce{AnimeID: animeID})

		w.Write(ans)

	}
}

// GetAnimeByListIDHandler ...
func GetAnimeByListIDHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		idsStrs := req.URL.Query()["id"]
		if len(idsStrs) == 0 {
			apiErrors.SetErrorInResponce(
				apiErrors.ErrBadRequest.SetNewMessage("need 'id'"),
				w,
				http.StatusBadRequest,
			)
			return
		} else if len(idsStrs) > 50 {
			apiErrors.SetErrorInResponce(
				apiErrors.ErrBadRequest.SetNewMessage("max 50 ids"),
				w,
				http.StatusBadRequest,
			)
			return
		}

		ids := make([]types.AnimeID, len(idsStrs))
		for idx, el := range idsStrs {
			id, err := strconv.Atoi(el)
			if err != nil {
				apiErrors.SetErrorInResponce(
					apiErrors.ErrBadRequest.SetNewMessage(el+" is not a number"),
					w,
					http.StatusBadRequest,
				)
				return
			}
			ids[idx] = types.AnimeID(id)
		}

		animeList, err := anime.GetAnimesByListID(ids)

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

		ans, _ := json.Marshal(animeList)
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

		ans, _ := json.Marshal(animesIDs)
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
			return
		}
		if count > 50 {
			apiErrors.SetErrorInResponce(
				apiErrors.ErrBadRequest.SetNewMessage("Invalid 'count', must less than 51"),
				w,
				http.StatusBadRequest)
			return
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

func SearchHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		toSearch := req.URL.Query().Get("search")
		if toSearch == "" {
			apiErrors.SetErrorInResponce(
				apiErrors.ErrBadRequest.SetNewMessage("search must be not null"),
				w,
				http.StatusBadRequest,
			)
			return
		} else if len(toSearch) < 3 {
			apiErrors.SetErrorInResponce(
				apiErrors.ErrBadRequest.SetNewMessage("search must have len > 2"),
				w,
				http.StatusBadRequest,
			)
			return
		}
		animes, err := anime.Search(toSearch)
		if err != nil {
			switch err.(type) {
			default:
				return
			}
		}

		ans, _ := json.Marshal(animes)
		w.Write(ans)
	}
}

func EditAnimeHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		_animeID, _ := strconv.Atoi(mux.Vars(req)["anime_id"])
		animeID := types.AnimeID(_animeID)

		reqData := EditAnimeRequest{}
		err := other.CheckRequestJSONData(w, req, &reqData)
		if err != nil {
			return
		}

		err = anime.EditAnime(animeID, reqData.NewAnimeModel())
		if err != nil {
			switch err.(type) {
			case *myerrors.ErrAnimeNotFound:
				apiErrors.SetErrorInResponce(&apiErrors.ErrAnimeNotFound, w, http.StatusNotFound)
				return
			}
			apiErrors.SetErrorInResponce(&apiErrors.ErrInternalServerError, w, http.StatusInternalServerError)
			return
		}

	}
}
