package middle

import (
	"encoding/base64"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/sean-tech/webservice/config"
	"github.com/sean-tech/webservice/encrypt"
	"github.com/sean-tech/webservice/service"
	"log"
	"sync"
	"time"
)

type ISecretManager interface {
	SetAesKeyStorage(storage IAesKeyStorage)
	GetAesKey(token string) (key string, ok bool)
	newAesKeyByChangedToken(newToken string, expiresTime time.Duration)
	deleteAesKeyByOriginToken(originToken string)
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

func (this *secretManagerImpl) GetAesKey(token string) (key string, ok bool) {
	return this.aesKeyStorage.Load(token)
}

func (this *secretManagerImpl) newAesKeyByChangedToken(newToken string, expiresTime time.Duration) {
	keyStr := hex.EncodeToString(encrypt.GetAes().GenerateKey())
	this.aesKeyStorage.Store(newToken, keyStr, expiresTime)
}

func (this *secretManagerImpl) deleteAesKeyByOriginToken(originToken string) {
	this.aesKeyStorage.Delete(originToken)
}

type SecretParams struct {
	Secret string	`json:"secret" validate:"required,base64"`
} 

/**
 * rsa拦截校验
 */
func (this *secretManagerImpl) InterceptRsa() gin.HandlerFunc {
	handler := func(ctx *gin.Context) {
		g := service.Gin{ctx}

		var sign = ctx.GetHeader("sign")
		var params SecretParams
		if err := g.Ctx.Bind(&params); err != nil {
			g.ResponseCode(service.STATUS_CODE_SECRET_CHECK_FAILED, nil)
			ctx.Abort()
			return
		}
		if err := service.ValidateParameter(params); err != nil {
			g.ResponseCode(service.STATUS_CODE_INVALID_PARAMS, nil)
			ctx.Abort()
			return
		}

		encrypted, err := base64.StdEncoding.DecodeString(params.Secret)
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
		signDatas, _ := base64.StdEncoding.DecodeString(sign)
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

		var params SecretParams
		if err := g.Ctx.Bind(&params); err != nil {
			g.ResponseCode(service.STATUS_CODE_SECRET_CHECK_FAILED, nil)
			ctx.Abort()
			return
		}
		if err := service.ValidateParameter(params); err != nil {
			g.ResponseCode(service.STATUS_CODE_INVALID_PARAMS, nil)
			ctx.Abort()
			return
		}

		rq := service.GetRequisition(ctx)
		key, ok := this.aesKeyStorage.Load(rq.Token)
		if !ok {
			g.ResponseCode(service.STATUS_CODE_SECRET_CHECK_FAILED, nil)
			ctx.Abort()
		}
		encrypted, err := base64.StdEncoding.DecodeString(params.Secret)
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
		ctx.Set(service.KEY_CTX_AES_KEY, keyBytes)
		// next
		ctx.Next()
	}
	return handler
}