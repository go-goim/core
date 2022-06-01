package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	redisv8 "github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	"github.com/go-goim/core/apps/user/internal/app"
	"github.com/go-goim/core/apps/user/internal/data"
	"github.com/go-goim/core/pkg/cache"
	"github.com/go-goim/core/pkg/consts"
	"github.com/go-goim/core/pkg/db"
	"github.com/go-goim/core/pkg/log"
	"github.com/go-goim/core/pkg/util"
)

var (
	userDao     *UserDao
	userDaoOnce sync.Once
)

type UserDao struct {
	rdb *redisv8.Client
}

func GetUserDao() *UserDao {
	userDaoOnce.Do(func() {
		userDao = &UserDao{
			rdb: app.GetApplication().Redis,
		}
	})
	return userDao
}

func (u *UserDao) GetUser(ctx context.Context, id int64) (*data.User, error) {
	user := &data.User{}
	tx := db.GetDBFromCtx(ctx).Where("id = ?", id).First(user)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, tx.Error
	}

	if user.IsDeleted() {
		return nil, nil
	}

	return user, nil
}

func (u *UserDao) getUserFromCache(ctx context.Context, uid string) (*data.User, error) {
	log.Debug("getUserFromCache", "uid", uid)
	user := &data.User{}
	key := fmt.Sprintf("user:%s", uid)
	val, err := cache.Get(ctx, key) // use default cache
	if err != nil {
		if err == cache.ErrCacheMiss {
			return nil, nil
		}

		return nil, err
	}

	err = json.Unmarshal(val, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserDao) setUserToCache(ctx context.Context, user *data.User) error {
	log.Debug("setUserToCache", "uid", user.ID)
	key := fmt.Sprintf("user:%s", user.UID)
	val, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return cache.Set(ctx, key, val, time.Duration(data.UserCacheExpire+util.RandIntn(data.UserCacheExpire/10))*time.Second)

}

// GetUserByUID get user by uid
func (u *UserDao) GetUserByUID(ctx context.Context, uid string) (*data.User, error) {
	user, err := u.getUserFromCache(ctx, uid)
	log.Debug("getUserFromCache result", "user", user, "err", err)
	if err != nil {
		return user, nil
	}

	if user != nil {
		return user, nil
	}

	user = &data.User{}
	tx := db.GetDBFromCtx(ctx).Where("uid = ?", uid).First(user)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, tx.Error
	}

	// put data to cache
	err = u.setUserToCache(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByEmail get user by email directly from db
func (u *UserDao) GetUserByEmail(ctx context.Context, email string) (*data.User, error) {
	user := &data.User{}
	tx := db.GetDBFromCtx(ctx).Where("email = ?", email).First(user)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, tx.Error
	}

	if user.IsDeleted() {
		return nil, nil
	}

	return user, nil
}

// GetUserByPhone get user by phone directly from db
func (u *UserDao) GetUserByPhone(ctx context.Context, phone string) (*data.User, error) {
	user := &data.User{}
	tx := db.GetDBFromCtx(ctx).Where("phone = ?", phone).First(user)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, tx.Error
	}

	if user.IsDeleted() {
		return nil, nil
	}

	return user, nil
}

// GetUserOnlineAgent get user online agent from redis
func (u *UserDao) GetUserOnlineAgent(ctx context.Context, uid string) (string, error) {
	key := consts.GetUserOnlineAgentKey(uid)
	val, err := u.rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redisv8.Nil {
			return "", nil
		}
		return "", err
	}

	return val, nil
}

func (u *UserDao) CreateUser(ctx context.Context, user *data.User) error {
	tx := db.GetDBFromCtx(ctx).Create(user)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (u *UserDao) UpdateUser(ctx context.Context, user *data.User) error {
	tx := db.GetDBFromCtx(ctx).Save(user)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

// UndoDelete undo delete user with new password
func (u *UserDao) UndoDelete(ctx context.Context, user *data.User) error {
	tx := db.GetDBFromCtx(ctx).Model(user).UpdateColumns(map[string]interface{}{
		"password":   user.Password,
		"updated_at": time.Now(),
		"status":     data.UserStatusNormal,
	})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (u *UserDao) ListUsers(ctx context.Context, uids ...string) ([]*data.User, error) {
	var users []*data.User
	if len(uids) == 0 {
		return users, nil
	}
	tx := db.GetDBFromCtx(ctx).Where("uid in (?)", uids).Find(&users)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return users, nil
}
