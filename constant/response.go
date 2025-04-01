package request

import (
	"net/http"
	"uploader/constant/errcode"

	"github.com/gin-gonic/gin"
)

var ResponseWrapper = new(responseWrapper)

type RespResult struct {
	errcode.CodeResult
	Data any `json:"data"`
}

type responseWrapper struct{}

func (*responseWrapper) Ok(c *gin.Context, res *RespResult) {
	c.JSON(http.StatusOK, res)
}

func (*responseWrapper) Fail(c *gin.Context, errCode errcode.CodeResult) {
	c.JSON(http.StatusOK, RespResult{
		CodeResult: errCode,
	})
}
