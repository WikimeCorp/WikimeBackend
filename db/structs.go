package db

import (
	"time"

	. "github.com/WikimeCorp/WikimeBackend/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Anime is presentation of a document from the `Anime` collection
type Anime struct {
	ID          AnimeID   `bson:"_id"`
	Title       string    `bson:"Title"`
	OriginTitle string    `bson:"OriginTitle"`
	Genres      []string  `bson:"Genres"`
	Description string    `bson:"Description"`
	Images      []string  `bson:"Images"`
	URLs        []string  `bson:"URLs"`
	Director    string    `bson:"Director"`
	DateAdded   time.Time `bson:"DateAdded"`
	RelaseDate  time.Time `bson:"RelaseDate"`
	Author      UserID    `bson:"Author"`
}

// Rating is presentation of a document from the `Rating` collection
type Rating struct {
	ID          AnimeID `bson:"_id"`
	Five        uint32
	Four        uint32
	Three       uint32
	Two         uint32
	One         uint32
	InFavorites uint32
	Average     float64
}

// User is presentation of a document from the `User` collection
type User struct {
	ID        UserID `bson:"_id"`
	Nickname  string
	Role      string
	Favorites []AnimeID
	Watched   []AnimeID
	Added     []AnimeID
	Rated     []struct {
		ID   AnimeID `bson:"AnimeId"`
		Rate AnimeRating
	}
}

// Comments is presentation of a document from the `Comments` collection
type Comments struct {
	ID       AnimeID `bson:"_id"`
	Comments []struct {
		id      primitive.ObjectID
		UserID  UserID
		Message string
	}
}

// idBaseStruct is presentation of a document from the `IdBase` collection. Only for inner use
type idBaseStruct struct {
	ID     string `bson:"_id"`
	LastID uint32 `bson:"LastId"`
}
