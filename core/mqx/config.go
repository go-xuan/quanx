package mqx

import "github.com/go-xuan/quanx/core/configx"

type Config struct {
	Source    string `json:"source" yaml:"source" default:"default"`
	Type      string `json:"type" yaml:"type"`
	Enable    bool   `json:"enable" yaml:"enable"`
	Endpoint  string `json:"endpoint" yaml:"endpoint"`
	AccessKey string `json:"accessKey" yaml:"accessKey"`
	SecretKey string `json:"secretKey" yaml:"secretKey"`
}

func (c *Config) Format() string {
	//TODO implement me
	panic("implement me")
}

func (c *Config) Reader() *configx.Reader {
	//TODO implement me
	panic("implement me")
}

func (c *Config) Execute() error {
	//TODO implement me
	panic("implement me")
}
