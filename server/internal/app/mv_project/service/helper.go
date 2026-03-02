package service

import (
	"mvtable/pkg/log"

	"go.uber.org/zap"
)

func hasPermission(userId, projectId string) bool {
	log.Info("check mv project permission", zap.String("userId", userId), zap.String("projectId", projectId))
	return true
}
