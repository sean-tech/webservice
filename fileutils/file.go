package fileutils

import (
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
)

func GetSize(f multipart.File) (int, error) {
	content, err := ioutil.ReadAll(f)
	return len(content), err
}

func GetExt(fileName string) string {
	return path.Ext(fileName)
}

func CheckExist(src string) bool {
	_, err := os.Stat(src)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func CheckPermission(src string) bool {
	_, err := os.Stat(src)
	return os.IsPermission(err)
}

func MKDirIfNotExist(src string) error {
	if exist := CheckExist(src); exist == false {
		if err := MKDir(src); err != nil {
			return err
		}
	}
	return nil
}

func MKDir(src string) error {
	err := os.MkdirAll(src, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func Open(name string, flag int, perm os.FileMode) (*os.File, error) {
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}

	return f, nil
}