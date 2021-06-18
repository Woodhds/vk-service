package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/rs/cors"
	"github.com/woodhds/vk.service/database"
	"github.com/woodhds/vk.service/vkclient"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	message "github.com/woodhds/vk.service/message"
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

		Json(rw, data)
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

		for _, id := range ids {
			for i := 1; i <= 4; i++ {
				go getMessages(conn, id, i)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}).Methods(http.MethodGet)

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

				Json(rw, res)
			}
		}

		if r.Method == http.MethodPost {
			u := &vkclient.VkUserMdodel{}
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

	r.HandleFunc("/repost", func(rw http.ResponseWriter, r *http.Request) {

		var d []message.VkRepostMessage
		json.NewDecoder(r.Body).Decode(&d)

		go func() {

			wallClient, _ := vkclient.NewWallClient(token, version)
			groupClient, _ := vkclient.NewGroupClient(token, version)

			data, _ := wallClient.GetById(&d, "is_member")

			for _, d := range data.Response.Groups {
				if d.IsMember == 0 {
					groupClient.Join(d.ID)
				}
			}

			for _, i := range data.Response.Items {
				wallClient.Repost(&message.VkRepostMessage{OwnerID: i.OwnerID, ID: i.ID})
			}

		}()

		rw.WriteHeader(http.StatusOK)

	}).Methods(http.MethodPost, http.MethodOptions)

	r.HandleFunc("/users/search", func(rw http.ResponseWriter, r *http.Request) {

		q := r.URL.Query().Get("q")

		if q == "" {
			return
		}

		client, _ := vkclient.NewUserClient(token, version)
		response, _ := client.Search(q)

		json.NewEncoder(rw).Encode(&response)

	}).Methods(http.MethodGet, http.MethodOptions)

	http.ListenAndServe(":4222", cors.Default().Handler(r))
}

func getMessages(conn *sql.DB, id int, page int) {

	wallClient, _ := vkclient.NewWallClient(token, version)

	data, e := wallClient.Get(&vkclient.WallGetRequest{OwnerId: id, Offset: (page - 1) * count, Count: count})

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

	ch := make(chan *message.VkMessageModel)

	if len(posts.Response.Items) > 0 {
		go func() {
			var c int
			for _, post := range posts.Response.Items {
				ch <- message.New(post, id, posts.Response.Groups)
				c++
			}
			close(ch)

			fmt.Printf("Fetched: %d\n", c)
		}()

		go func() {

			for m := range ch {
				_, sqlErr := conn.Exec(`
			insert or ignore into messages (Id, FromId, Date, Images, LikesCount, Owner, OwnerId, RepostedFrom, RepostsCount, Text, UserReposted) 
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) 
			ON CONFLICT(id, ownerId) DO UPDATE SET LikesCount=excluded.LikesCount, RepostsCount=excluded.RepostsCount, UserReposted=excluded.UserReposted`,
					m.ID, m.FromID, time.Time(*m.Date), strings.Join(m.Images, ";"), m.LikesCount, m.Owner, m.OwnerID, m.RepostedFrom, m.RepostsCount, m.Text, m.UserReposted)

				if sqlErr != nil {
					log.Print(sqlErr)
				}
			}
		}()
	}
}
