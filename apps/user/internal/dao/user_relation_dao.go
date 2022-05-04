package dao

import (
	"context"
	"sync"
	"time"

	"gorm.io/gorm"

	"github.com/yusank/goim/apps/user/internal/data"
	"github.com/yusank/goim/pkg/db"
)

type UserRelationDao struct{}

var (
	userRelationDao  *UserRelationDao
	userRelationOnce sync.Once
)

func GetUserRelationDao() *UserRelationDao {
	userRelationOnce.Do(func() {
		userRelationDao = &UserRelationDao{}
	})
	return userRelationDao
}

func (d *UserRelationDao) GetUserRelation(ctx context.Context, uid, friendUID string) (*data.UserRelation, error) {
	userRelation := &data.UserRelation{}
	err := db.GetDBFromCtx(ctx).Where("uid = ? AND friend_uid = ?", uid, friendUID).First(userRelation).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}

	return userRelation, nil
}

func (d *UserRelationDao) GetUserRelationList(ctx context.Context, uid string) ([]*data.UserRelation, error) {
	userRelationList := make([]*data.UserRelation, 0)
	err := db.GetDBFromCtx(ctx).Where("uid = ?", uid).Order("id").Find(&userRelationList).Error
	if err != nil {
		return nil, err
	}

	return userRelationList, nil
}

func (d *UserRelationDao) AddUserRelation(ctx context.Context, userRelation *data.UserRelation) error {
	return db.GetDBFromCtx(ctx).Create(userRelation).Error
}

func (d *UserRelationDao) UpdateUserRelationStatus(ctx context.Context, userRelation *data.UserRelation) error {
	tx := db.GetDBFromCtx(ctx).Model(userRelation).UpdateColumns(map[string]interface{}{
		"updated_at": time.Now(),
		"status":     userRelation.Status,
	})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (d *UserRelationDao) CountUserRelation(ctx context.Context, uid string) (int64, error) {
	var count int64
	err := db.GetDBFromCtx(ctx).Model(&data.UserRelation{}).Where("uid = ?", uid).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}
