// package mq init and pub/sub operation
package mq

import (
	messagev1 "github.com/yusank/goim/api/message/v1"
)

// Client defines mq pub/sub apis
type Client interface {
	Publish(msg *messagev1.MqMessage) error
	Subscribe(topic, consumerGroup string, concurrence int, cb SubscribeCallback)
}

// SubscribeCallback called by when new message coming,return error if got error and need to requeue.
type SubscribeCallback func(msg *messagev1.MqMessage) error
