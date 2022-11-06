package anime

import (
	"encoding/json"
	"errors"
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

func GetAnimeByIDHandler() func(http.ResponseWriter, *http.Request) {
	return getAnimeByIDEndpoint
}

func createAnimeEndpoint(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		log.Println("createAnimeEndpoint ParseForm error", err)
		return
	}

	animeReq := &AnimeCreateRequest{}

	err = json.NewDecoder(req.Body).Decode(animeReq)

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
