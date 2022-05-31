package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/transport/grpc"
	ggrpc "google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"

	responsepb "github.com/yusank/goim/api/transport/response"
	userv1 "github.com/yusank/goim/api/user/v1"
	"github.com/yusank/goim/apps/gateway/internal/app"
	"github.com/yusank/goim/apps/gateway/internal/dao"
	"github.com/yusank/goim/pkg/log"
	"github.com/yusank/goim/pkg/util"
)

type UserService struct {
	userDao         *dao.UserDao
	userServiceConn *ggrpc.ClientConn
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

func (s *UserService) QueryUserInfo(ctx context.Context, req *userv1.QueryUserRequest) (*userv1.User, error) {
	err := s.checkGrpcConn(ctx)
	if err != nil {
		return nil, err
	}

	rsp, err := userv1.NewUserServiceClient(s.userServiceConn).QueryUser(ctx, req)
	if err != nil {
		return nil, err
	}

	if !rsp.GetResponse().Success() {
		return nil, rsp.GetResponse()
	}

	return rsp.GetUser().ToUser(), nil
}

// Login check user login status and return user info
func (s *UserService) Login(ctx context.Context, req *userv1.UserLoginRequest) (*userv1.User, error) {
	err := s.checkGrpcConn(ctx)
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

	ddl, ok := ctx.Deadline()
	if ok {
		log.Debug("Login ctx deadline", "ddl", ddl)
	}

	ctx2, cancel := context.WithTimeout(ctx, 3*time.Second)
	rsp, err := userv1.NewUserServiceClient(s.userServiceConn).QueryUser(ctx2, queryReq)
	cancel()

	if err != nil {
		return nil, err
	}

	if !rsp.GetResponse().Success() {
		return nil, rsp.GetResponse()
	}

	user := rsp.GetUser()

	if user.GetPassword() != util.HashString(req.GetPassword()) {
		return nil, responsepb.Code_InvalidUsernameOrPassword.BaseResponse()
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

	return user.ToUser(), nil
}

// Register register user.
func (s *UserService) Register(ctx context.Context, req *userv1.CreateUserRequest) (*userv1.User, error) {
	err := s.checkGrpcConn(ctx)
	if err != nil {
		return nil, err
	}

	// do check user exist and create.
	rsp, err := userv1.NewUserServiceClient(s.userServiceConn).CreateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	if !rsp.GetResponse().Success() {
		return nil, rsp.GetResponse()
	}

	return rsp.GetUser().ToUser(), nil
}

// UpdateUser update user info.
func (s *UserService) UpdateUser(ctx context.Context, req *userv1.UpdateUserRequest) (*userv1.User, error) {
	err := s.checkGrpcConn(ctx)
	if err != nil {
		return nil, err
	}

	// do check user exist and update.
	rsp, err := userv1.NewUserServiceClient(s.userServiceConn).UpdateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	if !rsp.GetResponse().Success() {
		return nil, rsp.GetResponse()
	}

	return rsp.GetUser().ToUser(), nil
}

func (s *UserService) checkGrpcConn(ctx context.Context) error {
	if s.userServiceConn != nil {
		switch s.userServiceConn.GetState() {
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
		grpc.WithTimeout(5*time.Second))
	if err != nil {
		return err
	}

	s.userServiceConn = cc
	return nil
}
