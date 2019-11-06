package database

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/sean-tech/webservice/config"
	"time"
)

var client *redis.Client

type redisManagerImpl struct {}

/**
 * 开启redis并初始化客户端连接
 */
func (this *redisManagerImpl) Open() {
	client = redis.NewClient(&redis.Options{
		Addr:     config.Redis.Host,
		Password: config.Redis.Password,
		DB:       0,  // use default DB
		IdleTimeout:config.Redis.IdleTimeout,
	})
	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>
	// 初始化后通讯失败
	if err != nil {
		panic(err)
	}
}

/**
 * 存
 */
func (this *redisManagerImpl) Set(key string, value interface{}, expiration time.Duration) error {
	err := client.Set(key, value, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

/**
 * 取
 */
func (this *redisManagerImpl) Get(key string) (*string, error) {
	val, err := client.Get(key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return &val, nil
	}
}

/**
 * 删除key
 */
func (this *redisManagerImpl) Delete(key string) {
	client.Del(key)
}

/**
 * try lock
 */
func (this *redisManagerImpl) TryLock(key string, expiration time.Duration) (result bool) {
	// lock
	resp := client.SetNX(key, 1, expiration)
	lockSuccess, err := resp.Result()
	if err != nil || !lockSuccess {
		return false
	}
	return true
}

func (this *redisManagerImpl) ReleaseLock(key string) (result bool) {
	delResp := client.Del(key)
	unlockSuccess, err := delResp.Result()
	if err == nil && unlockSuccess > 0 {
		return true
	} else {
		return false
	}
}