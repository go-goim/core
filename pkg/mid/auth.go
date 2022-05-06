package mid

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

	"github.com/yusank/goim/pkg/log"
)

var (
	jwtHmacSecret = []byte("secret")
	expireTime    = time.Hour * 24
)

func SetJwtHmacSecret(secret string) {
	log.Debug("set jwt hmac secret", "secret", secret)
	jwtHmacSecret = []byte(secret)
}

type JwtClaims struct {
	UserID string `json:"uid"`
	jwt.RegisteredClaims
}

func (c *JwtClaims) Valid() error {
	if err := c.RegisteredClaims.Valid(); err != nil {
		return err
	}

	if c.UserID == "" {
		return fmt.Errorf("missing user id")
	}

	return nil
}

func newJwtClaims(userID string) *JwtClaims {
	return &JwtClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireTime)),
		},
	}
}

func NewJwtToken(userID string) (string, error) {
	claims := newJwtClaims(userID)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtHmacSecret)
}

func ParseJwtToken(tokenString string) (*JwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtHmacSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*JwtClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

func SetJwtToHeader(c *gin.Context, userID string) error {
	token, err := NewJwtToken(userID)
	if err != nil {
		return err
	}

	c.Header("Authorization", fmt.Sprintf("Bearer %s", token))
	return nil
}

func AuthJwtCookie(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.AbortWithStatus(401)
		return
	}

	if !strings.HasPrefix(token, "Bearer ") {
		_ = c.AbortWithError(401, fmt.Errorf("invalid token prefix")) // nolint: errcheck
		return
	}

	token = strings.TrimPrefix(token, "Bearer ")
	claims, err := ParseJwtToken(token)
	if err != nil {
		log.Error("parse jwt token failed", "err", err)
		_ = c.AbortWithError(401, err) // nolint: errcheck
		return
	}

	c.Set("uid", claims.UserID)
	c.Next()
}
