package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/WikimeCorp/WikimeBackend/applogic/authentication"
	"github.com/WikimeCorp/WikimeBackend/applogic/user"
	"github.com/WikimeCorp/WikimeBackend/dependencies"
	apiErrors "github.com/WikimeCorp/WikimeBackend/restapi/errors"
	"github.com/WikimeCorp/WikimeBackend/restapi/handlers/other"
	"github.com/WikimeCorp/WikimeBackend/types"
	"github.com/WikimeCorp/WikimeBackend/types/myerrors"
	"github.com/gorilla/mux"
)

// GetUserHandler return get user handler
func GetUserHandler(userIDGetter func(*http.Request) types.UserID) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		userID := userIDGetter(req)
		userObj, err := user.GetUser(userID)

		var errUserNotFound *myerrors.ErrUserNotFound

		if err != nil {
			if errors.As(err, &errUserNotFound) {
				apiErrors.SetErrorInResponce(&apiErrors.ErrUserNotFound, w, http.StatusNotFound)
				return
			}
		}

		addedAnime, err := user.GetAddedAnimeByUser(userID)
		if err != nil {
			apiErrors.SetErrorInResponce(&apiErrors.ErrInternalServerError, w, http.StatusInternalServerError)
			return
		}
		ansUser := User{UserModel: *userObj, Added: addedAnime}
		jsonAns, _ := json.Marshal(ansUser)

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
func GetCurrentUserHandler(userIDGetter func(*http.Request) types.UserID) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {

		userID := userIDGetter(req)
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

func ChangeRoleHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		currentUserID := req.Context().Value(dependencies.CtxUserID).(types.UserID)
		currentUser, _ := user.GetUser(currentUserID)

		userID, _ := strconv.Atoi(mux.Vars(req)["user_id"])
		secondUser, err := user.GetUser(types.UserID(userID))

		var errUserNotFound *myerrors.ErrUserNotFound

		if err != nil {
			if errors.As(err, &errUserNotFound) {
				apiErrors.SetErrorInResponce(&apiErrors.ErrUserNotFound, w, http.StatusUnauthorized)
				return
			}
		}

		authAns := authentication.CheckAll(
			authentication.CheckAdmin(currentUser),
			authentication.CheckFirstUserGreaterThenSecondUserByPriority(currentUser, secondUser),
		)()

		if authAns.Bool() == false {
			apiErrors.SetErrorInResponce(apiErrors.ErrForbidden.SetNewMessage(authAns.MessageIfFalse()), w, http.StatusForbidden)
			return
		}

		roleName := req.URL.Query().Get("roleName")
		if roleName == "" {
			apiErrors.SetErrorInResponce(
				apiErrors.ErrBadRequest.SetNewMessage("request must contain roleName parameter"),
				w,
				http.StatusBadRequest,
			)
			return
		}

		tmp := types.CheckRole(roleName)
		if tmp == false {
			apiErrors.SetErrorInResponce(
				apiErrors.ErrBadRequest.SetNewMessage("roleName must be one of: "+strings.Join(types.GetRoles(), ", ")),
				w,
				http.StatusBadRequest,
			)
			return
		}

		authAns = authentication.CheckFirstRoleGreaterThenSecondRole(currentUser.Role, types.Role(roleName))()
		if authAns.Bool() == false {
			apiErrors.SetErrorInResponce(apiErrors.ErrForbidden.SetNewMessage("your role less or equal"+roleName), w, http.StatusForbidden)
			return
		}

		err = user.ChangeRole(types.UserID(userID), types.Role(roleName))
		if err != nil {
			fmt.Println(err)
			apiErrors.SetErrorInResponce(&apiErrors.ErrInternalServerError, w, http.StatusInternalServerError)
			return
		}

	}
}

func ResetRoleHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		currentUserID := req.Context().Value(dependencies.CtxUserID).(types.UserID)
		currentUser, _ := user.GetUser(currentUserID)

		userID, _ := strconv.Atoi(mux.Vars(req)["user_id"])
		secondUser, err := user.GetUser(types.UserID(userID))

		var errUserNotFound *myerrors.ErrUserNotFound

		if err != nil {
			if errors.As(err, &errUserNotFound) {
				apiErrors.SetErrorInResponce(&apiErrors.ErrUserNotFound, w, http.StatusUnauthorized)
				return
			}
		}

		authAns := authentication.CheckAll(
			authentication.CheckAdmin(currentUser),
			authentication.CheckFirstUserGreaterThenSecondUserByPriority(currentUser, secondUser),
		)()

		if authAns.Bool() == false {
			apiErrors.SetErrorInResponce(apiErrors.ErrForbidden.SetNewMessage(authAns.MessageIfFalse()), w, http.StatusForbidden)
			return
		}

		err = user.ChangeRole(types.UserID(userID), types.DefaultRole)
		if err != nil {
			fmt.Println(err)
			apiErrors.SetErrorInResponce(&apiErrors.ErrInternalServerError, w, http.StatusInternalServerError)
			return
		}

	}
}

func GetModeratorsHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		res, err := user.GetUsersByRole(types.ModeratorRole)
		if err != nil {
			apiErrors.SetErrorInResponce(&apiErrors.ErrInternalServerError, w, http.StatusInternalServerError)
			return
		}

		ans := make([]ShortUser, 0, len(res))
		for _, el := range res {
			ans = append(ans, ShortUser{ID: el.UserID, Avatar: el.Avatar, Nickname: el.Nickname})
		}

		ansBytes, _ := json.Marshal(ans)
		w.Write(ansBytes)
	}
}

func GetAdminsHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		res, err := user.GetUsersByRole(types.AdminRole)
		if err != nil {
			apiErrors.SetErrorInResponce(&apiErrors.ErrInternalServerError, w, http.StatusInternalServerError)
			return
		}

		ans := make([]ShortUser, 0, len(res))
		for _, el := range res {
			ans = append(ans, ShortUser{ID: el.UserID, Avatar: el.Avatar, Nickname: el.Nickname})
		}

		ansBytes, _ := json.Marshal(ans)
		w.Write(ansBytes)
	}
}
