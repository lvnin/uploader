package upload

import (
	"mime/multipart"
	"path"
	"strings"

	"uploader/global"

	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"github.com/pkg/errors"
)

var HuaWeiObs = new(Obs)

type Obs struct{}

func NewHuaWeiObsClient() (client *obs.ObsClient, err error) {
	return obs.New(global.ServerConfig.HuaWeiObs.AccessKey, global.ServerConfig.HuaWeiObs.SecretKey, global.ServerConfig.HuaWeiObs.Endpoint)
}

func (o *Obs) UploadFile(file *multipart.FileHeader) (string, string, error) {
	// var open multipart.File
	open, err := file.Open()
	if err != nil {
		return "", "", err
	}
	defer open.Close()
	filename := file.Filename
	input := &obs.PutObjectInput{
		PutObjectBasicInput: obs.PutObjectBasicInput{
			ObjectOperationInput: obs.ObjectOperationInput{
				Bucket: global.ServerConfig.HuaWeiObs.Bucket,
				Key:    filename,
			},
			ContentType: file.Header.Get("content-type"),
		},
		Body: open,
	}

	var client *obs.ObsClient
	client, err = NewHuaWeiObsClient()
	if err != nil {
		return "", "", errors.Wrap(err, "获取华为对象存储对象失败!")
	}

	_, err = client.PutObject(input)
	if err != nil {
		return "", "", errors.Wrap(err, "文件上传失败!")
	}
	filepath := global.ServerConfig.HuaWeiObs.Path + "/" + filename
	return filepath, filename, err
}

func (o *Obs) UploadFormFile(formData UploadFormData) (UploadResult, error) {
	// var open multipart.File
	open, err := formData.File.Open()
	if err != nil {
		return UploadResult{}, err
	}
	defer open.Close()
	ext := path.Ext(formData.File.Filename)
	name := strings.TrimSuffix(formData.File.Filename, ext)
	filename := formData.File.Filename
	input := &obs.PutObjectInput{
		PutObjectBasicInput: obs.PutObjectBasicInput{
			ObjectOperationInput: obs.ObjectOperationInput{
				Bucket: global.ServerConfig.HuaWeiObs.Bucket,
				Key:    filename,
			},
			ContentType: formData.File.Header.Get("content-type"),
		},
		Body: open,
	}

	var client *obs.ObsClient
	client, err = NewHuaWeiObsClient()
	if err != nil {
		return UploadResult{}, errors.Wrap(err, "获取华为对象存储对象失败!")
	}

	_, err = client.PutObject(input)
	if err != nil {
		return UploadResult{}, errors.Wrap(err, "文件上传失败!")
	}
	filepath := global.ServerConfig.HuaWeiObs.Path + "/" + filename
	return UploadResult{
		Url:  filepath,
		Name: name,
		Tag:  ext,
		Key:  filename,
	}, err
}

func (o *Obs) DeleteFile(key string) error {
	client, err := NewHuaWeiObsClient()
	if err != nil {
		return errors.Wrap(err, "获取华为对象存储对象失败!")
	}
	input := &obs.DeleteObjectInput{
		Bucket: global.ServerConfig.HuaWeiObs.Bucket,
		Key:    key,
	}
	var output *obs.DeleteObjectOutput
	output, err = client.DeleteObject(input)
	if err != nil {
		return errors.Wrapf(err, "删除对象(%s)失败!, output: %v", key, output)
	}
	return nil
}

func (o *Obs) CovertFile(url string, target string) (string, error) {
	return "", nil
}

func (*Obs) PreviewFile(url string, expires int64) (string, error) {
	return "", nil
}
