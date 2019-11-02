package middle

import (
	"github.com/sean-tech/webservice/database"
	"github.com/sean-tech/webservice/logging"
	"sync"
	"time"
)

/** jwt token 存储接口 **/
type IJwtTokenStorage interface {
	Store(userId uint64, token string, expiresTime time.Duration)
	Load(userId uint64) (token string, ok bool)
	Delete(userId uint64)
}

var (
	memeoryTokenStorageOnce sync.Once
	memeoryTokenStorage IJwtTokenStorage
	redisTokenStorageOnce sync.Once
	redisTokenStorage IJwtTokenStorage
)
/**
 * 获取jwt token 内存存储实例
 */
func GetMemeoryTokenStorage() IJwtTokenStorage {
	memeoryTokenStorageOnce.Do(func() {
		memeoryTokenStorage = new(JwtMemeoryTokenStorage)
	})
	return memeoryTokenStorage
}

/**
 * 获取jwt token Redis存储实例
 */
func GetRedisTokenStorage() IJwtTokenStorage {
	redisTokenStorageOnce.Do(func() {
		redisTokenStorage = new(JwtRedisTokenStorage)
	})
	return redisTokenStorage
}



// 内存存储实现
type JwtMemeoryTokenStorage struct {
	// 用户token存储映射，用户当前token唯一性保证
	userCurrentTokenMap sync.Map
}

func (this *JwtMemeoryTokenStorage) Store(userId uint64, token string, expiresTime time.Duration) {
	this.userCurrentTokenMap.Store(userId, token)
	// 定时删除
	select {
	case <- time.After(expiresTime):
		this.Delete(userId)
	}
}

func (this *JwtMemeoryTokenStorage) Load(userId uint64) (token string, ok bool) {
	tokenInter, ok := this.userCurrentTokenMap.Load(userId)
	return tokenInter.(string), ok
}

func (this *JwtMemeoryTokenStorage) Delete(userId uint64) {
	this.userCurrentTokenMap.Delete(userId)
}



// redis存储实现
type JwtRedisTokenStorage struct {}

func (this *JwtRedisTokenStorage) Store(userId uint64, token string, expiresTime time.Duration) {
	err := database.GetRedisManager().Set(string(userId), token, expiresTime)
	if err != nil {
		logging.Error(err)
	}
}

func (this *JwtRedisTokenStorage) Load(userId uint64) (token string, ok bool) {
	tokenPointer, err := database.GetRedisManager().Get(string(userId))
	if err != nil {
		logging.Error(err)
		return "", false
	}
	return *tokenPointer, true
}

func (this *JwtRedisTokenStorage) Delete(userId uint64) {
	database.GetRedisManager().Delete(string(userId))
}