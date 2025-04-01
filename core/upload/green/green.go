package green

import (
	"uploader/global"
)

const (
	GreenStatsPass    = 1 // 通过
	GreenStatsFaliure = 0 // 未通过
)

type Green interface {
	GreenImage(url string) bool
	GreenErrorImageDefault() string
}

// NewOss OSS的实例化方法
func NewGreen() Green {
	switch global.ServerConfig.System.OssType {
	case "local":
		return &AliyunGreen{}
	case "qiniu":
		return &AliyunGreen{}
	case "tencent-cos":
		return &AliyunGreen{}
	case "aliyun-oss":
		return &AliyunGreen{}
	case "huawei-obs":
		return &AliyunGreen{}
	case "aws-s3":
		return &AliyunGreen{}
	default:
		return &AliyunGreen{}
	}
}
