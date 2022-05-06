package data

import (
	userv1 "github.com/yusank/goim/api/user/v1"
)

// UserRelation is the model of user relation table based on gorm, which is used for user relation management.
// UserRelation data stored in mysql.
type UserRelation struct {
	ID int64 `gorm:"primary_key"`
	// UID is the user uid of the user.
	UID string `gorm:"column:uid"`
	// FriendUID is the user uid of the friend.
	FriendUID string `gorm:"column:friend_uid"`
	// Status is the status of the relation.
	Status userv1.RelationStatus `gorm:"column:status"`
	// CreateAt is the creation time of the relation.
	CreateAt int64 `gorm:"column:create_at"`
	// UpdateAt is the update time of the relation.
	UpdateAt int64 `gorm:"column:update_at"`
}

func (UserRelation) TableName() string {
	return "user_relation"
}

const (
	UserMaxRelationCount = 2000 // UserMaxRelationCount is the max count of user relation.
)

func (ur *UserRelation) IsFriend() bool {
	return ur.Status == userv1.RelationStatus_FRIEND
}

func (ur *UserRelation) IsRequested() bool {
	return ur.Status == userv1.RelationStatus_REQUESTED
}

func (ur *UserRelation) IsStranger() bool {
	return ur.Status == userv1.RelationStatus_STRANGER
}

func (ur *UserRelation) IsBlocked() bool {
	return ur.Status == userv1.RelationStatus_BLOCKED
}

func (ur *UserRelation) SetStatus(status userv1.RelationStatus) {
	ur.Status = status
}

func (ur *UserRelation) SetFriend() {
	ur.SetStatus(userv1.RelationStatus_FRIEND)
}

func (ur *UserRelation) SetRequested() {
	ur.SetStatus(userv1.RelationStatus_REQUESTED)
}

func (ur *UserRelation) SetStranger() {
	ur.SetStatus(userv1.RelationStatus_STRANGER)
}

func (ur *UserRelation) SetBlocked() {
	ur.SetStatus(userv1.RelationStatus_BLOCKED)
}

func (ur *UserRelation) ToProtoUserRelation() *userv1.UserRelation {
	return &userv1.UserRelation{
		Uid:       ur.UID,
		FriendUid: ur.FriendUID,
		Status:    ur.Status,
		CreateAt:  ur.CreateAt,
		UpdateAt:  ur.UpdateAt,
	}
}
