package middle

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sean-tech/webservice/config"
	"github.com/sean-tech/webservice/service"
	"log"
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
		jwtManager = NewJwtManagerImpl()
	})
	return jwtManager
}



/**
 * jwt实现
 */
type jwtManagerImpl struct{
	tokenStorage IJwtTokenStorage
	publisher *Publisher
}

func NewJwtManagerImpl() *jwtManagerImpl {
	return &jwtManagerImpl{
		tokenStorage: GetRedisTokenStorage(),
		publisher:    NewPublisher("token", 10*time.Second, 1000),
	}
}

/**
 * 设置用户token存储媒介
 */
func (this *jwtManagerImpl) SetJwtTokenStorage(tokenStorage IJwtTokenStorage) {
	if tokenStorage == nil {
		 log.Fatal("jwt token storage must not be nil in SetJwtTokenStorage method")
		return
	}
	if this.tokenStorage != tokenStorage {
		this.tokenStorage = tokenStorage
	}
}

/**
 * 生成token
 */
func (this *jwtManagerImpl) GenerateToken(userId uint64, userName, password string, isAdministrotor bool) (string, error) {
	expireTime := time.Now().Add(config.Global.JwtExpiresTime)
	c := Claims{
		UserId:			userId,
		UserName:       userName,
		Password:       password,
		IsAdministrotor:isAdministrotor,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    config.Global.JwtIssuer,
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	token, err := tokenClaims.SignedString([]byte(config.Global.JwtSecret))
	if err == nil {
		this.tokenStorage.Store(userId, token, config.Global.JwtExpiresTime)
		this.publisher.Publish(map[string]interface{}{"userId":userId, "token":token, "expires":config.Global.JwtExpiresTime})
	}
	return token, err
}

/**
 * 解析token
 */
func (this *jwtManagerImpl) ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Global.JwtSecret), nil
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
	if claims.Issuer != config.Global.JwtIssuer {
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
		id, _ := service.GenerateId(config.App.WorkerId)
		requisition := &service.Requisition{
			ServiceId:    	uint64(id),
			ServicePaths:	make([]string, 5),
			UserId:       	claims.UserId,
			UserName:     	claims.UserName,
			Password:     	claims.Password,
			IsAdministrotor:claims.IsAdministrotor,
		}
		ctx.Set(service.KEY_CTX_REQUISITION, requisition)
		// next
		ctx.Next()
	}
	return handler
}