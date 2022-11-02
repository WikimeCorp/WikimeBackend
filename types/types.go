package types

type UserID uint32
type AnimeID uint32
type AnimeRating uint8
type VKUserID uint32

type SomeID interface {
	UserID | AnimeID
}

// OuterIDs is type for outer IDs
type OuterIDs interface {
	VKUserID
}

// JWTPayload is data in JWT token
type JWTPayload struct {
	UserID UserID `json:"sup"`
	Exp    int64  `json:"exp"`
}

// Role is user role type
type Role string

// Roles
const (
	UserRole      Role = "user"
	AdminRole          = "admin"
	ModeratorRole      = "moder"
)
