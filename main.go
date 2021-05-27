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
	"user-fetcher/database"
	message "user-fetcher/message"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

var token string
var version string
var count int
var clientId int

type VkUserMdodel struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

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
		rw.Header().Add("Access-control-allow-origin", "*")
		rw.Header().Add("Access-control-allow-method", "*")
		rw.Header().Add("Access-control-allow-headers", "*")

		if r.Method == http.MethodOptions {
			return
		}

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
			rw.Header().Add("Content-type", "application/json")
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
	})

	r.HandleFunc("/users", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Access-control-allow-origin", "*")
		rw.Header().Add("Access-control-allow-method", "*")
		rw.Header().Add("Access-control-allow-headers", "*")

		if r.Method == http.MethodOptions {
			return
		}

		if rows, e := conn.Query(`SELECT Id, Name, Avatar from VkUserModel`); e != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		} else {
			defer rows.Close()

			var res []VkUserMdodel

			for rows.Next() {
				u := VkUserMdodel{}
				rows.Scan(&u.Id, &u.Name, &u.Avatar)
				res = append(res, u)
			}

			if j, e := json.Marshal(res); e != nil {
				rw.WriteHeader(http.StatusBadRequest)
			} else {
				rw.Write(j)
				rw.Header().Add("Content-type", "application/json")
			}
		}
	})

	http.ListenAndServe(":4222", r)
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
	mutex.Unlock()
	fmt.Printf("Fetched: %d\n", c)
}
