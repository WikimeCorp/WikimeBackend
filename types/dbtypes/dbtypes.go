package dbtypes

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
	ReleaseDate time.Time `bson:"ReleaseDate"`
	Author      UserID    `bson:"Author"`
	Rating      *Rating   `bson:"Rating"`
}

// Rating is presentation of a document from the `Rating` collection
type Rating struct {
	Five        uint32  `bson:"Five"`
	Four        uint32  `bson:"Four"`
	Three       uint32  `bson:"Three"`
	Two         uint32  `bson:"Two"`
	One         uint32  `bson:"One"`
	InFavorites uint32  `bson:"InFavorites"`
	Average     float64 `bson:"Average"`
}

// User is presentation of a document from the `User` collection
type User struct {
	ID        UserID    `bson:"_id"`
	Nickname  string    `bson:"Nickname"`
	Role      string    `bson:"Role"`
	Favorites []AnimeID `bson:"Favorites"`
	Watched   []AnimeID `bson:"Watched"`
	Added     []AnimeID `bson:"Added"`
	Rated     []struct {
		ID   AnimeID     `bson:"AnimeId"`
		Rate AnimeRating `bson:"Rate"`
	} `bson:"Rated"`
}

// Comments is presentation of a document from the `Comments` collection
type Comments struct {
	ID       AnimeID   `bson:"_id"`
	Comments []Comment `bson:"Comments"`
}

// Comment is a comment struct
type Comment struct {
	ID      primitive.ObjectID `bson:"_id"`
	UserID  UserID             `bson:"UserId"`
	Message string             `bson:"Message"`
}
