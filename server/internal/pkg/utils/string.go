package utils

import (
	"math/rand"
	"strings"
)

// GenRandomNumber 生成随机数字验证码
func GenRandomNumber(length int) string {
	if length <= 0 {
		return ""
	}
	digits := []rune("0123456789")
	code := make([]rune, length)
	for i := 0; i < length; i++ {
		code[i] = digits[randomInt(0, 10)]
	}
	return string(code)
}

func randomInt(min, max int) int {
	return rand.Intn(max-min) + min
}

func SplitAndFilter(s, sep string) []string {
	result := []string{}
	for part := range strings.SplitSeq(s, sep) {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
