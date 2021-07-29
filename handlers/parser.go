package handlers

import (
	"database/sql"
	"fmt"
	"github.com/woodhds/vk.service/message"
	"github.com/woodhds/vk.service/notifier"
	"github.com/woodhds/vk.service/vkclient"
	"log"
	"net/http"
	"strings"
	"time"
)

func ParserHandler(conn *sql.DB, token string, version string, count int, notifier notifier.Notifier) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		defer func() {
			notifier.Success("Grab start")
		}()

		rows, _ := conn.Query(`select Id from VkUserModel`)
		var err error

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
				go getMessages(conn, id, i, token, version, count)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	})
}

func getMessages(conn *sql.DB, id int, page int, token string, version string, count int) {

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
