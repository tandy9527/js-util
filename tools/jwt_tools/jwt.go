package jwt_tools

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tandy9527/js-util/logger"
	"github.com/tandy9527/js-util/tools/str_tools"
)

type Claims struct {
	U  int64  `json:"u"`
	N  string `json:"n,omitempty"`
	IP string `json:"ip,omitempty"`
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
		U:  uid,
		N:  str_tools.RandLetterStr(8), // 随机字符串，防止重复
		IP: ip,                         //预留
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
	return tokenStr, nil
}

// 解析 JWT
func ParseToken(tokenStr, secret, ip string) (int64, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
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
		return claims.U, nil
	}

	return -4, fmt.Errorf("jwt ParseToken  token Expired")
}
