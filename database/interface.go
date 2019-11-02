package database

import (
	"github.com/jmoiron/sqlx"
	"sync"
	"time"
)

type IMysqlManager interface {
	Open()
	GetDbByUserName(userName string) (db *sqlx.DB, err error)
	GetAllDbs() (dbs []*sqlx.DB)
}

type IRedisManager interface {
	Open()
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (*string, error)
	Delete(key string)
	TryLock(key string, expiration time.Duration) (result bool)
	ReleaseLock(key string) (result bool)
}

var (
	mysqlManagerOnce sync.Once
	redisManagerOnce sync.Once
	mysqlManager IMysqlManager
	redisManager IRedisManager
)

func GetMysqlManager() IMysqlManager {
	mysqlManagerOnce.Do(func() {
		mysqlManager = new(mysqlManagerImpl)
	})
	return mysqlManager
}

func GetRedisManager() IRedisManager {
	redisManagerOnce.Do(func() {
		redisManager = new(redisManagerImpl)
	})
	return redisManager
}

