package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
	"github.com/woodhds/vk.service/database"
	"github.com/woodhds/vk.service/handlers"
	"github.com/woodhds/vk.service/notifier"
	"github.com/woodhds/vk.service/predictor"
	"log"
	"net/http"
)

var (
	token   string
	version string
	count   int
	host    string
)

func main() {

	flag.StringVar(&token, "token", "", "access token required")
	flag.StringVar(&version, "version", "", "version required")
	flag.IntVar(&count, "count", 10, "used count")
	flag.StringVar(&host, "host", "", "host save not setup")
	flag.Parse()

	if token == "" {
		panic("access token required")
	}

	log.Printf("Used token: %s", token)
	log.Printf("Used version: %s", version)
	log.Printf("Used count: %d", count)

	conn, err := sql.Open("sqlite3", "./data.db")

	if err != nil {
		log.Fatal(err)
		return
	}

	predictorClient, _ := predictor.NewClient(host)
	notifyService := notifier.NewNotifyService()
	messageQueryService := database.NewMessageQueryService(conn)

	defer conn.Close()

	database.Migrate(conn)

	router := mux.NewRouter()
	r := router.PathPrefix("/api").Subrouter()

	r.Path("/messages").Handler(handlers.MessagesHandler(messageQueryService, predictorClient)).Methods(http.MethodGet)
	r.Path("/like").Handler(handlers.LikeHandler(notifyService)).Methods(http.MethodPost)

	r.Path("/grab").Handler(handlers.ParserHandler(conn, token, version, count, notifyService)).Methods(http.MethodGet)

	r.Path("/users").Handler(handlers.UsersHandler(conn)).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)

	r.Path("/repost").Handler(handlers.RepostHandler(conn, token, version)).Methods(http.MethodPost, http.MethodOptions)

	r.Path("/users/search").Handler(handlers.UsersSearchHandler(token, version)).Methods(http.MethodGet, http.MethodOptions)
	r.Path("/messages/{ownerId:-?[0-9]+}/{id:[0-9]+}").Handler(handlers.MessageSaveHandler(predictorClient, conn)).Methods(http.MethodPost)
	r.Path("/notifications").Handler(notifier.NotificationHandler(notifyService)).Methods(http.MethodGet)

	fmt.Println(http.ListenAndServe(":4222", cors.Default().Handler(r)))
}
