package service

import (
	"encoding/json"
	"fmt"
	"github.com/sean-tech/webservice/config"
	"testing"
)

type User struct {
	UserId int						`json:"user_id" validate:"required,min=1"`
	UserName string					`json:"user_name" validate:"required,eorp"`
	Email string					`json:"email" validate:"required,email"`
}

type GoodsPayInfoParameter struct {
	GoodsId int			`json:"goods_id" validate:"required,min=1"`
	GoodsName string	`json:"goods_name" validate:"required,gte=1"`
	GoodsAmount int		`json:"goods_amount" validate:"required,min=1"`
	Remark string 		`json:"remark" validate:"gte=0"`
}

type GoodsPayParameter struct {
	UserInfo *User					`json:"user_info" validate:"required"`
	Goods []*GoodsPayInfoParameter	`json:"goods" validate:"required,gte=1,dive,required"`
	GoodsIds []int				`json:"goods_ids" validate:"required,gte=1,dive,min=1"`
}

// 正确数据
var parameter1 *GoodsPayParameter = &GoodsPayParameter{
	UserInfo: &User{
		UserId:   101,
		UserName: "12345678901",
		Email:    "123456@qq.com",
	},
	Goods:    []*GoodsPayInfoParameter{&GoodsPayInfoParameter{
		GoodsId:     1001,
		GoodsName:   "三只松鼠干果巧克力100g包邮",
		GoodsAmount: 1,
		Remark:      "",
	}},
	GoodsIds: []int{1},
}
// 错误数据
var parameter2 *GoodsPayParameter = &GoodsPayParameter{
	UserInfo: &User{
		UserId:   0, // 错误1
		UserName: "12345678901000", // 错误2
		Email:    "123456@qq.com",
	},
	Goods:    []*GoodsPayInfoParameter{&GoodsPayInfoParameter{
		GoodsId:     1001,
		GoodsName:   "三只松鼠干果巧克力100g包邮",
		GoodsAmount: 1,
		Remark:      "",
	}},
	GoodsIds: []int{1},
}
// 错误数据
var parameter3 *GoodsPayParameter = &GoodsPayParameter{
	UserInfo: &User{
		UserId:   101,
		UserName: "12345678901",
		Email:    "123456@qq.com",
	},
	Goods:    []*GoodsPayInfoParameter{}, // 错误
	GoodsIds: []int{1},
}
// 错误数据
var parameter4 *GoodsPayParameter = &GoodsPayParameter{
	UserInfo: &User{
		UserId:   101,
		UserName: "12345678901",
		Email:    "123456@qq.com",
	},
	Goods:    []*GoodsPayInfoParameter{&GoodsPayInfoParameter{
		GoodsId:     1001,
		GoodsName:   "三只松鼠干果巧克力100g包邮",
		GoodsAmount: 1,
		Remark:      "",
	}},
	GoodsIds: []int{0}, // 错误
}
// 错误数据
var parameter5 *GoodsPayParameter = &GoodsPayParameter{
	UserInfo: &User{
		UserId:   101,
		UserName: "12345678901",
		Email:    "123456@qq.com",
	},
	Goods:    []*GoodsPayInfoParameter{&GoodsPayInfoParameter{
		GoodsId:     1001,
		GoodsName:   "三只松鼠干果巧克力100g包邮",
		GoodsAmount: 1,
		Remark:      "",
	}},
	GoodsIds: nil, // 错误
}

func TestValidateParameter(t *testing.T) {

	var dataArray = []*GoodsPayParameter{parameter1, parameter2, parameter3, parameter4, parameter5}
	for idx, parameter := range dataArray {
		err := ValidateParameter(parameter)
		if err != nil {
			fmt.Printf("parameter%d validate failed: %v\n", idx, err)
		} else  {
			fmt.Printf("parameter%d validate success\n", idx)
		}
	}

}

func TestValdateParameter2(t *testing.T) {
	user := &User{
		UserId:   102,
		UserName: "18922311001",
		//UserName: "18922311001@qq.com",
		Email:    "102@qq.com",
	}
	err := ValidateParameter(user)
	if err != nil {
		t.Errorf("parameter validate failed: %v\n", err)
	} else  {
		t.Errorf("parameter validate success\n")
	}
}

func TestPutConfig(t *testing.T) {

	var (
		path = "sean.tech/webservice/config"
		endpoints = []string{"127.0.0.1:2379"}
	)

	var g = &config.GlobalConfig{
		RunMode:          "debug",
		JwtSecret:        "23347$040412",
		JwtIssuer:        "sean.tech/webservice",
		JwtExpiresTime:   36,
		RsaServerPubKey:  "",
		RsaServerPriKey:  "",
		RsaClientPubKey:  "",
		RedisHost:        "127.0.0.1:6379",
		RedisPassword:    "",
		RedisMaxIdle:     30,
		RedisMaxActive:   30,
		RedisIdleTimeout: 200,
		EtcdConfigPath:   "/sean.tech/webservice/config",
		EtcdRpcBasePath:  "/sean.tech/webservice/rpcservice",
		EtcdEndPoints:    []string{"127.0.0.1:2379"},
		KafkaAddress:     []string{"127.0.0.1:9092", "127.0.0.1:9093", "127.0.0.1:9094"},
	}
	if err := ValidateParameter(g); err != nil {
		t.Error(err)
	}
	jsonBytes, err := json.Marshal(g)
	if err != nil {
		t.Error(err)
	}
	err = config.Put(path + "/global", endpoints, string(jsonBytes))
	if err != nil {
		t.Error(err)
	}

	var cfg = &config.ModuleConfig{
		App:      config.AppConfig{
			Module:          "user",
			WorkerId:        0,
			RuntimeRootPath: "/Users/Lyra/Desktop/Go/",
			LogSavePath: "log/webservice/",
		},
		//Log:      config.LogConfig{
		//	LogSavePath: "log/webservice/",
		//	LogSaveName: "user",
		//	LogFileExt:  "log",
		//},
		//Upload:   config.UploadConfig{
		//	FilePrefixUrl: "http://127.0.0.1:8001",
		//	FileSavePath:  "uploadfiles/",
		//	FileMaxSize:   10,
		//	FileAllowExts: []string{".jpg", ".jpeg", ".png"},
		//},
		Server:   config.ServerConfig{
			HttpPort:             8811,
			RpcPort:              8812,
			ReadTimeout:          60,
			WriteTimeout:         60,
			RpcPerSecondConnIdle: 500,
		},
		Database: config.DatabaseConfig{
			Type:        "mysql",
			User:        "root",
			Password:    "admin2018",
			HostStr:     "",
			Hosts: map[int]string{0:"127.0.0.1:3306"},
			Name:        "svt_user",
			MaxIdle:     30,
			MaxOpen:     30,
			MaxLifetime: 200,
		},
		Redis:    config.RedisConfig{
			Host:        "127.0.0.1:6379",
			Password:    "",
			MaxIdle:     30,
			MaxActive:   30,
			IdleTimeout: 200,
		},
	}
	if err := ValidateParameter(cfg); err != nil {
		t.Error(err)
	}
	jsonBytes, err = json.Marshal(cfg)
	if err != nil {
		t.Error(err)
	}
	err = config.Put(path + "/192.168.1.52", endpoints, string(jsonBytes))
	if err != nil {
		t.Error(err)
	}
}