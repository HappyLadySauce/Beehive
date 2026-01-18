package passwd

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const (
	// PasswordChars 密码字符集：大小写字母和数字
	PasswordChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// GenerateRandomPassword 生成指定长度的随机密码
// 包含大小写字母和数字
func GenerateRandomPassword(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("password length must be greater than 0")
	}

	password := make([]byte, length)
	charsLen := big.NewInt(int64(len(PasswordChars)))

	for i := range password {
		// 生成随机索引
		randIndex, err := rand.Int(rand.Reader, charsLen)
		if err != nil {
			return "", fmt.Errorf("failed to generate random password: %w", err)
		}
		password[i] = PasswordChars[randIndex.Int64()]
	}

	return string(password), nil
}
