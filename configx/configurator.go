package configx

import (
	"strings"

	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/reflectx"
	log "github.com/sirupsen/logrus"
)

// Configurator 配置器接口，实现该接口的结构体可以作为配置器使用，用于读取和处理配置信息
// 配置器的工作流程：
// 1. 调用Readers()获取配置读取器列表
// 2. 按顺序使用读取器读取配置数据
// 3. 调用Valid()检查读取的配置是否有效
// 4. 配置有效则调用Execute()执行相关的业务逻辑
type Configurator interface {

	// GetClientName 获取客户端名称
	//GetClientName() string

	// Readers 获取配置读取器列表
	// 返回配置读取器列表，程序会按照读取器顺序依次尝试读取配置
	Readers() []Reader

	// Valid 验证配置是否有效
	// 用于决定配置器是否需要继续读取配置数据，通常在配置读取完成后检查关键配置项是否已正确设置
	Valid() bool

	// Execute 执行配置器逻辑
	// 用于根据读取到的配置值执行相关的业务逻辑，在配置读取完成且验证有效后调用此方法
	Execute() error
}

// LoadConfigurator 加载配置器
func LoadConfigurator(configurator Configurator) error {
	if configurator == nil {
		return nil
	}

	logger := log.WithField("configurator", reflectx.TypeOf(configurator).String())

	// 读取配置器
	location, err := ReadConfigurator(configurator)
	logger = logger.WithField("location", location)
	if err != nil {
		logger.WithError(err).Warn("read configurator failed")
		return errorx.Wrap(err, "read configurator failed")
	}

	// 执行配置器逻辑
	if err = configurator.Execute(); err != nil {
		logger.WithError(err).Warn("execute configurator failed")
		return errorx.Wrap(err, "execute configurator failed")
	}
	logger.Info("load configurator success")
	return nil
}

// ReadConfigurator 读取配置器，返回配置文件位置
func ReadConfigurator(configurator Configurator) (string, error) {
	if configurator == nil {
		return "nil", errorx.New("configurator is nil")
	} else if configurator.Valid() {
		return "self", nil
	}

	// 获取配置读取器
	readers := configurator.Readers()
	if len(readers) == 0 {
		return "", errorx.New("the configurator's reader is empty")
	}
	// 按照读取器的先后顺序依次读取配置
	var locations []string
	for _, reader := range readers {
		locations = append(locations, reader.Location())
		if err := ReaderRead(reader, configurator); err == nil && configurator.Valid() {
			return reader.Location(), nil
		}
	}
	return strings.Join(locations, ","), errorx.New("no available reader")
}
