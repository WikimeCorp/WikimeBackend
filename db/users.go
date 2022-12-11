package db

import (
	"errors"

	. "github.com/WikimeCorp/WikimeBackend/types"
	"github.com/WikimeCorp/WikimeBackend/types/dbtypes"
	inerr "github.com/WikimeCorp/WikimeBackend/types/myerrors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetUser gets a user by anime id
func GetUser(id UserID) (*dbtypes.User, error) {
	ans := &dbtypes.User{}

	err := usersCollection.FindOne(ctx, bson.M{"_id": id}).Decode(ans)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, &inerr.ErrUserNotFound{id}
		}
		return nil, err
	}

	return ans, err
}

// CheckUser checks if the user exists
func CheckUser(id UserID) (bool, error) {
	_, err := GetUser(id)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, err
	}

	return true, nil

}

func CreateUserDoc(nickname string, avatarPath string) (UserID, error) {
	userID, err := getNextID[UserID]("UserID")
	if err != nil {
		return 0, err
	}

	_, err = usersCollection.InsertOne(ctx, dbtypes.User{
		ID:        userID,
		Nickname:  nickname,
		Role:      string(UserRole),
		Favorites: []AnimeID{},
		Watched:   []AnimeID{},
		Avatar:    avatarPath,
		Rated: []struct {
			ID   AnimeID     `bson:"AnimeId"`
			Rate AnimeRating `bson:"Rate"`
		}{},
	})

	return userID, err
}

// EditNickname changes the nickname of the user with ID id :)
func EditNickname(id UserID, newNickname string) error {
	ans, err := usersCollection.UpdateByID(ctx, id, bson.M{"$set": bson.M{"Nickname": newNickname}})
	if err != nil {
		return err
	}

	if ans.MatchedCount == 0 {
		err = &inerr.ErrUserNotFound{id}
	}
	return err
}

func _pushToSetInUser(userID UserID, animeID AnimeID, list string) (bool, error) {
	anime := animeCollection.FindOne(ctx, bson.M{"_id": animeID})
	if errors.Is(anime.Err(), mongo.ErrNoDocuments) {
		return false, &inerr.ErrAnimeNotFound{animeID}
	}

	ans, err := usersCollection.UpdateByID(ctx, userID, bson.M{
		"$addToSet": bson.M{list: animeID},
	})
	if err != nil {
		return false, err
	}
	if ans.MatchedCount == 0 {
		return false, &inerr.ErrUserNotFound{userID}
	}
	if ans.ModifiedCount == 0 {
		return false, nil
	}
	return true, nil
}

func _pullFromSetInUser(userID UserID, animeID AnimeID, list string) (bool, error) {
	anime := animeCollection.FindOne(ctx, bson.M{"_id": animeID})
	if errors.Is(anime.Err(), mongo.ErrNoDocuments) {
		return false, &inerr.ErrAnimeNotFound{animeID}
	}

	ans, err := usersCollection.UpdateByID(ctx, userID, bson.M{
		"$pull": bson.M{list: animeID},
	})
	if err != nil {
		return false, err
	}
	if ans.MatchedCount == 0 {
		return false, &inerr.ErrUserNotFound{userID}
	}
	if ans.ModifiedCount == 0 {
		return false, nil
	}
	return true, nil
}

func addToFavorites(userID UserID, animeID AnimeID) (bool, error) {
	return _pushToSetInUser(userID, animeID, "Favorites")
}

func deleteFromFavorites(userID UserID, animeID AnimeID) (bool, error) {
	return _pullFromSetInUser(userID, animeID, "Favorites")
}

func addToWatched(userID UserID, animeID AnimeID) (bool, error) {
	return _pushToSetInUser(userID, animeID, "Watched")
}

func deleteFromWatched(userID UserID, animeID AnimeID) (bool, error) {
	return _pullFromSetInUser(userID, animeID, "Watched")
}

func addToAdded(userID UserID, animeID AnimeID) (bool, error) {
	return _pushToSetInUser(userID, animeID, "Added")
}

func addToRated(userID UserID, animeID AnimeID, rate AnimeRating) error {
	anime := animeCollection.FindOne(ctx, bson.M{"_id": animeID})
	if errors.Is(anime.Err(), mongo.ErrNoDocuments) {
		return &inerr.ErrAnimeNotFound{animeID}
	}

	ans, err := usersCollection.UpdateByID(ctx, userID, bson.M{
		"$addToSet": bson.M{"Rated": bson.M{"AnimeId": animeID, "Rate": rate}},
	})
	if err != nil {
		return err
	}
	if ans.MatchedCount == 0 {
		err = &inerr.ErrUserNotFound{userID}
	}
	return err

}

func checkInRated(animeID AnimeID, userID UserID) (bool, error) {
	err := usersCollection.FindOne(ctx, bson.M{"_id": userID, "Rated.AnimeId": animeID}).Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func RemoveUser(id UserID) error {
	_, err := usersCollection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func ChangeRole(userID UserID, role Role) error {
	_, err := usersCollection.UpdateOne(ctx, bson.M{"_id": userID}, bson.M{"$set": bson.M{"Role": string(role)}})
	if err != nil {
		return err
	}
	return nil
}

func GetUsersByRole(role Role) ([]*dbtypes.User, error) {
	cur, err := usersCollection.Find(ctx, bson.M{"Role": role})
	if err != nil {
		return nil, err
	}

	results := make([]*dbtypes.User, 0)

	for cur.Next(ctx) {

		elem := dbtypes.User{}

		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}
		results = append(results, &elem)
	}

	return results, nil
}
