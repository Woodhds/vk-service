package users

import (
	"context"
	"errors"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/woodhds/vk.service/database"
	pb "github.com/woodhds/vk.service/gen/users"
	vkClient "github.com/woodhds/vk.service/internal/vkclient"
	"google.golang.org/protobuf/types/known/emptypb"
)

type usersImplementation struct {
	pb.UnimplementedUsersServiceServer
	usersService database.UsersQueryService
	token        string
	version      string
	userClient   *vkClient.UserClient
}

func (i *usersImplementation) GetUsers(ctx context.Context, _ *emptypb.Empty) (*pb.GetUsersResponse, error) {
	if rows, e := i.usersService.GetFullUsers(ctx); e != nil {
		return &pb.GetUsersResponse{Users: toModel(rows)}, nil
	} else {
		return nil, e
	}
}

func (i *usersImplementation) Add(ctx context.Context, user *pb.VkUserProto) (*emptypb.Empty, error) {
	if user == nil || user.Id == 0 {
		return nil, errors.New("body is empty")
	}

	if e := i.usersService.InsertNew(int(user.GetId()), user.Name, user.Avatar.String(), ctx); e != nil {
		return nil, e
	}

	return nil, nil
}

func (i *usersImplementation) Delete(ctx context.Context, request *pb.DeleteUserRequest) (*emptypb.Empty, error) {
	e := i.usersService.Delete(int(request.Id), ctx)

	if e != nil {
		return nil, e
	}

	return nil, nil
}

func (i *usersImplementation) Search(ctx context.Context, request *pb.UserSearchRequest) (*pb.UserSearchResponse, error) {

	if request == nil || request.Search.GetValue() == "" {
		return nil, errors.New("search query is empty")
	}

	client, _ := vkClient.NewUserClient(i.token, i.version)

	if response, e := client.Search(request.Search.GetValue()); e != nil {
		return nil, e
	} else {
		return &pb.UserSearchResponse{Users: toModel(response)}, nil
	}
}

func NewUsersHandler(
	messagesQueryService database.UsersQueryService,
	token string,
	version string) pb.UsersServiceServer {
	return &usersImplementation{
		usersService: messagesQueryService,
		token:        token,
		version:      version,
	}
}

func toModel(rows []*vkClient.VkUserModel) []*pb.VkUserProto {
	result := make([]*pb.VkUserProto, len(rows), len(rows))

	for i := 0; i < len(rows); i++ {
		avatar, _ := runtime.StringValue(rows[i].Avatar)
		result[i] = &pb.VkUserProto{
			Id:     int32(rows[i].Id),
			Name:   rows[i].Name,
			Avatar: avatar,
		}
	}

	return result
}
