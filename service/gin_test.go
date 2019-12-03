package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sean-tech/webservice/config"
	"github.com/sean-tech/webservice/logging"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestGinServer(t *testing.T) {
	//default_file_path := "../web_config.ini"
	// config & logging
	//config.SetupFromLocal(default_file_path)
	var (
		path = "sean.tech/webservice/config"
		endpoints = "localhost:2379"
	)
	config.Setup(path, endpoints)
	logging.Setup()
	// server start
	HttpServerServe(ginApiRegister)
}

func ginApiRegister(engine *gin.Engine) {
	apiv1 := engine.Group("api/order/v1")
	{
		apiv1.POST("/bindtest", bindtest)
	}
}

func bindtest(ctx *gin.Context)  {
	date := ctx.Request.Header.Get("Date")
	fmt.Println(date)
	g := Gin{
		Ctx: ctx,
	}
	var parameter GoodsPayParameter
	err := g.BindParameter(&parameter)
	if err != nil {
		logging.Info(err)
		g.ResponseMsg(STATUS_CODE_INVALID_PARAMS, err.Error(), nil)
		return
	}
	var payMoney float64 = 0
	err = GoodsPay(ctx, &parameter, &payMoney)
	if err != nil {
		logging.Info(err)
		g.ResponseMsg(STATUS_CODE_ERROR, err.Error(), nil)
		return
	}
	var resp = make(map[string]string)
	resp["payMoney"] = fmt.Sprintf("%v", payMoney)
	g.ResponseCode(STATUS_CODE_SUCCESS, resp)
}

func GoodsPay(ctx context.Context, parameter *GoodsPayParameter, payMoney *float64) error {
	err := ValidateParameter(parameter)
	if err != nil {
		return err
	}
	*payMoney = 10.0
	return nil
}

func TestPostToGinServer(t *testing.T)  {
	var url = "http://localhost:8811/api/order/v1/bindtest"

	var user_info map[string]interface{} = make(map[string]interface{})
	user_info["user_id"] = 101
	user_info["user_name"] = "18922311056"
	user_info["email"] = "1028990481@qq.com"

	var goods1 map[string]interface{} = make(map[string]interface{})
	goods1["goods_id"] = 1001
	goods1["goods_name"] = "三只松鼠干果巧克力100g包邮"
	goods1["goods_amount"] = 1
	goods1["remark"] = ""
	var goods []interface{} = []interface{}{goods1}
	var goods_ids []int = []int{1}

	var parameter map[string]interface{} = make(map[string]interface{})
	parameter["user_info"] = user_info
	parameter["goods"] = goods
	parameter["goods_ids"] = goods_ids

	jsonStr, err := json.Marshal(parameter)
	if err != nil {
		fmt.Printf("to json error:%v\n", err)
		return
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	//defer resp.Body.Close()
	if err != nil {
		fmt.Printf("resp error:%v", err)
	} else {
		statuscode := resp.StatusCode
		hea := resp.Header
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))
		fmt.Println(statuscode)
		fmt.Println(hea)
	}
}