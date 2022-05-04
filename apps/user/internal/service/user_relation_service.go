package service

import (
	"context"
	"fmt"
	"sync"

	userv1 "github.com/yusank/goim/api/user/v1"
	"github.com/yusank/goim/apps/user/internal/dao"
	"github.com/yusank/goim/apps/user/internal/data"
)

// UserRelationService implements userv1.UserRelationService
type UserRelationService struct {
	userRelationDao *dao.UserRelationDao
	userDao         *dao.UserDao
	userv1.UnimplementedUserRelationServiceServer
}

var (
	_                       userv1.UserRelationServiceServer = &UserRelationService{}
	userRelationService     *UserRelationService
	userRelationServiceOnce sync.Once
)

func GetUserRelationService() *UserRelationService {
	userRelationServiceOnce.Do(func() {
		userRelationService = &UserRelationService{
			userRelationDao: dao.GetUserRelationDao(),
			userDao:         dao.GetUserDao(),
		}
	})
	return userRelationService
}

func (s *UserRelationService) AddFriend(ctx context.Context, req *userv1.AddFriendRequest) (*userv1.AddFriendResponse, error) {
	// TODO: check user exist
	friend, err := s.userDao.GetUserByUID(ctx, req.GetFriendUid())
	if err != nil {
		return nil, err
	}

	if friend == nil {
		return &userv1.AddFriendResponse{
			Status:  userv1.AddFriendStatus_NOT_FOUND,
			Message: userv1.AddFriendStatus_NOT_FOUND.String(),
		}, nil
	}

	ur, err := s.userRelationDao.GetUserRelation(ctx, req.GetUid(), req.GetFriendUid())
	if err != nil {
		return nil, err
	}

	// has no relation
	if ur == nil {
		ur = &data.UserRelation{
			UID:       req.GetUid(),
			FriendUID: req.GetFriendUid(),
			Status:    userv1.RelationStatus_REQUESTED,
		}

		if err = s.userRelationDao.AddUserRelation(ctx, ur); err != nil {
			return nil, err
		}

		return &userv1.AddFriendResponse{
			Status:       userv1.AddFriendStatus_SEND_REQUEST_SUCCESS,
			Message:      userv1.AddFriendStatus_SEND_REQUEST_SUCCESS.String(),
			UserRelation: ur.ToProtoUserRelation(),
		}, nil
	}

	switch {
	case ur.IsFriend():
		return &userv1.AddFriendResponse{
			Status:       userv1.AddFriendStatus_ALREADY_FRIEND,
			Message:      userv1.AddFriendStatus_ALREADY_FRIEND.String(),
			UserRelation: ur.ToProtoUserRelation(),
		}, nil

	case ur.IsRequested():
		return &userv1.AddFriendResponse{
			Status:       userv1.AddFriendStatus_ALREADY_REQUESTED,
			Message:      userv1.AddFriendStatus_ALREADY_REQUESTED.String(),
			UserRelation: ur.ToProtoUserRelation(),
		}, nil
		// need remove from black list first.
	case ur.IsBlocked():
		return &userv1.AddFriendResponse{
			Status:       userv1.AddFriendStatus_BLOCKED_BY_ME,
			Message:      userv1.AddFriendStatus_BLOCKED_BY_ME.String(),
			UserRelation: ur.ToProtoUserRelation(),
		}, nil
	case ur.IsStranger():
		// need update user relation to friend if the other user is friend, otherwise, just update to requested just update to requested
		ur2, err2 := s.userRelationDao.GetUserRelation(ctx, req.GetFriendUid(), req.GetUid())
		if err2 != nil {
			return nil, err2
		}

		// the other person has no relation with me
		if ur2 == nil || ur2.IsStranger() {
			ur.SetRequested()
			if err := s.userRelationDao.UpdateUserRelationStatus(ctx, ur); err != nil {
				return nil, err
			}
			return &userv1.AddFriendResponse{
				Status:       userv1.AddFriendStatus_SEND_REQUEST_SUCCESS,
				Message:      userv1.AddFriendStatus_SEND_REQUEST_SUCCESS.String(),
				UserRelation: ur.ToProtoUserRelation(),
			}, nil
		}

		// the other person has relation with me
		if ur2.IsFriend() {
			ur.SetFriend()
			if err := s.userRelationDao.UpdateUserRelationStatus(ctx, ur); err != nil {
				return nil, err
			}

			return &userv1.AddFriendResponse{
				Status:       userv1.AddFriendStatus_SUCCESS,
				Message:      userv1.AddFriendStatus_SUCCESS.String(),
				UserRelation: ur.ToProtoUserRelation(),
			}, nil
		}

		// the other person had blocked me
		if ur2.IsBlocked() {
			return &userv1.AddFriendResponse{
				Status:       userv1.AddFriendStatus_BLOCKED_BY_FRIEND,
				Message:      userv1.AddFriendStatus_BLOCKED_BY_FRIEND.String(),
				UserRelation: ur.ToProtoUserRelation(),
			}, nil
		}
	}

	return nil, fmt.Errorf("unknown user relation status: %d", ur.Status)
}

func (s *UserRelationService) GetUserRelation(ctx context.Context, req *userv1.GetUserRelationRequest) (
	*userv1.GetUserRelationResponse, error) {
	ur, err := s.userRelationDao.GetUserRelation(ctx, req.GetUid(), req.GetFriendUid())
	if err != nil {
		return nil, err
	}

	rsp := &userv1.GetUserRelationResponse{}
	if ur != nil {
		rsp.UserRelation = ur.ToProtoUserRelation()
	}

	return rsp, nil
}

func (s *UserRelationService) QueryUserRelationList(ctx context.Context, req *userv1.QueryUserRelationListRequest) (
	*userv1.QueryUserRelationListResponse, error) {
	urList, err := s.userRelationDao.GetUserRelationList(ctx, req.GetUid())
	if err != nil {
		return nil, err
	}

	rsp := &userv1.QueryUserRelationListResponse{}
	for _, ur := range urList {
		rsp.UserRelationList = append(rsp.UserRelationList, ur.ToProtoUserRelation())
	}

	return rsp, nil
}

func (s *UserRelationService) UpdateUserRelation(ctx context.Context, req *userv1.UpdateUserRelationRequest) (
	*userv1.UpdateUserRelationResponse, error) {
	ur, err := s.userRelationDao.GetUserRelation(ctx, req.GetUid(), req.GetFriendUid())
	if err != nil {
		return nil, err
	}

	if ur == nil {
		return &userv1.UpdateUserRelationResponse{
			Success: false,
			Message: "user relation not found",
		}, nil
	}

	if req.GetStatus() == ur.Status {
		// no need to update
		return &userv1.UpdateUserRelationResponse{
			Success: true,
		}, nil
	}

	switch req.GetStatus() {
	case userv1.RelationStatus_FRIEND:
		ur.SetFriend()
	case userv1.RelationStatus_STRANGER:
		ur.SetStranger()
	case userv1.RelationStatus_BLOCKED:
		ur.SetBlackList()
	default:
		return &userv1.UpdateUserRelationResponse{
			Success: false,
			Message: "unknown relation status",
		}, nil
	}

	if err := s.userRelationDao.UpdateUserRelationStatus(ctx, ur); err != nil {
		return &userv1.UpdateUserRelationResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &userv1.UpdateUserRelationResponse{
		Success: true,
	}, nil
}
