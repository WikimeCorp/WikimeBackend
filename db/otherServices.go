package db

import (
	"errors"

	. "github.com/WikimeCorp/WikimeBackend/types"
	"github.com/WikimeCorp/WikimeBackend/types/myerrors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CheckVkUserInDB(id VKUserID) (UserID, error) {
	ans := &innerUser[VKUserID]{}
	err := vkCollection.FindOne(ctx, bson.M{"_id": id}).Decode(ans)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return 0, myerrors.ErrNoDocumentsFromDB
		}

		return 0, err
	}

	return ans.InnerID, nil

}

func AddVkUser(vkID VKUserID, innerID UserID) error {
	_, err := vkCollection.InsertOne(ctx, bson.D{{"_id", vkID}, {"InnerId", innerID}})
	return err
}
