package main

import (
	"context"
	"github.com/woodhds/vk.service/api/handlers"
	"github.com/woodhds/vk.service/database"
	vkService "github.com/woodhds/vk.service/internal/app/vk-service"
	"github.com/woodhds/vk.service/internal/notifier"
	"github.com/woodhds/vk.service/internal/vkclient"
	"log"
	"os"
)

var (
	token   string
	version string
	count   int
	port    int
)

func main() {
	token = os.Getenv("TOKEN")
	version = os.Getenv("VERSION")
	ParseInt(&count, 50, os.Getenv("COUNT"))
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

	app := vkService.NewApp(messageQueryService, notifyService, usersQueryService, factory, messagesService, token, version, count)

	app.Initialize()

	app.Run(port)
}
