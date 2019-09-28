package services

import (
	"fmt"
	"github.com/rcrowley/go-metrics"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/serverplugin"
	"log"
	"sean.env/config"
	"sean.env/services/validation"
	"time"
)

/** 服务注册回调函数 **/
type ServiceRegisterFunc func(server *server.Server)
/**
 * 启动 服务server
 * registerFunc 服务注册回调函数
 */
func ServiceServe(registerFunc ServiceRegisterFunc) {
	address := fmt.Sprintf(":%d", config.ServerSetting.ServicePort)
	s := server.NewServer()
	RegisterPluginEtcd(s, address)
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
		EtcdServers:    config.EtcdSetting.EndPoints,
		BasePath:       config.EtcdSetting.BasePath,
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

var discovery client.ServiceDiscovery = nil

/**
 * 创建rpc调用客户端，基于Etcd服务发现
 */
func CreateXClient(servicePath string) client.XClient {
	if discovery == nil {
		discovery = client.NewEtcdDiscovery(config.EtcdSetting.BasePath, servicePath, config.EtcdSetting.EndPoints, nil)
	}
	option := client.DefaultOption
	option.Heartbeat = true
	option.HeartbeatInterval = time.Second
	option.ReadTimeout = config.ServerSetting.ReadTimeout
	option.WriteTimeout = config.ServerSetting.WriteTimeout
	xclient := client.NewXClient(servicePath, client.Failover, client.RoundRobin, discovery, option)
	return xclient
}

/**
 * 参数绑定验证
 */
func ValidParameter(parameter interface{}) error {
	valid := validation.Validation{}
	check, err := valid.Valid(parameter)
	if err != nil {
		return err
	}
	if !check {
		validation.MarkErrors(valid.Errors)
		return err
	}
	return nil
}