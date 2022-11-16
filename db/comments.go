package db

import (
	"errors"

	. "github.com/WikimeCorp/WikimeBackend/types"
	"github.com/WikimeCorp/WikimeBackend/types/dbtypes"
	inerr "github.com/WikimeCorp/WikimeBackend/types/myerrors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func createCommentsDoc(animeID AnimeID) error {
	anime, err := CheckAnime(animeID)
	if err != nil {
		return err
	}
	if !anime {
		return &inerr.ErrAnimeNotFound{animeID}
	}
	_, err = commentsCollection.InsertOne(ctx, dbtypes.Comments{ID: animeID, Comments: []dbtypes.Comment{}})
	return err
}

// AddComment adds a comment
//
// It must be guaranteed that the user with the `id` exists
func AddComment(animeID AnimeID, userID UserID, text string) (*CommentID, error) {
	comID := primitive.NewObjectID()
	ans, err := commentsCollection.UpdateByID(ctx, animeID,
		bson.M{
			"$push": bson.M{
				"Comments": dbtypes.Comment{comID, userID, text},
			},
		},
	)
	if err != nil {
		return nil, err
	}
	if ans.MatchedCount == 0 {
		return nil, &inerr.ErrAnimeNotFound{animeID}
	}
	tmp := CommentID(comID.Hex())
	return &tmp, nil
}

// DeleteCommentFromAnime removes comment from anime
//
// It must be guaranteed that the user with the `id` exists
func DeleteCommentFromAnime(animeID AnimeID, commentID primitive.ObjectID) error {
	ans, err := commentsCollection.UpdateByID(ctx, animeID,
		bson.M{
			"$pull": bson.M{
				"Comments": bson.M{"id": commentID},
			},
		},
	)
	if err != nil {
		return err
	}
	if ans.MatchedCount == 0 {
		return &inerr.ErrAnimeNotFound{animeID}
	}
	return nil
}

// DeleteCommentByID removes comment from anime
//
// It must be guaranteed that the user with the `id` exists
func DeleteCommentByID(commentID *CommentID) error {
	commentIDObj, err := primitive.ObjectIDFromHex(string(*commentID))
	ans, err := commentsCollection.UpdateMany(ctx,
		bson.M{},
		bson.M{
			"$pull": bson.M{
				"Comments": bson.M{"_id": commentIDObj},
			},
		},
	)

	if err != nil {
		return err
	}
	if ans.ModifiedCount == 0 {
		return inerr.ErrCommentNotFound
	}
	return nil
}

// GetComments returns anime comments with the id `animeid`
func GetComments(animeID AnimeID) ([]dbtypes.Comment, error) {
	comments := &dbtypes.Comments{}
	err := commentsCollection.FindOne(ctx, bson.M{"_id": animeID}).Decode(comments)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, &inerr.ErrAnimeNotFound{AnimeID: animeID}
	}
	return comments.Comments, nil
}
