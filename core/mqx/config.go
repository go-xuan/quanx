package mqx

type Config struct {
	Source    string `json:"source" yaml:"source" default:"default"`
	Type      string `json:"type" yaml:"type"`
	Enable    bool   `json:"enable" yaml:"enable"`
	Endpoint  string `json:"endpoint" yaml:"endpoint"`
	AccessKey string `json:"accessKey" yaml:"accessKey"`
	SecretKey string `json:"secretKey" yaml:"secretKey"`
}
