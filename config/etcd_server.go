package config

import (
	"context"
	"errors"
	"fmt"
	"github.com/coreos/etcd/clientv3"
)

const (
	PATH_GLOBAL = "/global"
)

func Put(path string, endpoints []string, value string) error {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: _dialTimeout,
	})
	if err != nil {
		return err
	}
	defer cli.Close()

	if resp, err := cli.Put(context.Background(), path, value, clientv3.WithPrevKV()); err != nil {
		return err
	} else {
		_ = resp
		//Log.Println(resp)
	}
	return nil
}

func Delete(path string, endpoints []string) error {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: _dialTimeout,
	})
	if err != nil {
		return err
	}
	defer cli.Close()

	if resp, err := cli.Delete(context.Background(), path); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
}

func GetConfigGlobal(path string, endpoints []string) (*GlobalConfig, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: _dialTimeout,
	})
	if err != nil {
		return nil, err
	}
	defer cli.Close()
	// global get
	global_path := path + PATH_GLOBAL
	if resp, err := cli.Get(context.Background(), global_path); err != nil {
		return nil, err
	} else {
		if len(resp.Kvs) != 1 {
			return nil, errors.New("gloabl config get error:kvs count not only 1")
		}
		kvs := resp.Kvs[0]
		var global *GlobalConfig
		if global, err = globalConfigWithJson(kvs.Value); err != nil {
			return nil, err
		}
		return global, nil
	}
}

func GetConfigModule(path string, endpoints []string, moduleName string, ip string) (*ModuleConfig, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: _dialTimeout,
	})
	if err != nil {
		return nil, err
	}
	defer cli.Close()
	// app get
	var cfg_path = path + "/" + moduleName
	if ip != "" && len(ip) > 0 {
		cfg_path = cfg_path + "/" + ip
	}
	if resp, err := cli.Get(context.Background(), cfg_path); err != nil {
		return nil, err
	} else {
		if len(resp.Kvs) != 1 {
			return nil, errors.New("app config get error:kvs count not only 1")
		}
		kvs := resp.Kvs[0]
		var cfg *ModuleConfig
		if cfg, err = appConfigWithJson(kvs.Value); err != nil {
			return nil, err
		}
		return cfg, err
	}
}