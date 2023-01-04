package groups

import (
	"context"
	"errors"
	"github.com/woodhds/vk.service/database"
	pb "github.com/woodhds/vk.service/gen/groups"
	"github.com/woodhds/vk.service/internal/vkclient"
	"google.golang.org/protobuf/types/known/emptypb"
)

type groupsImplementation struct {
	*pb.UnimplementedGroupsServiceServer
	queryService database.GroupsQueryService
	token        string
	version      string
}

func (gr *groupsImplementation) AddFavorite(ctx context.Context, request *pb.AddFavoriteGroupRequest) (*emptypb.Empty, error) {
	if len(request.Ids) == 0 {
		return nil, errors.New("empty list for add")
	}

	groupsClient, _ := vkclient.NewGroupClient(gr.token, gr.version)
	res := make([]int, len(request.Ids), len(request.Ids))

	for i := 0; i < len(request.Ids); i++ {
		res[i] = int(request.Ids[i])
	}

	groups, _ := groupsClient.Get(res)

	for i := 0; i < len(groups); i++ {
		if e := gr.queryService.Add(groups[i].ID, groups[i].Name, "", ctx); e != nil {
			return nil, e
		}
	}

	return nil, nil
}

func (gr *groupsImplementation) RemoveGroupFromFavorite(ctx context.Context, request *pb.RemoveGroupFromFavoriteRequest) (*emptypb.Empty, error) {
	for i := 0; i < len(request.Ids); i++ {
		if e := gr.queryService.Remove(int(request.Ids[i]), ctx); e != nil {
			return nil, e
		}
	}

	return nil, nil
}

func (gr *groupsImplementation) GetFavorites(ctx context.Context, request *pb.GetFavoritesRequest) (*pb.GetFavoriteResponse, error) {
	return nil, nil
}

func NewGroupsServer(queryService database.GroupsQueryService) pb.GroupsServiceServer {
	return &groupsImplementation{queryService: queryService}
}
