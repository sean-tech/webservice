package middle

import (
	"github.com/sean-tech/webservice/database"
	"github.com/sean-tech/webservice/logging"
	"sync"
	"time"
)

/** aes key 存储接口 **/
type IAesKeyStorage interface {
	Store(userId uint64,  key string, expiresTime time.Duration)
	Load(userId uint64) ( key string, ok bool)
	Delete(userId uint64)
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

func (this *secretRedisAesKeyStorage) Store(userId uint64, key string, expiresTime time.Duration) {
	err := database.GetGlobalRedis().Set(string(userId), key, expiresTime)
	if err != nil {
		logging.Error(err)
	}
}

func (this *secretRedisAesKeyStorage) Load(userId uint64) (key string, ok bool) {
	keyPointer, err := database.GetGlobalRedis().Get(string(userId))
	if err != nil {
		logging.Error(err)
		return "", false
	}
	return *keyPointer, true
}

func (this *secretRedisAesKeyStorage) Delete(userId uint64) {
	database.GetGlobalRedis().Delete(string(userId))
}