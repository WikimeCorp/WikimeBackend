package comments

import (
	"github.com/WikimeCorp/WikimeBackend/db"
	"github.com/WikimeCorp/WikimeBackend/types"
)

func CreateComment(animeID types.AnimeID, userID types.UserID, message string) (*types.CommentID, error) {
	comID, err := db.AddComment(animeID, userID, message)
	return comID, err
}

func GetComments(animeID types.AnimeID) ([]Comment, error) {
	comments, err := db.GetComments(animeID)
	ans := make([]Comment, len(comments))
	for idx, el := range comments {
		tmp := types.CommentID(el.ID.Hex())
		ans[idx] = Comment{Message: &el.Message, Author: &el.UserID, ID: &tmp}
	}
	return ans, err
}

func DeleteComment(commentID *types.CommentID) error {
	err := db.DeleteCommentByID(commentID)
	return err
}

func CheckAuthorComment(commentID *types.CommentID, userID types.UserID) (bool, error) {
	return db.CheckCommentAuthor(commentID, userID)
}

func GetComment(commentID *types.CommentID) (*Comment, error) {
	ansDB, err := db.GetComment(commentID)
	if err != nil {
		return nil, err
	}

	ans := Comment{Message: &ansDB.Message, Author: &ansDB.UserID, ID: commentID}
	return &ans, nil
}
