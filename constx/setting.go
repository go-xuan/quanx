package constx

import "path/filepath"

// default setting
const (
	DefaultPort           = 9996          // 默认端口
	DefaultConfigDir      = "conf"        // 配置文件路径
	DefaultConfigFilename = "config.yaml" // 主配置文件
	DefaultResourceDir    = "resource"    // 默认资源存放路径
)

// GetDefaultConfigPath 获取默认配置文件路径
func GetDefaultConfigPath() string {
	return filepath.Join(DefaultConfigDir, DefaultConfigFilename)
}
