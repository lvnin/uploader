package config

type System struct {
	Addr          int    `mapstructure:"addr" json:"addr" yaml:"addr"`                               // 端口值
	OssType       string `mapstructure:"oss-type" json:"oss-type" yaml:"oss-type"`                   // Oss类型
	UseMultipoint bool   `mapstructure:"use-multipoint" json:"use-multipoint" yaml:"use-multipoint"` // 多点登录拦截
	Locale        string `mapstructure:"locale" json:"locale" yaml:"locale"`
	TemporaryPath string `mapstructure:"temporary-path" json:"temporary-path" yaml:"temporary-path"`
}
