package db

import (
	"errors"

	. "github.com/WikimeCorp/WikimeBackend/types"
	"github.com/WikimeCorp/WikimeBackend/types/myerrors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func PutRefreshToken(token string, userID UserID) error {
	_, err := tokensCollecton.InsertOne(ctx, refreshToken{token, userID})
	return err
}

func GetAndRemoveRefreshToken(token string) error {
	err := tokensCollecton.FindOneAndDelete(ctx, bson.M{"_id": token}).Err()
	if errors.Is(err, mongo.ErrNoDocuments) {
		return myerrors.ErrNoDocumentsFromDB
	}
	return err
}
