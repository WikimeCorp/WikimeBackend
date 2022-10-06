package db

import (
	. "github.com/WikimeCorp/WikimeBackend/types"
	"go.mongodb.org/mongo-driver/bson"
)

func getNextID[T SomeID](idxName string) (answer T, err error) {
	ans := &idBaseStruct{}
	err = idBaseCollection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": idxName},
		bson.M{"$inc": bson.M{"LastId": 1}},
	).Decode(ans)

	answer = T(ans.LastID)

	return answer, err
}
