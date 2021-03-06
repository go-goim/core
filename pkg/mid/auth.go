package mid

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

	"github.com/go-goim/core/pkg/log"
	"github.com/go-goim/core/pkg/types"
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
	UserID types.ID `json:"uid"`
	jwt.RegisteredClaims
}

func (c *JwtClaims) Valid() error {
	if err := c.RegisteredClaims.Valid(); err != nil {
		return err
	}

	if c.UserID == 0 {
		return fmt.Errorf("missing user id")
	}

	return nil
}

func newJwtClaims(userID types.ID) *JwtClaims {
	return &JwtClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireTime)),
		},
	}
}

func NewJwtToken(userID types.ID) (string, error) {
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

func SetJwtToHeader(c *gin.Context, userID types.ID) error {
	token, err := NewJwtToken(userID)
	if err != nil {
		return err
	}

	c.Header("Authorization", fmt.Sprintf("Bearer %s", token))
	return nil
}

func AuthJwt(c *gin.Context) {
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

	if err != nil {
		log.Error("convert user id to int64 failed", "err", err, "user id", claims.UserID)
		_ = c.AbortWithError(401, fmt.Errorf("unknown uid format")) // nolint: errcheck
		return
	}

	c.Set(uidKey, claims.UserID)
	c.Next()
}

const (
	uidKey = "uid"
)

func GetUID(c *gin.Context) types.ID {
	if uid, ok := c.Get(uidKey); ok {
		return uid.(types.ID)
	}
	return 0
}
