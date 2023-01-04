package parser

import (
	"context"
	"fmt"
	"github.com/woodhds/vk.service/database"
	pb "github.com/woodhds/vk.service/gen/parser"
	"github.com/woodhds/vk.service/message"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"sync"
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
			if posts == nil {
				continue
			}
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

	go func() {
		wg.Wait()
		close(postsCh)
	}()

	return nil, nil
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
