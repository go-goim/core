package consts

import (
	"fmt"
	"time"
)

const (
	UserOnlineAgentKeyPrefix  = "userOnlineAgent:%d"  // userOnlineAgent:uid
	UserOfflineQueueKeyPrefix = "userOfflineQueue:%d" // userOfflineQueue:uid
	UserOnlineAgentKeyExpire  = time.Second * 30
	UserOfflineQueueKeyExpire = time.Hour * 24 * 7
	UserOfflineQueueMemberMax = 1000 // only store latest 1000 offline messages
)

func GetUserOnlineAgentKey(uid int64) string {
	return fmt.Sprintf(UserOnlineAgentKeyPrefix, uid)
}

func GetUserOfflineQueueKey(uid int64) string {
	return fmt.Sprintf(UserOfflineQueueKeyPrefix, uid)
}
