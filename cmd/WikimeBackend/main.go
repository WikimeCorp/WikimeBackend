package wikimebackend

import (
	"log"
	"net/http"
	"path"

	"github.com/WikimeCorp/WikimeBackend/config"
	"github.com/WikimeCorp/WikimeBackend/restapi/handlers/anime"
	"github.com/WikimeCorp/WikimeBackend/restapi/handlers/auth"
	"github.com/WikimeCorp/WikimeBackend/restapi/handlers/comments"
	"github.com/WikimeCorp/WikimeBackend/restapi/handlers/images"
	"github.com/WikimeCorp/WikimeBackend/restapi/handlers/other"
	"github.com/WikimeCorp/WikimeBackend/restapi/handlers/rating"
	"github.com/WikimeCorp/WikimeBackend/restapi/handlers/user"

	"github.com/WikimeCorp/WikimeBackend/restapi/middleware"
	"github.com/gorilla/mux"
)

func setupRouter() *mux.Router {
	router := mux.NewRouter()

	apiRouter := mux.NewRouter()

	// Routers
	userRouter := apiRouter.PathPrefix("/user").Subrouter()
	animeRouter := apiRouter.PathPrefix("/anime").Subrouter()
	commentsRouter := apiRouter.PathPrefix("/comments").Subrouter()
	authRouter := apiRouter.PathPrefix("/auth/").Subrouter()

	// User section
	userRouter.Handle(
		"/{user_id:[0-9]+}",
		http.HandlerFunc(user.GetUserHandler()),
	).Methods("GET")
	userRouter.Handle(
		"",
		middleware.NeedAuthentication(http.HandlerFunc(user.GetCurrentUserHandler())),
	).Methods("GET")
	userRouter.Handle(
		"/nickname",
		middleware.NeedAuthentication(http.HandlerFunc(user.ChangeNicknameHandler())),
	).Methods("PUT")
	userRouter.Handle(
		"/favorites",
		middleware.NeedAuthentication(http.HandlerFunc(user.AddToFavoritesHandler())),
	).Methods("POST")
	userRouter.Handle(
		"/favorites",
		middleware.NeedAuthentication(http.HandlerFunc(user.DeleteFromFavoritesHandler())),
	).Methods("DELETE")
	userRouter.Handle(
		"/watched",
		middleware.NeedAuthentication(http.HandlerFunc(user.AddToWatchedHandler())),
	).Methods("POST")
	userRouter.Handle(
		"/watched",
		middleware.NeedAuthentication(http.HandlerFunc(user.DeleteFromWatchedHandler())),
	).Methods("DELETE")

	// Anime section
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
		"/popular",
		anime.MostPopularHandler(),
	)
	animeRouter.Handle(
		"",
		http.HandlerFunc(anime.CreateAnimeHandler()),
	).Methods("POST")
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

	// Comment section
	commentsRouter.HandleFunc(
		"",
		comments.CreateAnimeEndpoint,
	).Methods("POST")
	animeRouter.HandleFunc(
		"/comments/{anime_id:[0-9]+}",
		comments.GetCommentByIDEndpoint,
	).Methods("GET")
	commentsRouter.HandleFunc(
		"/{comment_id:[0-9a-z]{24}}",
		comments.DeleteCommentEndpoint,
	).Methods("DELETE")

	// Rating section
	animeRouter.Handle(
		"/{anime_id:[0-9]+}/rating",
		middleware.NeedAuthentication(http.HandlerFunc(rating.SetRatingHandler())),
	).Methods("POST")

	// Auth section
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
