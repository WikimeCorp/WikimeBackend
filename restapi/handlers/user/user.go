package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/WikimeCorp/WikimeBackend/applogic/user"
	"github.com/WikimeCorp/WikimeBackend/dependencies"
	apiErrors "github.com/WikimeCorp/WikimeBackend/restapi/errors"
	"github.com/WikimeCorp/WikimeBackend/restapi/handlers/other"
	"github.com/WikimeCorp/WikimeBackend/types"
	"github.com/WikimeCorp/WikimeBackend/types/myerrors"
	"github.com/gorilla/mux"
)

// GetUserHandler return get user handler
func GetUserHandler() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		userID, _ := strconv.Atoi(mux.Vars(req)["user_id"])
		user, err := user.GetUser(types.UserID(userID))

		var errUserNotFound *myerrors.ErrUserNotFound

		if err != nil {
			if errors.As(err, &errUserNotFound) {
				apiErrors.SetErrorInResponce(&apiErrors.ErrUserNotFound, w, http.StatusUnauthorized)
				return
			}
		}

		jsonAns, _ := json.Marshal(user)

		w.Write(jsonAns)
	}
}

func ChangeNicknameHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		userID := req.Context().Value(dependencies.CtxUserID).(types.UserID)

		reqData := ChangeNicknameRequest{}
		err := other.CheckRequestJSONData(w, req, &reqData)
		if err != nil {
			return
		}
		err = user.SetNickname(userID, reqData.Nickname)
		if err != nil {
			apiErrors.SetErrorInResponce(&apiErrors.ErrInternalServerError, w, http.StatusInternalServerError)
			return
		}
	}
}

// GetCurrentUserHandler return get current user handler
func GetCurrentUserHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {

		userID := req.Context().Value(dependencies.CtxUserID).(types.UserID)
		user, err := user.GetUser(userID)

		var errUserNotFound *myerrors.ErrUserNotFound

		if err != nil {
			if errors.As(err, &errUserNotFound) {
				apiErrors.SetErrorInResponce(&apiErrors.ErrUserNotFound, w, http.StatusUnauthorized)
				return
			}
		}

		jsonAns, _ := json.Marshal(user)

		w.Write(jsonAns)
	}
}

func AddToFavoritesHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		userID := req.Context().Value(dependencies.CtxUserID).(types.UserID)

		reqData := AddToFavoritesRequest{}
		err := other.CheckRequestJSONData(w, req, &reqData)
		if err != nil {
			return
		}

		err = user.AddToFavorites(userID, *reqData.AnimeID)
		switch err.(type) {
		case *myerrors.ErrAnimeNotFound:
			apiErrors.SetErrorInResponce(&apiErrors.ErrAnimeNotFound, w, http.StatusNotFound)
			return
		}
	}
}

func AddToWatchedHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		userID := req.Context().Value(dependencies.CtxUserID).(types.UserID)

		reqData := AddToWatchedRequest{}
		err := other.CheckRequestJSONData(w, req, &reqData)
		if err != nil {
			return
		}

		err = user.AddToWatched(userID, *reqData.AnimeID)
		switch err.(type) {
		case *myerrors.ErrAnimeNotFound:
			apiErrors.SetErrorInResponce(&apiErrors.ErrAnimeNotFound, w, http.StatusNotFound)
			return
		}
	}
}

func DeleteFromWatchedHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		userID := req.Context().Value(dependencies.CtxUserID).(types.UserID)

		reqData := DeleteFromWatchedRequest{}
		err := other.CheckRequestJSONData(w, req, &reqData)
		if err != nil {
			return
		}

		err = user.DeleteFromWatched(userID, *reqData.AnimeID)
		switch err.(type) {
		case *myerrors.ErrAnimeNotFound:
			apiErrors.SetErrorInResponce(&apiErrors.ErrAnimeNotFound, w, http.StatusNotFound)
			return
		}
	}
}

func DeleteFromFavoritesHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		userID := req.Context().Value(dependencies.CtxUserID).(types.UserID)

		reqData := DeleteFromFavoritesRequest{}
		err := other.CheckRequestJSONData(w, req, &reqData)
		if err != nil {
			return
		}

		err = user.DeleteFromFavorites(userID, *reqData.AnimeID)
		switch err.(type) {
		case *myerrors.ErrAnimeNotFound:
			apiErrors.SetErrorInResponce(&apiErrors.ErrAnimeNotFound, w, http.StatusNotFound)
			return
		}
	}
}
