package services

import (
	"context"
	"github.com/sean-tech/webservice/config"
)

const KEY_SERVICE_INFO = "KEY_SERVICE_INFO"

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
func NewContext(parentCtx context.Context) context.Context {
	id, _ := GenerateId(config.AppSetting.WorkerId)
	info := &ServiceInfo{
		ServiceId:    uint64(id),
		ServicePaths: make([]string, 5),
		UserId:       0,
		UserName:     "",
		Password:     "",
		IsAdministrotor:false,
	}
	return context.WithValue(parentCtx, KEY_SERVICE_INFO, info)
}

/**
 * 信息获取，获取传输链上context绑定的用户服务信息
 */
func GetServiceInfo(ctx context.Context) *ServiceInfo {
	obj := ctx.Value(KEY_SERVICE_INFO)
	if info, ok := obj.(*ServiceInfo); ok {
		return  info
	}
	return nil
}

/**
 * 信息校验，token绑定的用户信息同参数传入信息校验，信息不一致说明恶意用户传他人数据渗透
 */
func ServiceInfoCheck(ctx context.Context, userId uint64, userName string) bool {
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