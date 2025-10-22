package configx

import (
	"github.com/go-xuan/utilx/anyx"
	"github.com/go-xuan/utilx/errorx"
	log "github.com/sirupsen/logrus"
)

// Configurator 配置器
type Configurator interface {
	Readers() []Reader // 配置读取器, 配置读取器的顺序会影响配置的读取顺序
	Execute() error    // 配置器运行, 根据配置内容执行相关逻辑
	Valid() bool       // 是否有效, 用于判断配置是否需要继续读取
}

// ConfiguratorReadAndExecute 读取配置并运行
func ConfiguratorReadAndExecute(configurator Configurator) error {
	var logger = log.WithField("type", anyx.TypeOf(configurator).String())

	location, err := ConfiguratorRead(configurator)
	logger = logger.WithField("location", location)
	if err != nil {
		return errorx.Wrap(err, "configurator read error")
	}

	if !configurator.Valid() {
		logger.Info("configurator is invalid")
		return nil
	}
	if err = configurator.Execute(); err != nil {
		logger.WithField("error", err.Error()).Error("configurator run error")
		return errorx.Wrap(err, "configurator execute error")
	}
	logger.Info("configurator run success")
	return nil
}

// ConfiguratorRead 读取配置
func ConfiguratorRead(configurator Configurator) (string, error) {
	if configurator.Valid() {
		return "", nil
	}

	readers := configurator.Readers()
	if len(readers) == 0 {
		return "", errorx.New("configurator reader is empty")
	}

	// 按照读取器的先后顺序依次读取配置
	for _, reader := range readers {
		if err := ReadWithReader(configurator, reader); err == nil && configurator.Valid() {
			return reader.Location(), nil
		}
	}

	return "", errorx.New("configurator all readers are invalid")
}
