package database

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/sean-tech/webservice/config"
	"log"
	"time"
)

type redisManagerImpl struct {client *redis.Client}

type _redis_link_type int
const (
	_ _redis_link_type = 0
	_redis_link_global = 1
	_redis_link_model = 2
)
func newRedisManagerImpl(link_type _redis_link_type) *redisManagerImpl {
	if link_type == _redis_link_global {
		return &redisManagerImpl{
			client: redis.NewClient(&redis.Options{
				Addr:     config.Redis.Host,
				Password: config.Redis.Password,
				DB:       0,  // use default DB
				IdleTimeout:config.Redis.IdleTimeout,
			}),
		}
	} else if link_type == _redis_link_model {
		return &redisManagerImpl{
			client: redis.NewClient(&redis.Options{
				Addr:     config.Global.RedisHost,
				Password: config.Global.RedisPassword,
				DB:       0,  // use default DB
				IdleTimeout:config.Global.RedisIdleTimeout,
			}),
		}
	} else {
		log.Fatal("redis impl new failed: link type is wrong")
		return nil
	}
}

/**
 * 开启redis并初始化客户端连接
 */
func (this *redisManagerImpl) Open() {
	pong, err := this.client.Ping().Result()
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
	err := this.client.Set(key, value, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

/**
 * 取
 */
func (this *redisManagerImpl) Get(key string) (*string, error) {
	val, err := this.client.Get(key).Result()
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
	this.client.Del(key)
}

/**
 * try lock
 */
func (this *redisManagerImpl) TryLock(key string, expiration time.Duration) (result bool) {
	// lock
	resp := this.client.SetNX(key, 1, expiration)
	lockSuccess, err := resp.Result()
	if err != nil || !lockSuccess {
		return false
	}
	return true
}

func (this *redisManagerImpl) ReleaseLock(key string) (result bool) {
	delResp := this.client.Del(key)
	unlockSuccess, err := delResp.Result()
	if err == nil && unlockSuccess > 0 {
		return true
	} else {
		return false
	}
}