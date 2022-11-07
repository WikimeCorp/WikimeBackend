package anime

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/WikimeCorp/WikimeBackend/applogic/anime"
	"github.com/WikimeCorp/WikimeBackend/dependencies"
	"github.com/WikimeCorp/WikimeBackend/types"
	"github.com/WikimeCorp/WikimeBackend/types/myerrors"
	"github.com/go-playground/validator/v10"

	apiErrors "github.com/WikimeCorp/WikimeBackend/restapi/errors"
	"github.com/gorilla/mux"

	"golang.org/x/sync/errgroup"
)

func getAnimeListEndpoint(w http.ResponseWriter, req *http.Request) {

}

func getAnimeByIDEndpoint(w http.ResponseWriter, req *http.Request) {
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

// GetAnimeByIDHandler ...
func GetAnimeByIDHandler() func(http.ResponseWriter, *http.Request) {
	return getAnimeByIDEndpoint
}

func createAnimeEndpoint(w http.ResponseWriter, req *http.Request) {
	animeReq := &AnimeCreateRequest{}

	err := json.NewDecoder(req.Body).Decode(animeReq)

	if err != nil {
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

// CreateAnimeHandler ...
func CreateAnimeHandler() func(http.ResponseWriter, *http.Request) {
	return createAnimeEndpoint
}

func getAnimeByListID(w http.ResponseWriter, req *http.Request) {
	animeListReq := &AnimeByListIDRequest{}
	err := json.NewDecoder(req.Body).Decode(animeListReq)

	if err != nil {
		apiErrors.SetErrorInResponce(&apiErrors.ErrBadJSONStruct, w, http.StatusBadRequest)
		return
	}

	err = dependencies.Validate.Struct(animeListReq)

	if err != nil {
		tmpErrors := make([]string, 0)
		for _, err := range err.(validator.ValidationErrors) {
			tmpErrors = append(tmpErrors, err.Error())
		}

		err := apiErrors.ErrValidate(tmpErrors)
		apiErrors.SetErrorInResponce(err, w, http.StatusBadRequest)
		return
	}

	errg, ctx := errgroup.WithContext(context.Background())

	results := make(chan types.Pair[int, *anime.Anime])

	for idx, id := range animeListReq.IDs {
		idx := idx
		id := id
		errg.Go(func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				animeAns, err := anime.GetAnimeByID(id)
				if err != nil {
					return err
				}

				select {
				case results <- types.Pair[int, *anime.Anime]{First: idx, Second: animeAns}:
					return nil
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		})
	}

	go func() {
		errg.Wait()
		close(results)
	}()

	animeListRes := AnimeByListIDResponce{Animes: make([]*anime.Anime, len(animeListReq.IDs))}

	for result := range results {
		animeListRes.Animes[result.First] = result.Second
	}

	err = errg.Wait()

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

	ans, _ := json.Marshal(animeListRes)
	w.Write(ans)
}

// GetAnimeByListIDHandler ...
func GetAnimeByListIDHandler() func(http.ResponseWriter, *http.Request) {
	return getAnimeByListID
}
