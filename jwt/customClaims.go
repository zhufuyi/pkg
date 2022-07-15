package jwt

import (
	"time"

	"github.com/golang-jwt/jwt"
)

/* Custom Claims 在StandardClaims的payload基础上添加自定义字段 */

// CustomClaims 自定义Claims
type CustomClaims struct {
	Uid  string `json:"uid"`
	Role string `json:"role"`
	jwt.StandardClaims
}

// GenerateToken 生成token
func GenerateToken(uid string, role ...string) (string, error) {
	if opt == nil {
		return "", initError
	}

	roleVal := ""
	if len(role) > 1 {
		roleVal = role[0]
	}
	claims := CustomClaims{
		uid,
		roleVal,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(opt.expire).Unix(),
			Issuer:    opt.issuer,
		},
	}

	token := jwt.NewWithClaims(opt.signingMethod, claims)
	return token.SignedString(opt.signingKey)
}

// VerifyToken 验证token
func VerifyToken(tokenString string) (*CustomClaims, error) {
	if opt == nil {
		return nil, initError
	}

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return opt.signingKey, nil
	})
	if err != nil {
		ve, ok := err.(*jwt.ValidationError)
		if ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, formatErr
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, expiredErr
			} else if ve.Errors&jwt.ValidationErrorUnverifiable != 0 {
				return nil, unverifiableErr
			} else if ve.Errors&jwt.ValidationErrorSignatureInvalid != 0 {
				return nil, signatureErr
			} else {
				return nil, ve
			}
		}
		return nil, signatureErr
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, signatureErr
}
