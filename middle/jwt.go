package middle

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sean-tech/webservice/config"
	"github.com/sean-tech/webservice/service"
	"sync"
	"time"
)

type Claims struct {
	UserId uint64 			`json:"userId"`
	UserName string 		`json:"userName"`
	Password string 		`json:"password"`
	IsAdministrotor bool 	`json:"isAdministrotor"`
	jwt.StandardClaims
}

type IJwtManager interface {
	SetJwtTokenStorage(tokenStorage IJwtTokenStorage)
	GenerateToken(userId uint64, userName, password string, isAdministrotor bool) (string, error)
	ParseToken(token string) (*Claims, error)
	InterceptCheck() gin.HandlerFunc
}

var (
	jwtManagerOnce sync.Once
	jwtManager IJwtManager
)
/**
 * 获取jwt管理实例
 */
func GetJwtManager() IJwtManager {
	jwtManagerOnce.Do(func() {
		jwtManager = new(jwtManagerImpl)
		jwtManager.SetJwtTokenStorage(GetRedisTokenStorage())
	})
	return jwtManager
}



var (
	// 密钥
	jwtSecret = []byte(config.AppSetting.JwtSecret)
	// 发行者
	jwtIssuer = config.AppSetting.JwtIssuer
	// 过期时间
	JwtExpiresTime = config.AppSetting.JwtExpiresTime
)

type jwtManagerImpl struct{
	tokenStorage IJwtTokenStorage
}

/**
 * 生成token
 */
func (this *jwtManagerImpl) SetJwtTokenStorage(tokenStorage IJwtTokenStorage) {
	if this.tokenStorage != tokenStorage {
		this.tokenStorage = tokenStorage
	}
}

/**
 * 生成token
 */
func (this *jwtManagerImpl) GenerateToken(userId uint64, userName, password string, isAdministrotor bool) (string, error) {
	expireTime := time.Now().Add(JwtExpiresTime)
	c := Claims{
		UserId:			userId,
		UserName:       userName,
		Password:       password,
		IsAdministrotor:isAdministrotor,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    jwtIssuer,
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	token, err := tokenClaims.SignedString(jwtSecret)
	if err == nil {
		this.tokenStorage.Store(userId, token, JwtExpiresTime)
	}
	return token, err
}

/**
 * 解析token
 */
func (this *jwtManagerImpl) ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if tokenClaims == nil {
		return nil, errors.New("token parse to nil")
	}
	if !tokenClaims.Valid {
		return nil, errors.New("token parsed not valid")
	}
	claims, ok := tokenClaims.Claims.(*Claims)
	if !ok {
		return nil, errors.New("token parse failed")
	}
	if claims.Issuer != jwtIssuer {
		return nil, errors.New("token parsed issuer not right")
	}
	return claims, nil
}

/**
 * jwt拦截校验
 */
func (this *jwtManagerImpl) InterceptCheck() gin.HandlerFunc {
	handler := func(ctx *gin.Context) {
		g := service.Gin{ctx}
		// token
		token := ctx.GetHeader("Authorization")
		if token == "" {
			g.ResponseCode(service.STATUS_CODE_AUTH_CHECK_TOKEN_FAILED, nil)
			ctx.Abort()
			return
		}
		// token parse
		claims, err := this.ParseToken(token)
		if err != nil {
			g.ResponseCode(service.STATUS_CODE_AUTH_CHECK_TOKEN_FAILED, nil)
			ctx.Abort()
			return
		}
		// time judge
		if time.Now().Unix() > claims.ExpiresAt {
			g.ResponseCode(service.STATUS_CODE_AUTH_CHECK_TOKEN_TIMEOUT, nil)
			ctx.Abort()
			return
		}
		// current token union judge
		if savedToken, ok := this.tokenStorage.Load(claims.UserId); !ok || savedToken != token {
			g.ResponseCode(service.STATUS_CODE_AUTH_CHECK_TOKEN_TIMEOUT, nil)
			ctx.Abort()
			return
		}
		// transfer
		ctx.Set(service.KEY_CTX_USERID, claims.UserId)
		ctx.Set(service.KEY_CTX_USERNAME, claims.UserName)
		ctx.Set(service.KEY_CTX_PASSWORD, claims.Password)
		ctx.Set(service.KEY_CTX_IS_ADMINISTROTOR, claims.IsAdministrotor)
		// next
		ctx.Next()
	}
	return handler
}