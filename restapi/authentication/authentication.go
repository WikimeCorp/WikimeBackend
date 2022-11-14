package authentication

import (
	"github.com/WikimeCorp/WikimeBackend/applogic/user"
	"github.com/WikimeCorp/WikimeBackend/types"
)

type responder func() bool
type checker func(...any) responder

// CheckAdminOrModeratorUserRole checks that the user is a moderator or admin
func CheckAdminOrModeratorUserRole(userID types.UserID) responder {
	user, _ := user.GetUser(userID)
	return func() bool {
		return AnyOne(CheckAdmin(user), CheckModeratorRole(user))
	}
}

// CheckAdmin ...
func CheckAdmin(user *user.UserModel) responder {
	return func() bool {
		return user.Role == types.AdminRole
	}
}

// CheckModeratorRole ...
func CheckModeratorRole(user *user.UserModel) responder {
	return func() bool {
		return user.Role == types.ModeratorRole
	}
}

// AnyOne checks that at least one checker is true
func AnyOne(firstChecker responder, checkers ...responder) bool {
	ans := firstChecker()
	for _, f := range checkers {
		if ans == true {
			return true
		}
		ans = f()
	}
	return ans
}
