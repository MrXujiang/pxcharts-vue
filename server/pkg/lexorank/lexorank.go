package lexorank

import "strings"

// 字母表
const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// MaxString 极大值字符串
const MaxString = "zzzzzzzzzz"

// MinString 极小值字符串
const MinString = "aaaaaaaaaa"

// Between 计算两个字符串之间的中间值
func Between(a, b string) string {
	// 确保a < b（字典序）
	if a >= b {
		// 如果a >= b，无法计算中间值，返回a后面加一个字符
		return a + "a"
	}

	maxLen := max(len(b), len(a))

	// 补齐到相同长度
	aPadded := padRight(a, maxLen, 'a')
	bPadded := padRight(b, maxLen, 'z')

	res := make([]rune, 0, maxLen)
	carry := 0

	for i := range maxLen {
		ai := strings.IndexRune(alphabet, rune(aPadded[i]))
		bi := strings.IndexRune(alphabet, rune(bPadded[i]))

		// 如果字符不在字母表中，使用默认值
		if ai == -1 {
			ai = 0
		}
		if bi == -1 {
			bi = 25
		}

		sum := ai + bi + carry
		mid := sum / 2

		// 确保 mid 在有效范围内 [0, 25]
		if mid < 0 {
			mid = 0
		} else if mid > 25 {
			mid = 25
		}

		res = append(res, rune(alphabet[mid]))

		// 计算下一轮的进位
		if sum%2 != 0 {
			carry = 26
		} else {
			carry = 0
		}
	}

	result := string(res)

	// 如果结果等于a或b，说明需要调整
	if result <= a || result >= b {
		// 找到第一个可以调整的位置
		runes := []rune(result)
		adjusted := false

		// 从右向左查找可以增加的位置
		for i := len(runes) - 1; i >= 0; i-- {
			idx := strings.IndexRune(alphabet, runes[i])
			if idx < 25 { // 如果不是最大值，可以增加
				runes[i] = rune(alphabet[idx+1])
				// 将右边的所有字符设置为'a'
				for j := i + 1; j < len(runes); j++ {
					runes[j] = 'a'
				}
				adjusted = true
				break
			}
		}

		// 如果无法调整（所有字符都是最大值），说明需要增加长度
		if !adjusted {
			result = result + "a"
		} else {
			result = string(runes)
		}
	}

	return result
}

// 辅助函数：右侧补齐字符
func padRight(s string, length int, pad rune) string {
	if len(s) >= length {
		return s
	}
	return s + strings.Repeat(string(pad), length-len(s))
}
