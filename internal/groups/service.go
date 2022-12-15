package groups

import (
	"context"
	pb "github.com/woodhds/vk.service/gen/groups"
	"google.golang.org/protobuf/types/known/emptypb"
)

type groupsImplementation struct {
	*pb.UnimplementedGroupsServiceServer
}

func (gr *groupsImplementation) AddFavorite(ctx context.Context, request *pb.AddFavoriteGroupRequest) (*emptypb.Empty, error) {
	return nil, nil
}

func (gr *groupsImplementation) RemoveGroupFromFavorite(ctx context.Context, request *pb.RemoveGroupFromFavoriteRequest) (*emptypb.Empty, error) {
	return nil, nil
}

func (gr *groupsImplementation) GetFavorites(ctx context.Context, request *pb.GetFavoritesRequest) (*pb.GetFavoriteResponse, error) {
	return nil, nil
}
