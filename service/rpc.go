package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/rcrowley/go-metrics"
	"github.com/sean-tech/webservice/config"
	"github.com/sean-tech/webservice/logging"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/serverplugin"
	"gopkg.in/go-playground/validator.v9"
	"log"
	"math"
	"sync"
	"time"
)

/** 服务注册回调函数 **/
type RpcRegisterFunc func(server *server.Server)
/**
 * 启动 服务server
 * registerFunc 服务注册回调函数
 */
func RpcServerServe(registerFunc RpcRegisterFunc) {
	address := fmt.Sprintf(":%d", config.Server.RpcPort)
	s := server.NewServer()
	RegisterPluginEtcd(s, address)
	RegisterPluginRateLimit(s)
	registerFunc(s)
	go func() {
		err := s.Serve("tcp", address)
		if err != nil {
			log.Fatalf("server start error : %v", err)
		}
	}()
}

/**
 * 注册插件，Etcd注册中心，服务发现
 */
func RegisterPluginEtcd(s *server.Server, serviceAddr string)  {
	plugin := &serverplugin.EtcdRegisterPlugin{
		ServiceAddress: "tcp@" + serviceAddr,
		EtcdServers:    config.Etcd.EndPoints,
		BasePath:       config.Etcd.BasePath,
		Metrics:        metrics.NewRegistry(),
		Services:       nil,
		UpdateInterval: time.Minute,
		Options:        nil,
	}
	err := plugin.Start()
	if err != nil {
		log.Fatal(err)
	}
	s.Plugins.Add(plugin)
}

/**
 * 注册插件，限流器，限制客户端连接数
 */
func RegisterPluginRateLimit(s *server.Server)  {
	var fillSpeed float64 = 1.0 / float64(config.Server.RpcPerSecondConnIdle)
	fillInterval := time.Duration(fillSpeed * math.Pow(10, 9))
	plugin := serverplugin.NewRateLimitingPlugin(fillInterval, config.Server.RpcPerSecondConnIdle)
	s.Plugins.Add(plugin)
}

var discoveryMap sync.Map
func getDiscovery(servicePath string) *client.ServiceDiscovery {
	if discovery, ok := discoveryMap.Load(servicePath); ok {
		return discovery.(*client.ServiceDiscovery)
	}
	discovery := client.NewEtcdDiscovery(config.Etcd.BasePath, servicePath, config.Etcd.EndPoints, nil)
	discoveryMap.Store(servicePath, &discovery)
	return &discovery
}

/**
 * 创建rpc调用客户端，基于Etcd服务发现
 */
func CreateRpcClient(servicePath string) client.XClient {
	option := client.DefaultOption
	option.Heartbeat = true
	option.HeartbeatInterval = time.Second
	option.ReadTimeout = config.Server.ReadTimeout
	option.WriteTimeout = config.Server.WriteTimeout
	xclient := client.NewXClient(servicePath, client.Failover, client.RoundRobin, *getDiscovery(servicePath), option)
	return xclient
}

/**
 * 参数绑定验证
 */
func ValidParameter(parameter interface{}) error {
	return ValidParameterWithRegisterFunc(parameter, "", nil)
}

/**
* 参数绑定验证，自定义验证函数注册
*/
func ValidParameterWithRegisterFunc(parameter interface{}, tag string, fn validator.Func) error {
	validate := validator.New()
	if len(tag) < 0  && fn != nil {
		validate.RegisterValidation(tag, fn)
	}
	err := validate.Struct(parameter)
	if err == nil {
		return nil
	}
	if _, ok := err.(*validator.InvalidValidationError); ok {
		logging.Warning(err)
		fmt.Println(err)
		return err
	}
	for _, err := range err.(validator.ValidationErrors) {
		info := fmt.Sprintf("validate err.Namespace:%s Field:%s StructNamespace:%s StructField:%s Tag:%s ActualTag:%s Kind:%v Type:%v Value:%v Param:%s",
			err.Namespace(), err.Field(), err.StructNamespace(), err.StructField(), err.Tag(), err.ActualTag(), err.Kind(), err.Type(), err.Value(), err.Param(),
		)
		logging.Warning(info)
		return errors.New(fmt.Sprintf("the value %s of parameter filed %s, type %s not fit tag %s", err.Value(), err.Field(), err.Type(), err.Tag()))
	}
	return nil
}



const key_service_info = "key_service_info"

/**
 * 服务信息
 */
type ServiceInfo struct {
	ServiceId uint64		`json:"serviceId"`
	ServicePaths []string  	`json:"servicePath"`
	UserId uint64 			`json:"userId"`
	UserName string 		`json:"userName"`
	Password string 		`json:"password"`
	IsAdministrotor bool 	`json:"isAdministrotor"`
}

/**
 * 新建context，并初始化info，绑定serviceId
 */
func NewServiceInfoContext(parentCtx context.Context) context.Context {
	id, _ := GenerateId(config.App.WorkerId)
	info := &ServiceInfo{
		ServiceId:    uint64(id),
		ServicePaths: make([]string, 5),
		UserId:       0,
		UserName:     "",
		Password:     "",
		IsAdministrotor:false,
	}
	return context.WithValue(parentCtx, key_service_info, info)
}

/**
 * 信息获取，获取传输链上context绑定的用户服务信息
 */
func GetServiceInfo(ctx context.Context) *ServiceInfo {
	obj := ctx.Value(key_service_info)
	if info, ok := obj.(*ServiceInfo); ok {
		return  info
	}
	return nil
}

/**
 * 信息校验，token绑定的用户信息同参数传入信息校验，信息不一致说明恶意用户传他人数据渗透
 */
func CheckServiceInfo(ctx context.Context, userId uint64, userName string) bool {
	info := GetServiceInfo(ctx)
	if info == nil {
		return false
	}
	if info.IsAdministrotor {
		return true
	}
	if info.UserId != userId || info.UserName != userName {
		return false
	}
	return true
}