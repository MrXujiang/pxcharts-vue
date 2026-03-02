package constants

import "time"

const (
	RedisResourceAction    = "resourceAction:"
	RedisAccessPasswordTTL = 24 * 60 * time.Minute
)

const (
	RedisRegisterVerifyCodePrefix = "email:verifyCode:"
	RedisRegisterVerifyCodeTTL    = 60 * time.Second // 秒
)
