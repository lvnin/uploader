package utils

import (
	"errors"
	request "uploader/constant"
	"uploader/constant/errcode"
	"uploader/global"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// CustomClaims 自定义 Payload 信息
type UserData struct {
	UserID   int    `json:"id"`
	Avatar   string `json:"avatar"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
}

type CustomClaims struct {
	User UserData
	jwt.StandardClaims
}

type JWT struct {
	signKey []byte // Jwt 密钥
}

var (
	ErrTokenExpired     = errors.New("token is expired")        // 令牌过期
	ErrTokenNotValidYet = errors.New("token not active yet")    // 令牌未生效
	ErrTokenMalformed   = errors.New("that's not even a token") // 令牌不完整
	ErrTokenInvalid     = errors.New("token is invalid")        // 无效令牌
)

func NewJWT() *JWT {
	return &JWT{
		signKey: []byte(global.ServerConfig.Jwt.SigningKey),
	}
}

func (j *JWT) CreateToken(claims CustomClaims) (token string, err error) {
	withClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return withClaims.SignedString(j.signKey)
}

func (j *JWT) ParseToken(token string) (*CustomClaims, error) {
	withClaims, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.signKey, nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, ErrTokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, ErrTokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, ErrTokenNotValidYet
			} else {
				return nil, ErrTokenInvalid
			}
		}

		return nil, ErrTokenInvalid
	}

	if withClaims == nil {
		return nil, ErrTokenInvalid
	}

	if claims, ok := withClaims.Claims.(*CustomClaims); ok { // 验证成功
		return claims, nil
	}

	return nil, ErrTokenInvalid

}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("x-token")
		if token == "" {
			request.ResponseWrapper.Fail(c, errcode.ErrorCode("ERRCODE_TOKEN_INVALID"))
			c.Abort()
			return
		}

		j := NewJWT()
		claims, err := j.ParseToken(token)

		if err != nil {
			if err == ErrTokenExpired {
				request.ResponseWrapper.Fail(c, errcode.ErrorCode("ERRCODE_TOKEN_EXPIRED"))
			} else if err == ErrTokenNotValidYet {
				request.ResponseWrapper.Fail(c, errcode.ErrorCode("ERRCODE_TOKEN_NOTVALIDYET"))
			} else if err == ErrTokenMalformed {
				request.ResponseWrapper.Fail(c, errcode.ErrorCode("ERRCODE_TOKEN_MALFORMED"))
			} else {
				request.ResponseWrapper.Fail(c, errcode.ErrorCode("ERRCODE_TOKEN_INVALID"))
			}

			c.Abort()
			return
		}

		c.Set("claims", claims)
	}
}

func GetToken(c *gin.Context) string {
	token, _ := c.Cookie("x-token")
	if token == "" {
		token = c.Request.Header.Get("x-token")
	}
	return token
}

func GetClaims(c *gin.Context) (*CustomClaims, error) {
	token := GetToken(c)
	j := NewJWT()
	claims, err := j.ParseToken(token)
	if err != nil {
		global.Logger.Error("从Gin的Context中获取从jwt解析信息失败, 请检查请求头是否存在x-token且claims是否为规定结构")
	}
	return claims, err
}

func GetUserID(c *gin.Context) int {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return 0
		} else {
			return cl.User.UserID
		}
	} else {
		waitUse := claims.(*CustomClaims)
		return waitUse.User.UserID
	}
}
