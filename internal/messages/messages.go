package messages

import (
	"context"
	"github.com/woodhds/vk.service/database"
	pb "github.com/woodhds/vk.service/gen/messages"
	vkClient "github.com/woodhds/vk.service/internal/vkclient"
	"github.com/woodhds/vk.service/message"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type messagesImplementation struct {
	pb.UnimplementedMessagesServiceServer
	messagesQueryService database.MessagesQueryService
	token                string
	version              string
	factory              database.ConnectionFactory
}

func (m *messagesImplementation) GetMessages(ctx context.Context, r *pb.GetMessagesRequest) (*pb.GetMessagesResponse, error) {
	search := r.GetSearch()

	data, e := m.messagesQueryService.GetMessages(search.GetValue(), ctx)

	if e != nil {
		return nil, e
	}

	response := mapToResponse(data)

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

	for _, d := range data.Response.Groups {
		if d.IsMember == 0 {
			groupClient.Join(d.ID)
		}
	}

	conn, _ := m.factory.GetConnection(ctx)
	defer conn.Close()
	for _, i := range data.Response.Items {
		if e := wallClient.Repost(&message.VkRepostMessage{OwnerID: i.OwnerID, ID: i.ID}); e == nil {
			if _, e := conn.ExecContext(ctx, "UPDATE messages SET UserReposted = true where Id = $1 and OwnerId = $2", i.ID, i.OwnerID); e != nil {
				return nil, e
			}
		} else {
			return nil, e
		}
	}

	return nil, nil
}

func mapToResponse(data []*message.VkCategorizedMessageModel) []*pb.VkMessageExt {
	n := len(data)
	res := make([]*pb.VkMessageExt, n, n)

	for i := 0; i < n; i++ {
		res[i] = &pb.VkMessageExt{
			Id:           int32(data[i].ID),
			FromId:       int32(data[i].FromID),
			Date:         timestamppb.New(time.Time(*data[i].Date)),
			Images:       data[i].Images,
			LikesCount:   int32(data[i].LikesCount),
			Owner:        data[i].Owner,
			OwnerId:      int32(data[i].OwnerID),
			RepostsCount: int32(data[i].RepostsCount),
			Text:         data[i].Text,
			UserReposted: data[i].UserReposted,
			Category:     "",
			IsAccept:     false,
		}
	}

	return res
}

func NewMessageHandler(
	messagesQueryService database.MessagesQueryService,
	token string,
	version string,
	factory database.ConnectionFactory) pb.MessagesServiceServer {
	return &messagesImplementation{
		messagesQueryService: messagesQueryService,
		token:                token,
		version:              version,
		factory:              factory,
	}
}
