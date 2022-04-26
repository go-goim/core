package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	redisv8 "github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	"github.com/yusank/goim/apps/user/internal/app"
	"github.com/yusank/goim/apps/user/internal/data"
	"github.com/yusank/goim/pkg/consts"
	"github.com/yusank/goim/pkg/db"
	"github.com/yusank/goim/pkg/log"
	"github.com/yusank/goim/pkg/util"
)

var (
	userDao     *UserDao
	userDaoOnce sync.Once
)

type UserDao struct {
	rdb *redisv8.Client
	// mysql.DB get from context, because we may need use transaction
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

func (u *UserDao) getUserFromRedis(ctx context.Context, uid string) (*data.User, error) {
	log.Debug("getUserFromRedis", "uid", uid)
	user := &data.User{}
	key := fmt.Sprintf("user:%s", uid)
	val, err := u.rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redisv8.Nil {
			return nil, nil
		}
		return nil, err
	}

	err = json.Unmarshal([]byte(val), user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByUID get user by uid
func (u *UserDao) GetUserByUID(ctx context.Context, uid string) (*data.User, error) {
	user, err := u.getUserFromRedis(ctx, uid)
	log.Debug("getUserFromRedis result", "user", user, "err", err)
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

	// put data to redis
	b, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}
	if err = u.rdb.Set(ctx, fmt.Sprintf("user:%s", uid), b,
		// expire in 24 hours + random(0-2.4 hours)
		time.Duration(data.UserCacheExpire+util.RandIntn(data.UserCacheExpire/10))*time.Second).Err(); err != nil {
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
