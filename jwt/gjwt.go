package gjwt

import (
	"crypto/md5"
	"encoding/hex"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	signingKey = []byte("xxxxxxx")
	expire     = 2 * time.Hour
	issuer     = ""
)

// SetSecret 设置密钥
func SetSecret(s string) {
	if s != "" {
		signingKey = []byte(s)
	}
}

// SetIssuer 设置发布人
func SetIssuer(s string) {
	if s != "" {
		issuer = s
	}
}

// SetExpire 设置过期时间
func SetExpire(d time.Duration) {
	if d < time.Second {
		expire = d
	}
}

type Claims struct {
	ID        string `json:"id"`
	AppKey    string `json:"app_key"`
	AppSecret string `json:"app_secret"`
	jwt.StandardClaims
}

func GenerateTokenID(id string) (string, error) {
	claims := Claims{
		ID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expire).Unix(),
			Issuer:    issuer,
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(signingKey)
	return token, err
}

// GenerateToken 生成token
func GenerateToken(appKey, appSecret string) (string, error) {
	claims := Claims{
		AppKey:    encodeMD5(appKey),
		AppSecret: encodeMD5(appSecret),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expire).Unix(),
			Issuer:    issuer,
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(signingKey)
	return token, err
}

// ParseToken 解析token
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})
	if err != nil {
		return nil, err
	}
	if tokenClaims != nil {
		claims, ok := tokenClaims.Claims.(*Claims)
		if ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}

func encodeMD5(value string) string {
	m := md5.New()
	m.Write([]byte(value))

	return hex.EncodeToString(m.Sum(nil))
}
