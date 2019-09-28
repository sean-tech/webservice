package dbutils
//
//
//import (
//	"encoding/json"
//	"github.com/gomodule/redigo/redis"
//	"sean.env/config"
//	"time"
//)
//
//var Redis *redis.Pool
//
//func Setup()  {
//	Redis = &redis.Pool{
//		Dial: func() (redis.Conn, error) {
//			c, err := redis.Dial("tcp", config.RedisSetting.Host)
//			if err != nil {
//				return nil, err
//			}
//			if config.RedisSetting.Password != "" {
//				if _, err := c.Do("AUTH", config.RedisSetting.Password); err != nil {
//					c.Close()
//					return nil, err
//				}
//			}
//			return c, nil
//		},
//		TestOnBorrow: func(c redis.Conn, t time.Time) error {
//			_, err := c.Do("PING")
//			return err
//		},
//		MaxIdle:         config.RedisSetting.MaxIdle,
//		MaxActive:       config.RedisSetting.MaxActive,
//		IdleTimeout:     config.RedisSetting.IdleTimeout,
//		Wait:            false,
//		MaxConnLifetime: 0,
//	}
//}
//
//func Set(key string, data interface{}, time int) (bool, error) {
//	connection := Redis.Get()
//	defer connection.Close()
//	value, err := json.Marshal(data)
//	if err != nil {
//		return false, err
//	}
//	reply, err := redis.Bool(connection.Do("SET", key, value))
//	connection.Do("EXPIRE", key, time)
//	return reply, err
//}
//
//func Exist(key string) bool {
//	connection := Redis.Get()
//	defer connection.Close()
//	exist, err := redis.Bool(connection.Do("EXISTS", key))
//	if err != nil {
//		return false
//	}
//	return exist
//}
//
//func Get(key string) ([]byte, error) {
//	connection := Redis.Get()
//	defer connection.Close()
//	reply, err := redis.Bytes(connection.Do("GET", key))
//	if err != nil {
//		return nil, err
//	}
//	return reply, nil
//}
//
//func Delete(key string) (bool, error) {
//	connection := Redis.Get()
//	defer connection.Close()
//	reply, err := redis.Bool(connection.Do("DEL", key))
//	return reply, err
//}
//
//func LikeDeletes(key string) error {
//	connection := Redis.Get()
//	defer connection.Close()
//
//	keys, err := redis.Strings(connection.Do("KEYS", "*"+key+"*"))
//	if err != nil {
//		return err
//	}
//	for _, key := range keys {
//		_, err := Delete(key)
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}