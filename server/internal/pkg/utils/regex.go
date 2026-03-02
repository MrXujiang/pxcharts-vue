package utils

import "regexp"

// ValidateEmail 校验邮箱
func ValidateEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

// ValidPhone 校验手机号
func ValidPhone(phoneNumber string) bool {
	phoneRegex := `^1[3-9]\d{9}$`
	re := regexp.MustCompile(phoneRegex)
	return re.MatchString(phoneNumber)
}
