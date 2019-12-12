package service

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/rcrowley/go-metrics"
	"github.com/sean-tech/webservice/config"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/serverplugin"
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
	cert, err := tls.LoadX509KeyPair(config.App.RuntimeRootPath + config.App.TLSCerPath, config.App.RuntimeRootPath + config.App.TLSKeyPath)
	if err != nil {
		log.Fatal(err)
		return
	}
	tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}
	s := server.NewServer(server.WithTLSConfig(tlsConfig))

	address := fmt.Sprintf(":%d", config.Server.RpcPort)
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
		EtcdServers:    config.Global.EtcdEndPoints,
		BasePath:       config.Global.EtcdRpcBasePath,
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
	discovery := client.NewEtcdDiscovery(config.Global.EtcdRpcBasePath, servicePath, config.Global.EtcdEndPoints, nil)
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
 * 服务信息
 */
type Requisition struct {
	ServiceId uint64		`json:"serviceId"`
	ServicePaths []string  	`json:"servicePath"`
	Token string			`json:"token"`
	UserId uint64 			`json:"userId"`
	UserName string 		`json:"userName"`
	Password string 		`json:"password"`
	IsAdministrotor bool 	`json:"isAdministrotor"`
}

/**
 * 信息获取，获取传输链上context绑定的用户请求调用信息
 */
func GetRequisition(ctx context.Context) *Requisition {
	obj := ctx.Value(KEY_CTX_REQUISITION)
	if info, ok := obj.(*Requisition); ok {
		return  info
	}
	return nil
}

/**
 * 信息校验，token绑定的用户信息同参数传入信息校验，信息不一致说明恶意用户传他人数据渗透
 */
func CheckServiceInfo(ctx context.Context, userId uint64, userName string) bool {
	info := GetRequisition(ctx)
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