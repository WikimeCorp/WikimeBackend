package comments

import "github.com/WikimeCorp/WikimeBackend/types"

type Comment struct {
	Message *string          `json:"message"`
	Author  *types.UserID    `json:"author"`
	ID      *types.CommentID `json:"id"`
}
