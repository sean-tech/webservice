package middle

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sean-tech/webservice/config"
	"github.com/sean-tech/webservice/services"
	"sync"
	"time"
)


var jwtSecret = []byte(config.AppSetting.JwtSecret)

type Claims struct {
	UserId uint64 	`json:"userId"`
	UserName string `json:"userName"`
	Password string `json:"password"`
	jwt.StandardClaims
}

// 用户token存储映射，用户当前token唯一性保证
var userCurrentTokenMap sync.Map

func GenerateToken(userId uint64, userName, password string) (string, error) {
	expireTime := time.Now().Add(3 * time.Hour)
	claims := Claims{
		UserId:			userId,
		UserName:       userName,
		Password:       password,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "Team",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	if err == nil {
		userCurrentTokenMap.Store(userId, token)
	}
	return token, err
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}

func Jwt() gin.HandlerFunc {
	handler := func(ctx *gin.Context) {
		g := services.Gin{ctx}
		// token
		token := ctx.GetHeader("Authorization")
		if token == "" {
			g.ResponseCode(services.STATUS_CODE_AUTH_CHECK_TOKEN_FAILED, nil)
			ctx.Abort()
			return
		}
		// token parse
		claims, err := ParseToken(token)
		if err != nil {
			g.ResponseCode(services.STATUS_CODE_AUTH_CHECK_TOKEN_FAILED, nil)
			ctx.Abort()
			return
		}
		// time judge
		if time.Now().Unix() > claims.ExpiresAt {
			g.ResponseCode(services.STATUS_CODE_AUTH_CHECK_TOKEN_TIMEOUT, nil)
			ctx.Abort()
			return
		}
		// current token union judge
		if savedToken, ok := userCurrentTokenMap.Load(claims.UserId); !ok || savedToken != token {
			g.ResponseCode(services.STATUS_CODE_AUTH_CHECK_TOKEN_TIMEOUT, nil)
			ctx.Abort()
			return
		}
		// transfer
		ctx.Set(services.KEY_CTX_USERID, claims.UserId)
		ctx.Set(services.KEY_CTX_USERNAME, claims.UserName)
		ctx.Set(services.KEY_CTX_PASSWORD, claims.Password)
		// next
		ctx.Next()
	}
	return handler
}