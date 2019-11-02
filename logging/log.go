package logging

import (
	"fmt"
	"github.com/robfig/cron"
	"github.com/sean-tech/webservice/config"
	"io"
	"log"
	"path/filepath"
	"runtime"
	"sync"
)

type Level int

const (
	level_debug Level = iota
	level_info
	level_warning
	level_error
	level_gin
)

var (
	logPrefix = ""
	DefaultPrefix = ""
	DefaultCallerDepth = 2

	levelFlags = []string{"DEBUG", "INFO", "WARN", "ERROR", "GIN"}
	lock sync.Mutex
	debugLogger *log.Logger
	infoLogger 	*log.Logger
	warnLogger 	*log.Logger
	errorLogger *log.Logger
	ginWriter io.Writer
)

func Setup() {

	if config.AppSetting.RunMode == "debug" {
		debugLogger = initLogger(level_debug)
	}
	infoLogger = initLogger(level_info)
	warnLogger = initLogger(level_warning)
	errorLogger = initLogger(level_error)
	logFileSliceTiming()
}

func Debug(v ...interface{})  {
	if config.AppSetting.RunMode == "debug" {
		debugLogger.Print(v)
		fmt.Println(v)
	}
}

func Info(v ...interface{})  {
	infoLogger.Print(v)
}

func Warning(v ...interface{})  {
	warnLogger.Print(v)
}

func Error(v ...interface{})  {
	errorLogger.Print(v)
}

func Fatal(v ...interface{})  {
	errorLogger.Print(v)
}

func initLogger(level Level) *log.Logger {
	var err error
	file, err := openLogFile(getLogFileName(levelFlags[level]), getLogFilePath())
	if err != nil {
		log.Fatalln(err)
	}
	lock.Lock()
	defer lock.Unlock()
	logger := log.New(file, DefaultPrefix, log.LstdFlags)
	logger.SetPrefix(getPrefix(level))
	return logger
}

func getPrefix(level Level) string {
	_, file, line, ok := runtime.Caller(DefaultCallerDepth)
	if ok {
		logPrefix = fmt.Sprintf("[%s]:[%s:%d]", levelFlags[level], filepath.Base(file), line)
	} else {
		logPrefix = fmt.Sprintf("[%s]", levelFlags[level])
	}
	return logPrefix
}

func logFileSliceTiming()  {
	c := cron.New()
	spec := "0 0 0 * * *"
	err := c.AddFunc(spec, func() {
		if config.AppSetting.RunMode == "debug" && fileTimePassDaySlice(levelFlags[level_debug]) {
			debugLogger = initLogger(level_debug)
		}
		if fileTimePassDaySlice(levelFlags[level_info]) {
			infoLogger = initLogger(level_info)
		}
		if fileTimePassDaySlice(levelFlags[level_warning]) {
			warnLogger = initLogger(level_warning)
		}
		if fileTimePassDaySlice(levelFlags[level_error]) {
			errorLogger = initLogger(level_error)
		}
		if ginWriterCallback != nil && fileTimePassDaySlice(levelFlags[level_gin]) {
			GinWriterGet(ginWriterCallback)
		}
	})
	if err != nil {
		log.Fatal(err)
	}
	c.Start()
}


type GinWriterCallback func(writer io.Writer)
var ginWriterCallback GinWriterCallback = nil
/**
 * 提供gin日志文件writer回调
 */
func GinWriterGet(callback GinWriterCallback)  {
	if callback == nil {
		return
	}
	if &ginWriterCallback != &callback {
		ginWriterCallback = callback
	}

	var err error
	ginWriter, err = openLogFile(getLogFileName(levelFlags[level_gin]), getLogFilePath())
	if err != nil {
		log.Fatalln(err)
	}
	ginWriterCallback(ginWriter)
}