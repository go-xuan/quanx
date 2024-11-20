package mqx

import (
	"time"
	
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

type Config struct {
	Source    string `json:"source" yaml:"source" default:"default"`
	Type      string `json:"type" yaml:"type"`
	Enable    bool   `json:"enable" yaml:"enable"`
	Endpoint  string `json:"endpoint" yaml:"endpoint"`
	AccessKey string `json:"accessKey" yaml:"accessKey"`
	SecretKey string `json:"secretKey" yaml:"secretKey"`
}

func (c *Config) NewProducer() (rocketmq.Producer, error) {
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
	if p, err := rocketmq.NewProducer(opts...); err != nil {
		return nil, err

	} else {
		return p, nil
	}
}

func (c *Config) NewPushConsumer() (rocketmq.PushConsumer, error) {
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
	if p, err := rocketmq.NewPushConsumer(opts...); err != nil {
		return nil, err
	} else {
		return p, nil
	}
}
func (c *Config) Credentials() primitive.Credentials {
	return primitive.Credentials{AccessKey: c.AccessKey, SecretKey: c.SecretKey}
}
