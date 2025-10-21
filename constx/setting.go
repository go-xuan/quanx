package constx

import "path/filepath"

const (
	DefaultConfigDir      = "conf"        // 配置文件路径
	DefaultConfigFilename = "config.yaml" // 主配置文件
	DefaultResourceDir    = "resource"    // 默认资源存放路径
)

func GetDefaultConfigPath() string {
	return filepath.Join(DefaultConfigDir, DefaultConfigFilename)
}
