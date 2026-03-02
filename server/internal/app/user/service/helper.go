package service

import (
	"context"
	"encoding/base64"
	"errors"
	model3 "mvtable/internal/app/team/model"
	model4 "mvtable/internal/app/user_team/model"
	"mvtable/internal/pkg/constants"
	"mvtable/internal/pkg/errorx"
	"mvtable/internal/storage/db"
	"mvtable/internal/storage/redis"
	"mvtable/pkg/jwt"
	"mvtable/pkg/log"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func genTokens(userID string) (token string, refreshToken string, err error) {
	// 生成JWT令牌
	token, err = jwt.GetInstance().GenerateAccessToken(userID)
	if err != nil {
		return "", "", err
	}

	// 生成刷新令牌
	refreshToken, err = jwt.GetInstance().GenerateRefreshToken(userID)
	if err != nil {
		return "", "", err
	}

	return token, refreshToken, err
}

func judgeVerifyCode(email, verifyCode string) error {
	key := constants.RedisRegisterVerifyCodePrefix + email
	cacheCode, err := redis.Get(context.Background(), key)
	if err != nil && !errors.Is(err, redis.Nil) {
		log.Error("get verify code error", zap.Error(err))
		return errorx.InternalServerError("注册失败")
	}
	if cacheCode == "" {
		return errorx.New(errorx.ErrOperationFailed, "验证码已失效")
	}
	if cacheCode != verifyCode {
		return errorx.New(errorx.ErrOperationFailed, "验证码错误")
	}

	return nil
}

func genNickname(s string) string {
	u := base64.URLEncoding.EncodeToString([]byte(s))
	return "用户" + u[:8]
}

func createTeam(tx *gorm.DB, userId string) error {
	// 创建个人团队
	team := &model3.Team{
		UserID: userId,
		Name:   "个人空间",
	}
	if err := db.Create(tx, team); err != nil {
		log.Error("create team error", zap.Error(err))
		return errorx.InternalServerError("注册失败")
	}

	// 创建团队关联关系
	if err := db.Create(tx, &model4.UserTeam{
		UserID:   userId,
		TeamID:   team.ID,
		Identity: constants.TeamIdentityCreator,
	}); err != nil {
		log.Error("create user team error", zap.Error(err))
		return errorx.InternalServerError("注册失败")
	}

	return nil
}
