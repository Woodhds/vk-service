package handlers

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/woodhds/vk.service/database"
	"github.com/woodhds/vk.service/internal/notifier"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/woodhds/vk.service/message"
)

func ParserHandler(factory database.ConnectionFactory, messageService VkMessagesService, count int, notifier *notifier.NotifyService, userQueryService database.UsersQueryService) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		defer func() {
			notifier.Success("Grab start")
		}()

		ids, _ := userQueryService.GetAll()

		postsCh := make(chan []*message.VkRepostMessage)

		ch := make(chan *message.VkMessageModel)

		go func() {
			for reposts := range postsCh {
				posts := messageService.GetById(reposts)
				var c int
				for _, post := range posts {
					ch <- post
					c++
				}

				fmt.Printf("Fetched: %d\n", c)
			}
		}()

		go func() {
			conn, _ := factory.GetConnection(context.Background())
			defer conn.Close()
			statement, _ := conn.PrepareContext(context.Background(), `
			insert into messages (Id, FromId, Date, Images, LikesCount, Owner, OwnerId, RepostsCount, Text, UserReposted) 
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) 
			ON CONFLICT(id, ownerId) DO UPDATE SET LikesCount=excluded.LikesCount, RepostsCount=excluded.RepostsCount, UserReposted=excluded.UserReposted, Images=excluded.Images`)

			for m := range ch {
				if m == nil {
					continue
				}

				_, sqlErr := statement.ExecContext(context.Background(),
					m.ID, m.FromID, time.Time(*m.Date), strings.Join(m.Images, ";"), m.LikesCount, m.Owner, m.OwnerID, m.RepostsCount, m.Text, m.UserReposted)

				if sqlErr != nil {
					log.Print(sqlErr)
				}
			}
		}()

		for _, id := range ids {
			for i := 1; i <= 4; i++ {
				go func(page int, co int, c chan []*message.VkRepostMessage) {
					c <- messageService.GetMessages(id, page, co)
				}(i, count, postsCh)
			}
		}

		httpClient := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}

		for i := 0; i < 10; i++ {
			go fetch(httpClient, i, postsCh)
		}
	})
}

func fetch(httpClient *http.Client, page int, postsCh chan []*message.VkRepostMessage) {
	res, err := httpClient.PostForm("https://wingri.ru/main/getPosts",
		url.Values{
			"page_num": []string{strconv.Itoa(page)},
			"our":      []string{},
			"city_id":  []string{strconv.Itoa(97)},
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