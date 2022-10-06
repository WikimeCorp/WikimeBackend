package types

type UserID uint32
type AnimeID uint32
type AnimeRating uint8

type SomeID interface {
	UserID | AnimeID
}
