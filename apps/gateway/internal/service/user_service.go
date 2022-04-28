package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-kratos/kratos/v2/transport/grpc"
	ggrpc "google.golang.org/grpc"

	userv1 "github.com/yusank/goim/api/user/v1"
	"github.com/yusank/goim/apps/gateway/internal/app"
	"github.com/yusank/goim/apps/gateway/internal/dao"
	"github.com/yusank/goim/pkg/conn/pool"
	"github.com/yusank/goim/pkg/conn/wrapper"
	"github.com/yusank/goim/pkg/util"
)

type UserService struct {
	userDao *dao.UserDao
}

var (
	userService     *UserService
	userServiceOnce sync.Once
)

func GetUserService() *UserService {
	userServiceOnce.Do(func() {
		userService = &UserService{
			userDao: dao.GetUserDao(),
		}
	})
	return userService
}

// Login check user login status and return user info
func (s *UserService) Login(ctx context.Context, req *userv1.UserLoginRequest) (*userv1.User, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, err
	}

	cc, err := s.loadConn(ctx)
	if err != nil {
		return nil, err
	}

	var queryReq = &userv1.QueryUserRequest{}
	switch {
	case req.GetEmail() != "":
		queryReq.User = &userv1.QueryUserRequest_Email{Email: req.GetEmail()}
	case req.GetPhone() != "":
		queryReq.User = &userv1.QueryUserRequest_Phone{Phone: req.GetPhone()}
	default:
		return nil, fmt.Errorf("invalid user login request")
	}

	user, err := userv1.NewUserServiceClient(cc).QueryUser(ctx, queryReq)
	if err != nil {
		return nil, err
	}

	if user.GetPassword() != util.Md5String(req.GetPassword()) {
		return nil, fmt.Errorf("invalid password or user not exist")
	}

	agentID, err := s.userDao.GetUserOnlineAgent(ctx, user.GetUid())
	if err != nil {
		return nil, err
	}

	if len(agentID) == 0 {
		// not login
		user.LoginStatus = userv1.LoginStatus_LOGIN
	} else {
		// already login
		user.LoginStatus = userv1.LoginStatus_ALREADY_LOGIN
		user.AgentId = &agentID
	}

	return user, nil
}

// Register register user.
func (s *UserService) Register(ctx context.Context, req *userv1.CreateUserRequest) (*userv1.User, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, err
	}

	cc, err := s.loadConn(ctx)
	if err != nil {
		return nil, err
	}

	var createReq = &userv1.CreateUserRequest{
		Name:     req.GetName(),
		Password: util.Md5String(req.GetPassword()),
	}
	switch {
	case req.GetEmail() != "":
		createReq.User = &userv1.CreateUserRequest_Email{Email: req.GetEmail()}
	case req.GetPhone() != "":
		createReq.User = &userv1.CreateUserRequest_Phone{Phone: req.GetPhone()}
	default:
		return nil, fmt.Errorf("invalid user register request")
	}

	// do check user exist and create.
	user, err := userv1.NewUserServiceClient(cc).CreateUser(ctx, createReq)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUser update user info.
func (s *UserService) UpdateUser(ctx context.Context, req *userv1.UpdateUserRequest) (*userv1.User, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, err
	}

	cc, err := s.loadConn(ctx)
	if err != nil {
		return nil, err
	}

	// do check user exist and update.
	user, err := userv1.NewUserServiceClient(cc).UpdateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) loadConn(ctx context.Context) (*ggrpc.ClientConn, error) {
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
