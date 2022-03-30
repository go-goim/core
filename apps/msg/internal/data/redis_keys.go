package data

import (
	"fmt"
	"time"
)

var (
	UserOnlineAgentKeyPrefix  = "userOnlineAgent:%s"  // userOnlineAgent:uid
	UserOfflineQueueKeyPrefix = "userOfflineQueue:%s" // userOfflineQueue:uid
	UserOfflineQueueKeyExpire = time.Hour * 24 * 7
	UserOfflineQueueMemberMax = 1000 // only store latest 1000 offline messages
)

func GetUserOnlineAgentKey(uid string) string {
	return fmt.Sprintf(UserOnlineAgentKeyPrefix, uid)
}

func GetUserOfflineQueueKey(uid string) string {
	return fmt.Sprintf(UserOfflineQueueKeyPrefix, uid)
}
