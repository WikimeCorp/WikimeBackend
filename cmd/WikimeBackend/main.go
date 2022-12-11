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
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func setupRouter() http.Handler {
	router := mux.NewRouter()

	apiRouter := mux.NewRouter()

	// Routers
	usersRouter := apiRouter.PathPrefix("/users").Subrouter()
	currentUserRouter := usersRouter.PathPrefix("/current").Subrouter()
	animeRouter := apiRouter.PathPrefix("/anime").Subrouter()
	commentsRouter := apiRouter.PathPrefix("/comments").Subrouter()
	authRouter := apiRouter.PathPrefix("/auth/").Subrouter()

	// Users section
	usersRouter.Handle(
		"/{user_id:[0-9]+}",
		http.HandlerFunc(user.GetUserHandler()),
	).Methods("GET")
	usersRouter.HandleFunc(
		"/moderators",
		user.GetModeratorsHandler(),
	).Methods("GET")
	usersRouter.HandleFunc(
		"/admins",
		user.GetAdminsHandler(),
	).Methods("GET")
	usersRouter.Handle(
		"/{user_id:[0-9]+}/role",
		middleware.NeedAuthentication(http.HandlerFunc(user.ChangeRoleHandler())),
	).Methods("PUT")
	usersRouter.Handle(
		"/{user_id:[0-9]+}/role",
		middleware.NeedAuthentication(http.HandlerFunc(user.ResetRoleHandler())),
	).Methods("DELETE")

	// User section
	currentUserRouter.Handle(
		"",
		middleware.NeedAuthentication(http.HandlerFunc(user.GetCurrentUserHandler())),
	).Methods("GET")
	currentUserRouter.Handle(
		"/nickname",
		middleware.NeedAuthentication(http.HandlerFunc(user.ChangeNicknameHandler())),
	).Methods("PUT")
	currentUserRouter.Handle(
		"/favorites",
		middleware.NeedAuthentication(http.HandlerFunc(user.AddToFavoritesHandler())),
	).Methods("POST")
	currentUserRouter.Handle(
		"/favorites",
		middleware.NeedAuthentication(http.HandlerFunc(user.DeleteFromFavoritesHandler())),
	).Methods("DELETE")
	currentUserRouter.Handle(
		"/watched",
		middleware.NeedAuthentication(http.HandlerFunc(user.AddToWatchedHandler())),
	).Methods("POST")
	currentUserRouter.Handle(
		"/watched",
		middleware.NeedAuthentication(http.HandlerFunc(user.DeleteFromWatchedHandler())),
	).Methods("DELETE")

	// Anime section
	animeRouter.HandleFunc(
		"",
		anime.SearchHandler(),
	).Queries("search", "{search}").Methods("GET")
	animeRouter.HandleFunc(
		"",
		anime.GetAnimesHangler(),
	).Methods("GET")
	animeRouter.Handle(
		"",
		middleware.NeedAuthentication(http.HandlerFunc(anime.CreateAnimeHandler())),
	).Methods("POST")
	animeRouter.HandleFunc(
		"/{anime_id:[0-9]+}",
		anime.GetAnimeByIDHandler(),
	).Methods("GET")
	animeRouter.HandleFunc(
		"/{anime_id:[0-9]+}",
		anime.EditAnimeHandler(),
	).Methods("PUT")
	animeRouter.HandleFunc(
		"/{anime_id:[0-9]+}/average",
		anime.SetAverageEndpoint,
	).Methods("PUT")
	animeRouter.HandleFunc(
		"/list",
		anime.GetAnimeByListIDHandler(),
	).Methods("GET")
	animeRouter.HandleFunc(
		"/popular",
		anime.MostPopularHandler(),
	).Queries("count", "{count:[0-9]+}").Methods("GET")

	// Images section
	router.PathPrefix("/images/").Handler(
		http.StripPrefix(
			"/images",
			http.FileServer(http.Dir(path.Join(config.Config.ImagePathDisk, config.Config.ImagesPathURI))),
		),
	).Methods("GET")
	animeRouter.Handle(
		"/{anime_id:[0-9]+}/images",
		middleware.NeedAuthentication(http.HandlerFunc(images.AddImageHandler())),
	).Methods("POST")
	animeRouter.Handle(
		"/{anime_id:[0-9]+}/poster",
		middleware.NeedAuthentication(http.HandlerFunc(images.SetPosterHandler())),
	).Methods("POST")
	currentUserRouter.Handle(
		"/avatar",
		middleware.NeedAuthentication(http.HandlerFunc(images.SetUserImageHandler())),
	).Methods("POST")
	animeRouter.Handle(
		"/{anime_id:[0-9]+}/images/{image:.{1,100}}",
		middleware.NeedAuthentication(http.HandlerFunc(images.DeleteImageFromAnimeHandler())),
	).Methods("DELETE")

	// Comment section
	commentsRouter.Handle(
		"",
		middleware.NeedAuthentication(http.HandlerFunc(comments.CreateCommentEndpoint)),
	).Methods("POST")
	animeRouter.HandleFunc(
		"/{anime_id:[0-9]+}/comments",
		comments.GetCommentByIDEndpoint,
	).Methods("GET")
	commentsRouter.Handle(
		"/{comment_id:[0-9a-z]{24}}",
		middleware.NeedAuthentication(http.HandlerFunc(comments.DeleteCommentEndpoint)),
	).Methods("DELETE")
	commentsRouter.HandleFunc(
		"/{comment_id:[0-9a-z]{24}}",
		comments.GetCommentEndpoint,
	).Methods("GET")

	// Rating section
	animeRouter.Handle(
		"/{anime_id:[0-9]+}/rating",
		middleware.NeedAuthentication(http.HandlerFunc(rating.SetRatingHandler())),
	).Methods("POST")

	// Auth section
	authRouter.HandleFunc("/vk", auth.OAuthVkHandler()).Methods("POST")

	router.NotFoundHandler = http.HandlerFunc(other.NotFoundEndpoint)

	router.PathPrefix("/").Handler(middleware.SetJSONHeader(apiRouter))
	return router

}

func Start() error {
	config := config.Config
	router := setupRouter()

	handler := router

	addr := config.Addr + ":" + config.Port
	log.Println(addr)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "PUT"})

	err := http.ListenAndServe(addr, middleware.PrintRequestURL(handlers.CORS(headersOk, methodsOk)(handler)))
	return err
}
