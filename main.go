package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/woodhds/vk.service/database"
	"github.com/woodhds/vk.service/vkclient"

	message "github.com/woodhds/vk.service/message"

	client "github.com/woodhds/vk.service/vkclient"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

var token string
var version string
var count int
var clientId int

func main() {

	flag.StringVar(&token, "token", "", "access token required")
	flag.StringVar(&version, "version", "", "version required")
	flag.IntVar(&count, "count", 10, "used count")
	flag.IntVar(&clientId, "clientId", 0, "clientId required")
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

	database.Migrate(conn)

	r := mux.NewRouter()

	r.HandleFunc("/messages", func(rw http.ResponseWriter, r *http.Request) {

		search := r.URL.Query().Get("search")

		res, e := conn.Query(`
			SELECT messages.Id, FromId, Date, Images, LikesCount, Owner, messages.OwnerId, RepostedFrom, RepostsCount, messages.Text, UserReposted
			FROM messages inner join messages_search as search  on messages.Id = search.Id AND  messages.OwnerId = search.OwnerId 
				where search.Text MATCH @search
				order by rank desc
				`, sql.Named("search", fmt.Sprintf(`"%s"`, search)))

		if e != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		var data []message.VkMessageModel

		for res.Next() {
			m := message.VkMessageModel{}
			e := res.Scan(&m.ID, &m.FromID, &m.Date, &m.Images, &m.LikesCount, &m.Owner, &m.OwnerID, &m.RepostedFrom, &m.RepostsCount, &m.Text, &m.UserReposted)
			if e == nil {
				data = append(data, m)
			}
		}
		defer res.Close()

		json, e := json.Marshal(data)

		if e == nil {
			rw.Write(json)
		}

		fmt.Println(search)

	})

	r.HandleFunc("/grab", func(rw http.ResponseWriter, r *http.Request) {
		rows, _ := conn.Query(`select Id from VkUserModel`)

		var ids []int
		for rows.Next() {
			var id int
			err = rows.Scan(&id)
			if err == nil {
				ids = append(ids, id)
			}
		}

		rows.Close()

		var wg sync.WaitGroup

		for _, id := range ids {
			for i := 1; i <= 4; i++ {
				wg.Add(1)
				go getMessages(conn, &wg, id, i)
				if err != nil {
					fmt.Println(err)
				}
			}
		}

		wg.Wait()
	})

	r.HandleFunc("/users", func(rw http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodGet {

			if rows, e := conn.Query(`SELECT Id, Name, Avatar from VkUserModel`); e != nil {
				rw.WriteHeader(http.StatusBadRequest)
				return
			} else {
				defer rows.Close()

				var res []vkclient.VkUserMdodel

				for rows.Next() {
					u := vkclient.VkUserMdodel{}
					rows.Scan(&u.Id, &u.Name, &u.Avatar)
					res = append(res, u)
				}

				if j, e := json.Marshal(res); e != nil {
					rw.WriteHeader(http.StatusBadRequest)
				} else {
					rw.Write(j)
				}
			}
		}

		if r.Method == http.MethodPost {
			u := &client.VkUserMdodel{}
			json.NewDecoder(r.Body).Decode(u)

			if u.Id == 0 {
				rw.WriteHeader(http.StatusBadRequest)
				return
			}

			_, e := conn.Exec(`INSERT INTO VkUserModel (Id, Avatar, Name) VALUES ($1, $2, $3)`, u.Id, u.Avatar, u.Name)

			if e != nil {
				rw.WriteHeader(http.StatusBadRequest)
			}
		}

	}).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)

	r.HandleFunc("/users/search", func(rw http.ResponseWriter, r *http.Request) {

		q := r.URL.Query().Get("q")

		if q == "" {
			return
		}

		client, _ := vkclient.NewUserClient(token, version)
		response, _ := client.Search(q)

		json.NewEncoder(rw).Encode(&response)

	}).Methods(http.MethodGet, http.MethodOptions)

	http.ListenAndServe(":4222", handlers.ContentTypeHandler(handlers.CORS()(r), "application/json"))
}

func getMessages(conn *sql.DB, wg *sync.WaitGroup, id int, page int) {
	defer wg.Done()

	wallClient, _ := client.NewWallClient(token, version)

	data, e := wallClient.Get(&client.WallGetRequest{OwnerId: id, Offset: (page - 1) * count, Count: count})

	if e != nil {
		fmt.Println(e)
		return
	}

	var reposts []message.VkRepostMessage

	for _, v := range data.Response.Items {
		if len(v.CopyHistory) > 0 {
			reposts = append(reposts, message.VkRepostMessage{OwnerID: v.CopyHistory[0].OwnerID, ID: v.CopyHistory[0].ID})
		}
	}

	posts, _ := wallClient.GetById(&reposts)

	time.Sleep(time.Millisecond * 500)

	var c int
	var m *message.VkMessageModel

	if len(posts.Response.Items) > 0 {

		tran, _ := conn.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelDefault})
		for _, post := range posts.Response.Items {

			m = message.New(post, id, posts.Response.Groups)
			c = c + 1
			_, sqlErr := tran.Exec(`
			insert or ignore into messages (Id, FromId, Date, Images, LikesCount, Owner, OwnerId, RepostedFrom, RepostsCount, Text, UserReposted) 
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) 
			ON CONFLICT(id, ownerId) DO UPDATE SET LikesCount=excluded.LikesCount, RepostsCount=excluded.RepostsCount, UserReposted=excluded.UserReposted`,
				m.ID, m.FromID, time.Time(*m.Date), strings.Join(m.Images, ";"), m.LikesCount, m.Owner, m.OwnerID, m.RepostedFrom, m.RepostsCount, m.Text, m.UserReposted)

			if sqlErr != nil {
				log.Print(sqlErr)
			}
		}

		tran.Commit()
	}

	fmt.Printf("Fetched: %d\n", c)
}
