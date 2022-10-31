package db

import (
	. "github.com/WikimeCorp/WikimeBackend/types"
)

// idBaseStruct is presentation of a document from the `IdBase` collection. Only for inner use
type idBaseStruct struct {
	ID     string `bson:"_id"`
	LastID uint32 `bson:"LastId"`
}

type genres struct {
	Genres []string `bson:"Genres"`
}

type innerUser[OuterID OuterIDs] struct {
	OuterID OuterID
	InnerID UserID
}

type refreshToken struct {
	Token  string `bson:"_id"`
	UserID UserID `bson:"UserID"`
}
