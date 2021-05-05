package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
	message "user-fetcher/message"

	_ "github.com/mattn/go-sqlite3"
)

var token string
var version string
var count int

func migrate(conn *sql.DB) {
	log.Println("Start Migrate")
	createUserStm := `CREATE TABLE IF NOT EXISTS VkUserModel 
	(
		Id INTEGER,
		PRIMARY KEY(Id)
	)
	`
	log.Println("Create VkUserModel")
	log.Println(createUserStm)
	_, createRes := conn.Exec(createUserStm)
	if createRes != nil {
		log.Fatal(createRes)
	}

	log.Println("Created VkUserModel")

	createUserStm = `
	CREATE TABLE IF NOT EXISTS messages(
		Id Integer,
		FromId Integer,
		Date DateTime,
		Images TExt,
		LikesCount integer,
		Owner Text,
		OwnerId Integer,
		RepostedFrom integer,
		RepostsCount Integer,
		Text text,
		Primary Key(Id, OwnerId) 
		)`
	_, crecreateRes := conn.Exec(createUserStm)

	if crecreateRes != nil {
		log.Fatal(createRes)
	}

	log.Println("Create Fulltext search table")

	createUserStm = `
	CREATE VIRTUAL TABLE IF NOT EXISTS messages_search USING fts5(Id, OwnerId, Text)
	`

	_, crecreateRes = conn.Exec(createUserStm)

	if crecreateRes != nil {
		log.Fatal("Error occured during creating virtual table: ", createUserStm)
		panic(crecreateRes)
	}

	log.Println("Created messages")

	log.Println("Create on create trigger for full text search")
	createUserStm = `
	CREATE TRIGGER IF NOT EXISTS TR_messages_AI AFTER INSERT on messages
	BEGIN
		INSERT INTO messages_search (Id, OwnerId, Text) VALUES (new.Id, new.OwnerId, new.Text);
	END;
	`
	_, crecreateRes = conn.Exec(createUserStm)

	if crecreateRes != nil {
		log.Fatalln("Error creating TRIGGER on message: ", crecreateRes)
		panic(crecreateRes)
	}

	log.Println("Trigger created")

	log.Println("Stop migrate")
}

func main() {

	flag.StringVar(&token, "token", "", "access token required")
	flag.StringVar(&version, "version", "", "version required")
	flag.IntVar(&count, "count", 10, "used count")
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

	defer conn.Close()

	migrate(conn)
	rows, _ := conn.Query(`select Id from VkUserModel`)

	var ids []int
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err == nil {
			ids = append(ids, id)
		}
	}

	var mutex sync.Mutex
	var wg sync.WaitGroup

	httpClient := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}

	for _, id := range ids {
		for i := 1; i <= 4; i++ {
			wg.Add(1)
			go getMessages(&mutex, conn, httpClient, &wg, id, i)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	wg.Wait()
}

func getMessages(mutex *sync.Mutex, conn *sql.DB, httpClient *http.Client, wg *sync.WaitGroup, id int, page int) {
	defer wg.Done()
	mutex.Lock()

	url := fmt.Sprintf(`https://api.vk.com/method/wall.get?owner_id=%d&offset=%d&filter=owner&count=%d&extended=1&access_token=%s&v=%s`,
		id, (page-1)*count, count, token, version)

	fmt.Printf("Request URL: %s", url)

	resp, e := httpClient.Get(url)

	if e != nil {
		fmt.Println(e)
		return
	}

	if resp == nil {
		fmt.Println("RESPONSE NULL")
		return
	}

	defer resp.Body.Close()

	var data message.VkWallResponse
	err := json.NewDecoder(resp.Body).Decode(&data)

	if err != nil {
		fmt.Println(err)
	}

	time.Sleep(time.Millisecond * 500)

	var yt bytes.Buffer
	for _, dataItem := range data.Response.Items {
		if len(dataItem.CopyHistory) > 0 {
			yt.WriteString(fmt.Sprintf("%d_%d,", dataItem.CopyHistory[0].OwnerID, dataItem.CopyHistory[0].ID))
		}
	}

	url = fmt.Sprintf(`https://api.vk.com/method/wall.getById?posts=%s&extended=1&access_token=%s&v=%s`, yt.String(), token, version)

	resp, _ = httpClient.Get(url)

	var posts message.VkResponse

	e = json.NewDecoder(resp.Body).Decode(&posts)
	if e != nil {
		fmt.Println(e)
	}

	time.Sleep(time.Millisecond * 500)

	var c int
	var m *message.VkMessageModel

	if len(posts.Response.Items) > 0 {

		tran, _ := conn.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelDefault})
		for _, post := range posts.Response.Items {

			m = vkMessageModel(post, id, posts.Response.Groups)

			c = c + 1

			sqlResult, sqlErr := tran.Exec(`insert or ignore into messages (Id, FromId, Date, Images, LikesCount, Owner, OwnerId, RepostedFrom, RepostsCount, Text) 
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
				m.ID, m.FromID, time.Time(*m.Date), strings.Join(m.Images, ";"), m.LikesCount, m.Owner, m.OwnerID, m.RepostedFrom, m.RepostsCount, m.Text)

			if sqlErr != nil {
				log.Print(sqlErr)
			} else {
				rowAffected, _ := sqlResult.RowsAffected()
				log.Printf("Row affected: %d", rowAffected)
			}
		}
		tran.Commit()
	}
	mutex.Unlock()
	fmt.Printf("Fetched: %d\n", c)
}

func vkMessageModel(post *message.VkMessage, id int, groups []*message.VkGroup) *message.VkMessageModel {
	model := &message.VkMessageModel{
		ID:           post.ID,
		FromID:       post.FromID,
		Date:         post.Date,
		Images:       []string{},
		LikesCount:   post.Likes.Count,
		Owner:        "",
		OwnerID:      post.OwnerID,
		RepostedFrom: id,
		RepostsCount: post.Reposts.Count,
		Text:         post.Text,
	}

	for _, i := range post.Attachments {
		if len(i.Photo.Sizes) > 2 {
			model.Images = append(model.Images, i.Photo.Sizes[3].Url)
		}
	}

	for _, g := range groups {
		if g.ID == -post.OwnerID {
			model.Owner = g.Name
		}
	}

	return model
}
