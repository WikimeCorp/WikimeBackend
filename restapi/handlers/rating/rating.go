package rating

import (
	"log"
	"net/http"
	"strconv"

	"github.com/WikimeCorp/WikimeBackend/applogic/anime"
	"github.com/WikimeCorp/WikimeBackend/dependencies"
	apiErrors "github.com/WikimeCorp/WikimeBackend/restapi/errors"
	"github.com/WikimeCorp/WikimeBackend/restapi/handlers/other"
	"github.com/WikimeCorp/WikimeBackend/types"
	"github.com/WikimeCorp/WikimeBackend/types/myerrors"
	"github.com/gorilla/mux"
)

// SetRatingHandler is handler for set rating req
func SetRatingHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		userID := req.Context().Value(dependencies.CtxUserID).(types.UserID)
		_animeID, _ := strconv.Atoi(mux.Vars(req)["anime_id"])
		animeID := types.AnimeID(_animeID)

		reqData := SetRatingRequest{}
		err := other.CheckRequestJSONData(w, req, &reqData)
		if err != nil {
			return
		}

		err = anime.SetRating(userID, animeID, *reqData.Rating)
		if err != nil {
			switch err.(type) {
			case *myerrors.ErrAnimeNotFound:
				apiErrors.SetErrorInResponce(&apiErrors.ErrAnimeNotFound, w, http.StatusNotFound)
				return
			default:
				log.Println("Error: ", err)
				apiErrors.SetErrorInResponce(&apiErrors.ErrInternalServerError, w, http.StatusInternalServerError)
				return
			}
		}
	}
}
