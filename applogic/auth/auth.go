package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path"
	"strconv"

	"github.com/WikimeCorp/WikimeBackend/config"
	"github.com/WikimeCorp/WikimeBackend/db"
	. "github.com/WikimeCorp/WikimeBackend/types"
	"github.com/WikimeCorp/WikimeBackend/types/myerrors"
)

func VkAuth(authToken string) (jwtToken string, err error) {
	innerUserID, err := getInnerUserIDForVk(authToken)

	if err != nil {
		return "", err
	}

	jwtToken = generateJWT(innerUserID)
	return jwtToken, nil
}

func getInnerUserIDForVk(token string) (UserID, error) {
	vkUserID, err := getVkUserID(token)
	if err != nil {
		return 0, err
	}

	userID, err := db.CheckVkUserInDB(vkUserID)
	if err != nil {
		if errors.Is(err, myerrors.ErrNoDocumentsFromDB) {
			userID, err = registerUser(vkUserID)
			if err != nil {
				return 0, err
			}
		} else {
			return 0, err
		}
	}

	return userID, nil
}

func registerUser[T OuterIDs](outerID T) (userID UserID, err error) {
	userID, err = db.CreateUserDoc("", path.Join(config.Config.ImagesPathURI, config.Config.DefaultUserAvatarPath))
	if err != nil {
		return 0, err
	}

	defer func() {
		if err != nil {
			db.RemoveUser(userID)
			userID = 0
		}
	}()

	err = db.EditNickname(userID, "id"+strconv.Itoa(int(userID)))
	if err != nil {
		return 0, err
	}
	switch any(outerID).(type) {
	case VKUserID:
		err = db.AddVkUser(VKUserID(outerID), userID)
	}

	if err != nil {
		return 0, err
	}

	return userID, nil
}

func getVkUserID(token string) (VKUserID, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.vk.com/method/users.get?access_token=%s&v=%s",
		token,
		config.Config.VKAPIVersion))

	if err != nil {
		return 0, err
	}

	ans := struct {
		Error *struct {
			ErrorCode     int    `json:"error_code"`
			ErrorMsg      string `json:"error_msg"`
			RequestParams []struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			} `json:"request_params"`
		} `json:"error"`
		Response []struct {
			ID int `json:"id"`
		} `json:"response"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&ans)
	if ans.Error != nil {
		if ans.Error.ErrorCode == 5 {
			return 0, myerrors.ErrOuterTokenIsWrong
		}

		return 0, err
	}
	return VKUserID(ans.Response[0].ID), nil // maybe need to add a response check
}
