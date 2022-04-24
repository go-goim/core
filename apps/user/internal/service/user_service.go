package service

import (
	"context"
	"sync"

	userv1 "github.com/yusank/goim/api/user/v1"
	"github.com/yusank/goim/apps/user/internal/dao"
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
