package logging

import (
	"fmt"
	"github.com/robfig/cron"
	"log"
	"path/filepath"
	"runtime"
	"sean.env/config"
	"sync"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARNING
	ERROR
)

var (
	logPrefix = ""
	DefaultPrefix = ""
	DefaultCallerDepth = 2

	levelFlags = []string{"DEBUG", "INFO", "WARN", "ERROR"}
	lock sync.Mutex
	debugLogger *log.Logger
	infoLogger 	*log.Logger
	warnLogger 	*log.Logger
	errorLogger *log.Logger
)

func Setup() {

	debugLogger = initLogger(DEBUG)
	infoLogger = initLogger(INFO)
	warnLogger = initLogger(WARNING)
	errorLogger = initLogger(ERROR)
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
	c := cron.New(cron.WithSeconds())
	spec := "0 0 0 * * *"
	_, err := c.AddFunc(spec, func() {
		if fileTimePassDaySlice(levelFlags[DEBUG]) {
			debugLogger = initLogger(DEBUG)
		}
		if fileTimePassDaySlice(levelFlags[INFO]) {
			infoLogger = initLogger(INFO)
		}
		if fileTimePassDaySlice(levelFlags[WARNING]) {
			warnLogger = initLogger(WARNING)
		}
		if fileTimePassDaySlice(levelFlags[ERROR]) {
			errorLogger = initLogger(ERROR)
		}
	})
	if err != nil {
		log.Fatal(err)
	}
	c.Start()
}