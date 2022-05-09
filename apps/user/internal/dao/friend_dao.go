package dao

import (
	"context"
	"sync"
	"time"

	"gorm.io/gorm"

	"github.com/yusank/goim/apps/user/internal/data"
	"github.com/yusank/goim/pkg/db"
)

type FriendDao struct{}

var (
	friendDao     *FriendDao
	friendDaoOnce sync.Once
)

func GetUserRelationDao() *FriendDao {
	friendDaoOnce.Do(func() {
		friendDao = &FriendDao{}
	})
	return friendDao
}

func (d *FriendDao) GetFriend(ctx context.Context, uid, friendUID string) (*data.Friend, error) {
	userRelation := &data.Friend{}
	err := db.GetDBFromCtx(ctx).Where("uid = ? AND friend_uid = ?", uid, friendUID).First(userRelation).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}

	return userRelation, nil
}

func (d *FriendDao) GetFriends(ctx context.Context, uid string) ([]*data.Friend, error) {
	userRelationList := make([]*data.Friend, 0)
	err := db.GetDBFromCtx(ctx).Where("uid = ?", uid).Order("id").Find(&userRelationList).Error
	if err != nil {
		return nil, err
	}

	return userRelationList, nil
}

func (d *FriendDao) CreateFriend(ctx context.Context, friend *data.Friend) error {
	friend.CreatedAt = time.Now().Unix()
	friend.UpdatedAt = time.Now().Unix()

	return db.GetDBFromCtx(ctx).Create(friend).Error
}

func (d *FriendDao) UpdateFriendStatus(ctx context.Context, userRelation *data.Friend) error {
	tx := db.GetDBFromCtx(ctx).Model(userRelation).UpdateColumns(map[string]interface{}{
		"updated_at": time.Now(),
		"status":     userRelation.Status,
	})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (d *FriendDao) CountFriends(ctx context.Context, uid string) (int64, error) {
	var count int64
	err := db.GetDBFromCtx(ctx).Model(&data.Friend{}).Where("uid = ?", uid).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}
