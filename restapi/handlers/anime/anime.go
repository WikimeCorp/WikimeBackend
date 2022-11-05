package anime

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/WikimeCorp/WikimeBackend/applogic/anime"
	"github.com/WikimeCorp/WikimeBackend/types"
	"github.com/WikimeCorp/WikimeBackend/types/myerrors"

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
