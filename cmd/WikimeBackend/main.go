package wikimebackend

import (
	"log"
	"net/http"

	"github.com/WikimeCorp/WikimeBackend/config"
	"github.com/WikimeCorp/WikimeBackend/restapi/handlers/anime"
	"github.com/WikimeCorp/WikimeBackend/restapi/handlers/auth"
	"github.com/WikimeCorp/WikimeBackend/restapi/handlers/other"
	"github.com/WikimeCorp/WikimeBackend/restapi/handlers/user"

	"github.com/WikimeCorp/WikimeBackend/restapi/middleware"
	"github.com/gorilla/mux"
)

func setupRouter() *mux.Router {
	router := mux.NewRouter()

	// User section
	userRouter := router.PathPrefix("/user/").Subrouter()
	userRouter.Handle("/{user_id:[0-9]+}",
		http.HandlerFunc(user.GetUserHandler()),
	).Methods("GET")
	userRouter.Handle("/",
		middleware.NeedAuthorization(http.HandlerFunc(user.GetCurrentUserHandler())),
	).Methods("GET")

	// Anime section
	animeRouter := router.PathPrefix("/anime/").Subrouter()
	animeRouter.HandleFunc("/{anime_id:[0-9]+}",
		anime.GetAnimeByIDHandler(),
	).Methods("GET")

	// Auth section
	authRouter := router.PathPrefix("/auth/").Subrouter()
	authRouter.HandleFunc("/vk", auth.OAuthVkHandler()).Methods("GET")

	router.NotFoundHandler = http.HandlerFunc(other.NotFoundEndpoint)

	return router

}

func Start() error {
	config := config.Config
	router := setupRouter()

	handler := middleware.SetJSONHeader(router)

	addr := config.Addr + ":" + config.Port
	log.Println(addr)
	err := http.ListenAndServe(addr, handler)
	return err
}
