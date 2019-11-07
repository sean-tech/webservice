package config

import (
	"github.com/go-ini/ini"
	"log"
	"strconv"
	"strings"
	"time"
)

type appConfig struct {
	RunMode 		string
	Module 			string
	WorkerId 		int64
	RuntimeRootPath string
	JwtSecret 		string
	JwtIssuer 		string
	JwtExpiresTime 	time.Duration
}
var App = &appConfig{}

type logConfig struct {
	LogSavePath string
	LogSaveName string
	LogFileExt 	string
	TimeFormat 	string
}
var Log = &logConfig{}

type uploadConfig struct {
	FilePrefixUrl	string
	FileSavePath 	string
	FileMaxSize 	int
	FileAllowExts 	[]string
}
var Upload = &uploadConfig{}

type serverConfig struct {
	HttpPort              int
	RpcPort               int
	ReadTimeout           time.Duration
	WriteTimeout          time.Duration
	RpcPerSecondConnIdle  int64
}
var Server = &serverConfig{}

type databaseConfig struct {
	Type 		string
	User 		string
	Password 	string
	HostStr 	string
	Hosts 		map[int]string
	Name 		string
	MaxIdle 	int
	MaxOpen 	int
	MaxLifetime time.Duration
}
var Database = &databaseConfig{}

type redisConfig struct {
	Host        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}
var Redis = &redisConfig{}

type etcdConfig struct {
	KeyPrefix 	string
	BasePath 	string
	EndPointStr string
	EndPoints 	[]string
}
var Etcd = &etcdConfig{}

type kafkaConfig struct {
	AddressStr 	string
	Address 	[]string
}
var Kafka = &kafkaConfig{}


func Setup(configFilePath string) {
	// load
	Cfg, err := ini.Load(configFilePath)
	if err != nil {
		log.Fatalf("Fail to parse 'config.ini': %v", err)
	}

	// App convert
	err = Cfg.Section("app").MapTo(App)
	if err != nil {
		log.Fatalf("Cfg.MapTo App err: %v", err)
	}
	App.JwtExpiresTime = App.JwtExpiresTime * time.Hour

	// Log convert
	err = Cfg.Section("log").MapTo(Log)
	if err != nil {
		log.Fatalf("Cfg.MapTo Log err: %v", err)
	}

	// Upload convert
	err = Cfg.Section("upload").MapTo(Upload)
	if err != nil {
		log.Fatalf("Cfg.MapTo Upload err: %v", err)
	}
	Upload.FileMaxSize = Upload.FileMaxSize * 1024 * 1024

	// Server convert
	err = Cfg.Section("server").MapTo(Server)
	if err != nil {
		log.Fatalf("Cfg.MapTo Server err: %v", err)
	}
	Server.ReadTimeout = Server.ReadTimeout * time.Second
	Server.WriteTimeout = Server.ReadTimeout * time.Second

	// Database convert
	err = Cfg.Section("database").MapTo(Database)
	if err != nil {
		log.Fatalf("Cfg.MapTo Database err: %v", err)
	}
	Database.MaxLifetime = Database.MaxLifetime * time.Second
	idHosts := strings.Split(Database.HostStr, ", ")
	Database.Hosts = make(map[int]string, len(idHosts))
	for _, idHost := range idHosts {
		seps := strings.Split(idHost, "-")
		id, _ := strconv.Atoi(seps[0])
		host := seps[1]
		Database.Hosts[id] = host
	}

	// Redis convert
	err = Cfg.Section("redis").MapTo(Redis)
	if err != nil {
		log.Fatalf("Cfg.MapTo Redis err: %v", err)
	}
	Redis.IdleTimeout = Redis.IdleTimeout * time.Second

	// etcdConfig convert
	err = Cfg.Section("etcd").MapTo(Etcd)
	if err != nil {
		log.Fatalf("Cfg.MapTo Redis err: %v", err)
	}
	Etcd.EndPoints = strings.Split(Etcd.EndPointStr, ", ")

	// kafkaConfig convert
	err = Cfg.Section("kafka").MapTo(Kafka)
	if err != nil {
		log.Fatalf("Cfg.MapTo Redis err: %v", err)
	}
	Kafka.Address = strings.Split(Kafka.AddressStr, ", ")
}