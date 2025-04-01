package api

import (
	"strconv"
	"time"
	request "uploader/constant"
	"uploader/constant/errcode"
	"uploader/core/upload"
	"uploader/global"
	"uploader/model"
	"uploader/service"
	"uploader/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UploadApi struct{}

// @function: UploadSingle
// @description: 上传单文件
// @param: file *multipart.FileHeader
// @return: err error
func (s *UploadApi) UploadSingle(c *gin.Context) {
	uploadType := c.DefaultPostForm("type", "file")
	renamed := c.DefaultPostForm("renamed", "true") == "true"
	uploadID, _ := strconv.Atoi(c.DefaultPostForm("id", "0"))
	noGreen := c.DefaultPostForm("noGreen", "false") == "true"

	_, header, err := c.Request.FormFile("file")
	if err != nil {
		global.Logger.Error(errcode.ErrorCode("ERRCODE_FILE_RECEIVE_FAILURE").Message, zap.Error(err))
		request.ResponseWrapper.Fail(c, errcode.ErrorCode("ERRCODE_FILE_RECEIVE_FAILURE"))
		return
	}

	formData := upload.UploadFormData{
		UserID:  utils.GetUserID(c),
		Type:    uploadType,
		ID:      uploadID,
		Renamed: renamed,
		NoGreen: noGreen,
		File:    header,
	}
	var result upload.UploadResult
	result, err = upload.UploadService.UploadFormFile(formData)
	if err != nil {
		global.Logger.Error(errcode.ErrorCode("ERRCODE_DBLINK_UPDATE_FAILURE").Message, zap.Error(err))
		request.ResponseWrapper.Fail(c, errcode.ErrorCode("ERRCODE_DBLINK_UPDATE_FAILURE"))
		return
	}
	errCode := errcode.ErrorCode("ERRCODE_UPLOAD_SUCCESS")
	if !result.Green {
		errCode = errcode.ErrorCode("ERRCODE_FILE_UPLOAD_NOT_COMPLIANCE")
	}

	request.ResponseWrapper.Ok(c, &request.RespResult{
		CodeResult: errCode,
		Data: struct {
			Result []upload.UploadResult `json:"result"`
		}{
			Result: []upload.UploadResult{
				result,
			},
		},
	})
}

// @function: UploadMulti
// @description: 上传多文件
// @param: file []*multipart.FileHeader
// @return: err error
func (s *UploadApi) UploadMulti(c *gin.Context) {
	uploadType := c.DefaultPostForm("type", "file")
	uploadID, _ := strconv.Atoi(c.DefaultPostForm("id", "0"))
	renamed := c.DefaultPostForm("renamed", "true") == "true"
	noGreen := c.DefaultPostForm("noGreen", "true") == "true"

	form, err := c.MultipartForm()
	if err != nil {
		global.Logger.Error(errcode.ErrorCode("ERRCODE_FILE_RECEIVE_FAILURE").Message, zap.Error(err))
		request.ResponseWrapper.Fail(c, errcode.ErrorCode("ERRCODE_FILE_RECEIVE_FAILURE"))
		return
	}

	files := form.File["file"]

	var results []upload.UploadResult

	for _, f := range files {
		formData := upload.UploadFormData{
			UserID:  utils.GetUserID(c),
			Type:    uploadType,
			ID:      uploadID,
			Renamed: renamed,
			NoGreen: noGreen,
			File:    f,
		}
		var result upload.UploadResult
		result, err = upload.UploadService.UploadFormFile(formData)
		if err != nil {
			global.Logger.Error(errcode.ErrorCode("ERRCODE_DBLINK_UPDATE_FAILURE").Message, zap.Error(err))
			request.ResponseWrapper.Fail(c, errcode.ErrorCode("ERRCODE_DBLINK_UPDATE_FAILURE"))
			return
		}

		results = append(results, result)
	}

	errCode := errcode.ErrorCode("ERRCODE_UPLOAD_SUCCESS")
	if !noGreen {
		for _, v := range results {
			if !v.Green {
				errCode = errcode.ErrorCode("ERRCODE_FILE_UPLOAD_NOT_COMPLIANCE")
				break
			}
		}
	}

	request.ResponseWrapper.Ok(c, &request.RespResult{
		CodeResult: errCode,
		Data: struct {
			Result []upload.UploadResult `json:"result"`
		}{
			Result: results,
		},
	})
}

// @function: UploadDelete
// @description: 删除上传文件
// @param: param request.DeleteFileReq
// @return: err error
func (s *UploadApi) UploadDelete(c *gin.Context) {
	var param request.UrlsReq
	if err := c.ShouldBindJSON(&param); err != nil {
		request.ResponseWrapper.Fail(c, errcode.ErrorCode("ERRCODE_PARAMETER_FAILURE"))
		return
	}

	for _, url := range param.Urls {
		upload.UploadService.DeleteFile(url)
	}

	request.ResponseWrapper.Ok(c, &request.RespResult{
		CodeResult: errcode.ErrorCode("ERRCODE_DELETE_SUCCESS"),
	})
}

// @function: UploadCovert
// @description: 转换上传文件
// @param: param request.UploadCovertReq
// @return: err error
func (s *UploadApi) UploadCovert(c *gin.Context) {
	var param request.UploadCovertReq
	if err := c.ShouldBindJSON(&param); err != nil {
		request.ResponseWrapper.Fail(c, errcode.ErrorCode("ERRCODE_PARAMETER_FAILURE"))
		return
	}

	url, err := upload.UploadService.CovertFile(param.Url, param.Target)
	if err != nil {
		request.ResponseWrapper.Fail(c, errcode.ErrorCode("ERRCODE_FILE_COVERT_FAILURE"))
		return
	}

	request.ResponseWrapper.Ok(c, &request.RespResult{
		CodeResult: errcode.ErrorCode("ERRCODE_COVERT_SUCCESS"),
		Data: request.UrlReq{
			Url: url,
		},
	})
}

// @function: UploadPreview
// @description: 预览上传文件
// @param: param request.UploadPreviewReq
// @return: err error
func (s *UploadApi) UploadPreview(c *gin.Context) {
	var param request.UploadPreviewReq
	if err := c.ShouldBindQuery(&param); err != nil {
		request.ResponseWrapper.Fail(c, errcode.ErrorCode("ERRCODE_PARAMETER_FAILURE"))
		return
	}

	url, err := upload.UploadService.PreviewFile(param.Url, param.Expires)
	if err != nil {
		request.ResponseWrapper.Fail(c, errcode.ErrorCode("ERRCODE_FILE_PREVIEW_FAILURE"))
		return
	}

	request.ResponseWrapper.Ok(c, &request.RespResult{
		CodeResult: errcode.ErrorCode("ERRCODE_PREVIEW_SUCCESS"),
		Data: request.UploadPreviewResp{
			Url:       url,
			ExpiresIn: time.Unix(time.Now().Unix()+param.Expires, 0),
		},
	})
}

// @function: UploadGenerateContract
// @description: 生成合同
// @param: param request.UploadPreviewReq
// @return: err error
func (s *UploadApi) UploadGenerateContract(c *gin.Context) {
	var param model.ContractModel
	if err := c.ShouldBindJSON(&param); err != nil {
		request.ResponseWrapper.Fail(c, errcode.ErrorCode("ERRCODE_PARAMETER_FAILURE"))
		return
	}

	result, err := service.ServiceGroupApp.ContractService.CovertContract(param)
	if err != nil {
		global.Logger.Error(err.Error())
		request.ResponseWrapper.Fail(c, errcode.ErrorCode("ERRCODE_GENERATE_CONTRACT_FAILURE"))
		return
	}

	request.ResponseWrapper.Ok(c, &request.RespResult{
		CodeResult: errcode.ErrorCode("ERRCODE_GENERATE_SUCCESS"),
		Data:       result,
	})
}
