package middle

import (
	"github.com/sean-tech/webservice/database"
	"github.com/sean-tech/webservice/logging"
	"sync"
	"time"
)

/** aes key 存储接口 **/
type IAesKeyStorage interface {
	Store(token string,  key string, expiresTime time.Duration)
	Load(token string) ( key string, ok bool)
	Delete(token string)
}

var (
	_aesKeyStorageOnce sync.Once
	_aesKeyStorage *secretRedisAesKeyStorage
)

func GetRedisAesKeyStorage() IAesKeyStorage {
	_aesKeyStorageOnce.Do(func() {
		_aesKeyStorage = new(secretRedisAesKeyStorage)
	})
	return _aesKeyStorage
}


// redis存储实现
type secretRedisAesKeyStorage struct {}

func (this *secretRedisAesKeyStorage) Store(token string, key string, expiresTime time.Duration) {
	err := database.GetGlobalRedis().Set(token, key, expiresTime)
	if err != nil {
		logging.Error(err)
	}
}

func (this *secretRedisAesKeyStorage) Load(token string) (key string, ok bool) {
	key, err := database.GetGlobalRedis().Get(token)
	if err != nil {
		logging.Error(err)
		return "", false
	}
	if key == "" {
		return key,false
	}
	return key, true
}

func (this *secretRedisAesKeyStorage) Delete(token string) {
	database.GetGlobalRedis().Delete(token)
}