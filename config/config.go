package config

import (
	"github.com/go-ini/ini"
	"log"
	"strconv"
	"strings"
	"time"
)

type appSetting struct {
	RunMode string
	Module string
	WorkerId int64
	RuntimeRootPath string
	JwtSecret string
	JwtIssuer string
}
var AppSetting = &appSetting{}

type logSetting struct {
	LogSavePath string
	LogSaveName string
	LogFileExt string
	TimeFormat string
}
var LogSetting = &logSetting{}

type uploadSetting struct {
	FilePrefixUrl string
	FileSavePath string
	FileMaxSize int
	FileAllowExts []string
}
var UploadSetting = &uploadSetting{}

type serverSetting struct {
	ApiPort int
	ServicePort int
	ReadTimeout time.Duration
	WriteTimeout time.Duration
}
var ServerSetting = &serverSetting{}

type databaseSetting struct {
	Type string
	User string
	Password string
	HostStr string
	Hosts map[int]string
	Name string
	MaxIdle int
	MaxOpen int
	MaxLifetime time.Duration
}
var DatabaseSetting = &databaseSetting{}

type redis struct {
	Host        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}
var RedisSetting = &redis{}

type etcdSetting struct {
	KeyPrefix string
	BasePath string
	EndPointStr string
	EndPoints []string
}
var EtcdSetting = &etcdSetting{}

type kafkaSetting struct {
	AddressStr string
	Address []string
}
var KafkaSetting = &kafkaSetting{}


func Setup(configFilePath string) {
	// load
	Cfg, err := ini.Load(configFilePath)
	if err != nil {
		log.Fatalf("Fail to parse 'config.ini': %v", err)
	}

	// AppSetting convert
	err = Cfg.Section("app").MapTo(AppSetting)
	if err != nil {
		log.Fatalf("Cfg.MapTo AppSetting err: %v", err)
	}

	// LogSetting convert
	err = Cfg.Section("log").MapTo(LogSetting)
	if err != nil {
		log.Fatalf("Cfg.MapTo LogSetting err: %v", err)
	}

	// UploadSetting convert
	err = Cfg.Section("upload").MapTo(UploadSetting)
	if err != nil {
		log.Fatalf("Cfg.MapTo UploadSetting err: %v", err)
	}
	UploadSetting.FileMaxSize = UploadSetting.FileMaxSize * 1024 * 1024

	// ServerSetting convert
	err = Cfg.Section("server").MapTo(ServerSetting)
	if err != nil {
		log.Fatalf("Cfg.MapTo ServerSetting err: %v", err)
	}
	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.ReadTimeout * time.Second

	// DatabaseSetting convert
	err = Cfg.Section("database").MapTo(DatabaseSetting)
	if err != nil {
		log.Fatalf("Cfg.MapTo DatabaseSetting err: %v", err)
	}
	DatabaseSetting.MaxLifetime = DatabaseSetting.MaxLifetime * time.Second
	idHosts := strings.Split(DatabaseSetting.HostStr, ", ")
	DatabaseSetting.Hosts = make(map[int]string, len(idHosts))
	for _, idHost := range idHosts {
		seps := strings.Split(idHost, "-")
		id, _ := strconv.Atoi(seps[0])
		host := seps[1]
		DatabaseSetting.Hosts[id] = host
	}

	// RedisSetting convert
	err = Cfg.Section("redis").MapTo(RedisSetting)
	if err != nil {
		log.Fatalf("Cfg.MapTo RedisSetting err: %v", err)
	}
	RedisSetting.IdleTimeout = RedisSetting.IdleTimeout * time.Second

	// etcdSetting convert
	err = Cfg.Section("etcd").MapTo(EtcdSetting)
	if err != nil {
		log.Fatalf("Cfg.MapTo RedisSetting err: %v", err)
	}
	EtcdSetting.EndPoints = strings.Split(EtcdSetting.EndPointStr, ", ")

	// kafkaSetting convert
	err = Cfg.Section("kafka").MapTo(KafkaSetting)
	if err != nil {
		log.Fatalf("Cfg.MapTo RedisSetting err: %v", err)
	}
	KafkaSetting.Address = strings.Split(KafkaSetting.AddressStr, ", ")
}