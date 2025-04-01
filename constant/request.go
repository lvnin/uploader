package request

import (
	"time"
	"uploader/model"
)

type PageReq struct {
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"pageSize" form:"pageSize"`
	Keyword  string `json:"keyword" form:"keyword"`
}

type UrlReq struct {
	Url string `json:"url" form:"url"`
}

type UrlsReq struct {
	Urls []string `json:"urls" form:"urls"`
}

type UploadCovertReq struct {
	Url    string `json:"url" form:"url"`
	Target string `json:"target" form:"target"`
}

type UploadPreviewReq struct {
	Url     string `json:"url" form:"url"`
	Expires int64  `json:"expires" form:"expires"`
}

type UploadPreviewResp struct {
	Url       string    `json:"url"`
	ExpiresIn time.Time `json:"expiresIn"`
}

type GenerateContractReq struct {
	Source string              `json:"source" from:"source"` // 来源文件
	Target string              `json:"target" form:"target"` // 生成文件
	Params model.ContractModel `json:"params" form:"params"`
}
