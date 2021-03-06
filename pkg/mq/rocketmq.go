package mq

import (
	"context"
	"fmt"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

type ProducerConfig struct {
	Retry int
	Addr  []string
}

func (pc *ProducerConfig) validate() error {
	if pc.Retry <= 0 {
		pc.Retry = 3
	}

	if len(pc.Addr) == 0 {
		return fmt.Errorf("invalid producer config")
	}

	return nil
}

func NewProducer(cfg *ProducerConfig) (Producer, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	p, err := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver(cfg.Addr)),
		producer.WithRetry(cfg.Retry),
	)

	if err != nil {
		return nil, err
	}

	return p, nil
}

type ConsumerConfig struct {
	Addr        []string
	Subscriber  Subscriber
	Concurrence int
}

type Subscriber interface {
	Group() string
	Topic() string
	Consume(context.Context, ...*primitive.MessageExt) (consumer.ConsumeResult, error)
}

func (pc *ConsumerConfig) validate() error {
	if pc.Concurrence <= 0 {
		pc.Concurrence = 1
	}

	if len(pc.Addr) == 0 {
		return fmt.Errorf("invalid producer config")
	}

	if pc.Subscriber == nil {
		return fmt.Errorf("must specify subscriber")
	}

	return nil
}

func NewConsumer(cfg *ConsumerConfig) (Consumer, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	c, err := rocketmq.NewPushConsumer(
		consumer.WithGroupName(cfg.Subscriber.Group()),
		consumer.WithNsResolver(primitive.NewPassthroughResolver(cfg.Addr)),
	)
	if err != nil {
		return nil, err
	}

	for i := 0; i < cfg.Concurrence; i++ {
		if err = c.Subscribe(cfg.Subscriber.Topic(), consumer.MessageSelector{}, cfg.Subscriber.Consume); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func NewMessage(topic string, data []byte) *primitive.Message {
	return primitive.NewMessage(topic, data)
}
