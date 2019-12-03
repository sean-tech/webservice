package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sean-tech/webservice/config"
	"github.com/sean-tech/webservice/data"
	"github.com/sean-tech/webservice/encrypt"
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
	//engine.StaticFS(config.Upload.FileSavePath, http.Dir(GetUploadFilePath()))
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
 * 响应数据，成功，原数据转json返回
 */
func (g *Gin) ResponseData(data interface{}) {
	g.ResponseCode(STATUS_CODE_SUCCESS, data)
	return
}

/**
 * 响应数据，成功，aes加密返回
 */
func (g *Gin) ResponseAesJsonData(data []byte) {
	if keyBytes, exist := g.Ctx.Get(KEY_CTX_AES_KEY); exist == true {
		if secretBytes, err := encrypt.GetAes().EncryptCBC(data, keyBytes.([]byte)); err == nil {
			g.ResponseCode(STATUS_CODE_SUCCESS, base64.StdEncoding.EncodeToString(secretBytes))
			return
		}
	}
	g.ResponseCode(STATUS_CODE_SUCCESS, data)
	return
}

/**
 * 响应数据，成功，rsa加密返回
 */
func (g *Gin) ResponseRsaJsonData(data []byte) {
	if secretBytes, err := encrypt.GetRsa().Encrypt(config.Global.RsaClientPubKey, data); err == nil {
		if signBytes, err :=encrypt.GetRsa().Sign(config.Global.RsaServerPriKey, data); err == nil {
			g.Ctx.JSON(http.StatusOK, gin.H{
				"code" : STATUS_CODE_SUCCESS,
				"msg" :  STATUS_MSG_SUCCESS,
				"data" : base64.StdEncoding.EncodeToString(secretBytes),
				"sign" : base64.StdEncoding.EncodeToString(signBytes),
			})
			return
		}
	}
	g.ResponseCode(STATUS_CODE_SUCCESS, data)
	return
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
 * 响应数据，自定义error
 */
func (g *Gin) ResponseError(err error) {
	if e, ok := err.(*data.Error); ok {
		g.ResponseCode(StatusCode(e.Code), nil)
		return
	}
	g.ResponseMsg(STATUS_CODE_FAILED, err.Error(), nil)

	return
}



const (
	KEY_CTX_REQUISITION 		= "KEY_CTX_REQUISITION"
	KEY_CTX_PARAMS_JSON 		= "KEY_CTX_PARAMS_JSON"
	KEY_CTX_AES_KEY 			= "KEY_CTX_AES_KEY"
)

/**
 * 参数绑定
 */
func (g *Gin) BindParameter(parameter interface{}) error {
	//err := g.Ctx.Bind(parameter)
	//if err != nil {
	//	return err
	//}
	//return nil
	paramJsonBytes, exist := g.Ctx.Get(KEY_CTX_PARAMS_JSON)
	if !exist {
		return data.NewError(STATUS_CODE_INVALID_PARAMS, "参数json获取失败")
	}
	if err := json.Unmarshal(paramJsonBytes.([]byte), parameter); err != nil {
		return data.NewError(STATUS_CODE_INVALID_PARAMS, err.Error())
	}
	return nil
}

/**
 * 信息获取，获取传输链上context绑定的用户请求调用信息
 */
func (g *Gin) GetRequisition() *Requisition {
	rq := GetRequisition(g.Ctx)
	if rq != nil {
		return rq
	}
	id, _ := GenerateId(config.App.WorkerId)
	rq = &Requisition{
		ServiceId:    uint64(id),
		ServicePaths: make([]string, 5),
		UserId:       0,
		UserName:     "",
		Password:     "",
		IsAdministrotor:false,
	}
	g.Ctx.Set(KEY_CTX_REQUISITION, rq)
	return rq
}

/**
 * 文件上传处理函数
 */
//func (g *Gin) UploadFile() (fileUrl, filePath string, ok bool) {
//
//	data := make(map[string]string)
//
//	file, fileHeader, err := g.Ctx.Request.FormFile("file")
//	if err != nil {
//		logging.Warning(err)
//		g.ResponseMsg(STATUS_CODE_ERROR, err.Error(), data)
//		return "", "", false
//	}
//	if fileHeader == nil {
//		g.ResponseCode(STATUS_CODE_INVALID_PARAMS, data)
//		return "", "", false
//	}
//
//	fileName := GetUploadFileName(fileHeader.Filename)
//	fullPath := GetUploadFileFullPath()
//	savePath := GetUploadFilePath()
//	src := fullPath + fileName
//	if !CheckUploadFileExt(src) || !CheckUploadFileSize(file) {
//		g.ResponseCode(STATUS_CODE_UPLOAD_FILE_CHECK_FORMAT_WRONG, nil)
//		return "", "", false
//	}
//	if err := CheckUploadFile(fullPath); err != nil {
//		logging.Warning(err)
//		g.ResponseCode(STATUS_CODE_UPLOAD_FILE_CHECK_FAILED, nil)
//		return "", "", false
//	}
//	if err := g.Ctx.SaveUploadedFile(fileHeader, src); err != nil {
//		logging.Warning(err)
//		g.ResponseCode(STATUS_CODE_UPLOAD_FILE_SAVE_FAILED, nil)
//		return "", "", false
//	}
//	fileUrl = GetUploadFileFullUrl(fileName)
//	filePath = savePath + fileName
//	return fileUrl, filePath, true
//}


