package rocketx

import (
	"context"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/core/mqx"
	"github.com/go-xuan/quanx/os/errorx"
)

type RocketClient struct {
	conf         *mqx.Config
	producer     rocketmq.Producer
	pushConsumer rocketmq.PushConsumer
	pullConsumer rocketmq.PullConsumer
}

func (c *RocketClient) PrepareProducer() {
	if c.producer == nil {

	}
}

func (c *RocketClient) PreparePushConsumer() {
	if c.pushConsumer == nil {

	}
}

func (c *RocketClient) PreparePullConsumer() {
	if c.pullConsumer == nil {

	}
}

// Publish 发布消息
func (c *RocketClient) Publish(ctx context.Context, topic, tag string, body []byte) error {
	c.PrepareProducer()
	var msg = &primitive.Message{Topic: topic, Body: body}
	if tag != "" {
		msg.WithTag(tag)
	}
	if result, err := c.producer.SendSync(ctx, msg); err != nil {
		return errorx.Wrap(err, "rocketmq publish error")
	} else {
		log.Info("rocketMQ publish success: ", result.String())
	}
	return nil
}

type ConsumerFunc func(context.Context, ...*primitive.MessageExt) (consumer.ConsumeResult, error)

// PushSubscribe 根据tag订阅消息
func (c *RocketClient) PushSubscribe(ctx context.Context, topic, tag string, f ConsumerFunc) error {
	c.PreparePushConsumer()
	selector := consumer.MessageSelector{}
	if tag != "" || topic == "*" {
		selector.Type = consumer.TAG
		selector.Expression = tag
	}
	if err := c.pushConsumer.Subscribe(topic, selector, f); err != nil {
		return errorx.Wrap(err, "rocketMQ pushConsumer subscribe failed")
	}
	return nil
}

// PullSubscribe 根据tag订阅消息
func (c *RocketClient) PullSubscribe(ctx context.Context, topic, tag string, f ConsumerFunc) error {
	c.PreparePushConsumer()
	selector := consumer.MessageSelector{}
	if tag != "" || topic == "*" {
		selector.Type = consumer.TAG
		selector.Expression = tag
	}
	if err := c.pullConsumer.Subscribe(topic, selector); err != nil {
		return errorx.Wrap(err, "rocketMQ pullConsumer subscribe failed")
	}
	return nil
}
