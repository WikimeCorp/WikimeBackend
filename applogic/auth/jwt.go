package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/WikimeCorp/WikimeBackend/config"
	"github.com/WikimeCorp/WikimeBackend/types"
	. "github.com/WikimeCorp/WikimeBackend/types"
	"github.com/WikimeCorp/WikimeBackend/types/myerrors"
)

var (
	hashAlgName = "HS256"
	hashNewFunc = sha256.New

	jwtLifeTime, _ = strconv.Atoi(config.Config.JWTLifeTime)
)

var headerJWTBytes, _ = json.Marshal(struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}{hashAlgName, "JWT"})
var headerJWTBase64 = base64.RawURLEncoding.EncodeToString(headerJWTBytes)

var secretKeyForHash = []byte(config.Config.SecretKeyForHash)

func hashFunc(bytes []byte) []byte {
	h := hmac.New(hashNewFunc, secretKeyForHash)
	h.Write(bytes)
	return h.Sum(nil)
}

func generateJWT(userID UserID) string {
	payload := JWTPayload{
		UserID: userID,
		Exp:    time.Now().UTC().Add(time.Duration(jwtLifeTime) * time.Second).Unix(),
	}

	payloadBytes, _ := json.Marshal(payload)
	payloadBase64 := base64.RawURLEncoding.EncodeToString(payloadBytes)

	tmp := headerJWTBase64 + "." + payloadBase64

	return tmp + "." + base64.RawURLEncoding.EncodeToString(hashFunc([]byte(tmp)))
}

func checkValidJWTSignature(jwtToken string) (bool, error) {
	strs := strings.Split(jwtToken, ".")
	if len(strs) != 3 {
		return false, myerrors.ErrIncorrectPartsNumberOfJWTToken
	}
	header := strs[0]
	payload := strs[1]
	signature := strs[2]

	hash := hashFunc([]byte(header + "." + payload))
	if base64.RawURLEncoding.EncodeToString(hash) != signature {
		return false, nil
	}

	return true, nil
}

func generateRefreshToken(userID UserID) (string, error) {
	return "", nil
}

func checkExp(payload *types.JWTPayload) bool {
	nowTime := time.Now().UTC().Unix()
	if payload.Exp < nowTime {
		return false
	}
	return true
}

// CheckJWTValid checks the JWT token for validity
func CheckJWTValid(jwtToken string) (*types.JWTPayload, error) {
	goodSignature, err := checkValidJWTSignature(jwtToken)
	if err != nil || goodSignature == false {
		return nil, myerrors.ErrJWTTokenInvalidSignature
	}

	payloadStrInBase64 := strings.Split(jwtToken, ".")[1]
	payloadBytes, err := base64.RawURLEncoding.DecodeString(payloadStrInBase64)
	payload := types.JWTPayload{}

	_ = json.Unmarshal(payloadBytes, &payload)

	if checkExp(&payload) == false {
		return nil, myerrors.ErrJWTTokenTimeout
	}

	return &payload, nil

}
