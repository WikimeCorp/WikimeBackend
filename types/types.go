package types

type UserID uint32
type AnimeID uint32
type AnimeRating uint8
type VKUserID uint32

type SomeID interface {
	UserID | AnimeID
}

type OuterIDs interface {
	VKUserID
}

type JWTPayload struct {
	UserID   UserID `json:"sup"`
	DeadLine int64  `json:"exp"`
}
