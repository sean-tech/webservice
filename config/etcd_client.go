package config

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"github.com/sean-tech/webservice/utils"
	"log"
	"time"
)

var (
	_dialTimeout    = 5 * time.Second
	_requestTimeout = 2 * time.Second
)

func loadEtcdConfig(path string, endpoints []string, moduleName string, isAdmin bool)  {
	if len(path) <= 0 {
		log.Fatal("failed to load config from Etcd:path is nil")
	}
	if len(endpoints) <= 0 {
		log.Fatal("failed to load config from Etcd:_endpoints is nil")
	}

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: _dialTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}
	//defer cli.Close()

	loadGlobal(cli, path)
	switch isAdmin {
	case true:
		loadModuleWithName(cli, path, moduleName)
		break
	case false:
		loadModuleWithIp(cli, path, moduleName)
		break
	}
}

func loadGlobal(cli *clientv3.Client, path string) {
	// global get
	global_path := path + PATH_GLOBAL
	if resp, err := cli.Get(context.Background(), global_path); err != nil {
		log.Fatal(err)
	} else {
		if len(resp.Kvs) != 1 {
			log.Fatal("gloabl config get error:kvs count not only 1")
		}
		kvs := resp.Kvs[0]
		var g *GlobalConfig
		if g, err = globalConfigWithJson(kvs.Value); err != nil {
			log.Fatal(err)
		}
		Global = g
	}
	// global watch
	go watchGlobal(cli, global_path)
}

func watchGlobal(cli *clientv3.Client, path string) {
	rch := cli.Watch(context.Background(), path, clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			if ev.Type == clientv3.EventTypePut {
				if g, err := globalConfigWithJson(ev.Kv.Value); err == nil {
					Global = g
				}
			}
		}
	}
}

func loadModuleWithIp(cli *clientv3.Client, path string, moduleName string) {
	// module get
	var cfg_path = ""
	for _, ip := range utils.Ip.GetLocalIP() {
		this_path := path + "/" + moduleName + "/" + ip
		if resp, err := cli.Get(context.Background(), this_path); err != nil {
			log.Fatal(err)
		} else {
			if len(resp.Kvs) != 1 {
				continue
			}
			cfg_path = this_path
			kvs := resp.Kvs[0]
			var cfg *ModuleConfig
			if cfg, err = appConfigWithJson(kvs.Value); err != nil {
				log.Fatal(err)
			}
			cfg.bestow()
		}
	}
	if cfg_path == "" {
		log.Fatal("app config get error:not found config path")
	}
	go watchModuleCfg(cli, cfg_path)
}

func loadModuleWithName(cli *clientv3.Client, path string, moduleName string) {
	cfg_path := path + "/" + moduleName
	if resp, err := cli.Get(context.Background(), cfg_path); err != nil {
		log.Fatal(err)
	} else {
		if len(resp.Kvs) != 1 {
			log.Fatal("module config get error:kvs count not only 1")
		}
		kvs := resp.Kvs[0]
		var cfg *ModuleConfig
		if cfg, err = appConfigWithJson(kvs.Value); err != nil {
			log.Fatal(err)
		}
		cfg.bestow()
	}
	go watchModuleCfg(cli, cfg_path)
}

func watchModuleCfg(cli *clientv3.Client, path string) {
	rch := cli.Watch(context.Background(), path, clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			if ev.Type == clientv3.EventTypePut {
				if cfg, err := appConfigWithJson(ev.Kv.Value); err == nil {
					cfg.bestow()
				}
			}
		}
	}
}




