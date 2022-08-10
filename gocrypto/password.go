package gocrypto

import "golang.org/x/crypto/bcrypt"

// HashAndSaltPassword 对密码加盐哈希
func HashAndSaltPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), err
}

// VerifyPassword 验证密码和密文是否匹配
func VerifyPassword(password string, hashed string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)) == nil
}
