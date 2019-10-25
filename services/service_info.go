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

func GetServiceInfo(ctx context.Context) *ServiceInfo {
	obj := ctx.Value(KEY_SERVICE_INFO)
	if info, ok := obj.(*ServiceInfo); ok {
		return  info
	}
	return nil
}