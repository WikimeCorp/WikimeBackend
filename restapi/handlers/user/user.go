package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/WikimeCorp/WikimeBackend/applogic/user"
	apiErrors "github.com/WikimeCorp/WikimeBackend/restapi/errors"
	"github.com/WikimeCorp/WikimeBackend/types"
	"github.com/WikimeCorp/WikimeBackend/types/myerrors"
	"github.com/gorilla/mux"
)

func getUser(w http.ResponseWriter, req *http.Request) {
	userID, _ := strconv.Atoi(mux.Vars(req)["user_id"])
	user, err := user.GetUser(types.UserID(userID))

	var errUserNotFound *myerrors.ErrUserNotFound

	if err != nil {
		if errors.As(err, &errUserNotFound) {
			apiErrors.SetErrorInResponce(&apiErrors.ErrUserNotFound, w, http.StatusUnauthorized)
			return
		}
	}

	jsonAns, _ := json.Marshal(user)

	w.Write(jsonAns)
}

// GetUserHandler return get user handler
func GetUserHandler() func(w http.ResponseWriter, req *http.Request) {
	return getUser
}
