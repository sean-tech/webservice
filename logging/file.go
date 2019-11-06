package logging

import (
	"fmt"
	"os"
	"github.com/sean-tech/webservice/config"
	"github.com/sean-tech/webservice/fileutils"
	"strings"
	"time"
)

func getLogFilePath() string {
	return fmt.Sprintf("%s%s", config.App.RuntimeRootPath, config.Log.LogSavePath)
}

func getLastDayLogFileName(levelFlag string) string {
	lastDayTime := time.Now().AddDate(0, 0, -1)
	return fmt.Sprintf("%s_%s_%s.%s",
		config.Log.LogSaveName,
		strings.ToLower(levelFlag),
		lastDayTime.Format(config.Log.TimeFormat),
		config.Log.LogFileExt,
	)
}

func getLogFileName(levelFlag string) string {
	return fmt.Sprintf("%s_%s.%s",
		config.Log.LogSaveName,
		strings.ToLower(levelFlag),
		config.Log.LogFileExt,
	)
}

func openLogFile(fileName, filePath string) (*os.File, error) {

	src := filePath
	perm := fileutils.CheckPermission(src)
	if perm == true {
		return nil, fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	err := fileutils.MKDirIfNotExist(src)
	if err != nil {
		return nil, fmt.Errorf("file.IsNotExistMkDir src: %s, err: %v", src, err)
	}

	f, err := fileutils.Open(src + fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("Fail to OpenFile :%v", err)
	}

	return f, nil
}

func fileTimePassDaySlice(levelFlag string) bool {
	// 默认日志文件不存在，无初始化文件，不继续处理
	currentLogFileExist := fileutils.CheckExist(getLogFilePath() + getLogFileName(levelFlag))
	if !currentLogFileExist {
		return false
	}
	// 昨日日志文件存在，说明已处理，不继续处理
	lastDayLogFileExist := fileutils.CheckExist(getLogFilePath() + getLastDayLogFileName(levelFlag))
	if lastDayLogFileExist {
		return false
	}
	// 把当前日志文件重命名为昨日日志文件
	originalPath := getLogFilePath() + getLogFileName(levelFlag)
	newPath := getLogFilePath() + getLastDayLogFileName(levelFlag)
	err := os.Rename(originalPath, newPath)
	if err != nil {
		Error(err)
		return false
	}
	return true
}








