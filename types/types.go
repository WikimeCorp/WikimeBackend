package types

import "golang.org/x/exp/maps"

type UserID uint32
type AnimeID uint32
type AnimeRating uint8
type VKUserID uint32
type CommentID string

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
var (
	RootRole      Role = "root"
	AdminRole     Role = "admin"
	ModeratorRole Role = "moderator"
	UserRole      Role = "user"
	DefaultRole        = UserRole
)

func CheckRole(role string) bool {
	_, ok := rolesPriority[role]
	if ok == false {
		return false
	}
	return true
}

func GetRoles() []string {
	return maps.Keys(rolesPriority)
}

var rolesPriority = map[string]int{
	string(RootRole):      0,
	string(AdminRole):     1,
	string(ModeratorRole): 2,
	string(UserRole):      3,
}

func (r *Role) GetPriority() int {
	return rolesPriority[string(*r)]
}

type Pair[first, second any] struct {
	First  first
	Second second
}
