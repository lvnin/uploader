package errcode

import "uploader/global"

type CodeResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func ErrorCode(s string) CodeResult {
	if global.ServerConfig.System.Locale == "zh_CN" {
		return ErrCode_zh_CN[s]
	}

	return CodeResult{}
}
