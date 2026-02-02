package nacosx

import (
	"reflect"

	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/marshalx"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
	log "github.com/sirupsen/logrus"
)

// 客户端
var _client *Client

// GetClient 获取客户端
func GetClient() *Client {
	if _client == nil {
		panic("nacos client not initialized, please check the relevant config")
	}
	return _client
}

// Initialized 是否初始化
func Initialized() bool {
	return _client != nil
}

// NewClient 创建nacos客户端
func NewClient(config *Config) (*Client, error) {
	client := &Client{config: config}
	var param = config.ClientParam()
	var err error
	switch config.Mode {
	case OnlyConfig: // 仅初始化配置中心
		if client.configClient, err = config.ConfigClient(param); err != nil {
			return nil, errorx.Wrap(err, "init nacos config client failed")
		}
	case OnlyNaming: // 仅初始化服务发现
		if client.namingClient, err = config.NamingClient(param); err != nil {
			return nil, errorx.Wrap(err, "init nacos naming client failed")
		}
	case ConfigAndNaming: // 初始化配置中心和服务发现
		if client.configClient, err = config.ConfigClient(param); err != nil {
			return nil, errorx.Wrap(err, "init nacos config client failed")
		}
		if client.namingClient, err = config.NamingClient(param); err != nil {
			return nil, errorx.Wrap(err, "init nacos naming client failed")
		}
	}
	return client, nil
}

// Client nacos客户端
type Client struct {
	config       *Config                     // nacos配置
	configClient config_client.IConfigClient // nacos配置中心客户端
	namingClient naming_client.INamingClient // nacos服务发现客户端
}

// GetGroup 获取分组
func (c *Client) GetGroup() string {
	return c.config.Group
}

// GetConfigClient 获取配置中心客户端
func (c *Client) GetConfigClient() config_client.IConfigClient {
	if c.configClient == nil {
		panic("the nacos config client has not been initialized")
	}
	return c.configClient
}

// GetNamingClient 获取服务发现客户端
func (c *Client) GetNamingClient() naming_client.INamingClient {
	if c.namingClient == nil {
		panic("the nacos naming client has not been initialized")
	}
	return c.namingClient
}

// PublishConfig 发布nacos配置
func (c *Client) PublishConfig(param vo.ConfigParam) error {
	if _, err := c.GetConfigClient().PublishConfig(param); err != nil {
		return errorx.Wrap(err, "nacos config client publish config failed")
	}
	return nil
}

// GetConfig 获取nacos配置
func (c *Client) GetConfig(param vo.ConfigParam, publishIfNotExist bool) (string, error) {
	content, err := c.GetConfigClient().GetConfig(param)
	if err != nil {
		return "", errorx.Wrap(err, "get nacos config failed")
	}
	// 如果配置不存在，则发布配置
	if content == "" && publishIfNotExist {
		if err = c.PublishConfig(param); err != nil {
			return "", errorx.Wrap(err, "publish nacos config failed")
		}
	}
	return content, nil
}

// DeleteConfig 删除nacos配置
func (c *Client) DeleteConfig(param vo.ConfigParam) error {
	if _, err := c.GetConfigClient().DeleteConfig(param); err != nil {
		return errorx.Wrap(err, "nacos config client delete config failed")
	}
	return nil
}

// ReadConfig 读取nacos配置
func (c *Client) ReadConfig(config any, param vo.ConfigParam) ([]byte, error) {
	// 配置值必须是指针类型，否则不允许读取
	if typeOf := reflect.TypeOf(config); typeOf.Kind() != reflect.Pointer {
		return nil, errorx.New("the scanned object must be of pointer type")
	}
	// 读取配置文件内容
	content, err := c.GetConfig(param, false)
	if err != nil {
		return nil, errorx.Wrap(err, "read nacos config failed")
	}
	if content == "" {
		return nil, errorx.New("read nacos config empty")
	}
	// 解析配置文件内容
	data := []byte(content)
	if err = marshalx.Apply(string(param.Type)).Unmarshal(data, config); err != nil {
		return nil, errorx.Wrap(err, "unmarshal nacos config failed")
	}
	return data, nil
}

// ListenConfig 监听nacos配置
func (c *Client) ListenConfig(config any, param vo.ConfigParam) error {
	// 配置监听响应方法
	param.OnChange = func(namespace, group, dataId, data string) {
		logger := log.WithField("dataId", dataId).
			WithField("group", group).
			WithField("namespace", namespace).
			WithField("data", data)
		logger.Info("the nacos config data has changed !!!")
		if err := marshalx.Apply(dataId).Unmarshal([]byte(data), config); err != nil {
			logger.WithError(err).Error("update nacos config failed")
		}
	}
	if err := c.GetConfigClient().ListenConfig(param); err != nil {
		return errorx.Wrap(err, "listen nacos config failed")
	}
	return nil
}

// CancelListenConfig 取消监听nacos配置
func (c *Client) CancelListenConfig(param vo.ConfigParam) error {
	if err := c.GetConfigClient().CancelListenConfig(param); err != nil {
		return errorx.Wrap(err, "cancel listen nacos config failed")
	}
	return nil
}

// SearchConfig 搜索nacos配置
func (c *Client) SearchConfig(param vo.SearchConfigParam) (*model.ConfigPage, error) {
	page, err := c.GetConfigClient().SearchConfig(param)
	if err != nil {
		return nil, errorx.Wrap(err, "search nacos config failed")
	}
	return page, nil
}
