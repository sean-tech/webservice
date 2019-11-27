package service

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sean-tech/webservice/config"
	"github.com/sean-tech/webservice/data"
	"github.com/sean-tech/webservice/logging"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

/** 服务注册回调函数 **/
type GinRegisterFunc func(engine *gin.Engine)

/**
 * 启动 api server
 * handler: 接口实现serveHttp的对象
 */
func HttpServerServe(registerFunc GinRegisterFunc) {
	// gin
	gin.SetMode(config.Global.RunMode)
	gin.DisableConsoleColor()
	logging.GinWriterGet(func(writer io.Writer) {
		gin.DefaultWriter = io.MultiWriter(writer, os.Stdout)
		logging.Debug(writer)
	})
	engine := gin.Default()
	engine.StaticFS(config.Upload.FileSavePath, http.Dir(GetUploadFilePath()))
	registerFunc(engine)
	// server
	s := http.Server{
		Addr:           fmt.Sprintf(":%d", config.Server.HttpPort),
		Handler:        engine,
		ReadTimeout:    config.Server.ReadTimeout,
		WriteTimeout:   config.Server.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Fatal(fmt.Sprintf("Listen: %v\n", err))
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
 * 响应数据，根据code默认msg
 */
func (g *Gin) ResponseCode(statusCode StatusCode, data interface{}) {
	g.ResponseMsg(statusCode, statusCode.Msg(), data)
	return
}

/**
 * 响应数据，自定义msg
 */
func (g *Gin) ResponseMsg(statusCode StatusCode, msg string, data interface{}) {
	g.Ctx.JSON(http.StatusOK, gin.H{
		"code" : statusCode,
		"msg" :  msg,
		"data" : data,
	})
	return
}

/**
 * 响应数据，自定义msg
 */
func (g *Gin) ResponseError(err data.Error) {
	g.ResponseCode(StatusCode(err.Code), nil)
	return
}

/**
 * 参数绑定
 */
func (g *Gin) BindParameter(parameter interface{}) error {
	err := g.Ctx.Bind(parameter)
	if err != nil {
		return err
	}
	return nil
}

const (
	KEY_CTX_USERID 				= "KEY_CTX_USERID"
	KEY_CTX_USERNAME 			= "KEY_CTX_USERNAME"
	KEY_CTX_PASSWORD 			= "KEY_CTX_PASSWORD"
	KEY_CTX_IS_ADMINISTROTOR 	= "KEY_CTX_IS_ADMINISTROTOR"
	KEY_CTX_PARAMS_JSON 		= "KEY_CTX_PARAMS_JSON"
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
	paramJsonBytes, exist := g.Ctx.Get(KEY_CTX_PARAMS_JSON)
	if exist {
		serviceInfo.Params = paramJsonBytes.([]byte)
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

	fileName := GetUploadFileName(fileHeader.Filename)
	fullPath := GetUploadFileFullPath()
	savePath := GetUploadFilePath()
	src := fullPath + fileName
	if !CheckUploadFileExt(src) || !CheckUploadFileSize(file) {
		g.ResponseCode(STATUS_CODE_UPLOAD_FILE_CHECK_FORMAT_WRONG, nil)
		return "", "", false
	}
	if err := CheckUploadFile(fullPath); err != nil {
		logging.Warning(err)
		g.ResponseCode(STATUS_CODE_UPLOAD_FILE_CHECK_FAILED, nil)
		return "", "", false
	}
	if err := g.Ctx.SaveUploadedFile(fileHeader, src); err != nil {
		logging.Warning(err)
		g.ResponseCode(STATUS_CODE_UPLOAD_FILE_SAVE_FAILED, nil)
		return "", "", false
	}
	fileUrl = GetUploadFileFullUrl(fileName)
	filePath = savePath + fileName
	return fileUrl, filePath, true
}


