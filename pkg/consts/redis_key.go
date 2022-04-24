package consts

import (
	"fmt"
	"time"
)

const (
	UserOnlineAgentKeyPrefix  = "userOnlineAgent:%s"  // userOnlineAgent:uid
	UserOfflineQueueKeyPrefix = "userOfflineQueue:%s" // userOfflineQueue:uid
	UserOnlineAgentKeyExpire  = time.Second * 30
	UserOfflineQueueKeyExpire = time.Hour * 24 * 7
	UserOfflineQueueMemberMax = 1000 // only store latest 1000 offline messages
)

func GetUserOnlineAgentKey(uid string) string {
	return fmt.Sprintf(UserOnlineAgentKeyPrefix, uid)
}

func GetUserOfflineQueueKey(uid string) string {
	return fmt.Sprintf(UserOfflineQueueKeyPrefix, uid)
}
