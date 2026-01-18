package passwd

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const (
	// SaltLength 盐值长度（字节）
	SaltLength = 16
	// BcryptCost bcrypt 加密成本
	BcryptCost = 10
)

// GenerateSalt 生成随机盐值
func GenerateSalt() (string, error) {
	salt := make([]byte, SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}
	return hex.EncodeToString(salt), nil
}

// HashPassword 使用盐值加密密码
// 使用 bcrypt 算法，将盐值与密码结合后生成哈希
func HashPassword(password, salt string) (string, error) {
	// 将盐值与密码结合
	saltedPassword := password + salt

	// 使用 bcrypt 生成哈希
	hash, err := bcrypt.GenerateFromPassword([]byte(saltedPassword), BcryptCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hash), nil
}

// VerifyPassword 验证密码是否正确
func VerifyPassword(password, salt, hashedPassword string) bool {
	// 将盐值与密码结合
	saltedPassword := password + salt

	// 使用 bcrypt 验证
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(saltedPassword))
	return err == nil
}
