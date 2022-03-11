package mq

import (
	"context"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

/*
 * Producer
 */

// Producer defines mq message producer
type Producer interface {
	rocketmq.Producer
}

/*
 * Consumer
 */

// Consumer defines mq push consumer
type Consumer interface {
	rocketmq.PushConsumer
}

type SubscribeCallback func(context.Context, ...*primitive.MessageExt) (consumer.ConsumeResult, error)
