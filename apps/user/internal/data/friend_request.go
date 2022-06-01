package data

import (
	friendpb "github.com/go-goim/api/user/friend/v1"
)

// FriendRequest is the model of fiend request table based on gorm, which is used for add friend request.
// FriendRequest data stored in mysql.
type FriendRequest struct {
	ID        int64                        `gorm:"primary_key"`
	UID       string                       `gorm:"column:uid"`
	FriendUID string                       `gorm:"column:friend_uid"`
	Status    friendpb.FriendRequestStatus `gorm:"column:status"`
	CreatedAt int64                        `gorm:"column:created_at"`
	UpdatedAt int64                        `gorm:"column:updated_at"`
}

func (FriendRequest) TableName() string {
	return "friend_request"
}

func (fr *FriendRequest) IsRequested() bool {
	return fr.Status == friendpb.FriendRequestStatus_REQUESTED
}

func (fr *FriendRequest) IsAccepted() bool {
	return fr.Status == friendpb.FriendRequestStatus_ACCEPTED
}

func (fr *FriendRequest) IsRejected() bool {
	return fr.Status == friendpb.FriendRequestStatus_REJECTED
}

func (fr *FriendRequest) SetStatus(status friendpb.FriendRequestStatus) {
	fr.Status = status
}

func (fr *FriendRequest) SetRequested() {
	fr.SetStatus(friendpb.FriendRequestStatus_REQUESTED)
}

func (fr *FriendRequest) SetAccepted() {
	fr.SetStatus(friendpb.FriendRequestStatus_ACCEPTED)
}

func (fr *FriendRequest) SetRejected() {
	fr.SetStatus(friendpb.FriendRequestStatus_REJECTED)
}

func (fr *FriendRequest) ToProto() *friendpb.FriendRequest {
	return &friendpb.FriendRequest{
		Uid:       fr.UID,
		FriendUid: fr.FriendUID,
		Status:    fr.Status,
		CreatedAt: fr.CreatedAt,
		UpdatedAt: fr.UpdatedAt,
	}
}
