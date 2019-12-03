package config

import (
	"encoding/json"
	"time"
)

type ModuleConfig struct {
	App      AppConfig      `json:"app" validate:"required"`
	//Log      LogConfig      `json:"log" validate:"required"`
	//Upload   UploadConfig   `json:"upload" validate:"required"`
	Server   ServerConfig   `json:"server" validate:"required"`
	Database DatabaseConfig `json:"database" validate:"required"`
	Redis    RedisConfig    `json:"redis" validate:"required"`
}

func (this *ModuleConfig) bestow()  {
	App = &this.App
	//Log = &this.Log
	//Upload = &this.Upload
	Server = &this.Server
	Database = &this.Database
	Redis = &this.Redis
}

func globalConfigWithJson(value []byte) (*GlobalConfig, error) {
	var global = &GlobalConfig{}
	if err := json.Unmarshal(value, global); err != nil {
		return nil, err
	}
	global.JwtExpiresTime = Global.JwtExpiresTime * time.Hour
	global.RedisIdleTimeout = Global.RedisIdleTimeout * time.Second
	return global, nil
}

func appConfigWithJson(value []byte) (*ModuleConfig, error) {
	var cfg = new(ModuleConfig)
	if err := json.Unmarshal(value, cfg); err != nil {
		return nil, err
	}
	//cfg.Upload.FileMaxSize = cfg.Upload.FileMaxSize * 1024 * 1024
	cfg.Server.ReadTimeout = cfg.Server.ReadTimeout * time.Second
	cfg.Server.WriteTimeout = cfg.Server.ReadTimeout * time.Second
	cfg.Database.MaxLifetime = cfg.Database.MaxLifetime * time.Second
	cfg.Redis.IdleTimeout = cfg.Redis.IdleTimeout * time.Second
	return cfg, nil
}