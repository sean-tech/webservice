package fileutils

import (
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path"
	"github.com/sean-tech/webservice/config"
	"github.com/sean-tech/webservice/cryptutils"
	"strings"
)

func GetUploadFileFullUrl(name string) string {
	return config.UploadSetting.FilePrefixUrl + "/" + config.UploadSetting.FileSavePath + name
}

func GetUploadFileName(name string) string {
	ext := path.Ext(name)
	fileName := strings.TrimSuffix(name, ext)
	fileName = cryptutils.Md5.EncodeWithTimestamp(fileName)
	return fileName + ext
}

func GetUploadFilePath() string {
	return config.UploadSetting.FileSavePath
}

func GetUploadFileFullPath() string {
	return config.AppSetting.RuntimeRootPath + config.UploadSetting.FileSavePath
}

func CheckUploadFileExt(fileName string) bool {
	ext := GetExt(fileName)
	for _, allowExt := range config.UploadSetting.FileAllowExts {
		if strings.ToUpper(allowExt) == strings.ToUpper(ext) {
			return true
		}
	}
	return false
}

func CheckUploadFileSize(f multipart.File) bool {
	size, err := GetSize(f)
	if err != nil {
		log.Println(err)
		return false
	}

	return size <= config.UploadSetting.FileMaxSize
}

func CheckUploadFile(src string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.Getwd err: %v", err)
	}

	err = MKDirIfNotExist(dir + "/" + src)
	if err != nil {
		return fmt.Errorf("file.IsNotExistMkDir err: %v", err)
	}

	perm := CheckPermission(src)
	if perm == true {
		return fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	return nil
}
