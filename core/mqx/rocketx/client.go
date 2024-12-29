package rocketx

import (
	"context"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
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

// PrepareProducer 准备生产者
func (c *RocketClient) PrepareProducer() {
	if c.producer == nil {
		opts := NewProducerOptions(c.conf)
		if p, err := rocketmq.NewProducer(opts...); err != nil {
			log.Error("new producer err:", err.Error())
			panic(err)
		} else {
			c.producer = p
		}
	}
}

// PreparePushConsumer 准备PushConsumer消费者
func (c *RocketClient) PreparePushConsumer() {
	if c.pushConsumer == nil {
		opts := NewConsumerOptions(c.conf)
		if pushConsumer, err := rocketmq.NewPushConsumer(opts...); err != nil {
			log.Error("new push consumer err:", err.Error())
			panic(err)
		} else {
			c.pushConsumer = pushConsumer
		}
	}
}

// PreparePullConsumer 准备PullConsumer消费者
func (c *RocketClient) PreparePullConsumer() {
	if c.pullConsumer == nil {
		opts := NewConsumerOptions(c.conf)
		if pullConsumer, err := rocketmq.NewPullConsumer(opts...); err != nil {
			log.Error("new pull consumer err:", err.Error())
			panic(err)
		} else {
			c.pullConsumer = pullConsumer
		}
	}
}

func NewProducerOptions(c *mqx.Config) []producer.Option {
	var opts []producer.Option
	opts = append(opts,
		producer.WithNameServer([]string{c.Endpoint}),
		producer.WithGroupName(""),
		producer.WithInstanceName(""),
		producer.WithNamespace(""),
		producer.WithSendMsgTimeout(time.Millisecond*200),
		producer.WithRetry(0),
	)
	if c.AccessKey != "" && c.SecretKey != "" {
		opts = append(opts, producer.WithCredentials(
			primitive.Credentials{
				AccessKey: c.AccessKey,
				SecretKey: c.SecretKey,
			}))
	}
	return opts
}

func NewConsumerOptions(c *mqx.Config) []consumer.Option {
	var opts []consumer.Option
	opts = append(opts,
		consumer.WithNameServer([]string{c.Endpoint}),
		consumer.WithGroupName(""),
		consumer.WithNamespace(""),
		consumer.WithRetry(0),
	)
	if c.AccessKey != "" && c.SecretKey != "" {
		opts = append(opts, consumer.WithCredentials(
			primitive.Credentials{
				AccessKey: c.AccessKey,
				SecretKey: c.SecretKey,
			}))
	}
	return opts
}

// Publish 发布消息
func (c *RocketClient) Publish(ctx context.Context, body []byte, topic string, tag string) error {
	c.PrepareProducer()
	var msg = &primitive.Message{Topic: topic, Body: body}
	if tag != "" {
		msg.WithTag(tag)
	}
	if result, err := c.producer.SendSync(ctx, msg); err != nil {
		return errorx.Wrap(err, "rocketmq publish msg error")
	} else {
		log.Info("rocketmq publish msg success: ", result.String())
	}
	return nil
}

type ConsumerFunc func(context.Context, ...*primitive.MessageExt) (consumer.ConsumeResult, error)

// PushSubscribe push订阅消息
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

// PullSubscribe pull订阅消息
func (c *RocketClient) PullSubscribe(ctx context.Context, topic, tag string) error {
	c.PreparePullConsumer()
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
