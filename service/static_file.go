package service

//import (
//	"fmt"
//	"github.com/sean-tech/webservice/config"
//	"github.com/sean-tech/webservice/fileutils"
//	"log"
//	"mime/multipart"
//	"os"
//	"path"
//	"strings"
//)
//
//func GetUploadFileFullUrl(name string) string {
//	return config.Upload.FilePrefixUrl + "/" + config.Upload.FileSavePath + name
//}
//
//func GetUploadFileName(name string) string {
//	ext := path.Ext(name)
//	fileName := strings.TrimSuffix(name, ext)
//	//fileName = encrypt.GetMd5().EncryptWithTimestamp([]byte(fileName), 0)
//	return fileName + ext
//}
//
//func GetUploadFilePath() string {
//	return config.Upload.FileSavePath
//}
//
//func GetUploadFileFullPath() string {
//	return config.App.RuntimeRootPath + config.Upload.FileSavePath
//}
//
//func CheckUploadFileExt(fileName string) bool {
//	ext := fileutils.GetExt(fileName)
//	for _, allowExt := range config.Upload.FileAllowExts {
//		if strings.ToUpper(allowExt) == strings.ToUpper(ext) {
//			return true
//		}
//	}
//	return false
//}
//
//func CheckUploadFileSize(f multipart.File) bool {
//	size, err := fileutils.GetSize(f)
//	if err != nil {
//		log.Println(err)
//		return false
//	}
//
//	return size <= config.Upload.FileMaxSize
//}
//
//func CheckUploadFile(src string) error {
//	dir, err := os.Getwd()
//	if err != nil {
//		return fmt.Errorf("os.Getwd err: %v", err)
//	}
//
//	err = fileutils.MKDirIfNotExist(dir + "/" + src)
//	if err != nil {
//		return fmt.Errorf("file.IsNotExistMkDir err: %v", err)
//	}
//
//	perm := fileutils.CheckPermission(src)
//	if perm == true {
//		return fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
//	}
//
//	return nil
//}
