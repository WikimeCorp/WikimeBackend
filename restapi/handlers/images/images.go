package images

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/WikimeCorp/WikimeBackend/applogic/anime"
	"github.com/WikimeCorp/WikimeBackend/applogic/authentication"
	"github.com/WikimeCorp/WikimeBackend/applogic/images"
	"github.com/WikimeCorp/WikimeBackend/applogic/user"
	"github.com/WikimeCorp/WikimeBackend/config"
	"github.com/WikimeCorp/WikimeBackend/dependencies"
	apiErrors "github.com/WikimeCorp/WikimeBackend/restapi/errors"
	"github.com/WikimeCorp/WikimeBackend/types"
	"github.com/WikimeCorp/WikimeBackend/types/myerrors"
	"github.com/gorilla/mux"
)

var maxUploadesFileSize int64 = config.Config.MaxUploadedFileSize * 1024 * 1024

func fileFromMultipartForm(w http.ResponseWriter, req *http.Request) (io.ReadSeeker, string, error) {
	// Check file size
	req.Body = http.MaxBytesReader(w, req.Body, maxUploadesFileSize)
	if err := req.ParseMultipartForm(maxUploadesFileSize); err != nil {
		// Maybe need and for too big image size
		apiErrors.SetErrorInResponce(apiErrors.ErrBadRequest.SetNewMessage(err.Error()), w, http.StatusBadRequest)
		return nil, "", err
	}

	image, _, err := req.FormFile("file")
	if err != nil {
		apiErrors.SetErrorInResponce(&apiErrors.ErrBadRequest, w, http.StatusBadRequest)
		return nil, "", err
	}

	tail, err := checkFileType(w, image)
	if err != nil {
		return nil, "", err
	}

	return image, tail, nil
}

func checkFileType(w http.ResponseWriter, file io.ReadSeeker) (string, error) {
	buff := make([]byte, 512)
	_, err := file.Read(buff)
	if err != nil {
		apiErrors.SetErrorInResponce(&apiErrors.ErrInternalServerError, w, http.StatusInternalServerError)
		return "", err
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		apiErrors.SetErrorInResponce(&apiErrors.ErrInternalServerError, w, http.StatusInternalServerError)
		return "", err
	}

	filetype := http.DetectContentType(buff)
	if filetype == "image/png" {
		return "png", nil
	} else if filetype == "image/jpeg" {
		return "jpg", nil
	} else {
		apiErrors.SetErrorInResponce(&apiErrors.ErrBadImageFormat, w, http.StatusBadRequest)
		return "", &apiErrors.ErrBadImageFormat
	}
}

func fileFromRaw(w http.ResponseWriter, req *http.Request) (io.ReadSeeker, string, error) {
	tmp, _ := io.ReadAll(req.Body)
	file := bytes.NewReader(tmp)

	tail, err := checkFileType(w, file)
	if err != nil {
		return nil, "", err
	}
	return file, tail, nil
}

func addImageEndpoint(w http.ResponseWriter, req *http.Request) {
	_animeID, _ := strconv.Atoi(mux.Vars(req)["anime_id"])
	animeID := types.AnimeID(_animeID)

	animeObj, err := anime.GetAnimeByID(animeID)
	if err != nil {
		switch err.(type) {
		case *myerrors.ErrAnimeNotFound:
			apiErrors.SetErrorInResponce(&apiErrors.ErrAnimeNotFound, w, http.StatusBadRequest)
		default:
			apiErrors.SetErrorInResponce(&apiErrors.ErrInternalServerError, w, http.StatusInternalServerError)
		}
		return
	}

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
		authentication.CheckAll(
			authentication.CheckModeratorRole(user),
			authentication.CheckUserCreatedAnime(user, animeObj),
		),
	)

	if authAns.Bool() == false {
		apiErrors.SetErrorInResponce(apiErrors.ErrForbidden.SetNewMessage(authAns.MessageIfFalse()), w, http.StatusForbidden)
		return
	}

	image, tail, err := fileFromMultipartForm(w, req) //fileFromRaw(w, req)
	if err != nil {
		return
	}

	err = images.AddImageToAnime(animeID, image, tail)
	if err != nil {
		apiErrors.SetErrorInResponce(&apiErrors.ErrInternalServerError, w, http.StatusInternalServerError)
	}
}

// AddImageHandler ...
func AddImageHandler() func(http.ResponseWriter, *http.Request) {
	return addImageEndpoint
}

func setPosterEndpoint(w http.ResponseWriter, req *http.Request) {
	_animeID, _ := strconv.Atoi(mux.Vars(req)["anime_id"])
	animeID := types.AnimeID(_animeID)

	animeObj, err := anime.GetAnimeByID(animeID)
	if err != nil {
		switch err.(type) {
		case *myerrors.ErrAnimeNotFound:
			apiErrors.SetErrorInResponce(&apiErrors.ErrAnimeNotFound, w, http.StatusBadRequest)
		default:
			apiErrors.SetErrorInResponce(&apiErrors.ErrInternalServerError, w, http.StatusInternalServerError)
		}
		return
	}

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
		authentication.CheckAll(
			authentication.CheckModeratorRole(user),
			authentication.CheckUserCreatedAnime(user, animeObj),
		),
	)

	if authAns.Bool() == false {
		apiErrors.SetErrorInResponce(apiErrors.ErrForbidden.SetNewMessage(authAns.MessageIfFalse()), w, http.StatusForbidden)
		return
	}

	image, tail, err := fileFromMultipartForm(w, req) //fileFromRaw(w, req)
	if err != nil {
		return
	}

	err = images.SetPoster(animeID, image, tail)
	if err != nil {
		apiErrors.SetErrorInResponce(&apiErrors.ErrInternalServerError, w, http.StatusInternalServerError)
	}
}

// SetPosterHandler is handler for changing anime poster
func SetPosterHandler() func(http.ResponseWriter, *http.Request) {
	return setPosterEndpoint
}

func SetUserImageHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		userID := req.Context().Value(dependencies.CtxUserID).(types.UserID)

		image, tail, err := fileFromMultipartForm(w, req) //fileFromRaw(w, req)
		if err != nil {
			return
		}
		err = images.SetAvatar(userID, image, tail)
		if err != nil {
			apiErrors.SetErrorInResponce(&apiErrors.ErrInternalServerError, w, http.StatusInternalServerError)
		}

	}
}

func DeleteImageFromAnimeHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		_animeID, _ := strconv.Atoi(mux.Vars(req)["anime_id"])
		animeID := types.AnimeID(_animeID)

		animeObj, err := anime.GetAnimeByID(animeID)
		if err != nil {
			switch err.(type) {
			case *myerrors.ErrAnimeNotFound:
				apiErrors.SetErrorInResponce(&apiErrors.ErrAnimeNotFound, w, http.StatusBadRequest)
			default:
				apiErrors.SetErrorInResponce(&apiErrors.ErrInternalServerError, w, http.StatusInternalServerError)
			}
			return
		}
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
			authentication.CheckAll(
				authentication.CheckModeratorRole(user),
				authentication.CheckUserCreatedAnime(user, animeObj),
			),
		)

		if authAns.Bool() == false {
			apiErrors.SetErrorInResponce(apiErrors.ErrForbidden.SetNewMessage(authAns.MessageIfFalse()), w, http.StatusForbidden)
			return
		}
		image := mux.Vars(req)["image"]
		fmt.Println("image: ", image)
		err = images.DeleteImageFromAnime(animeID, image)
		if err != nil {
			if errors.Is(err, myerrors.ErrImageNotFound) {
				apiErrors.SetErrorInResponce(&apiErrors.ErrNotFound, w, http.StatusNotFound)
				return
			} else if errors.Is(err, myerrors.ErrForbiden) {
				apiErrors.SetErrorInResponce(&apiErrors.ErrForbidden, w, http.StatusForbidden)
				return
			}
			fmt.Println(err)
			apiErrors.SetErrorInResponce(&apiErrors.ErrInternalServerError, w, http.StatusInternalServerError)
			return
		}
	}
}
