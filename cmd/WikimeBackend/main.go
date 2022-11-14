package wikimebackend

import (
	"log"
	"net/http"
	"path"

	"github.com/WikimeCorp/WikimeBackend/config"
	"github.com/WikimeCorp/WikimeBackend/restapi/handlers/anime"
	"github.com/WikimeCorp/WikimeBackend/restapi/handlers/auth"
	"github.com/WikimeCorp/WikimeBackend/restapi/handlers/images"
	"github.com/WikimeCorp/WikimeBackend/restapi/handlers/other"
	"github.com/WikimeCorp/WikimeBackend/restapi/handlers/user"

	"github.com/WikimeCorp/WikimeBackend/restapi/middleware"
	"github.com/gorilla/mux"
)

func setupRouter() *mux.Router {
	router := mux.NewRouter()

	apiRouter := mux.NewRouter()

	// User section
	userRouter := apiRouter.PathPrefix("/user/").Subrouter()
	userRouter.Handle("/{user_id:[0-9]+}",
		http.HandlerFunc(user.GetUserHandler()),
	).Methods("GET")
	userRouter.Handle("/",
		middleware.NeedAuthorization(http.HandlerFunc(user.GetCurrentUserHandler())),
	).Methods("GET")

	// Anime section
	animeRouter := apiRouter.PathPrefix("/anime").Subrouter()
	animeRouter.HandleFunc(
		"",
		anime.GetAnimesHangler(),
	).Methods("GET")
	animeRouter.HandleFunc(
		"/{anime_id:[0-9]+}",
		anime.GetAnimeByIDHandler(),
	).Methods("GET")
	animeRouter.HandleFunc(
		"/list",
		anime.GetAnimeByListIDHandler(),
	).Methods("GET")
	animeRouter.HandleFunc(
		"",
		anime.CreateAnimeHandler(),
	).Methods("POST") // Add auth check
	animeRouter.HandleFunc(
		"/{anime_id:[0-9]+}",
		anime.SetAverageEndpoint,
	).Methods("PUT")

	// Images section
	router.PathPrefix("/images/").Handler(
		http.StripPrefix(
			"/images",
			http.FileServer(http.Dir(path.Join(config.Config.ImagePathDisk, config.Config.ImagesPathURI))),
		),
	).Methods("GET")
	animeRouter.HandleFunc(
		"/{anime_id:[0-9]+}/image",
		images.AddImageHandler(),
	).Methods("POST")
	animeRouter.HandleFunc(
		"/{anime_id:[0-9]+}/poster",
		images.SetPosterHandler(),
	).Methods("POST")

	// Auth section
	authRouter := apiRouter.PathPrefix("/auth/").Subrouter()
	authRouter.HandleFunc("/vk", auth.OAuthVkHandler()).Methods("POST")

	apiRouter.NotFoundHandler = http.HandlerFunc(other.NotFoundEndpoint)

	router.PathPrefix("/").Handler(middleware.SetJSONHeader(apiRouter))
	return router

}

func Start() error {
	config := config.Config
	router := setupRouter()

	handler := router

	addr := config.Addr + ":" + config.Port
	log.Println(addr)
	err := http.ListenAndServe(addr, handler)
	return err
}
