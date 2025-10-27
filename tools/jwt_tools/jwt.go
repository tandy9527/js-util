package jwt_tools

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tandy9527/js-util/logger"
)

var jwtSecret = []byte("eyJ1aWQiOjEwMDAsImlwIjoiMS4yLjMuNCIsIm5vbmNlIjoiYTJiYzEyMyIsImlhdCI6MTczODUyNzQ0NywiZXhwIjoxNzM4NTI3NjI3LCJpc3MiOiJTbG90TG9iYnkifQ") // 🔒 只保存在服务端

type Claims struct {
	UID   int64  `json:"uid"`
	Nonce string `json:"nonce,omitempty"`
	jwt.RegisteredClaims
}

// 生成 JWT（有效期3分钟）
func GenerateOneToken(uid int64) (string, error) {
	now := time.Now()
	expire := now.Add(30 * time.Second)

	claims := Claims{
		UID:   uid,
		Nonce: GenerateJTI(), // 随机字符串，防止重复
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expire),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(jwtSecret)
	if err != nil {
		logger.Errorf("GenerateOneToken error: %v", err)
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString([]byte(tokenStr)), nil
}

// 解析 JWT
func ParseToken(tokenStr string) (int64, error) {
	decoded, err := base64.RawURLEncoding.DecodeString(tokenStr)
	if err != nil {
		logger.Errorf("ParseToken decoded token error: %v", err)
		return -1, err
	}
	token, err := jwt.ParseWithClaims(string(decoded), &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return -1, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// 检查是否过期（jwt 库会自动检查）
		return claims.UID, nil
	}

	return -1, fmt.Errorf("invalid token")
}

// GenerateJTI 生成安全随机 JTI（32 位十六进制字符串）
func GenerateJTI() string {
	b := make([]byte, 16) // 16 字节 = 128 位
	rand.Read(b)
	return hex.EncodeToString(b)
}
