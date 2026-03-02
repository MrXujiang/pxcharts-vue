package service

import (
	model2 "mvtable/internal/app/user_team/model"
	"mvtable/internal/pkg/constants"
	"mvtable/internal/storage/db"
	"mvtable/pkg/log"
	"slices"

	"go.uber.org/zap"
)

func hasPermission(userId, teamId string, args ...constants.TeamIdentity) bool {
	userTeam, err := db.Get[model2.UserTeam](db.GetDB(), map[string]any{
		"user_id": userId,
		"team_id": teamId,
	})
	if err != nil {
		log.Error("get user team error", zap.Error(err))
		return false
	}

	if userTeam == nil {
		return false
	}

	if !slices.Contains(args, userTeam.Identity) {
		return false
	}

	return true
}
