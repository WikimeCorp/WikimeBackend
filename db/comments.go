package db

import (
	"errors"
	"fmt"

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
	_, err = commentsCollection.InsertOne(ctx, dbtypes.Comments{ID: animeID})
	return err
}

// AddComment adds a comment
//
// It must be guaranteed that the user with the `id` exists
func AddComment(animeID AnimeID, userID UserID, text string) error {
	ans, err := commentsCollection.UpdateByID(ctx, animeID,
		bson.M{
			"$push": bson.M{
				"Comments": dbtypes.Comment{primitive.NewObjectID(), userID, text},
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
func DeleteCommentByID(commentID primitive.ObjectID) error {
	ans, err := commentsCollection.UpdateMany(ctx,
		bson.M{},
		bson.M{
			"$pull": bson.M{
				"Comments": bson.M{"id": commentID},
			},
		},
	)

	if err != nil {
		return err
	}
	if ans.ModifiedCount == 0 {
		return fmt.Errorf("comment with id %v not fount", commentID.Hex())
	}
	return nil
}

// GetComments returns anime comments with the id `animeid`
func GetComments(animeID AnimeID) error {
	comments := &dbtypes.Comments{}
	err := commentsCollection.FindOne(ctx, bson.M{"_id": animeID}).Decode(comments)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return &inerr.ErrAnimeNotFound{AnimeID: animeID}
	}
	return err
}
