package comments

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/WikimeCorp/WikimeBackend/applogic/comments"
	"github.com/WikimeCorp/WikimeBackend/restapi/handlers/other"
	"github.com/WikimeCorp/WikimeBackend/types"
	"github.com/WikimeCorp/WikimeBackend/types/myerrors"
	"github.com/gorilla/mux"

	apiErrors "github.com/WikimeCorp/WikimeBackend/restapi/errors"
)

func CreateAnimeEndpoint(w http.ResponseWriter, req *http.Request) {
	userID := 0 // Add getting user id from context

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
	_commentID, _ := mux.Vars(req)["comment_id"]
	commentID := types.CommentID(_commentID)

	err := comments.DeleteComment(&commentID)
	if err != nil {
		if err == myerrors.ErrCommentNotFound {
			apiErrors.SetErrorInResponce(&apiErrors.ErrCommentNotFound, w, http.StatusNotFound)
			return
		}
		return
	}

}
