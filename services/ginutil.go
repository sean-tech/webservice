package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"github.com/sean-tech/webservice/config"
	"github.com/sean-tech/webservice/fileutils"
	"github.com/sean-tech/webservice/logging"
	"time"
)

/**
 * 启动 api server
 * handler: 接口实现serveHttp的对象
 */
func ServerServe(handler http.Handler) {
	// server
	s := http.Server{
		Addr:           fmt.Sprintf(":%d", config.ServerSetting.ApiPort),
		Handler:        handler,
		ReadTimeout:    config.ServerSetting.ReadTimeout,
		WriteTimeout:   config.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Fatal("Listen: %s\n", err)
		}
	}()
	// signal
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<- quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}

type Gin struct {
	Ctx *gin.Context
}

/**
 * 参数绑定
 */
func (g *Gin) Bind(parameter interface{}) error {
	err := g.Ctx.Bind(parameter)
	if err != nil {
		return err
	}
	return nil
}

/**
 * 参数绑定验证
 */
func (g *Gin) BindAndValid(parameter interface{}) error {
	err := g.Ctx.Bind(parameter)
	if err != nil {
		return err
	}
	valid := validation.Validation{}
	check, err := valid.Valid(parameter)
	if err != nil {
		return err
	}
	if !check {
		for _, err := range valid.Errors {
			errorStr := "参数:" + err.Key + " 校验失败:" + err.Message
			logging.Info(errorStr)
			return errors.New(errorStr)
		}
	}
	return nil
}

type ValidMidHandler func(parameter interface{})  error
/**
 * 参数绑定验证，带中间处理调用
 */
func (g *Gin) BindAndValidWithMidHandler(parameter interface{}, midhandler ValidMidHandler) error {
	err := g.BindAndValid(parameter)
	if err != nil {
		return err
	}
	return midhandler(parameter)
}

/**
 * 响应数据，根据code默认msg
 */
func (g *Gin) ResponseCode(statusCode StatusCode, data interface{})  {
	g.ResponseMsg(statusCode, statusCode.Msg(), data)
	return
}

/**
 * 响应数据，自定义msg
 */
func (g *Gin) ResponseMsg(statusCode StatusCode, msg string, data interface{})  {
	g.Ctx.JSON(http.StatusOK, gin.H{
		"code" : statusCode,
		"msg" :  msg,
		"data" : data,
	})
	return
}


const (
	KEY_CTX_USERID 				= "KEY_CTX_USERID"
	KEY_CTX_USERNAME 			= "KEY_CTX_USERNAME"
	KEY_CTX_PASSWORD 			= "KEY_CTX_PASSWORD"
	KEY_CTX_IS_ADMINISTROTOR 	= "KEY_CTX_IS_ADMINISTROTOR"
)
/**
 * 对ServiceInfo赋值
 */
func (g *Gin) BindServiceInfo(serviceCtx context.Context)  {
	serviceInfo := GetServiceInfo(serviceCtx)
	userId, exist := g.Ctx.Get(KEY_CTX_USERID)
	if exist {
		serviceInfo.UserId = userId.(uint64)
	}
	userName, exist := g.Ctx.Get(KEY_CTX_USERNAME)
	if exist {
		serviceInfo.UserName = userName.(string)
	}
	password, exist := g.Ctx.Get(KEY_CTX_PASSWORD)
	if exist {
		serviceInfo.Password = password.(string)
	}
	isAdministrotor, exist := g.Ctx.Get(KEY_CTX_IS_ADMINISTROTOR)
	if exist {
		serviceInfo.IsAdministrotor = isAdministrotor.(bool)
	}
}

/**
 * 文件上传处理函数
 */
func (g *Gin) UploadFile() (fileUrl, filePath string, ok bool) {

	data := make(map[string]string)

	file, fileHeader, err := g.Ctx.Request.FormFile("file")
	if err != nil {
		logging.Warning(err)
		g.ResponseMsg(STATUS_CODE_ERROR, err.Error(), data)
		return "", "", false
	}
	if fileHeader == nil {
		g.ResponseCode(STATUS_CODE_INVALID_PARAMS, data)
		return "", "", false
	}

	fileName := fileutils.GetUploadFileName(fileHeader.Filename)
	fullPath := fileutils.GetUploadFileFullPath()
	savePath := fileutils.GetUploadFilePath()
	src := fullPath + fileName
	if !fileutils.CheckUploadFileExt(src) || !fileutils.CheckUploadFileSize(file) {
		g.ResponseCode(STATUS_CODE_UPLOAD_FILE_CHECK_FORMAT_WRONG, nil)
		return "", "", false
	}
	if err := fileutils.CheckUploadFile(fullPath); err != nil {
		logging.Warning(err)
		g.ResponseCode(STATUS_CODE_UPLOAD_FILE_CHECK_FAILED, nil)
		return "", "", false
	}
	if err := g.Ctx.SaveUploadedFile(fileHeader, src); err != nil {
		logging.Warning(err)
		g.ResponseCode(STATUS_CODE_UPLOAD_FILE_SAVE_FAILED, nil)
		return "", "", false
	}
	fileUrl = fileutils.GetUploadFileFullUrl(fileName)
	filePath = savePath + fileName
	return fileUrl, filePath, true
}


