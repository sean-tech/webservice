package config

import (
	"encoding/json"
	"fmt"
	"github.com/sean-tech/webservice/utils"
	"testing"
)

var (
	path = "sean.tech/webservice/config"
	endpoints = []string{"localhost:2379"}
)

func TestPutConfig(t *testing.T) {

	var g = &GlobalConfig{
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
		EtcdEndPoints:    []string{"localhost:2379"},
		KafkaAddress:     []string{"localhost:9092", "localhost:9093", "localhost:9094"},
	}
	jsonBytes, err := json.Marshal(g)
	if err != nil {
		t.Error(err)
	}
	err = Put(path + "/global", endpoints, string(jsonBytes))
	if err != nil {
		t.Error(err)
	}

	var cfg = &ModuleConfig{
		App:      AppConfig{
			Module:          "user",
			WorkerId:        0,
			RuntimeRootPath: "/Users/Lyra/Desktop/Go/",
		},
		//Log:      LogConfig{
		//	LogSavePath: "Log/webservice/",
		//	LogSaveName: "user",
		//	LogFileExt:  "Log",
		//},
		//Upload:   UploadConfig{
		//	FilePrefixUrl: "http://127.0.0.1:8001",
		//	FileSavePath:  "uploadfiles/",
		//	FileMaxSize:   10,
		//	FileAllowExts: []string{".jpg", ".jpeg", ".png"},
		//},
		Server:   ServerConfig{
			HttpPort:             8811,
			RpcPort:              8812,
			ReadTimeout:          60,
			WriteTimeout:         60,
			RpcPerSecondConnIdle: 500,
		},
		Database: DatabaseConfig{
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
		Redis:    RedisConfig{
			Host:        "127.0.0.1:6379",
			Password:    "",
			MaxIdle:     30,
			MaxActive:   30,
			IdleTimeout: 200,
		},
	}
	jsonBytes, err = json.Marshal(cfg)
	if err != nil {
		t.Error(err)
	}
	err = Put(path + "/192.168.1.52", endpoints, string(jsonBytes))
	if err != nil {
		t.Error(err)
	}
}

func TestLoad(t *testing.T) {
	g, err := GetConfigGlobal(path, endpoints)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(g)
	app, err := GetConfigModule(path, endpoints, "192.168.1.52")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(app)
}

func TestDelete(t *testing.T) {
	Delete(path, endpoints)
	fmt.Println(utils.Ip.GetLocalIP())
}