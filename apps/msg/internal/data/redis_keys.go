package data

import (
	"fmt"
	"time"
)

var (
	UserOnlineAgentKeyPrefix  = "userOnlineAgent:%s"  // userOnlineAgent:uid
	UserOfflineQueueKeyPrefix = "userOfflineQueue:%s" // userOnlineAgent:uid
	UserOnlineAgentKeyExpire  = time.Second * 15      // default expire time
)

func GetUserOnlineAgentKey(uid string) string {
	return fmt.Sprintf(UserOnlineAgentKeyPrefix, uid)
}

func GetUserOfflineQueueKey(uid string) string {
	return fmt.Sprintf(UserOfflineQueueKeyPrefix, uid)
}
