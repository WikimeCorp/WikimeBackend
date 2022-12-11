package comments

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/WikimeCorp/WikimeBackend/applogic/authentication"
	"github.com/WikimeCorp/WikimeBackend/applogic/comments"
	"github.com/WikimeCorp/WikimeBackend/applogic/user"
	"github.com/WikimeCorp/WikimeBackend/dependencies"
	"github.com/WikimeCorp/WikimeBackend/restapi/handlers/other"
	"github.com/WikimeCorp/WikimeBackend/types"
	"github.com/WikimeCorp/WikimeBackend/types/myerrors"
	"github.com/gorilla/mux"

	apiErrors "github.com/WikimeCorp/WikimeBackend/restapi/errors"
)

func CreateCommentEndpoint(w http.ResponseWriter, req *http.Request) {
	userID := req.Context().Value(dependencies.CtxUserID).(types.UserID)

	commentReq := CreateCommentRequest{}
	err := other.CheckRequestJSONData(w, req, &commentReq)
	if err != nil {
		return
	}

	comID, err := comments.CreateComment(commentReq.AnimeID, types.UserID(userID), commentReq.Message)
	switch err.(type) {
	case *myerrors.ErrAnimeNotFound:
		apiErrors.SetErrorInResponce(&apiErrors.ErrAnimeNotFound, w, http.StatusNotFound)
		return
	default:
	}

	ans := CreateCommentResponce{CommentID: *comID}

	ansJSON, _ := json.Marshal(ans)
	w.Write(ansJSON)
}

func GetCommentByIDEndpoint(w http.ResponseWriter, req *http.Request) {
	animeID, _ := strconv.Atoi(mux.Vars(req)["anime_id"])

	comment, err := comments.GetComments(types.AnimeID(animeID))

	switch err.(type) {
	case *myerrors.ErrAnimeNotFound:
		apiErrors.SetErrorInResponce(&apiErrors.ErrAnimeNotFound, w, http.StatusNotFound)
		return
	default:
	}

	ansJSON, _ := json.Marshal(comment)
	w.Write(ansJSON)

}

func DeleteCommentEndpoint(w http.ResponseWriter, req *http.Request) {
	userID := req.Context().Value(dependencies.CtxUserID).(types.UserID)
	user, err := user.GetUser(userID)
	if err != nil {
		switch err.(type) {
		default:
			apiErrors.SetErrorInResponce(&apiErrors.ErrInternalServerError, w, http.StatusInternalServerError)
		}
		return
	}

	_commentID, _ := mux.Vars(req)["comment_id"]
	commentID := types.CommentID(_commentID)

	commentObj, err := comments.GetComment(&commentID)
	if err != nil {
		if errors.Is(err, myerrors.ErrCommentNotFound) {
			apiErrors.SetErrorInResponce(&apiErrors.ErrCommentNotFound, w, http.StatusNotFound)
			return
		}
	}

	authAns := authentication.AnyOne(
		authentication.CheckAdmin(user),
		authentication.CheckModeratorRole(user),
		authentication.CheckUserCreatedComment(user, commentObj),
	)

	if authAns.Bool() == false {
		apiErrors.SetErrorInResponce(apiErrors.ErrForbidden.SetNewMessage(authAns.MessageIfFalse()), w, http.StatusForbidden)
		return
	}

	err = comments.DeleteComment(&commentID)
	if err != nil {
		if err == myerrors.ErrCommentNotFound {
			apiErrors.SetErrorInResponce(&apiErrors.ErrCommentNotFound, w, http.StatusNotFound)
			return
		}
		return
	}

}

func GetCommentEndpoint(w http.ResponseWriter, req *http.Request) {
	_commentID, _ := mux.Vars(req)["comment_id"]
	commentID := types.CommentID(_commentID)

	ans, err := comments.GetComment(&commentID)
	if err != nil {
		if errors.Is(err, myerrors.ErrCommentNotFound) {
			apiErrors.SetErrorInResponce(&apiErrors.ErrCommentNotFound, w, http.StatusNotFound)
			return
		}
	}

	ansBytes, _ := json.Marshal(ans)
	w.Write(ansBytes)
}
