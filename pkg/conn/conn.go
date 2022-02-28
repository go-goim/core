package conn

import messagev1 "github.com/yusank/goim/api/message/v1"

// todo 想办法把 ws server(for 循环读数据)和维持连接的两个go routine 合并成一个,统一处理
type Conn interface {
	PushMessage(*messagev1.PushMessageReq) error
	Key() string
	Ping() <-chan struct{}
	Close() error
}
