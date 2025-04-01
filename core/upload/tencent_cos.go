package upload

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"uploader/global"

	"github.com/tencentyun/cos-go-sdk-v5"
	"go.uber.org/zap"
)

type TencentCOS struct{}

// UploadFile upload file to COS
func (*TencentCOS) UploadFile(file *multipart.FileHeader) (string, string, error) {
	client := NewClient()
	f, openError := file.Open()
	if openError != nil {
		global.Logger.Error("function file.Open() Filed", zap.Any("err", openError.Error()))
		return "", "", errors.New("function file.Open() Filed, err:" + openError.Error())
	}
	defer f.Close() // 创建文件 defer 关闭
	fileKey := fmt.Sprintf("%d%s", time.Now().Unix(), file.Filename)

	_, err := client.Object.Put(context.Background(), global.ServerConfig.TencentCOS.PathPrefix+"/"+fileKey, f, nil)
	if err != nil {
		panic(err)
	}
	return global.ServerConfig.TencentCOS.BaseURL + "/" + global.ServerConfig.TencentCOS.PathPrefix + "/" + fileKey, fileKey, nil
}

func (*TencentCOS) UploadFormFile(formData UploadFormData) (UploadResult, error) {
	client := NewClient()
	f, openError := formData.File.Open()
	if openError != nil {
		global.Logger.Error("function file.Open() Filed", zap.Any("err", openError.Error()))
		return UploadResult{}, errors.New("function file.Open() Filed, err:" + openError.Error())
	}
	defer f.Close() // 创建文件 defer 关闭
	ext := path.Ext(formData.File.Filename)
	name := strings.TrimSuffix(formData.File.Filename, ext)
	fileKey := fmt.Sprintf("%d%s", time.Now().Unix(), formData.File.Filename)
	_, err := client.Object.Put(context.Background(), global.ServerConfig.TencentCOS.PathPrefix+"/"+fileKey, f, nil)
	if err != nil {
		panic(err)
	}
	return UploadResult{
		Url:  global.ServerConfig.TencentCOS.BaseURL + "/" + global.ServerConfig.TencentCOS.PathPrefix + "/" + fileKey,
		Name: name,
		Tag:  ext,
		Key:  fileKey,
	}, nil
}

// DeleteFile delete file form COS
func (*TencentCOS) DeleteFile(key string) error {
	client := NewClient()
	name := global.ServerConfig.TencentCOS.PathPrefix + "/" + key
	_, err := client.Object.Delete(context.Background(), name)
	if err != nil {
		global.Logger.Error("function bucketManager.Delete() Filed", zap.Any("err", err.Error()))
		return errors.New("function bucketManager.Delete() Filed, err:" + err.Error())
	}
	return nil
}

func (*TencentCOS) CovertFile(url string, target string) (string, error) {
	return "", nil
}

func (*TencentCOS) PreviewFile(url string, expires int64) (string, error) {
	return "", nil
}

// NewClient init COS client
func NewClient() *cos.Client {
	urlStr, _ := url.Parse("https://" + global.ServerConfig.TencentCOS.Bucket + ".cos." + global.ServerConfig.TencentCOS.Region + ".myqcloud.com")
	baseURL := &cos.BaseURL{BucketURL: urlStr}
	client := cos.NewClient(baseURL, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  global.ServerConfig.TencentCOS.SecretID,
			SecretKey: global.ServerConfig.TencentCOS.SecretKey,
		},
	})
	return client
}
