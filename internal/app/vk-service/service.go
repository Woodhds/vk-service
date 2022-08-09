package vk_service

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"github.com/woodhds/vk.service/api/handlers"
	"github.com/woodhds/vk.service/database"
	vkMessages "github.com/woodhds/vk.service/gen/messages"
	"github.com/woodhds/vk.service/internal/notifier"
	"github.com/woodhds/vk.service/internal/predictor"
	"net/http"
)

type App struct {
	router              *mux.Router
	predictor           predictor.Predictor
	messageQueryService database.MessagesQueryService
	notifyService       *notifier.NotifyService
	usersQueryService   database.UsersQueryService
	token               string
	version             string
	count               int
	factory             database.ConnectionFactory
	messagesService     handlers.VkMessagesService
	grpcMux             *runtime.ServeMux
}

func (app *App) Initialize() {
	router := mux.NewRouter()
	app.router = router.PathPrefix("/api").Subrouter()
	app.grpcMux = runtime.NewServeMux()

	app.initializeRoutes()

}

func (app *App) Run(port int) {
	fmt.Println(http.ListenAndServe(fmt.Sprintf(":%d", port), cors.Default().Handler(app.router)))
}

func NewApp(
	predictor predictor.Predictor,
	messageQueryService database.MessagesQueryService,
	notifyService *notifier.NotifyService,
	usersQueryService database.UsersQueryService,
	factory database.ConnectionFactory,
	messagesService handlers.VkMessagesService,
	token string,
	version string,
	count int) *App {
	return &App{
		router:              nil,
		predictor:           predictor,
		messageQueryService: messageQueryService,
		notifyService:       notifyService,
		usersQueryService:   usersQueryService,
		token:               token,
		version:             version,
		count:               count,
		factory:             factory,
		messagesService:     messagesService,
	}
}

func (app *App) initializeRoutes() {
	app.router.Path("/like").Handler(handlers.LikeHandler(app.notifyService)).Methods(http.MethodPost)
	app.router.Path("/grab").Handler(handlers.ParserHandler(app.factory, app.messagesService, app.count, app.notifyService, app.usersQueryService)).Methods(http.MethodGet)
	app.router.Path("/users").Handler(handlers.UsersHandler(app.usersQueryService, app.notifyService)).Methods(http.MethodGet, http.MethodPost, http.MethodOptions, http.MethodDelete)
	app.router.Path("/repost").Handler(handlers.RepostHandler(app.factory, app.token, app.version)).Methods(http.MethodPost, http.MethodOptions)
	app.router.Path("/users/search").Handler(handlers.UsersSearchHandler(app.token, app.version)).Methods(http.MethodGet, http.MethodOptions)
	app.router.Path("/messages/{ownerId:-?[0-9]+}/{id:[0-9]+}").Handler(handlers.MessageSaveHandler(app.predictor, app.factory)).Methods(http.MethodPost)
	app.router.Path("/notifications").Handler(notifier.NotificationHandler(app.notifyService)).Methods(http.MethodGet)

	vkMessages.RegisterMessagesServiceHandlerServer(context.Background(), app.grpcMux, handlers.NewMessageHandler(app.messageQueryService, app.predictor))
	app.router.PathPrefix("").Handler(app.grpcMux)
}
