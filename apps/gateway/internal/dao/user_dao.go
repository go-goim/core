package dao

import (
	"context"
	"sync"

	redisv8 "github.com/go-redis/redis/v8"

	"github.com/yusank/goim/apps/gateway/internal/app"
	"github.com/yusank/goim/pkg/consts"
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
