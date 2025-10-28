package str_tools

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"slices"

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
