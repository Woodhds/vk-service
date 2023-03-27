package main

import (
	"context"
	"github.com/woodhds/vk.service/database"
	vkService "github.com/woodhds/vk.service/internal/app/vk-service"
	"github.com/woodhds/vk.service/internal/parser"
	sweeper "github.com/woodhds/vk.service/internal/sweeper"
	"github.com/woodhds/vk.service/internal/vkclient"
	"log"
	"os"
	"os/signal"
	"time"
)

var (
	token            string
	version          string
	count            int
	port             int
	connectionString string
)

func main() {
	token = os.Getenv("TOKEN")
	version = os.Getenv("VERSION")
	ParseInt(&count, 50, os.Getenv("COUNT"))
	ParseInt(&port, 4222, os.Getenv("PORT"))
	connectionString = os.Getenv("CONNECTION_STRING")

	if token == "" {
		panic("access token required")
	}

	log.Printf("Used token: %s", token)
	log.Printf("Used version: %s", version)
	log.Printf("Used count: %d", count)

	if connectionString == "" {
		connectionString = ":memory?Pooling=True&MaxPoolSize=100&Cache=shared"
	}
	factory, err := database.NewConnectionFactory(connectionString)

	if err != nil {
		log.Fatal(err)
		return
	}

	messageQueryService := database.NewMessageQueryService(factory)
	usersQueryService := database.NewUserQueryService(factory)
	groupsQueryService := database.NewGroupsQueryService(factory)
	wallClient, _ := vkclient.NewWallClient(token, version)
	messagesService := parser.NewMessageService(wallClient)

	conn, _ := factory.GetConnection(context.Background())

	database.Migrate(conn)
	deleteCtx, cl := context.WithTimeout(context.TODO(), time.Second*15)
	defer cl()

	sw := sweeper.NewSweeper(factory)
	go sw.Run(deleteCtx)

	if e := conn.Close(); e != nil {
		log.Println(e)
	}

	app := vkService.NewApp(messageQueryService, usersQueryService, groupsQueryService, factory, messagesService, token, version, count)

	app.Initialize()

	go app.Run(port)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	app.Stop(ctx)
	os.Exit(0)
}
