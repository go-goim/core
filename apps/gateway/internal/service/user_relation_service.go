package service

import (
	"context"
	"sync"

	"github.com/go-kratos/kratos/v2/transport/grpc"
	ggrpc "google.golang.org/grpc"

	apiresp "github.com/yusank/goim/api/transport/response"
	relationv1 "github.com/yusank/goim/api/user/relation/v1"
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

func (s *UserRelationService) AddFriend(ctx context.Context, req *relationv1.AddFriendRequest) (*relationv1.AddFriendResponse, error) {
	cc, err := s.loadConn(ctx)
	if err != nil {
		return nil, err
	}

	return relationv1.NewUserRelationServiceClient(cc).AddFriend(ctx, req)
}

func (s *UserRelationService) ListUserRelation(ctx context.Context, req *relationv1.QueryUserRelationListRequest) (
	*relationv1.QueryUserRelationListResponse, error) {

	cc, err := s.loadConn(ctx)
	if err != nil {
		return nil, err
	}

	return relationv1.NewUserRelationServiceClient(cc).QueryUserRelationList(ctx, req)
}

func (s *UserRelationService) AcceptFriend(ctx context.Context, req *relationv1.RelationPair) (*apiresp.BaseResponse, error) {
	cc, err := s.loadConn(ctx)
	if err != nil {
		return nil, err
	}

	updateReq := &relationv1.UpdateUserRelationRequest{
		RelationPair: req,
		Action:       relationv1.UpdateUserRelationAction_ACCEPT,
	}

	return relationv1.NewUserRelationServiceClient(cc).UpdateUserRelation(ctx, updateReq)
}

func (s *UserRelationService) RejectFriend(ctx context.Context, req *relationv1.RelationPair) (*apiresp.BaseResponse, error) {
	cc, err := s.loadConn(ctx)
	if err != nil {
		return nil, err
	}

	updateReq := &relationv1.UpdateUserRelationRequest{
		RelationPair: req,
		Action:       relationv1.UpdateUserRelationAction_REJECT,
	}

	return relationv1.NewUserRelationServiceClient(cc).UpdateUserRelation(ctx, updateReq)
}

func (s *UserRelationService) BlockFriend(ctx context.Context, req *relationv1.RelationPair) (*apiresp.BaseResponse, error) {
	cc, err := s.loadConn(ctx)
	if err != nil {
		return nil, err
	}

	updateReq := &relationv1.UpdateUserRelationRequest{
		RelationPair: req,
		Action:       relationv1.UpdateUserRelationAction_BLOCK,
	}

	return relationv1.NewUserRelationServiceClient(cc).UpdateUserRelation(ctx, updateReq)
}

func (s *UserRelationService) UnblockFriend(ctx context.Context, req *relationv1.RelationPair) (*apiresp.BaseResponse, error) {
	cc, err := s.loadConn(ctx)
	if err != nil {
		return nil, err
	}

	updateReq := &relationv1.UpdateUserRelationRequest{
		RelationPair: req,
		Action:       relationv1.UpdateUserRelationAction_UNBLOCK,
	}

	return relationv1.NewUserRelationServiceClient(cc).UpdateUserRelation(ctx, updateReq)
}

func (s *UserRelationService) DeleteFriend(ctx context.Context, req *relationv1.RelationPair) (*apiresp.BaseResponse, error) {
	cc, err := s.loadConn(ctx)
	if err != nil {
		return nil, err
	}

	updateReq := &relationv1.UpdateUserRelationRequest{
		RelationPair: req,
		Action:       relationv1.UpdateUserRelationAction_DELETE,
	}

	return relationv1.NewUserRelationServiceClient(cc).UpdateUserRelation(ctx, updateReq)
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
