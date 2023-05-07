package messages

import (
	"context"
	"github.com/woodhds/vk.service/database"
	pb "github.com/woodhds/vk.service/gen/messages"
	"github.com/woodhds/vk.service/internal/parser"
	vkClient "github.com/woodhds/vk.service/internal/vkclient"
	"github.com/woodhds/vk.service/message"
	"google.golang.org/protobuf/types/known/emptypb"
)

type messagesImplementation struct {
	*pb.UnimplementedMessagesServiceServer
	messagesQueryService database.MessagesQueryService
	token                string
	version              string
	factory              database.ConnectionFactory
	messageService       parser.VkMessagesService
}

func (m *messagesImplementation) GetMessages(ctx context.Context, r *pb.GetMessagesRequest) (*pb.GetMessagesResponse, error) {
	search := r.GetSearch()

	data, e := m.messagesQueryService.GetMessages(search.GetValue(), ctx)

	if e != nil {
		return nil, e
	}

	getById := make([]*message.VkRepostMessage, len(data), len(data))
	for i, datum := range data {
		getById[i] = &message.VkRepostMessage{
			OwnerID: datum.OwnerID,
			ID:      datum.ID,
		}
	}

	messages := m.messageService.GetById(getById)

	response := mapToResponse(messages)

	return &pb.GetMessagesResponse{
		Messages: response,
	}, nil
}

func (m *messagesImplementation) Repost(ctx context.Context, request *pb.RepostMessageRequest) (*emptypb.Empty, error) {
	wallClient, _ := vkClient.NewWallClient(m.token, m.version)
	groupClient, _ := vkClient.NewGroupClient(m.token, m.version)

	getByIdRequest := make([]*message.VkRepostMessage, len(request.Messages))

	for i := 0; i < len(request.Messages); i++ {
		getByIdRequest[i] = &message.VkRepostMessage{
			OwnerID: int(request.Messages[i].OwnerId),
			ID:      int(request.Messages[i].Id),
		}
	}

	data, _ := wallClient.GetById(getByIdRequest, "is_member")

	for _, d := range data.Groups {
		if d.IsMember == 0 {
			groupClient.Join(d.ID)
		}
	}

	conn, _ := m.factory.GetConnection(ctx)
	defer conn.Close()
	for _, i := range data.Items {
		e := wallClient.Repost(&message.VkRepostMessage{OwnerID: i.OwnerID, ID: i.ID})
		return nil, e
	}

	return nil, nil
}

func mapToResponse(data []*message.VkMessageModel) []*pb.VkMessageExt {
	n := len(data)
	res := make([]*pb.VkMessageExt, n, n)

	for i := 0; i < n; i++ {
		res[i] = data[i].ToDto()
	}

	return res
}

func NewMessageHandler(
	messagesQueryService database.MessagesQueryService,
	token string,
	version string,
	factory database.ConnectionFactory,
	messageService parser.VkMessagesService) pb.MessagesServiceServer {
	return &messagesImplementation{
		messagesQueryService: messagesQueryService,
		token:                token,
		version:              version,
		factory:              factory,
		messageService:       messageService,
	}
}
