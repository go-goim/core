package dao

import (
	"context"
	"sort"
	"sync"
	"time"

	"gorm.io/gorm"

	"github.com/go-goim/core/apps/user/internal/data"
	"github.com/go-goim/core/pkg/cache"
	"github.com/go-goim/core/pkg/db"
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

// GetFriendStatusFromCache get friend status from cache.
// cache key: sort(uid, friend_uid), so that there is no duplicated key, only one record between two users.
// cache value: 1 as constant.
func (d *FriendDao) GetFriendStatusFromCache(ctx context.Context, uid, friendUID string) (bool, error) {
	keys := []string{uid, friendUID}
	sort.Strings(keys)
	key := "friend_status:" + keys[0] + ":" + keys[1]
	_, err := cache.Get(ctx, key)
	if err != nil {
		if err == cache.ErrCacheMiss {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// SetFriendStatusToCache set friend status to cache.
func (d *FriendDao) SetFriendStatusToCache(ctx context.Context, uid, friendUID string) error {
	keys := []string{uid, friendUID}
	sort.Strings(keys)
	key := "friend_status:" + keys[0] + ":" + keys[1]
	return cache.Set(ctx, key, []byte("1"), 0) // 0 means no expire time.
}

func (d *FriendDao) DeleteFriendStatusFromCache(ctx context.Context, uid, friendUID string) error {
	keys := []string{uid, friendUID}
	sort.Strings(keys)
	key := "friend_status:" + keys[0] + ":" + keys[1]
	return cache.Delete(ctx, key)
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
