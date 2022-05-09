package service

import (
	"context"
	"sync"

	responsepb "github.com/yusank/goim/api/transport/response"
	friendpb "github.com/yusank/goim/api/user/friend/v1"
	"github.com/yusank/goim/apps/user/internal/dao"
	"github.com/yusank/goim/apps/user/internal/data"
)

// FriendService implements friendpb.FriendServiceServer
type FriendService struct {
	friendDao        *dao.FriendDao
	friendRequestDao *dao.FriendRequestDao
	userDao          *dao.UserDao
	friendpb.UnimplementedFriendServiceServer
}

var (
	_                 friendpb.FriendServiceServer = &FriendService{}
	friendService     *FriendService
	friendServiceOnce sync.Once
)

func GetUserRelationService() *FriendService {
	friendServiceOnce.Do(func() {
		friendService = &FriendService{
			friendDao:        dao.GetUserRelationDao(),
			friendRequestDao: dao.GetFriendRequestDao(),
			userDao:          dao.GetUserDao(),
		}
	})
	return friendService
}

func (s *FriendService) AddFriend(ctx context.Context, req *friendpb.BaseFriendRequest) (*friendpb.AddFriendResponse, error) {
	friendUser, err := s.userDao.GetUserByUID(ctx, req.GetFriendUid())
	if err != nil {
		return nil, err
	}

	rsp := &friendpb.AddFriendResponse{
		Response: responsepb.OK,
	}

	if friendUser == nil {
		rsp.Response = responsepb.ErrUserNotFound
		return rsp, nil
	}

	me, err := s.friendDao.GetFriend(ctx, req.GetUid(), req.GetFriendUid())
	if err != nil {
		return nil, err
	}

	friend, err := s.friendDao.GetFriend(ctx, req.GetFriendUid(), req.GetUid())
	if err != nil {
		return nil, err
	}

	// friend had blocked me or me had blocked friend
	if !s.canAddFriend(ctx, me, friend, rsp) {
		return rsp, nil
	}

	ok, err := s.canAddAutomatically(ctx, me, friend, rsp)
	if err != nil {
		return nil, err
	}

	// had added friend
	if ok {
		return rsp, nil
	}

	// send friend request
	err = s.sendFriendRequest(ctx, req, me, friend, rsp)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *FriendService) canAddFriend(ctx context.Context, me, friend *data.Friend, rsp *friendpb.AddFriendResponse) bool {
	// check if me blocked the friend
	if friend != nil && friend.IsBlocked() {
		rsp.Status = friendpb.AddFriendStatus_BLOCKED_BY_FRIEND
		return false
	}

	// check if me has blocked the friend
	if me != nil && me.IsBlocked() {
		rsp.Status = friendpb.AddFriendStatus_BLOCKED_BY_ME
		return false
	}

	return true
}

func (s *FriendService) canAddAutomatically(ctx context.Context, me, friend *data.Friend, rsp *friendpb.AddFriendResponse) (bool, error) {
	if friend == nil || friend.IsStranger() {
		return false, nil
	}

	// checked friend is not blocked me
	if me == nil {
		// create me -> friend relation
		me = &data.Friend{
			UID:       me.UID,
			FriendUID: friend.UID,
			Status:    friendpb.FriendStatus_FRIEND,
		}

		if err := s.friendDao.CreateFriend(ctx, me); err != nil {
			return false, err
		}

		rsp.Status = friendpb.AddFriendStatus_ADD_FRIEND_SUCCESS
		return true, nil
	}

	me.SetFriend()
	if err := s.friendDao.UpdateFriendStatus(ctx, me); err != nil {
		return false, err
	}

	rsp.Status = friendpb.AddFriendStatus_ADD_FRIEND_SUCCESS
	return true, nil
}

// friend has not blocked me and has no relation with me(no data or status is stranger)
// me has not blocked the friend and may have relation with the friend(no data or status in [friend, stranger])
func (s *FriendService) sendFriendRequest(ctx context.Context, req *friendpb.BaseFriendRequest,
	me, friend *data.Friend, rsp *friendpb.AddFriendResponse) error {
	// load old friend request
	fr, err := s.friendRequestDao.GetFriendRequest(ctx, req.GetUid(), req.GetFriendUid())
	if err != nil {
		return err
	}

	// if the friend request is not exist, create new one
	if fr == nil {
		fr = &data.FriendRequest{
			UID:       req.GetUid(),
			FriendUID: req.GetFriendUid(),
			Status:    friendpb.FriendRequestStatus_REQUESTED,
		}

		if err := s.friendRequestDao.CreateFriendRequest(ctx, fr); err != nil {
			return err
		}

		rsp.Status = friendpb.AddFriendStatus_SEND_REQUEST_SUCCESS
		return nil
	}

	// if the friend request is exist, check the status
	if fr.IsRequested() {
		rsp.Status = friendpb.AddFriendStatus_ALREADY_SENT_REQUEST
		return nil
	}

	if fr.IsAccepted() {
		// me and friend were friends before, no relation now, send friend request again
		fr.SetRequested()
		if err := s.friendRequestDao.UpdateFriendRequest(ctx, fr); err != nil {
			return err
		}

		rsp.Status = friendpb.AddFriendStatus_SEND_REQUEST_SUCCESS

		// if me status is friend, update friend status to stranger
		if me != nil && me.IsFriend() {
			me.SetStranger()
			if err := s.friendDao.UpdateFriendStatus(ctx, me); err != nil {
				return err
			}
		}
	}

	if fr.IsRejected() {
		// reject the friend request, resend friend request
		fr.SetRequested()
		if err := s.friendRequestDao.UpdateFriendRequest(ctx, fr); err != nil {
			return err
		}

		rsp.Status = friendpb.AddFriendStatus_SEND_REQUEST_SUCCESS

		// if me status is friend, update friend status to stranger
		if me != nil && me.IsFriend() {
			me.SetStranger()
			if err := s.friendDao.UpdateFriendStatus(ctx, me); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *FriendService) GetUserRelation(ctx context.Context, req *friendpb.GetUserRelationRequest) (
	*friendpb.GetUserRelationResponse, error) {
	ur, err := s.friendDao.GetFriend(ctx, req.GetUid(), req.GetFriendUid())
	if err != nil {
		return nil, err
	}

	rsp := &friendpb.GetUserRelationResponse{}
	if ur != nil {
		rsp.UserRelation = ur.ToProtoFriend()
	}

	return rsp, nil
}

func (s *FriendService) QueryUserRelationList(ctx context.Context, req *friendpb.QueryUserRelationListRequest) (
	*friendpb.QueryUserRelationListResponse, error) {
	urList, err := s.friendDao.GetFriends(ctx, req.GetUid())
	if err != nil {
		return nil, err
	}

	var (
		rsp           = &friendpb.QueryUserRelationListResponse{}
		friendUIDList = make([]string, len(urList))
		friendMap     = make(map[string]*data.User)
	)
	for i, ur := range urList {
		rsp.UserRelationList = append(rsp.UserRelationList, ur.ToProtoFriend())
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

func (s *FriendService) UpdateUserRelation(ctx context.Context, req *friendpb.UpdateUserRelationRequest) (
	*responsepb.BaseResponse, error) {
	pair := req.GetRelationPair()
	ur, err := s.friendDao.GetFriend(ctx, pair.GetUid(), pair.GetFriendUid())
	if err != nil {
		return nil, err
	}

	if ur == nil {
		return responsepb.ErrRelationNotFound, nil
	}

	newStatus, valid := req.GetAction().CheckActionAndGetNewStatus(ur.Status)
	if !valid {
		return responsepb.ErrInvalidUpdateRelationAction, nil
	}

	ur.SetStatus(newStatus)
	if err := s.friendDao.UpdateFriendStatus(ctx, ur); err != nil {
		return responsepb.ErrUnknown.SetMsg(err.Error()), nil
	}

	return responsepb.OK, nil
}
