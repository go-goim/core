package service

import (
	"context"
	"sync"

	"github.com/go-kratos/kratos/v2/transport/grpc"
	ggrpc "google.golang.org/grpc"

	responsepb "github.com/yusank/goim/api/transport/response"
	friendpb "github.com/yusank/goim/api/user/friend/v1"
	"github.com/yusank/goim/apps/gateway/internal/app"
	"github.com/yusank/goim/pkg/conn/pool"
	"github.com/yusank/goim/pkg/conn/wrapper"
)

type UserRelationService struct {
}

var (
	userRelationService     *UserRelationService
	userRelationServiceOnce sync.Once
)

func GetUserRelationService() *UserRelationService {
	userRelationServiceOnce.Do(func() {
		userRelationService = &UserRelationService{}
	})
	return userRelationService
}

func (s *UserRelationService) AddFriend(ctx context.Context, req *friendpb.BaseFriendRequest) (*friendpb.AddFriendResponse, error) {
	cc, err := s.loadConn(ctx)
	if err != nil {
		return nil, err
	}

	return friendpb.NewFriendServiceClient(cc).AddFriend(ctx, req)
}

func (s *UserRelationService) ListUserRelation(ctx context.Context, req *friendpb.QueryFriendListRequest) (
	*friendpb.QueryFriendListResponse, error) {

	cc, err := s.loadConn(ctx)
	if err != nil {
		return nil, err
	}

	return friendpb.NewFriendServiceClient(cc).QueryFriendList(ctx, req)
}

func (s *UserRelationService) AcceptFriend(ctx context.Context, req *friendpb.BaseFriendRequest) (*responsepb.BaseResponse, error) {
	cc, err := s.loadConn(ctx)
	if err != nil {
		return nil, err
	}

	updateReq := &friendpb.ConfirmFriendRequestReq{
		Info:   req,
		Action: friendpb.ConfirmFriendRequestAction_ACCEPT,
	}

	return friendpb.NewFriendServiceClient(cc).ConfirmFriendRequest(ctx, updateReq)
}

func (s *UserRelationService) RejectFriend(ctx context.Context, req *friendpb.BaseFriendRequest) (*responsepb.BaseResponse, error) {
	cc, err := s.loadConn(ctx)
	if err != nil {
		return nil, err
	}

	updateReq := &friendpb.ConfirmFriendRequestReq{
		Info:   req,
		Action: friendpb.ConfirmFriendRequestAction_REJECT,
	}

	return friendpb.NewFriendServiceClient(cc).ConfirmFriendRequest(ctx, updateReq)
}

func (s *UserRelationService) BlockFriend(ctx context.Context, req *friendpb.BaseFriendRequest) (*responsepb.BaseResponse, error) {
	cc, err := s.loadConn(ctx)
	if err != nil {
		return nil, err
	}

	updateReq := &friendpb.UpdateFriendStatusRequest{
		Info:   req,
		Status: friendpb.FriendStatus_BLOCKED,
	}

	return friendpb.NewFriendServiceClient(cc).UpdateFriendStatus(ctx, updateReq)
}

func (s *UserRelationService) UnblockFriend(ctx context.Context, req *friendpb.BaseFriendRequest) (*responsepb.BaseResponse, error) {
	cc, err := s.loadConn(ctx)
	if err != nil {
		return nil, err
	}

	updateReq := &friendpb.UpdateFriendStatusRequest{
		Info:   req,
		Status: friendpb.FriendStatus_UNBLOCKED,
	}

	return friendpb.NewFriendServiceClient(cc).UpdateFriendStatus(ctx, updateReq)
}

func (s *UserRelationService) DeleteFriend(ctx context.Context, req *friendpb.BaseFriendRequest) (*responsepb.BaseResponse, error) {
	cc, err := s.loadConn(ctx)
	if err != nil {
		return nil, err
	}

	updateReq := &friendpb.UpdateFriendStatusRequest{
		Info:   req,
		Status: friendpb.FriendStatus_STRANGER,
	}

	return friendpb.NewFriendServiceClient(cc).UpdateFriendStatus(ctx, updateReq)
}

func (s *UserRelationService) loadConn(ctx context.Context) (*ggrpc.ClientConn, error) {
	var ck = "discovery://dc1/goim.user.service"
	c := pool.Get(ck)
	if c != nil {
		wc := c.(*wrapper.GrpcWrapper)
		return wc.ClientConn, nil
	}

	cc, err := grpc.DialInsecure(ctx,
		grpc.WithDiscovery(app.GetApplication().Register),
		grpc.WithEndpoint(ck))
	if err != nil {
		return nil, err
	}

	pool.Add(wrapper.WrapGrpc(context.Background(), cc, ck))

	return cc, nil
}
