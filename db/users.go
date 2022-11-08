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

func CreateUserDoc(nickname string) (UserID, error) {
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
		Added:     []AnimeID{},
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

func _pushToSet(userID UserID, animeID AnimeID, list string) error {
	anime := animeCollection.FindOne(ctx, bson.M{"_id": animeID})
	if errors.Is(anime.Err(), mongo.ErrNoDocuments) {
		return &inerr.ErrAnimeNotFound{animeID}
	}

	ans, err := usersCollection.UpdateByID(ctx, userID, bson.M{
		"$addToSet": bson.M{list: animeID},
	})
	if err != nil {
		return err
	}
	if ans.MatchedCount == 0 {
		err = &inerr.ErrUserNotFound{userID}
	}
	return err
}

func addToFavorites(userID UserID, animeID AnimeID) error {
	return _pushToSet(userID, animeID, "Favorites")
}

func deleteFromFavorites(userID UserID, animeID AnimeID) error {
	_, err := usersCollection.UpdateByID(ctx, userID, bson.M{"$pull": bson.M{"Favorites": animeID}})
	return err
}

func addToWatched(userID UserID, animeID AnimeID) error {
	return _pushToSet(userID, animeID, "Watched")
}

func addToAdded(userID UserID, animeID AnimeID) error {
	return _pushToSet(userID, animeID, "Added")
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
