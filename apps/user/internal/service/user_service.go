package service

import (
	"context"
	"fmt"
	"sync"

	userv1 "github.com/yusank/goim/api/user/v1"
	"github.com/yusank/goim/apps/user/internal/dao"
	"github.com/yusank/goim/apps/user/internal/data"
)

type UserService struct {
	userDao *dao.UserDao
	userv1.UnimplementedUserServiceServer
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

func (s *UserService) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	user, err := s.userDao.GetUserByUID(ctx, req.GetUid())
	if err != nil {
		return nil, err
	}

	return &userv1.GetUserResponse{
		User: user.ToProtoUser(),
	}, nil
}

// Login check user login status and return user info
func (s *UserService) Login(ctx context.Context, req *userv1.UserLoginRequest) (*userv1.User, error) {
	var (
		filter  string
		getFunc func(ctx context.Context, filter string) (*data.User, error)
	)

	switch {
	case req.GetEmail() != "":
		filter = req.GetEmail()
		getFunc = s.userDao.GetUserByEmail
	case req.GetPhone() != "":
		filter = req.GetPhone()
		getFunc = s.userDao.GetUserByPhone
	default:
		// TODO: define error
		return nil, fmt.Errorf("invalid login request")
	}

	user, err := getFunc(ctx, filter)
	if err != nil {
		return nil, err
	}

	protoUser := user.ToProtoUser()
	agentID, err := s.userDao.GetUserOnlineAgent(ctx, user.UID)
	if err != nil {
		return protoUser, err
	}

	if len(agentID) == 0 {
		// not login
		protoUser.LoginStatus = userv1.LoginStatus_LOGIN_STATUS_LOGIN
	} else {
		// already login
		protoUser.LoginStatus = userv1.LoginStatus_LOGIN_STATUS_ALREADY_LOGIN
		protoUser.AgentId = agentID
	}

	return protoUser, nil
}
