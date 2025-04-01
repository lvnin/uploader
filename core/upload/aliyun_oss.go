package upload

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"path"
	"strings"
	"time"

	"uploader/global"
	"uploader/utils"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"go.uber.org/zap"
)

type AliyunOSS struct{}

func (*AliyunOSS) UploadFile(file *multipart.FileHeader) (string, string, error) {
	bucket, err := NewBucket()
	if err != nil {
		global.Logger.Error("function AliyunOSS.NewBucket() Failed", zap.Any("err", err.Error()))
		return "", "", errors.New("function AliyunOSS.NewBucket() Failed, err:" + err.Error())
	}

	// 读取本地文件。
	f, openError := file.Open()
	if openError != nil {
		global.Logger.Error("function file.Open() Failed", zap.Any("err", openError.Error()))
		return "", "", errors.New("function file.Open() Failed, err:" + openError.Error())
	}
	defer f.Close() // 创建文件 defer 关闭
	// 上传阿里云路径 文件名格式 自己可以改 建议保证唯一性
	// yunFileTmpPath := filepath.Join("uploads", time.Now().Format(time.DateOnly)) + "/" + file.Filename
	// 读取文件后缀
	ext := path.Ext(file.Filename)
	// 读取文件名并加密
	name := strings.TrimSuffix(file.Filename, ext)
	name = utils.MD5V([]byte(name))
	// 拼接新文件名
	filename := name + "_" + time.Now().Format(time.DateOnly) + ext

	yunFileTmpPath := global.ServerConfig.AliyunOSS.Bucket.BasePath + "/" + filename

	// 上传文件流。
	err = bucket.PutObject(yunFileTmpPath, f)
	if err != nil {
		global.Logger.Error("function formUploader.Put() Failed", zap.Any("err", err.Error()))
		return "", "", errors.New("function formUploader.Put() Failed, err:" + err.Error())
	}

	return global.ServerConfig.AliyunOSS.Bucket.BucketUrl + "/" + yunFileTmpPath, filename, nil
}

func (*AliyunOSS) UploadFormFile(formData UploadFormData) (UploadResult, error) {
	bucket, err := NewBucket()
	if err != nil {
		global.Logger.Error("function AliyunOSS.NewBucket() Failed", zap.Any("err", err.Error()))
		return UploadResult{}, errors.New("function AliyunOSS.NewBucket() Failed, err:" + err.Error())
	}

	// 读取上传文件
	f, openError := formData.File.Open()
	if openError != nil {
		global.Logger.Error("function file.Open() Failed", zap.Any("err", openError.Error()))
		return UploadResult{}, errors.New("function file.Open() Failed, err:" + openError.Error())
	}
	defer f.Close() // 创建文件 defer 关闭
	// 上传阿里云路径 文件名格式 自己可以改 建议保证唯一性
	// yunFileTmpPath := filepath.Join("uploads", time.Now().Format(time.DateOnly)) + "/" + file.Filename
	// 读取文件后缀
	filename := ""
	fileSize := formData.File.Size
	ext := path.Ext(formData.File.Filename)
	name := strings.TrimSuffix(formData.File.Filename, ext)
	dateTag := strings.ReplaceAll(time.Now().Format(time.DateTime), " ", "_")
	if formData.Renamed {
		// 需要重命名文件
		// 读取文件名并加密（格式：md5_格式化时间.ext）
		filename = utils.MD5V([]byte(name+utils.GetRandSlatString(5))) + "_" + dateTag + ext
	} else {
		// 判断是否已存在文件
		isExist, _ := bucket.IsObjectExist(fmt.Sprintf("%s/%s/%d/%s", global.ServerConfig.AliyunOSS.Bucket.BasePath,
			formData.Type, formData.ID, name+ext))
		if isExist {
			// 已存在相同文件需要改名
			filename = name + "_" + utils.GetRandSlatString(5) + "_" + dateTag + ext
		} else {
			filename = name + ext
		}
	}

	yunFileTmpPath := fmt.Sprintf("%s/%s/%d/%s", global.ServerConfig.AliyunOSS.Bucket.BasePath,
		formData.Type, formData.ID, filename)

	// 上传文件流。
	err = bucket.PutObject(yunFileTmpPath, f)
	if err != nil {
		global.Logger.Error("function formUploader.Put() Failed", zap.Any("err", err.Error()))
		return UploadResult{}, errors.New("function formUploader.Put() Failed, err:" + err.Error())
	}

	return UploadResult{
		Url:  global.ServerConfig.AliyunOSS.Bucket.BucketUrl + "/" + yunFileTmpPath,
		Name: name,
		Size: fileSize,
		Tag:  ext,
		Key:  yunFileTmpPath,
	}, nil
}

func (*AliyunOSS) DeleteFile(key string) error {
	bucket, err := NewBucket()
	if err != nil {
		global.Logger.Error("function AliyunOSS.NewBucket() Failed", zap.Any("err", err.Error()))
		return errors.New("function AliyunOSS.NewBucket() Failed, err:" + err.Error())
	}

	// 删除单个文件。objectName表示删除OSS文件时需要指定包含文件后缀在内的完整路径，例如abc/efg/123.jpg。
	// 如需删除文件夹，请将objectName设置为对应的文件夹名称。如果文件夹非空，则需要将文件夹下的所有object删除后才能删除该文件夹。
	err = bucket.DeleteObject(key)
	if err != nil {
		global.Logger.Error("function bucketManager.Delete() Filed", zap.Any("err", err.Error()))
		return errors.New("function bucketManager.Delete() Filed, err:" + err.Error())
	}

	return nil
}

// CovertFile - 转换文件
// @param {CovertFileData} covertData
// @returns error
func (*AliyunOSS) CovertFile(url string, target string) (string, error) {
	bucket, err := NewBucket()
	if err != nil {
		global.Logger.Error("function AliyunOSS.NewBucket() Failed", zap.Any("err", err.Error()))
		return "", err
	}

	// 判断转换之后的文件是否存在
	ext := path.Ext(url)
	key := strings.TrimSuffix(strings.ReplaceAll(url, global.ServerConfig.AliyunOSS.Bucket.BucketUrl+"/", ""), ext)

	covertKey := key + target
	isExist, _ := bucket.IsObjectExist(covertKey)
	if isExist {
		// 需要删除
		err = bucket.DeleteObject(covertKey)
		if err != nil {
			global.Logger.Error("function bucketManager.Delete() Filed", zap.Any("err", err.Error()))
			return "", err
		}
	}

	animationStyle := fmt.Sprintf("doc/convert,target_%s,source_%s",
		target, strings.ReplaceAll(ext, ".", ""))

	// 开始转换格式
	bucketNameEncoded := base64.URLEncoding.EncodeToString([]byte(global.ServerConfig.AliyunOSS.Bucket.BucketName))
	targetKeyEncoded := base64.URLEncoding.EncodeToString([]byte(covertKey))

	process := fmt.Sprintf("%s|sys/saveas,b_%v,o_%v/notify,topic_QXVkaW9Db252ZXJ0", animationStyle, bucketNameEncoded, targetKeyEncoded)
	_, err = bucket.AsyncProcessObject(key+ext, process)
	if err != nil {
		log.Fatalf("[AliyunOSS.CovertFile] Failed to process object: %s", err)
		return "", err
	}

	return covertKey, nil
}

func (*AliyunOSS) PreviewFile(url string, expires int64) (string, error) {
	bucket, err := NewBucket()
	if err != nil {
		global.Logger.Error("function AliyunOSS.NewBucket() Failed", zap.Any("err", err.Error()))
		return "", err
	}

	// 判断转换之后的文件是否存在
	key := strings.ReplaceAll(url, global.ServerConfig.AliyunOSS.Bucket.BucketUrl+"/", "")

	return bucket.SignURL(key, oss.HTTPGet, expires, oss.Process("imm/previewdoc,copy_0"))
}

func NewBucket() (*oss.Bucket, error) {
	// 创建OSSClient实例。
	client, err := oss.New(global.ServerConfig.AliyunOSS.Bucket.Endpoint, global.ServerConfig.AliyunOSS.Bucket.AccessKeyId, global.ServerConfig.AliyunOSS.Bucket.AccessKeySecret)
	if err != nil {
		return nil, err
	}

	// 获取存储空间。
	bucket, err := client.Bucket(global.ServerConfig.AliyunOSS.Bucket.BucketName)
	if err != nil {
		return nil, err
	}

	return bucket, nil
}
