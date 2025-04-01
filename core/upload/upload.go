package upload

import (
	"errors"
	"fmt"
	"strings"
	"time"
	request "uploader/constant"
	"uploader/constant/errcode"
	"uploader/core/upload/green"
	"uploader/global"

	"mime/multipart"

	"gorm.io/gorm"
)

// ============================================
// 上传数据结构
// ============================================
type UploadFormData struct {
	UserID  int                   `json:"user_id" form:"user_id"` // 用户ID
	Type    string                `json:"type" form:"type"`       // 上传类型
	ID      int                   `json:"id" form:"id"`           // 上传ID
	Renamed bool                  `json:"renamed" form:"renamed"` // 是否改名
	NoGreen bool                  `json:"noGreen" form:"noGreen"` // 是否进行审查
	File    *multipart.FileHeader `json:"file" form:"file"`       // 文件
}

type UploadResult struct {
	ID        uint      `json:"id" gorm:"primarykey"`        // 主键ID
	Name      string    `json:"name" gorm:"comment:文件名"`     // 文件名
	Size      int64     `json:"size" gorm:"comment:文件大小"`    // 文件大小
	Url       string    `json:"url" gorm:"comment:文件地址"`     // 文件地址
	Tag       string    `json:"tag" gorm:"comment:文件标签"`     // 文件标签
	Key       string    `json:"key" gorm:"comment:编号"`       // 编号
	Green     bool      `json:"green" gorm:"comment:是否审查通过"` // 是否审查通过
	CreatedAt time.Time `json:"created_at"`                  // 创建时间
	UpdatedAt time.Time `json:"updated_at"`                  // 更新时间
	Status    uint8     `json:"status"`                      // 状态
}

func (UploadResult) TableName() string {
	return "upload_records"
}

type UploadFile struct {
	ID           uint `gorm:"primarykey" json:"ID"` // 主键ID
	FileName     string
	FileMd5      string
	FilePath     string
	ExaFileChunk []UploadFileChunk
	ChunkTotal   int
	IsFinish     bool
	UpdatedAt    time.Time
	CreatedAt    time.Time
	Status       uint8
}

func (UploadFile) TableName() string {
	return "upload_files"
}

// file chunk struct, 切片结构体
type UploadFileChunk struct {
	ID              uint `gorm:"primarykey" json:"ID"` // 主键ID
	ExaFileID       uint
	FileChunkNumber int
	FileChunkPath   string
	UpdatedAt       time.Time
	CreatedAt       time.Time
	Status          uint8
}

func (UploadFileChunk) TableName() string {
	return "upload_file_chunks"
}

// ============================================
// OSS 对象存储接口
// ============================================
type OSS interface {
	UploadFile(file *multipart.FileHeader) (string, string, error)
	UploadFormFile(formData UploadFormData) (UploadResult, error)
	DeleteFile(key string) error
	CovertFile(url string, target string) (string, error)
	PreviewFile(url string, expires int64) (string, error)
}

// NewOss OSS的实例化方法
func NewOss() OSS {
	switch global.ServerConfig.System.OssType {
	case "local":
		return &Local{}
	case "qiniu":
		return &Qiniu{}
	case "tencent-cos":
		return &TencentCOS{}
	case "aliyun-oss":
		return &AliyunOSS{}
	case "huawei-obs":
		return HuaWeiObs
	case "aws-s3":
		return &AwsS3{}
	default:
		return &Local{}
	}
}

// ============================================
// 数据处理
// ============================================
var UploadService = new(uploadService)

type uploadService struct{}

// @function: RecordFile
// @description: 创建文件上传记录
// @param: file UploadResult
// @return: error
func (service *uploadService) RecordFile(file UploadResult) error {
	now := time.Now()
	file.CreatedAt = now
	file.UpdatedAt = now
	file.Status = 1
	if global.DB != nil {
		return global.DB.Create(&file).Error
	}
	return nil
}

// @function: FindFile
// @description: 查询文件记录
// @param: id uint
// @return: UploadResult, error
func (service *uploadService) FindFile(id uint) (file UploadResult, err error) {
	if global.DB != nil {
		err = global.DB.Where("id = ?", id).First(&file).Error
	}

	return file, err
}

// @function: FindFileByUrl
// @description: 查询文件记录
// @param: url string
// @return: UploadResult, error
func (service *uploadService) FindFileByUrl(url string) (file UploadResult, err error) {
	if global.DB != nil {
		err = global.DB.Where("url = ?", url).First(&file).Error
	} else {
		file.Key = url[strings.Index(url, "/"):]
	}

	return file, err
}

// @function: DeleteFile
// @description: 删除文件记录
// @param: url string
// @return: err error
func (service *uploadService) DeleteFile(url string) (err error) {
	var file UploadResult
	if global.DB != nil {
		file, err = service.FindFileByUrl(url)
		if err != nil {
			return
		}
	} else {
		file.Key = strings.ReplaceAll(url, global.ServerConfig.AliyunOSS.Bucket.BucketUrl+"/", "")
	}
	oss := NewOss()
	if err = oss.DeleteFile(file.Key); err != nil {
		return errors.New(errcode.ErrorCode("ERRCODE_FILE_DELETE_FAILURE").Message)
	}

	if global.DB != nil {
		err = global.DB.Where("id = ?", file.ID).Unscoped().Delete(&file).Error
	}

	return err
}

// EditFileName 编辑文件名或者备注
func (service *uploadService) EditFileName(file UploadResult) (err error) {
	if global.DB == nil {
		return
	}
	var fileFromDb UploadResult
	return global.DB.Where("id = ?", file.ID).First(&fileFromDb).Update("name", file.Name).Error
}

// @function: GetFileRecordInfoList
// @description: 分页获取数据
// @param: info request.PageInfo
// @return: list interface{}, total int64, err error
func (service *uploadService) GetFileRecordInfoList(info request.PageReq) (list interface{}, total int64, err error) {
	if global.DB == nil {
		return
	}

	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	keyword := info.Keyword
	db := global.DB.Model(&UploadResult{})
	var fileLists []UploadResult
	if len(keyword) > 0 {
		db = db.Where("name LIKE ?", "%"+keyword+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Order("updated_at desc").Find(&fileLists).Error
	return fileLists, total, err
}

// @function: UploadFile
// @description: 根据配置文件判断是文件上传到本地或者七牛云
// @param: header *multipart.FileHeader, noSave string
// @return: file UploadResult, err error
func (service *uploadService) UploadFile(header *multipart.FileHeader, noSave string) (file UploadResult, err error) {
	oss := NewOss()
	filePath, key, uploadErr := oss.UploadFile(header)
	if uploadErr != nil {
		return UploadResult{}, uploadErr
	}
	if noSave == "0" {
		s := strings.Split(header.Filename, ".")
		f := UploadResult{
			Url:  filePath,
			Name: header.Filename,
			Tag:  s[len(s)-1],
			Key:  key,
		}

		fmt.Printf("[uploadService.UploadFile]Url:%s, Name:%s, Tag:%s, Key:%s", f.Url, f.Name, f.Tag, f.Key)

		return f, service.RecordFile(f)
	}
	return
}

// @function: UploadFormFile
// @description: 根据配置文件判断是文件上传到本地或者七牛云
// @param: formData UploadFormData
// @return: result UploadResult, err error
func (service *uploadService) UploadFormFile(formData UploadFormData) (result UploadResult, err error) {
	oss := NewOss()

	result, err = oss.UploadFormFile(formData)
	if err != nil {
		return result, err
	}

	if !formData.NoGreen {
		// 需要进行图片审核（非图片直接通过）
		ossGreen := green.NewGreen()
		result.Green = ossGreen.GreenImage(result.Url)
		if !result.Green {
			service.DeleteFile(result.Url)
			result.Url = ossGreen.GreenErrorImageDefault() // 用审查不过默认图片代替
		}
	} else {
		result.Green = true
	}

	return result, service.RecordFile(result)
}

// @function: FindOrCreateFile
// @description: 上传文件时检测当前文件属性，如果没有文件则创建，有则返回文件的当前切片
// @param: fileMd5 string, fileName string, chunkTotal int
// @return: file model.ExaFileModel, err error
func (service *uploadService) FindOrCreateFile(fileMd5 string, fileName string, chunkTotal int) (file UploadFile, err error) {
	if global.DB == nil {
		return file, fmt.Errorf("[UploadService.FindOrCreateFile] DB init failed")
	}

	var cfile UploadFile
	cfile.FileMd5 = fileMd5
	cfile.FileName = fileName
	cfile.ChunkTotal = chunkTotal

	if errors.Is(global.DB.Where("file_md5 = ? AND is_finish = ?", fileMd5, true).First(&file).Error, gorm.ErrRecordNotFound) {
		err = global.DB.Where("file_md5 = ? AND file_name = ?", fileMd5, fileName).Preload("ExaFileChunk").FirstOrCreate(&file, cfile).Error
		return file, err
	}
	cfile.IsFinish = true
	cfile.FilePath = file.FilePath
	err = global.DB.Create(&cfile).Error
	return cfile, err
}

// @function: CreateFileChunk
// @description: 创建文件切片记录
// @param: id uint, fileChunkPath string, fileChunkNumber int
// @return: error
func (service *uploadService) CreateFileChunk(id uint, fileChunkPath string, fileChunkNumber int) error {
	if global.DB == nil {
		return fmt.Errorf("[UploadService.CreateFileChunk] DB init failed")
	}

	var chunk UploadFileChunk
	chunk.FileChunkPath = fileChunkPath
	chunk.ExaFileID = id
	chunk.FileChunkNumber = fileChunkNumber
	err := global.DB.Create(&chunk).Error
	return err
}

// @function: DeleteFileChunk
// @description: 删除文件切片记录
// @param: fileMd5 string, fileName string, filePath string
// @return: error
func (service *uploadService) DeleteFileChunk(fileMd5 string, filePath string) error {
	if global.DB == nil {
		return fmt.Errorf("[UploadService.DeleteFileChunk] DB init failed")
	}

	var chunks []UploadFileChunk
	var file UploadFile
	err := global.DB.Where("file_md5 = ? ", fileMd5).First(&file).
		Updates(map[string]interface{}{
			"IsFinish":  true,
			"file_path": filePath,
		}).Error
	if err != nil {
		return err
	}
	err = global.DB.Where("exa_file_id = ?", file.ID).Delete(&chunks).Unscoped().Error
	return err
}

// @function: CovertFile
// @description: 转换文件
// @param: url string
// @return: err error
func (service *uploadService) CovertFile(url string, target string) (string, error) {
	oss := NewOss()
	return oss.CovertFile(url, target)
}

// @function: PreviewFile
// @description: 预览文件
// @param: url string
// @param: expires int64
// @return: err error
func (service *uploadService) PreviewFile(url string, expires int64) (string, error) {
	oss := NewOss()
	return oss.PreviewFile(url, expires)
}
