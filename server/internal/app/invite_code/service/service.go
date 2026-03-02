package service

import (
	"crypto/rand"
	"go.uber.org/zap"
	"math/big"
	"mvtable/internal/app/invite_code/model"
	"mvtable/internal/app/invite_code/repo"
	"mvtable/internal/pkg/errorx"
	"mvtable/internal/storage/db"
	"mvtable/pkg/log"
)

type InviteCodeService struct{}

func NewInviteCodeService() *InviteCodeService {
	return &InviteCodeService{}
}

func generateInviteCode(length int) (string, error) {
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	code := make([]byte, length)
	for i := range code {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		code[i] = charset[n.Int64()]
	}
	return string(code), nil
}

func (s *InviteCodeService) AdminBatchCreat(userId string, req *model.BatchCreateInviteCodeReq) ([]string, error) {
	var (
		maxRetryPerCode = 5 // 每个邀请码最多重试5次
		inviteCodes     = make([]string, 0, req.Count)
	)

	if req.Count <= 0 {
		return inviteCodes, errorx.BadRequest("数量必须大于0")
	}

	if req.Count > 20 {
		return inviteCodes, errorx.BadRequest("单次创建数量不能超过20个")
	}

	for i := 0; i < req.Count; i++ {
		var (
			value string
			err   error
			retry int
		)

		for retry = 0; retry < maxRetryPerCode; retry++ {
			value, err = generateInviteCode(8)
			if err != nil {
				log.Error("generate invite code failed", zap.Error(err))
				return inviteCodes, errorx.InternalServerError("创建失败")
			}
			if err = db.GetDB().Create(&model.InviteCode{
				Value:  value,
				UserID: userId,
			}).Error; err != nil {
				if db.IsUniqueViolation(err) {
					continue
				} else {
					log.Error("create invite code failed", zap.Error(err))
					return inviteCodes, errorx.InternalServerError("创建失败")
				}
			}

			inviteCodes = append(inviteCodes, value)
			break
		}
	}
	return inviteCodes, nil
}

func (s *InviteCodeService) AdminGetList(req *model.GetInviteCodeListReq) (*model.GetInviteCodeListRes, error) {
	list, total, err := repo.List(req.Page, req.Size, req.IsUsed)

	if err != nil {
		log.Error("获取邀请码列表失败", zap.Error(err))
		return nil, errorx.InternalServerError("获取失败")
	}

	return &model.GetInviteCodeListRes{
		Total: total,
		List:  list,
	}, nil
}
