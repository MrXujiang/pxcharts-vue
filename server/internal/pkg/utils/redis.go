package utils

import "strings"

func GetRedisKey(prefix string, args ...string) string {
	return prefix + strings.Join(args, ":")
}
