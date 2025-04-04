package upload

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"path"
	"strings"
	"time"

	"uploader/global"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"go.uber.org/zap"
)

type Qiniu struct{}

// @object: *Qiniu
// @function: UploadFile
// @description: 上传文件
// @param: file *multipart.FileHeader
// @return: string, string, error
func (*Qiniu) UploadFile(file *multipart.FileHeader) (string, string, error) {
	putPolicy := storage.PutPolicy{Scope: global.ServerConfig.Qiniu.Bucket}
	mac := qbox.NewMac(global.ServerConfig.Qiniu.AccessKey, global.ServerConfig.Qiniu.SecretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := qiniuConfig()
	formUploader := storage.NewFormUploader(cfg)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{Params: map[string]string{"x:name": "github logo"}}

	f, openError := file.Open()
	if openError != nil {
		global.Logger.Error("function file.Open() Filed", zap.Any("err", openError.Error()))

		return "", "", errors.New("function file.Open() Filed, err:" + openError.Error())
	}
	defer f.Close()                                                  // 创建文件 defer 关闭
	fileKey := fmt.Sprintf("%d%s", time.Now().Unix(), file.Filename) // 文件名格式 自己可以改 建议保证唯一性
	putErr := formUploader.Put(context.Background(), &ret, upToken, fileKey, f, file.Size, &putExtra)
	if putErr != nil {
		global.Logger.Error("function formUploader.Put() Filed", zap.Any("err", putErr.Error()))
		return "", "", errors.New("function formUploader.Put() Filed, err:" + putErr.Error())
	}
	return global.ServerConfig.Qiniu.ImgPath + "/" + ret.Key, ret.Key, nil
}

func (*Qiniu) UploadFormFile(formData UploadFormData) (UploadResult, error) {
	putPolicy := storage.PutPolicy{Scope: global.ServerConfig.Qiniu.Bucket}
	mac := qbox.NewMac(global.ServerConfig.Qiniu.AccessKey, global.ServerConfig.Qiniu.SecretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := qiniuConfig()
	formUploader := storage.NewFormUploader(cfg)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{Params: map[string]string{"x:name": "github logo"}}

	f, openError := formData.File.Open()
	if openError != nil {
		global.Logger.Error("function file.Open() Filed", zap.Any("err", openError.Error()))

		return UploadResult{}, errors.New("function file.Open() Filed, err:" + openError.Error())
	}
	defer f.Close()
	ext := path.Ext(formData.File.Filename)
	name := strings.TrimSuffix(formData.File.Filename, ext)                   // 创建文件 defer 关闭
	fileKey := fmt.Sprintf("%d%s", time.Now().Unix(), formData.File.Filename) // 文件名格式 自己可以改 建议保证唯一性
	putErr := formUploader.Put(context.Background(), &ret, upToken, fileKey, f, formData.File.Size, &putExtra)
	if putErr != nil {
		global.Logger.Error("function formUploader.Put() Filed", zap.Any("err", putErr.Error()))
		return UploadResult{}, errors.New("function formUploader.Put() Filed, err:" + putErr.Error())
	}
	return UploadResult{
		Url:  global.ServerConfig.Qiniu.ImgPath + "/" + ret.Key,
		Name: name,
		Tag:  ext,
		Key:  ret.Key,
	}, nil
}

// @object: *Qiniu
// @function: DeleteFile
// @description: 删除文件
// @param: key string
// @return: error
func (*Qiniu) DeleteFile(key string) error {
	mac := qbox.NewMac(global.ServerConfig.Qiniu.AccessKey, global.ServerConfig.Qiniu.SecretKey)
	cfg := qiniuConfig()
	bucketManager := storage.NewBucketManager(mac, cfg)
	if err := bucketManager.Delete(global.ServerConfig.Qiniu.Bucket, key); err != nil {
		global.Logger.Error("function bucketManager.Delete() Filed", zap.Any("err", err.Error()))
		return errors.New("function bucketManager.Delete() Filed, err:" + err.Error())
	}
	return nil
}

func (*Qiniu) CovertFile(url string, target string) (string, error) {
	return "", nil
}

func (*Qiniu) PreviewFile(url string, expires int64) (string, error) {
	return "", nil
}

// @object: *Qiniu
// @function: qiniuConfig
// @description: 根据配置文件进行返回七牛云的配置
// @return: *storage.Config
func qiniuConfig() *storage.Config {
	cfg := storage.Config{
		UseHTTPS:      global.ServerConfig.Qiniu.UseHTTPS,
		UseCdnDomains: global.ServerConfig.Qiniu.UseCdnDomains,
	}
	switch global.ServerConfig.Qiniu.Zone { // 根据配置文件进行初始化空间对应的机房
	case "ZoneHuadong":
		cfg.Zone = &storage.ZoneHuadong
	case "ZoneHuabei":
		cfg.Zone = &storage.ZoneHuabei
	case "ZoneHuanan":
		cfg.Zone = &storage.ZoneHuanan
	case "ZoneBeimei":
		cfg.Zone = &storage.ZoneBeimei
	case "ZoneXinjiapo":
		cfg.Zone = &storage.ZoneXinjiapo
	}
	return &cfg
}
