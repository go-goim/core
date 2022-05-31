package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/transport/grpc"
	ggrpc "google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"

	responsepb "github.com/go-goim/goim/api/transport/response"
	friendpb "github.com/go-goim/goim/api/user/friend/v1"
	"github.com/go-goim/goim/apps/gateway/internal/app"
)

type FriendService struct {
	friendServiceConn *ggrpc.ClientConn
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

func (s *FriendService) AddFriend(ctx context.Context, req *friendpb.AddFriendRequest) (*friendpb.AddFriendResult, error) {
	err := s.checkGrpcConn(ctx)
	if err != nil {
		return nil, err
	}

	rsp, err := friendpb.NewFriendServiceClient(s.friendServiceConn).AddFriend(ctx, req)
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

	err := s.checkGrpcConn(ctx)
	if err != nil {
		return nil, err
	}

	rsp, err := friendpb.NewFriendServiceClient(s.friendServiceConn).QueryFriendList(ctx, req)
	if err != nil {
		return nil, err
	}

	if rsp.Response.Success() {
		return rsp.GetFriendList(), nil
	}

	return nil, rsp.GetResponse()
}

func (s *FriendService) AcceptFriend(ctx context.Context, req *friendpb.ConfirmFriendRequestReq) error {
	if err := req.Validate(); err != nil {
		return responsepb.NewBaseResponseWithMessage(responsepb.Code_InvalidParams, err.Error())
	}

	return s.confirmFriendRequest(ctx, req)
}

func (s *FriendService) RejectFriend(ctx context.Context, req *friendpb.ConfirmFriendRequestReq) error {
	if err := req.Validate(); err != nil {
		return responsepb.NewBaseResponseWithMessage(responsepb.Code_InvalidParams, err.Error())
	}

	return s.confirmFriendRequest(ctx, req)
}

func (s *FriendService) confirmFriendRequest(ctx context.Context, req *friendpb.ConfirmFriendRequestReq) error {
	err := s.checkGrpcConn(ctx)
	if err != nil {
		return err
	}

	rsp, err := friendpb.NewFriendServiceClient(s.friendServiceConn).ConfirmFriendRequest(ctx, req)
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
	err := s.checkGrpcConn(ctx)
	if err != nil {
		return err
	}

	updateReq := &friendpb.UpdateFriendStatusRequest{
		Info:   req,
		Status: status,
	}

	rsp, err := friendpb.NewFriendServiceClient(s.friendServiceConn).UpdateFriendStatus(ctx, updateReq)
	if err != nil {
		return err
	}

	if !rsp.Success() {
		return rsp
	}

	return nil
}

func (s *FriendService) checkGrpcConn(ctx context.Context) error {
	if s.friendServiceConn != nil {
		switch s.friendServiceConn.GetState() {
		case connectivity.Idle:
			return nil
		case connectivity.Connecting:
			return nil
		case connectivity.Ready:
			return nil
		default:
			// reconnect
		}
	}

	var ck = fmt.Sprintf("discovery://dc1/%s", app.GetApplication().Config.SrvConfig.UserService)

	cc, err := grpc.DialInsecure(ctx,
		grpc.WithDiscovery(app.GetApplication().Register),
		grpc.WithEndpoint(ck),
		grpc.WithTimeout(time.Second*5))
	if err != nil {
		return err
	}

	s.friendServiceConn = cc
	return nil
}
