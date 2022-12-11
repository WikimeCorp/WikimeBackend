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

func CheckCommentAuthor(commentID *CommentID, userID UserID) (bool, error) {
	commID, err := primitive.ObjectIDFromHex(string(*commentID))
	err = commentsCollection.FindOne(ctx, bson.M{"Comments": bson.M{"$elemMatch": bson.M{"_id": commID, "UserId": userID}}}).Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func GetComment(commentID *CommentID) (*dbtypes.Comment, error) {
	commID, _ := primitive.ObjectIDFromHex(string(*commentID))
	matchStage := bson.D{{"$match", bson.D{{"Comments", bson.D{{"$elemMatch", bson.D{{"_id", commID}}}}}}}}
	unwindStage := bson.D{{"$unwind", "$Comments"}}
	match2Stage := bson.D{{"$match", bson.D{{"Comments._id", commID}}}}
	finalStage := bson.D{{"$replaceRoot", bson.D{{"newRoot", "$Comments"}}}}

	cur, err := commentsCollection.Aggregate(ctx, mongo.Pipeline{matchStage, unwindStage, match2Stage, finalStage})
	if err != nil {
		return nil, err
	}
	results := make([]*dbtypes.Comment, 0)
	for cur.Next(ctx) {
		elem := dbtypes.Comment{}
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}
		results = append(results, &elem)
	}
	if len(results) == 0 {
		return nil, inerr.ErrCommentNotFound
	}
	ans := results[0]

	return ans, nil
}
