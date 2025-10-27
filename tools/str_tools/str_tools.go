package str_tools

import (
	"math/rand"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

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
