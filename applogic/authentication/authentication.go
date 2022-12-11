package authentication

import (
	"strings"

	"github.com/WikimeCorp/WikimeBackend/applogic/anime"
	"github.com/WikimeCorp/WikimeBackend/applogic/comments"
	"github.com/WikimeCorp/WikimeBackend/applogic/user"
	"github.com/WikimeCorp/WikimeBackend/types"
)

type responder func() AuthAnswer
type checker func(...any) responder

var AuthErrUserNotAdmin = "user not admin"
var AuthErrUserNotRoot = "user not admin"
var AuthErrUserNotModerator = "user not moderator"
var AuthErrUserNotAuthor = "user is not author of this object"
var AuthErrUserHasNotHigherPriority = "user does not have a higher priority"
var AuthErrFirstRoleLessThanSecondRole = "first(your) role less than second role"

type AuthAnswer interface {
	Bool() bool
	MessageIfFalse() string
	messagePointer() *string
}

type answer struct {
	ans     bool
	message *string
}

func (a *answer) Bool() bool {
	return a.ans
}

func (a *answer) MessageIfFalse() string {
	return *a.message
}

func (a *answer) messagePointer() *string {
	return a.message
}

// CheckAdminOrModeratorUserRole checks that the user is a moderator or admin
func CheckAdminOrModeratorUserRole(user *user.UserModel) responder {
	return func() AuthAnswer {
		tmp := AnyOne(CheckAdmin(user), CheckModeratorRole(user))
		ans := answer{tmp.Bool(), tmp.messagePointer()}
		return &ans
	}
}

func CheckUserCreatedAnime(userObj *user.UserModel, animeObj *anime.Anime) responder {
	return func() AuthAnswer {
		ans := answer{userObj.UserID == animeObj.Author, &AuthErrUserNotAuthor}
		return &ans
	}
}

func CheckRootRole(user *user.UserModel) responder {
	return func() AuthAnswer {
		return &answer{user.Role == types.RootRole, &AuthErrUserNotRoot}
	}
}

func CheckUserCreatedComment(userObj *user.UserModel, commentObj *comments.Comment) responder {
	return func() AuthAnswer {
		ans := answer{userObj.UserID == *commentObj.Author, &AuthErrUserNotAuthor}
		return &ans
	}
}

func firstRoleGreaterThanSecondRole(firstRole, secondRole types.Role) bool {
	// < instead of > because the lower the priority, the higher the role
	return firstRole.GetPriority() < secondRole.GetPriority()
}

func CheckFirstUserGreaterThenSecondUserByPriority(firstUserObj, secondUserObj *user.UserModel) responder {
	return func() AuthAnswer {
		return &answer{firstRoleGreaterThanSecondRole(firstUserObj.Role, secondUserObj.Role), &AuthErrUserHasNotHigherPriority}
	}
}

func CheckFirstRoleGreaterThenSecondRole(firstRole, secondRole types.Role) responder {
	return func() AuthAnswer {
		return &answer{firstRoleGreaterThanSecondRole(firstRole, secondRole), &AuthErrFirstRoleLessThanSecondRole}
	}
}

// CheckAdmin ...
func CheckAdmin(user *user.UserModel) responder {
	return func() AuthAnswer {
		return &answer{user.Role == types.AdminRole || CheckRootRole(user)().Bool(), &AuthErrUserNotAdmin}
	}
}

// CheckModeratorRole ...
func CheckModeratorRole(user *user.UserModel) responder {
	return func() AuthAnswer {
		return &answer{user.Role == types.ModeratorRole, &AuthErrUserNotModerator}
	}
}

// AnyOne checks that at least one checker is true
func AnyOne(firstChecker responder, checkers ...responder) AuthAnswer {
	ans := firstChecker()
	err := ""
	errs := make([]string, 0, len(checkers)+1)
	errs = append(errs, ans.MessageIfFalse())
	for _, f := range checkers {
		if ans.Bool() == true {
			return &answer{true, &err}
		}
		ans = f()
		errs = append(errs, ans.MessageIfFalse())
	}
	if ans.Bool() == false {
		err = "no one: " + strings.Join(errs, ", ") + ";"
		return &answer{false, &err}
	}
	return &answer{true, &err}

}

func CheckAll(firstChecker responder, checkers ...responder) responder {
	return func() AuthAnswer {
		errs := make([]string, 0, len(checkers)+1)
		errStr := ""
		ans := answer{true, &errStr}
		checkers = append(checkers, firstChecker)
		for _, f := range checkers {
			_ans := f()
			if _ans.Bool() == false {
				ans.ans = false
				errs = append(errs, _ans.MessageIfFalse())
			}
		}
		if ans.Bool() == false {
			errStr = "no: " + strings.Join(errs, ", ") + ";"
			return &answer{false, &errStr}
		}
		return &ans
	}

}
