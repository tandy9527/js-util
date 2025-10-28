package jwt_tools

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tandy9527/js-util/logger"
	"github.com/tandy9527/js-util/tools/str_tools"
)

var jwtSecret = []byte("eyJ1aWQiOjEwMDAsImlwIjoiMS4yLjMuNCIsIm5vbmNlIjoiYTJiYzEyMyIsImlhdCI6MTczODUyNzQ0NywiZXhwIjoxNzM4NTI3NjI3LCJpc3MiOiJTbG90TG9iYnkifQ")

type Claims struct {
	UID   int64  `json:"uid"`
	Nonce string `json:"nonce,omitempty"`
	IP    string `json:"ip,omitempty"`
	jwt.RegisteredClaims
}

// 生成 JWT
// expire 过期时间,单位:s
// uid 用户ID
// secret 秘钥
// ip 请求ip,预留
func GenerateToken(uid int64, secret, ip string, expire int) (string, error) {

	if str_tools.IsEmpty(secret) {
		return "", fmt.Errorf("secret is empty")
	}
	now := time.Now()

	claims := Claims{
		UID:   uid,
		Nonce: str_tools.RandLetterStr(32), // 随机字符串，防止重复
		IP:    ip,                          //预留
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(expire) * time.Second)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		logger.Errorf("GenerateOneToken error: %v", err)
		return "", err
	}
	return str_tools.Base64Encode(tokenStr), nil
}

// 解析 JWT
func ParseToken(tokenStr, ip string) (int64, error) {
	decoded := str_tools.Base64Decode(tokenStr)
	if str_tools.IsEmpty(decoded) {
		return -1, fmt.Errorf("jwt ParseToken invalid token")
	}

	token, err := jwt.ParseWithClaims(decoded, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return -2, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// 检查 IP 是否匹配
		if str_tools.IsNotEmpty(claims.IP) && claims.IP != ip {
			return -3, fmt.Errorf("jwt ParseToken ip mismatch")
		}
		// 检查是否过期（jwt 库会自动检查）
		return claims.UID, nil
	}

	return -4, fmt.Errorf("jwt ParseToken  token Expired")
}
