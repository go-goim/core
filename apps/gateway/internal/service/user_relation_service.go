package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-kratos/kratos/v2/transport/grpc"
	ggrpc "google.golang.org/grpc"

	userv1 "github.com/yusank/goim/api/user/v1"
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

func (s *UserRelationService) AddFriend(ctx context.Context, req *userv1.AddFriendRequest) (*userv1.AddFriendResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, err
	}

	cc, err := s.loadConn(ctx)
	if err != nil {
		return nil, err
	}

	return userv1.NewUserRelationServiceClient(cc).AddFriend(ctx, req)
}

func (s *UserRelationService) ListUserRelation(ctx context.Context, req *userv1.QueryUserRelationListRequest) (
	*userv1.QueryUserRelationListResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, err
	}

	cc, err := s.loadConn(ctx)
	if err != nil {
		return nil, err
	}

	return userv1.NewUserRelationServiceClient(cc).QueryUserRelationList(ctx, req)
}

func (s *UserRelationService) AcceptFriend(ctx context.Context, req *userv1.AcceptFriendRequest) error {
	if err := req.ValidateAll(); err != nil {
		return err
	}

	cc, err := s.loadConn(ctx)
	if err != nil {
		return err
	}

	updateReq := &userv1.UpdateUserRelationRequest{
		Uid:       req.GetUid(),
		FriendUid: req.GetFriendUid(),
		Action:    userv1.UpdateUserRelationAction_ACCEPT,
	}

	rsp, err := userv1.NewUserRelationServiceClient(cc).UpdateUserRelation(ctx, updateReq)
	if err != nil {
		return err
	}

	if !rsp.Success {
		return fmt.Errorf("update user relation failed. detail: %s", rsp.Message)
	}

	return nil
}

func (s *UserRelationService) RejectFriend(ctx context.Context, req *userv1.RejectFriendRequest) error {
	if err := req.ValidateAll(); err != nil {
		return err
	}

	cc, err := s.loadConn(ctx)
	if err != nil {
		return err
	}

	updateReq := &userv1.UpdateUserRelationRequest{
		Uid:       req.GetUid(),
		FriendUid: req.GetFriendUid(),
		Action:    userv1.UpdateUserRelationAction_REJECT,
	}

	rsp, err := userv1.NewUserRelationServiceClient(cc).UpdateUserRelation(ctx, updateReq)
	if err != nil {
		return err
	}

	if !rsp.Success {
		return fmt.Errorf("update user relation failed. detail: %s", rsp.Message)
	}

	return nil
}

func (s *UserRelationService) BlockFriend(ctx context.Context, req *userv1.BlockFriendRequest) error {
	if err := req.ValidateAll(); err != nil {
		return err
	}

	cc, err := s.loadConn(ctx)
	if err != nil {
		return err
	}

	updateReq := &userv1.UpdateUserRelationRequest{
		Uid:       req.GetUid(),
		FriendUid: req.GetFriendUid(),
		Action:    userv1.UpdateUserRelationAction_BLOCK,
	}

	rsp, err := userv1.NewUserRelationServiceClient(cc).UpdateUserRelation(ctx, updateReq)
	if err != nil {
		return err
	}

	if !rsp.Success {
		return fmt.Errorf("update user relation failed. detail: %s", rsp.Message)
	}

	return nil
}

func (s *UserRelationService) UnblockFriend(ctx context.Context, req *userv1.UnblockFriendRequest) error {
	if err := req.ValidateAll(); err != nil {
		return err
	}

	cc, err := s.loadConn(ctx)
	if err != nil {
		return err
	}

	updateReq := &userv1.UpdateUserRelationRequest{
		Uid:       req.GetUid(),
		FriendUid: req.GetFriendUid(),
		Action:    userv1.UpdateUserRelationAction_UNBLOCK,
	}

	rsp, err := userv1.NewUserRelationServiceClient(cc).UpdateUserRelation(ctx, updateReq)
	if err != nil {
		return err
	}

	if !rsp.Success {
		return fmt.Errorf("update user relation failed. detail: %s", rsp.Message)
	}

	return nil
}

func (s *UserRelationService) DeleteFriend(ctx context.Context, req *userv1.RemoveFriendRequest) error {
	if err := req.ValidateAll(); err != nil {
		return err
	}

	cc, err := s.loadConn(ctx)
	if err != nil {
		return err
	}

	updateReq := &userv1.UpdateUserRelationRequest{
		Uid:       req.GetUid(),
		FriendUid: req.GetFriendUid(),
		Action:    userv1.UpdateUserRelationAction_DELETE,
	}

	rsp, err := userv1.NewUserRelationServiceClient(cc).UpdateUserRelation(ctx, updateReq)
	if err != nil {
		return err
	}

	if !rsp.Success {
		return fmt.Errorf("update user relation failed. detail: %s", rsp.Message)
	}

	return nil
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
