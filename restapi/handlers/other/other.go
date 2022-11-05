package other

import (
	"net/http"

	apiErrors "github.com/WikimeCorp/WikimeBackend/restapi/errors"
)

func NotFoundEndpoint(w http.ResponseWriter, req *http.Request) {
	apiErrors.SetErrorInResponce(&apiErrors.ErrNotFound, w, http.StatusNotFound)
}
