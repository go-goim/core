package service

import (
	"context"
	"fmt"
	"sync"

	apiresp "github.com/yusank/goim/api/transport/response"
	relationv1 "github.com/yusank/goim/api/user/relation/v1"
	"github.com/yusank/goim/apps/user/internal/dao"
	"github.com/yusank/goim/apps/user/internal/data"
)

// UserRelationService implements relationv1.UserRelationService
type UserRelationService struct {
	userRelationDao *dao.UserRelationDao
	userDao         *dao.UserDao
	relationv1.UnimplementedUserRelationServiceServer
}

var (
	_                       relationv1.UserRelationServiceServer = &UserRelationService{}
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

func (s *UserRelationService) AddFriend(ctx context.Context, req *relationv1.AddFriendRequest) (*relationv1.AddFriendResponse, error) {
	friend, err := s.userDao.GetUserByUID(ctx, req.GetFriendUid())
	if err != nil {
		return nil, err
	}

	rsp := &relationv1.AddFriendResponse{
		Response: apiresp.OK,
	}

	if friend == nil {
		rsp.Response = apiresp.ErrUserNotFound
		return rsp, nil
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
			Status:    relationv1.RelationStatus_REQUESTED,
		}

		if err = s.userRelationDao.AddUserRelation(ctx, ur); err != nil {
			return nil, err
		}

		rsp.Status = relationv1.AddFriendStatus_SEND_REQUEST_SUCCESS
		rsp.UserRelation = ur.ToProtoUserRelation()
		return rsp, nil
	}

	rsp.UserRelation = ur.ToProtoUserRelation()
	switch {
	case ur.IsFriend():
		rsp.Status = relationv1.AddFriendStatus_ALREADY_FRIEND
		return rsp, nil
	case ur.IsRequested():
		rsp.Status = relationv1.AddFriendStatus_ALREADY_REQUESTED
		return rsp, nil
		// need remove from black list first.
	case ur.IsBlocked():
		rsp.Status = relationv1.AddFriendStatus_BLOCKED_BY_ME
		return rsp, nil
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
			rsp.Status = relationv1.AddFriendStatus_SEND_REQUEST_SUCCESS
			return rsp, nil
		}

		// the other person has relation with me
		if ur2.IsFriend() {
			ur.SetFriend()
			if err := s.userRelationDao.UpdateUserRelationStatus(ctx, ur); err != nil {
				return nil, err
			}

			rsp.Status = relationv1.AddFriendStatus_SUCCESS
			return rsp, nil
		}

		// the other person had blocked me
		if ur2.IsBlocked() {
			rsp.Status = relationv1.AddFriendStatus_BLOCKED_BY_FRIEND
			return rsp, nil
		}
	}

	rsp.Status = relationv1.AddFriendStatus_FAILED
	rsp.Response = apiresp.ErrUnknown.SetMsg(fmt.Sprintf("unknown status: %d", ur.Status))
	return rsp, nil
}

func (s *UserRelationService) GetUserRelation(ctx context.Context, req *relationv1.GetUserRelationRequest) (
	*relationv1.GetUserRelationResponse, error) {
	ur, err := s.userRelationDao.GetUserRelation(ctx, req.GetUid(), req.GetFriendUid())
	if err != nil {
		return nil, err
	}

	rsp := &relationv1.GetUserRelationResponse{}
	if ur != nil {
		rsp.UserRelation = ur.ToProtoUserRelation()
	}

	return rsp, nil
}

func (s *UserRelationService) QueryUserRelationList(ctx context.Context, req *relationv1.QueryUserRelationListRequest) (
	*relationv1.QueryUserRelationListResponse, error) {
	urList, err := s.userRelationDao.GetUserRelationList(ctx, req.GetUid())
	if err != nil {
		return nil, err
	}

	var (
		rsp           = &relationv1.QueryUserRelationListResponse{}
		friendUIDList = make([]string, len(urList))
		friendMap     = make(map[string]*data.User)
	)
	for i, ur := range urList {
		rsp.UserRelationList = append(rsp.UserRelationList, ur.ToProtoUserRelation())
		friendUIDList[i] = ur.FriendUID
	}

	// get friend info
	friendInfoList, err := s.userDao.ListUsers(ctx, friendUIDList...)
	if err != nil {
		return nil, err
	}

	for i, friendInfo := range friendInfoList {
		friendMap[friendInfo.UID] = friendInfoList[i]
	}

	for _, ur := range rsp.UserRelationList {
		if friendInfo, ok := friendMap[ur.FriendUid]; ok {
			ur.FriendName = friendInfo.Name
			ur.FriendAvatar = friendInfo.Avatar
		}
	}

	return rsp, nil
}

func (s *UserRelationService) UpdateUserRelation(ctx context.Context, req *relationv1.UpdateUserRelationRequest) (
	*apiresp.BaseResponse, error) {
	pair := req.GetRelationPair()
	ur, err := s.userRelationDao.GetUserRelation(ctx, pair.GetUid(), pair.GetFriendUid())
	if err != nil {
		return nil, err
	}

	if ur == nil {
		return apiresp.ErrRelationNotFound, nil
	}

	newStatus, valid := req.GetAction().CheckActionAndGetNewStatus(ur.Status)
	if !valid {
		return apiresp.ErrInvalidUpdateRelationAction, nil
	}

	ur.SetStatus(newStatus)
	if err := s.userRelationDao.UpdateUserRelationStatus(ctx, ur); err != nil {
		return apiresp.ErrUnknown.SetMsg(err.Error()), nil
	}

	return apiresp.OK, nil
}
