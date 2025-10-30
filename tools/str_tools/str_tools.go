package str_tools

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"slices"
	"strings"

	"github.com/tandy9527/js-util/logger"
)

const EMPTY = ""

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// StrSplingInt 字符串拼接int64
func StrSplingInt(s string, n int64) string {
	return s + fmt.Sprintf("%d", n)
}

func IsEmpty(s string) bool {
	return s == ""
}

func IsNotEmpty(s string) bool {
	return !IsEmpty(s)
}
func IsAllEmpty(ss ...string) bool {
	return !slices.ContainsFunc(ss, IsNotEmpty)
}
func IsAllNotEmpty(ss ...string) bool {
	return !slices.ContainsFunc(ss, IsEmpty)
}

// RandLetterString 生成随机英文字母字符串
func RandLetterStr(n int) string {
	result := make([]byte, n)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

// 生成指定长度的随机字符串
func RandNumStr(min, max int) string {
	length := rand.Intn(max-min+1) + min
	return RandLetterStr(length)
}

// Base64Encode base64编码
func Base64Encode(str string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(str))
}

// Base64Decode base64解码
func Base64Decode(str string) string {
	decoded, err := base64.RawURLEncoding.DecodeString(str)
	if err != nil {
		logger.Errorf("ParseToken decoded token error: %v", err)
		return ""
	}
	return string(decoded)
}

// JSON 转map
func JsonToMao(jsonstr string) map[string]any {
	if IsEmpty(jsonstr) {
		panic("json string is empty")
	}
	data := make(map[string]any)
	err := json.Unmarshal([]byte(jsonstr), &data)
	if err != nil {
		panic(fmt.Sprintf("json to map error: %v", err.Error()))
	}
	return data
}

// GetRealIP 获取真实客户端 IP，支持多层代理
func GetRealIP(r *http.Request) string {
	// 先从 X-Forwarded-For 获取
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		// 可能有多个 IP，用逗号分隔
		ips := strings.Split(xff, ",")
		for _, ip := range ips {
			ip = strings.TrimSpace(ip)
			if ip != "" && !isPrivateIP(ip) {
				return ip // 返回第一个非内网 IP
			}
		}
	}

	// 再尝试 X-Real-IP
	xrip := r.Header.Get("X-Real-IP")
	if xrip != "" && !isPrivateIP(xrip) {
		return xrip
	}

	// 最后取 RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

// isPrivateIP 判断是否是私网 IP
func isPrivateIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	privateBlocks := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"127.0.0.0/8",
		"::1/128",
	}

	for _, block := range privateBlocks {
		_, cidr, _ := net.ParseCIDR(block)
		if cidr.Contains(ip) {
			return true
		}
	}
	return false
}
