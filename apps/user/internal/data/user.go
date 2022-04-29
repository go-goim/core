package data

import (
	userv1 "github.com/yusank/goim/api/user/v1"
)

// User is the model of user table based on gorm, which contains user basic info.
// User data stored in mysql.
type User struct {
	ID       int64  `gorm:"primary_key"`
	UID      string `gorm:"type:varchar(64);unique_index;not null"`
	Name     string `gorm:"type:varchar(32);not null"`
	Password string `gorm:"type:varchar(32);not null"`
	Email    string `gorm:"type:varchar(32);index;not null"`
	Phone    string `gorm:"type:varchar(32);index;not null"`
	Avatar   string `gorm:"type:varchar(128);not null"`
	Status   int    `gorm:"type:tinyint(1);not null"`
	CreateAt int64  `gorm:"type:bigint(20);not null;autoCreateTime"`
	UpdateAt int64  `gorm:"type:bigint(20);not null;autoUpdateTime"`
}

func (User) TableName() string {
	return "user"
}

const (
	UserStatusNormal int = iota
	UserStatusDeleted
)

const (
	UserCacheExpire = 60 * 60 * 24 // 1 day
)

func (u *User) IsDeleted() bool {
	return u.Status == UserStatusDeleted
}

func (u *User) ToProtoUserInternal() *userv1.UserInternal {
	return &userv1.UserInternal{
		Uid:      u.UID,
		Name:     u.Name,
		Email:    u.Email,
		Phone:    u.Phone,
		Avatar:   u.Avatar,
		Password: u.Password,
		CreateAt: u.CreateAt,
		UpdateAt: u.UpdateAt,
	}
}
