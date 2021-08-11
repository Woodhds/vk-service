package handlers

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/woodhds/vk.service/message"
	"github.com/woodhds/vk.service/notifier"
	"github.com/woodhds/vk.service/vkclient"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func ParserHandler(conn *sql.DB, token string, version string, count int, notifier *notifier.NotifyService) http.Handler {
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

		postsCh := make(chan []*message.VkRepostMessage, 10)

		ch := make(chan *message.VkMessageModel)

		wallClient, _ := vkclient.NewWallClient(token, version)

		go func() {
			for reposts := range postsCh {
				posts, _ := wallClient.GetById(reposts)
				var c int
				for _, post := range posts.Response.Items {
					ch <- message.New(post, posts.Response.Groups)
					c++
				}

				fmt.Printf("Fetched: %d\n", c)
			}
		}()

		go func() {

			for m := range ch {
				_, sqlErr := conn.Exec(`
			insert into messages (Id, FromId, Date, Images, LikesCount, Owner, OwnerId, RepostsCount, Text, UserReposted) 
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) 
			ON CONFLICT(id, ownerId) DO UPDATE SET LikesCount=excluded.LikesCount, RepostsCount=excluded.RepostsCount, UserReposted=excluded.UserReposted`,
					m.ID, m.FromID, time.Time(*m.Date), strings.Join(m.Images, ";"), m.LikesCount, m.Owner, m.OwnerID, m.RepostsCount, m.Text, m.UserReposted)

				if sqlErr != nil {
					log.Print(sqlErr)
				}
			}
		}()

		for _, id := range ids {
			for i := 1; i <= 4; i++ {
				go getMessages(id, i, count, wallClient, postsCh)
				if err != nil {
					fmt.Println(err)
				}
			}
		}

		httpClient := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}

		for i := 0; i < 10; i++ {
			go fetch(httpClient, i, postsCh)
		}
	})
}

func getMessages(id int, page int, count int, wallClient *vkclient.WallClient, postsCh chan []*message.VkRepostMessage) {

	data, e := wallClient.Get(&vkclient.WallGetRequest{OwnerId: id, Offset: (page - 1) * count, Count: count})

	if e != nil {
		fmt.Println(e)
		return
	}

	var messages []*message.VkRepostMessage

	for _, v := range data.Response.Items {
		if len(v.CopyHistory) > 0 {
			messages = append(messages, &message.VkRepostMessage{OwnerID: v.CopyHistory[0].OwnerID, ID: v.CopyHistory[0].ID})
		}
	}

	postsCh <- messages
}

func fetch(httpClient *http.Client, page int, postsCh chan []*message.VkRepostMessage) {
	res, err := httpClient.PostForm("https://wingri.ru/main/getPosts",
		url.Values{
			"page_num": []string{strconv.Itoa(page)},
			"our":      []string{},
			"city_id":  []string{strconv.Itoa(24)},
		},
	)

	if err != nil {
		log.Print(err)
	}

	doc, _ := goquery.NewDocumentFromReader(res.Body)
	var reposts []*message.VkRepostMessage
	doc.Find(".grid-item .post_container .post_footer a").Each(func(i int, selection *goquery.Selection) {
		attr, exists := selection.Attr("href")
		if exists {
			arr := strings.Split(strings.Replace(attr, "https://vk.com/wall", "", 1), "_")
			owner, ownerE := strconv.Atoi(arr[0])
			id, idE := strconv.Atoi(arr[1])

			if ownerE == nil && idE == nil {
				reposts = append(reposts, &message.VkRepostMessage{
					OwnerID: owner,
					ID:      id,
				})
			}
		}
	})

	postsCh <- reposts
}
