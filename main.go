package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/woodhds/vk.service/api/handlers"
	"github.com/woodhds/vk.service/database"
	"github.com/woodhds/vk.service/internal/notifier"
	"github.com/woodhds/vk.service/internal/predictor"
	"github.com/woodhds/vk.service/internal/vkclient"
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
	ParseInt(&count, 50, os.Getenv("COUNT"))
	host = os.Getenv("HOST")
	ParseInt(&port, 4222, os.Getenv("PORT"))

	if token == "" {
		panic("access token required")
	}

	log.Printf("Used token: %s", token)
	log.Printf("Used version: %s", version)
	log.Printf("Used count: %d", count)

	connectionString := os.Getenv("DATABASE_URL")
	factory, err := database.NewConnectionFactory(connectionString)

	if err != nil {
		log.Fatal(err)
		return
	}

	predictorClient, _ := predictor.NewClient(host)
	notifyService := notifier.NewNotifyService()
	messageQueryService := database.NewMessageQueryService(factory)
	wallClient, _ := vkclient.NewWallClient(token, version)
	usersQueryService, _ := database.NewUserQueryService(factory)
	messagesService := handlers.NewMessageService(wallClient)

	conn, _ := factory.GetConnection(context.Background())

	database.Migrate(conn)

	if e := conn.Close(); e != nil {
		log.Println(e)
	}

	router := mux.NewRouter()
	r := router.PathPrefix("/api").Subrouter()

	r.Path("/messages").Handler(handlers.MessagesHandler(messageQueryService, predictorClient)).Methods(http.MethodGet)
	r.Path("/like").Handler(handlers.LikeHandler(notifyService)).Methods(http.MethodPost)

	r.Path("/grab").Handler(handlers.ParserHandler(factory, messagesService, count, notifyService, usersQueryService)).Methods(http.MethodGet)

	r.Path("/users").Handler(handlers.UsersHandler(usersQueryService, notifyService)).Methods(http.MethodGet, http.MethodPost, http.MethodOptions, http.MethodDelete)

	r.Path("/repost").Handler(handlers.RepostHandler(factory, token, version)).Methods(http.MethodPost, http.MethodOptions)

	r.Path("/users/search").Handler(handlers.UsersSearchHandler(token, version)).Methods(http.MethodGet, http.MethodOptions)
	r.Path("/messages/{ownerId:-?[0-9]+}/{id:[0-9]+}").Handler(handlers.MessageSaveHandler(predictorClient, factory)).Methods(http.MethodPost)
	r.Path("/notifications").Handler(notifier.NotificationHandler(notifyService)).Methods(http.MethodGet)

	fmt.Println(http.ListenAndServe(fmt.Sprintf(":%d", port), cors.Default().Handler(r)))
}
