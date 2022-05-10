package service

import (
	"context"
	"sync"

	"github.com/go-kratos/kratos/v2/transport/grpc"
	ggrpc "google.golang.org/grpc"

	friendpb "github.com/yusank/goim/api/user/friend/v1"
	"github.com/yusank/goim/apps/gateway/internal/app"
	"github.com/yusank/goim/pkg/conn/pool"
	"github.com/yusank/goim/pkg/conn/wrapper"
)

type FriendService struct {
}

var (
	userRelationService     *FriendService
	userRelationServiceOnce sync.Once
)

func GetUserRelationService() *FriendService {
	userRelationServiceOnce.Do(func() {
		userRelationService = &FriendService{}
	})
	return userRelationService
}

func (s *FriendService) AddFriend(ctx context.Context, req *friendpb.BaseFriendRequest) (*friendpb.AddFriendResult, error) {
	cc, err := s.loadConn(ctx)
	if err != nil {
		return nil, err
	}

	rsp, err := friendpb.NewFriendServiceClient(cc).AddFriend(ctx, req)
	if err != nil {
		return nil, err
	}

	if !rsp.Response.Success() {
		return nil, rsp.GetResponse()
	}

	return rsp.GetResult(), nil
}

func (s *FriendService) ListUserRelation(ctx context.Context, req *friendpb.QueryFriendListRequest) (
	[]*friendpb.Friend, error) {

	cc, err := s.loadConn(ctx)
	if err != nil {
		return nil, err
	}

	rsp, err := friendpb.NewFriendServiceClient(cc).QueryFriendList(ctx, req)
	if err != nil {
		return nil, err
	}

	if rsp.Response.Success() {
		return rsp.GetFriendList(), nil
	}

	return nil, rsp.GetResponse()
}

func (s *FriendService) AcceptFriend(ctx context.Context, req *friendpb.BaseFriendRequest) error {
	return s.confirmFriendRequest(ctx, req, friendpb.ConfirmFriendRequestAction_ACCEPT)
}

func (s *FriendService) RejectFriend(ctx context.Context, req *friendpb.BaseFriendRequest) error {
	return s.confirmFriendRequest(ctx, req, friendpb.ConfirmFriendRequestAction_REJECT)
}

func (s *FriendService) confirmFriendRequest(ctx context.Context, req *friendpb.BaseFriendRequest,
	action friendpb.ConfirmFriendRequestAction) error {
	cc, err := s.loadConn(ctx)
	if err != nil {
		return err
	}

	updateReq := &friendpb.ConfirmFriendRequestReq{
		Info:   req,
		Action: action,
	}

	rsp, err := friendpb.NewFriendServiceClient(cc).ConfirmFriendRequest(ctx, updateReq)
	if err != nil {
		return err
	}

	if !rsp.Success() {
		return rsp
	}

	return nil
}

func (s *FriendService) BlockFriend(ctx context.Context, req *friendpb.BaseFriendRequest) error {
	return s.updateFriendStatus(ctx, req, friendpb.FriendStatus_BLOCKED)
}

func (s *FriendService) UnblockFriend(ctx context.Context, req *friendpb.BaseFriendRequest) error {
	return s.updateFriendStatus(ctx, req, friendpb.FriendStatus_UNBLOCKED)
}

func (s *FriendService) DeleteFriend(ctx context.Context, req *friendpb.BaseFriendRequest) error {
	return s.updateFriendStatus(ctx, req, friendpb.FriendStatus_STRANGER)
}

func (s *FriendService) updateFriendStatus(ctx context.Context, req *friendpb.BaseFriendRequest, status friendpb.FriendStatus) error { // nolint: lll
	cc, err := s.loadConn(ctx)
	if err != nil {
		return err
	}

	updateReq := &friendpb.UpdateFriendStatusRequest{
		Info:   req,
		Status: status,
	}

	rsp, err := friendpb.NewFriendServiceClient(cc).UpdateFriendStatus(ctx, updateReq)
	if err != nil {
		return err
	}

	if !rsp.Success() {
		return rsp
	}

	return nil
}

func (s *FriendService) loadConn(ctx context.Context) (*ggrpc.ClientConn, error) {
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
