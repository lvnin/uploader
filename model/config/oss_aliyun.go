package config

type AliyunOSS struct {
	Bucket AliyunOSSBucket `mapstructure:"bucket" json:"bucket" yaml:"bucket"`
	Green  AliyunOSSGreen  `mapstructure:"green" json:"green" yaml:"green"`
}

type AliyunOSSBucket struct {
	Endpoint        string `mapstructure:"endpoint" json:"endpoint" yaml:"endpoint"`
	AccessKeyId     string `mapstructure:"access-key-id" json:"access-key-id" yaml:"access-key-id"`
	AccessKeySecret string `mapstructure:"access-key-secret" json:"access-key-secret" yaml:"access-key-secret"`
	BucketName      string `mapstructure:"bucket-name" json:"bucket-name" yaml:"bucket-name"`
	BucketUrl       string `mapstructure:"bucket-url" json:"bucket-url" yaml:"bucket-url"`
	BasePath        string `mapstructure:"base-path" json:"base-path" yaml:"base-path"`
}

type AliyunOSSGreen struct {
	AccessKey       string   `mapstructure:"access-key" json:"access-key" yaml:"access-key"`
	AccessKeySecret string   `mapstructure:"access-key-secret" json:"access-key-secret" yaml:"access-key-secret"`
	Region          string   `mapstructure:"region" json:"region" yaml:"region"`
	Endpoint        string   `mapstructure:"endpoint" json:"endpoint" yaml:"endpoint"`
	SpareRegion     string   `mapstructure:"spare-region" json:"spare-region" yaml:"spare-region"`
	SpareEndpoint   string   `mapstructure:"spare-endpoint" json:"spare-endpoint" yaml:"spare-endpoint"`
	Service         string   `mapstructure:"service" json:"service" yaml:"service"`
	ConnectTimeout  int      `mapstructure:"connect-timeout" json:"connect-timeout" yaml:"connect-timeout"`
	ReadTimeout     int      `mapstructure:"read-timeout" json:"read-timeout" yaml:"read-timeout"`
	ErrorImagePath  string   `mapstructure:"error-image-path" json:"error-image-path" yaml:"error-image-path"`
	Score           int      `mapstructure:"score" json:"score" yaml:"score"`
	Scenes          []string `mapstructure:"scenes" json:"scenes" yaml:"scenes"`
}
