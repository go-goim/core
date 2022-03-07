package data

import (
	"fmt"
	"time"
)

var (
	UserOnlineAgentKeyPrefix = "userOnlineAgent:%s" // userOnlineAgent:uid
	UserOnlineAgentKeyExpire = time.Second * 15     // default expire time
)

func GetUserOnlineAgentKey(uid string) string {
	return fmt.Sprintf(UserOnlineAgentKeyPrefix, uid)
}
