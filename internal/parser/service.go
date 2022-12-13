package parser

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/woodhds/vk.service/database"
	pb "github.com/woodhds/vk.service/gen/parser"
	"github.com/woodhds/vk.service/message"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

type parserImplementation struct {
	*pb.UnimplementedParserServiceServer
	factory          database.ConnectionFactory
	messageService   VkMessagesService
	count            int
	userQueryService database.UsersQueryService
}

func (impl *parserImplementation) Parse(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	cc, _ := impl.factory.GetConnection(ctx)
	wg := sync.WaitGroup{}

	ids, _ := impl.userQueryService.GetAll()

	postsCh := make(chan []*message.VkRepostMessage)

	ch := make(chan *message.VkMessageModel)

	go func() {
		for reposts := range postsCh {
			posts := impl.messageService.GetById(reposts)
			var c int
			for _, post := range posts {
				ch <- post
				c++
			}

			fmt.Printf("Fetched: %d\n", c)
		}
		close(ch)
	}()

	go func() {
		defer cc.Close()
		for m := range ch {
			if m == nil {
				continue
			}

			sqlErr := m.Save(cc, context.Background())

			if sqlErr != nil {
				log.Print(sqlErr)
			}
		}

		log.Print("Channel closed")
	}()

	for _, id := range ids {
		for i := 1; i <= 4; i++ {
			wg.Add(1)
			go func(userId int, page int, co int, c chan []*message.VkRepostMessage) {
				defer wg.Done()
				c <- impl.messageService.GetMessages(userId, page, co)
			}(id, i, impl.count, postsCh)
		}
	}

	httpClient := &http.Client{
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
		Timeout:   time.Second * 30,
	}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go fetch(httpClient, i, postsCh, &wg)
	}

	go func() {
		wg.Wait()
		close(postsCh)
	}()

	return nil, nil
}

func fetch(httpClient *http.Client, page int, postsCh chan []*message.VkRepostMessage, wg *sync.WaitGroup) {
	defer wg.Done()
	res, err := httpClient.PostForm("https://wingri.ru/main/getPosts",
		url.Values{
			"page_num": []string{strconv.Itoa(page)},
			"our":      []string{},
			"city_id":  []string{strconv.Itoa(5)},
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

func NewParserServer(
	factory database.ConnectionFactory,
	service VkMessagesService,
	count int,
	userQueryService database.UsersQueryService) pb.ParserServiceServer {
	return &parserImplementation{
		factory:          factory,
		messageService:   service,
		count:            count,
		userQueryService: userQueryService,
	}
}
