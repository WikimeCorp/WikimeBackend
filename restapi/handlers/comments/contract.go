package comments

import "github.com/WikimeCorp/WikimeBackend/types"

type CreateCommentRequest struct {
	Message string        `validate:"required" json:"message"`
	AnimeID types.AnimeID `validate:"required" json:"animeId"`
}

type CreateCommentResponce struct {
	CommentID types.CommentID `json:"commentId"`
}
