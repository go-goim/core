package conn

import messagev1 "github.com/yusank/goim/api/message/v1"

type Conn interface {
	PushMessage(*messagev1.PushMessage) error
	Key() string
	Ping() <-chan struct{}
	Close() error
}
