package middle

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sean-tech/webservice/config"
	"github.com/sean-tech/webservice/encrypt"
	"github.com/sean-tech/webservice/service"
	"log"
	"sync"
	"time"
)

func init() {
	GetSecretManager().startSubscribeToken()
}

type ISecretManager interface {
	SetAesKeyStorage(storage IAesKeyStorage)
	startSubscribeToken()
	GetAesKey(userId uint64) string
	InterceptRsa() gin.HandlerFunc
	InterceptAes() gin.HandlerFunc
}

var (
	_secretManagerOnce sync.Once
	_secretManager *secretManagerImpl
)

func GetSecretManager() ISecretManager {
	_secretManagerOnce.Do(func() {
		_secretManager = new(secretManagerImpl)
		_secretManager.SetAesKeyStorage(GetRedisAesKeyStorage())
	})
	return _secretManager
}

type secretManagerImpl struct {
	aesKeyStorage IAesKeyStorage
}

func (this *secretManagerImpl) SetAesKeyStorage(storage IAesKeyStorage) {
	if storage == nil {
		log.Fatal("aes key storage must not be nil in SetAesKeyStorage method")
		return
	}
	if this.aesKeyStorage != storage {
		this.aesKeyStorage = storage
	}
}

func (this *secretManagerImpl) startSubscribeToken() {
	sub := SubscribeTopic("token")
	go func() {
		for message := range sub.Message {
			if valMap, ok := message.(map[string]interface{}); ok {
				userId := valMap["userId"].(int64)
				expiresTime := valMap["expires"].(time.Duration)
				keyStr := hex.EncodeToString(encrypt.GetAes().GenerateKey())
				this.aesKeyStorage.Store(uint64(userId), keyStr, expiresTime)
			}
		}
	}()
}

func (this *secretManagerImpl) GetAesKey(userId uint64) string {
	key, ok := this.aesKeyStorage.Load(userId)
	if !ok {
		key = hex.EncodeToString(encrypt.GetAes().GenerateKey())
		this.aesKeyStorage.Store(uint64(userId), key, config.Global.JwtExpiresTime)
	}
	return key
}

/**
 * rsa拦截校验
 */
func (this *secretManagerImpl) InterceptRsa() gin.HandlerFunc {
	handler := func(ctx *gin.Context) {
		g := service.Gin{ctx}
		secretStr := ctx.Param("secret")
		encrypted, err := base64.StdEncoding.DecodeString(secretStr)
		if err != nil {
			g.ResponseCode(service.STATUS_CODE_SECRET_CHECK_FAILED, nil)
			ctx.Abort()
			return
		}
		jsonBytes, err := encrypt.GetRsa().Decrypt(config.Global.RsaServerPriKey, encrypted)
		if err != nil {
			g.ResponseCode(service.STATUS_CODE_SECRET_CHECK_FAILED, nil)
			ctx.Abort()
			return
		}
		var jsonMap map[string]interface{} = make(map[string]interface{})
		json.Unmarshal(jsonBytes, &jsonMap)
		var sign = jsonMap["sign"].(string)
		signDatas, _ := base64.StdEncoding.DecodeString(sign)
		delete(jsonMap, "sign")
		jsonBytes, _ = json.Marshal(jsonMap)
		err = encrypt.GetRsa().Verify(config.Global.RsaClientPubKey, jsonBytes, signDatas)
		if err != nil {
			g.ResponseCode(service.STATUS_CODE_SECRET_CHECK_FAILED, nil)
			ctx.Abort()
			return
		}
		ctx.Set(service.KEY_CTX_PARAMS_JSON, jsonBytes)
		// next
		ctx.Next()
	}
	return handler
}

/**
 * aes拦截校验
 */
func (this *secretManagerImpl) InterceptAes() gin.HandlerFunc {
	handler := func(ctx *gin.Context) {
		g := service.Gin{ctx}

		userId := ctx.GetInt64(service.KEY_CTX_USERID)
		key, ok := this.aesKeyStorage.Load(uint64(userId))
		if !ok {
			g.ResponseCode(service.STATUS_CODE_SECRET_CHECK_FAILED, nil)
			ctx.Abort()
		}
		secretStr := ctx.Param("secret")
		encrypted, err := base64.StdEncoding.DecodeString(secretStr)
		if err != nil {
			g.ResponseCode(service.STATUS_CODE_SECRET_CHECK_FAILED, nil)
			ctx.Abort()
			return
		}
		keyBytes, err := hex.DecodeString(key)
		if err != nil {
			g.ResponseCode(service.STATUS_CODE_SECRET_CHECK_FAILED, nil)
			ctx.Abort()
			return
		}
		jsonBytes, err := encrypt.GetAes().DecryptCBC(encrypted, keyBytes)
		if err != nil {
			g.ResponseCode(service.STATUS_CODE_SECRET_CHECK_FAILED, nil)
			ctx.Abort()
			return
		}
		ctx.Set(service.KEY_CTX_PARAMS_JSON, jsonBytes)
		// next
		ctx.Next()
	}
	return handler
}