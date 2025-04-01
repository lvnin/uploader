package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"uploader/global"
)

var FileHelper = new(fileHelper)

type fileHelper struct{}

// IsFileExists - 判断文件是否存在
// @param {string} filePath
// @returns bool
func (f *fileHelper) IsFileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}

	return false
}

// DeleteFile - 删除文件
// @param {string} filePath
// @returns error
func (f *fileHelper) DeleteFile(filePath string) error {
	if f.IsFileExists(filePath) {
		return os.Remove(filePath)
	}

	return nil
}

// DownloadFile - 下载文件
// @param {string} url
// @returns string, error
func (f *fileHelper) DownloadFile(url string) (string, error) {
	fileName := path.Base(url)
	filePath := global.ServerConfig.System.TemporaryPath + "/" + fileName

	if f.IsFileExists(filePath) {
		return filePath, nil
	}

	res, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	reader := bufio.NewReader(res.Body)

	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}

	writer := bufio.NewWriter(file)

	_, err = io.Copy(writer, reader)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func (f *fileHelper) GetFileHeader(filePath string) (*multipart.FileHeader, error) {
	// open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// create a buffer to hold the file in memory
	var buff bytes.Buffer
	buffWriter := io.Writer(&buff)

	// create a new form and create a new file field
	formWriter := multipart.NewWriter(buffWriter)
	formPart, err := formWriter.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		return nil, err
	}

	// copy the content of the file to the form's file field
	if _, err := io.Copy(formPart, file); err != nil {
		return nil, err
	}

	// close the form writer after the copying process is finished
	// I don't use defer in here to avoid unexpected EOF error
	formWriter.Close()

	// transform the bytes buffer into a form reader
	buffReader := bytes.NewReader(buff.Bytes())
	formReader := multipart.NewReader(buffReader, formWriter.Boundary())

	// read the form components with max stored memory of 1MB
	multipartForm, err := formReader.ReadForm(1 << 20)
	if err != nil {
		return nil, err
	}

	// return the multipart file header
	files, exists := multipartForm.File["file"]
	if !exists || len(files) == 0 {
		return nil, fmt.Errorf("multipart file not exists")
	}

	return files[0], nil
}
