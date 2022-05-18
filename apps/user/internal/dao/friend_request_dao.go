package dao

import (
	"context"
	"sync"
	"time"

	"gorm.io/gorm"

	"github.com/yusank/goim/apps/user/internal/data"
	"github.com/yusank/goim/pkg/db"
)

type FriendRequestDao struct {
}

var (
	friendRequestDao     *FriendRequestDao
	friendRequestDaoOnce sync.Once
)

func GetFriendRequestDao() *FriendRequestDao {
	friendRequestDaoOnce.Do(func() {
		friendRequestDao = &FriendRequestDao{}
	})
	return friendRequestDao
}

func (d *FriendRequestDao) CreateFriendRequest(ctx context.Context, fr *data.FriendRequest) error {
	fr.CreatedAt = time.Now().Unix()
	fr.UpdatedAt = time.Now().Unix()
	// get db from context
	return db.GetDBFromCtx(ctx).Create(fr).Error
}

func (d *FriendRequestDao) GetFriendRequest(ctx context.Context, uid, friendUID string) (*data.FriendRequest, error) {
	var fr data.FriendRequest
	if err := db.GetDBFromCtx(ctx).Where("uid = ? AND friend_uid = ?", uid, friendUID).First(&fr).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &fr, nil
}

func (d *FriendRequestDao) GetFriendRequestByID(ctx context.Context, id int64) (*data.FriendRequest, error) {
	var fr data.FriendRequest
	if err := db.GetDBFromCtx(ctx).Where("id = ?", id).First(&fr).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &fr, nil
}

func (d *FriendRequestDao) GetFriendRequests(ctx context.Context, uid string) ([]*data.FriendRequest, error) {
	var frs []*data.FriendRequest
	if err := db.GetDBFromCtx(ctx).Where("uid = ?", uid).Find(&frs).Error; err != nil {
		return nil, err
	}

	return frs, nil
}

func (d *FriendRequestDao) UpdateFriendRequest(ctx context.Context, fr *data.FriendRequest) error {
	return db.GetDBFromCtx(ctx).Model(fr).UpdateColumns(map[string]interface{}{
		"status":     fr.Status,
		"updated_at": fr.UpdatedAt,
	}).Error
}
