package service

import (
	"context"
	"fmt"
	"sync"

	userv1 "github.com/yusank/goim/api/user/v1"
	"github.com/yusank/goim/apps/user/internal/dao"
	"github.com/yusank/goim/apps/user/internal/data"
	"github.com/yusank/goim/pkg/util"
)

// UserService implements userv1.UserServiceServer
type UserService struct {
	userDao *dao.UserDao
	userv1.UnimplementedUserServiceServer
}

var (
	_               userv1.UserServiceServer = &UserService{}
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

func (s *UserService) GetUser(ctx context.Context, req *userv1.GetUserInfoRequest) (*userv1.User, error) {
	user, err := s.userDao.GetUserByUID(ctx, req.GetUid())
	if err != nil {
		return nil, err
	}

	if user == nil || user.IsDeleted() {
		return nil, fmt.Errorf("user not found")
	}

	return user.ToProtoUser(), nil
}

func (s *UserService) QueryUser(ctx context.Context, req *userv1.QueryUserRequest) (*userv1.User, error) {
	user, err := s.loadUserByEmailOrPhone(ctx, req.GetEmail(), req.GetPhone())
	if err != nil {
		return nil, err
	}

	if user == nil || user.IsDeleted() {
		return nil, fmt.Errorf("user not found")
	}

	return user.ToProtoUser(), nil
}

func (s *UserService) loadUserByEmailOrPhone(ctx context.Context, email, phone string) (*data.User, error) {
	var (
		value   string
		getFunc func(ctx context.Context, v string) (*data.User, error)
	)

	switch {
	case email != "":
		value = email
		getFunc = s.userDao.GetUserByEmail
	case phone != "":
		value = phone
		getFunc = s.userDao.GetUserByPhone
	default:
		return nil, fmt.Errorf("invalid query user request")
	}

	user, err := getFunc(ctx, value)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) CreateUser(ctx context.Context, req *userv1.CreateUserRequest) (*userv1.User, error) {
	user, err := s.loadUserByEmailOrPhone(ctx, req.GetEmail(), req.GetPhone())
	if err != nil {
		return nil, err
	}

	if user == nil {
		user = &data.User{
			UID:      util.UUID(),
			Email:    req.GetEmail(),
			Phone:    req.GetPhone(),
			Password: req.GetPassword(),
		}

		err = s.userDao.CreateUser(ctx, user)
		if err != nil {
			return nil, err
		}

		return user.ToProtoUser(), nil
	}

	// user exists
	if user.IsDeleted() {
		// undo delete
		// 这里会出现被删除用户 undo 后使用的是旧密码的情况,需要更新密码
		user.Password = req.GetPassword()
		err = s.userDao.UndoDelete(ctx, user)
		if err != nil {
			return nil, err
		}

		return user.ToProtoUser(), nil
	}

	return nil, fmt.Errorf("user exists")

}

func (s *UserService) UpdateUser(ctx context.Context, req *userv1.UpdateUserRequest) (*userv1.User, error) {
	user, err := s.userDao.GetUserByUID(ctx, req.GetUid())
	if err != nil {
		return nil, err
	}

	if user == nil || user.IsDeleted() {
		return nil, fmt.Errorf("user not found")
	}

	if req.GetEmail() != "" {
		user.Email = req.GetEmail()
	}

	if req.GetPhone() != "" {
		user.Phone = req.GetPhone()
	}

	if req.GetName() != "" {
		user.Password = req.GetName()
	}

	if req.GetAvatar() != "" {
		user.Avatar = req.GetAvatar()
	}

	err = s.userDao.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return user.ToProtoUser(), nil
}
