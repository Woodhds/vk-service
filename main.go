package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/rs/cors"
	"github.com/woodhds/vk.service/config"
	"github.com/woodhds/vk.service/database"
	"github.com/woodhds/vk.service/handlers"
	"github.com/woodhds/vk.service/notifier"
	"github.com/woodhds/vk.service/predictor"
	"github.com/woodhds/vk.service/vkclient"
	"log"
	"net/http"
	"os"
)

var (
	token   string
	version string
	count   int
	host    string
	port    int
)

func main() {
	token = os.Getenv("TOKEN")
	version = os.Getenv("VERSION")
	config.ParseInt(&count, 50, os.Getenv("COUNT"))
	host = os.Getenv("HOST")
	config.ParseInt(&port, 4222, os.Getenv("PORT"))

	if token == "" {
		panic("access token required")
	}

	log.Printf("Used token: %s", token)
	log.Printf("Used version: %s", version)
	log.Printf("Used count: %d", count)

	conn, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatal(err)
		return
	}

	predictorClient, _ := predictor.NewClient(host)
	notifyService := notifier.NewNotifyService()
	messageQueryService := database.NewMessageQueryService(conn)
	wallClient, _ := vkclient.NewWallClient(token, version)

	defer conn.Close()

	database.Migrate(conn)

	router := mux.NewRouter()
	r := router.PathPrefix("/api").Subrouter()

	r.Path("/messages").Handler(handlers.MessagesHandler(messageQueryService, predictorClient)).Methods(http.MethodGet)
	r.Path("/like").Handler(handlers.LikeHandler(notifyService)).Methods(http.MethodPost)

	r.Path("/grab").Handler(handlers.ParserHandler(conn, wallClient, count, notifyService)).Methods(http.MethodGet)

	r.Path("/users").Handler(handlers.UsersHandler(conn)).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)

	r.Path("/repost").Handler(handlers.RepostHandler(conn, token, version)).Methods(http.MethodPost, http.MethodOptions)

	r.Path("/users/search").Handler(handlers.UsersSearchHandler(token, version)).Methods(http.MethodGet, http.MethodOptions)
	r.Path("/messages/{ownerId:-?[0-9]+}/{id:[0-9]+}").Handler(handlers.MessageSaveHandler(predictorClient, conn)).Methods(http.MethodPost)
	r.Path("/notifications").Handler(notifier.NotificationHandler(notifyService)).Methods(http.MethodGet)
	r.Path("/predict/{ownerId:-?[0-9]+}/{id:[0-9]+}").Handler(handlers.PredictHandler(predictorClient, messageQueryService, wallClient)).Methods(http.MethodPost)

	fmt.Println(http.ListenAndServe(fmt.Sprintf(":%d", port), cors.Default().Handler(r)))
}
