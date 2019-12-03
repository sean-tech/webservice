package config

import (
	"flag"
	"github.com/go-ini/ini"
	"github.com/sean-tech/webservice/fileutils"
	"log"
	"strconv"
	"strings"
	"time"
)

type GlobalConfig struct{
	RunMode 			string			`json:"run_mode" validate:"required,oneof=debug release"`
	// jwt
	JwtSecret 			string			`json:"jwt_secret" validate:"required,gte=1"`
	JwtIssuer 			string			`json:"jwt_issuer" validate:"required,gte=1"`
	JwtExpiresTime 		time.Duration	`json:"jwt_expires_time" validate:"required,gte=1"`
	// rsa
	RsaServerPubKey 	string			`json:"rsa_server_pub_key"`
	RsaServerPriKey 	string			`json:"rsa_server_pri_key"`
	RsaClientPubKey 	string			`json:"rsa_client_pub_key"`
	// redis
	RedisHost 			string			`json:"redis_host" validate:"required,tcp_addr"`
	RedisPassword 		string			`json:"redis_password" validate:"gte=0"`
	RedisMaxIdle 		int				`json:"redis_max_idle" validate:"required,min=1"`
	RedisMaxActive 		int				`json:"redis_max_active" validate:"required,min=1"`
	RedisIdleTimeout 	time.Duration	`json:"redis_idle_timeout" validate:"required,gte=1"`
	// etcd
	EtcdConfigPath 		string			`json:"etcd_config_path" validate:"required,gte=1"`
	EtcdRpcBasePath 	string			`json:"etcd_rpc_base_path" validate:"required,gte=1"`
	EtcdEndPoints 		[]string		`json:"etcd_end_points" validate:"required,gte=1,dive,tcp_addr"`
	// kafka
	KafkaAddress 		[]string		`json:"kafka_address" validate:"required,gte=1,dive,tcp_addr"`
}
var Global = &GlobalConfig{}

type AppConfig struct{
	Module 			string	`json:"module" validate:"required,gte=1"`
	WorkerId 		int64	`json:"worker_id" validate:"min=0"`
	RuntimeRootPath string	`json:"runtime_root_path" validate:"required,gt=1"`
	LogSavePath string		`json:"log_save_path" validate:"required,gt=1"`
}
var App = &AppConfig{}

//type LogConfig struct{
//	LogSavePath string	`json:"log_save_path" validate:"required,gt=1"`
//	LogSaveName string	`json:"log_save_name" validate:"required,gte=1"`
//	LogFileExt 	string	`json:"log_file_ext" validate:"required,oneof=log txt"`
//}
//var Log = &LogConfig{}
//
//type UploadConfig struct{
//	FilePrefixUrl	string	`json:"file_prefix_url"`
//	FileSavePath 	string	`json:"file_save_path"`
//	FileMaxSize 	int		`json:"file_max_size"`
//	FileAllowExts 	[]string`json:"file_allow_exts"`
//}
//var Upload = &UploadConfig{}

type ServerConfig struct{
	HttpPort              int			`json:"http_port" validate:"required,min=1,max=10000"`
	RpcPort               int			`json:"rpc_port" validate:"required,min=1,max=10000"`
	ReadTimeout           time.Duration	`json:"read_timeout" validate:"required,gte=1"`
	WriteTimeout          time.Duration	`json:"write_timeout" validate:"required,gte=1"`
	RpcPerSecondConnIdle  int64			`json:"rpc_per_second_conn_idle" validate:"required,gte=1"`
}
var Server = &ServerConfig{}

type DatabaseConfig struct{
	Type 		string			`json:"type" validate:"required,oneof=mysql"`
	User 		string			`json:"user" validate:"required,gte=1"`
	Password 	string			`json:"password" validate:"required,gte=1"`
	HostStr 	string			`json:"host_str"`
	Hosts 		map[int]string	`json:"hosts" validate:"required,gte=1,dive,keys,min=0,endkeys,tcp_addr"`
	Name 		string			`json:"name" validate:"required,gte=1"`
	MaxIdle 	int				`json:"max_idle" validate:"required,min=1"`
	MaxOpen 	int				`json:"max_open" validate:"required,min=1"`
	MaxLifetime time.Duration	`json:"max_lifetime" validate:"required,gte=1"`
}
var Database = &DatabaseConfig{}

type RedisConfig struct{
	Host        string			`json:"host" validate:"required,tcp_addr"`
	Password    string			`json:"password" validate:"gte=0"`
	MaxIdle     int				`json:"max_idle" validate:"required,min=1"`
	MaxActive   int				`json:"max_active" validate:"required,min=1"`
	IdleTimeout time.Duration	`json:"idle_timeout" validate:"required,gte=1"`
}
var Redis = &RedisConfig{}



func loadLocalConfig(configFilePath string) {
	// load
	Cfg, err := ini.Load(configFilePath)
	if err != nil {
		log.Fatalf("Fail to parse 'config.ini': %v", err)
	}

	// Global convert
	err = Cfg.Section("global").MapTo(Global)
	if err != nil {
		log.Fatalf("Cfg.MapTo Global err: %v", err)
	}
	Global.JwtExpiresTime = Global.JwtExpiresTime * time.Hour
	Global.RedisIdleTimeout = Global.RedisIdleTimeout * time.Second
	// rsa
	rsaServerPubKeyBuf, err := fileutils.ReadFile(Global.RsaServerPubKey)
	if err != nil {
		log.Fatal(err)
	}
	Global.RsaServerPubKey = string(rsaServerPubKeyBuf)
	rsaServerPriKeyBuf, err := fileutils.ReadFile(Global.RsaServerPriKey)
	if err != nil {
		log.Fatal(err)
	}
	Global.RsaServerPriKey = string(rsaServerPriKeyBuf)
	rsaClientPubKeyBuf, err := fileutils.ReadFile(Global.RsaClientPubKey)
	if err != nil {
		log.Fatal(err)
	}
	Global.RsaClientPubKey = string(rsaClientPubKeyBuf)

	// App convert
	err = Cfg.Section("app").MapTo(App)
	if err != nil {
		log.Fatalf("Cfg.MapTo App err: %v", err)
	}

	//// Log convert
	//err = Cfg.Section("log").MapTo(Log)
	//if err != nil {
	//	log.Fatalf("Cfg.MapTo Log err: %v", err)
	//}
	//
	//// Upload convert
	//err = Cfg.Section("upload").MapTo(Upload)
	//if err != nil {
	//	log.Fatalf("Cfg.MapTo Upload err: %v", err)
	//}
	//Upload.FileMaxSize = Upload.FileMaxSize * 1024 * 1024

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
	Database.HostStr = strings.Replace(Database.HostStr, " ", "", -1)
	idHosts := strings.Split(Database.HostStr, ",")
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
}


/**
 * 初始化config，通过etcd注册中心
 */
func Setup(defaultEtcdPath, defaultEndPointsStr string) {
	// cfp Etcd path
	etcd_path_usage := "please use -etcdpath to pointing at Etcd path for config center to init webservice config."
	etcd_path := flag.String("etcdpath", defaultEtcdPath, etcd_path_usage)
	*etcd_path = strings.Replace(*etcd_path, " ", "", -1)
	// cfp Etcd endpoints
	etcd_endpoints_str_usage := "please use -endpoints to pointing at Etcd endpoints for config center to init webservice config. if more, use ',' to separate"
	etcd_endpoints_str := flag.String("endpoints", defaultEndPointsStr, etcd_endpoints_str_usage)
	*etcd_endpoints_str = strings.Replace(*etcd_endpoints_str, " ", "", -1)

	flag.Parse()

	if etcd_path == nil || len(*etcd_path) <= 0 {
		log.Fatal(etcd_path_usage)
	}
	if etcd_endpoints_str == nil || len(*etcd_endpoints_str) <= 0 {
		log.Fatal(etcd_endpoints_str_usage)
	}
	loadEtcdConfig(*etcd_path, strings.Split(*etcd_endpoints_str, ","))
}

/**
 * 初始化config，通过本地配置文件
 */
func SetupFromLocal(defaultFilePath string) {
	// cfp local
	config_file_path_usage := "please use -cfp to pointing at local config file _path for webservice"
	config_file_path := flag.String("cfp", defaultFilePath, config_file_path_usage)
	*config_file_path = strings.Replace(*config_file_path, " ", "", -1)

	flag.Parse()

	if config_file_path == nil || len(*config_file_path) <= 0 {
		log.Fatal(config_file_path_usage)
	}
	loadLocalConfig(*config_file_path)
}