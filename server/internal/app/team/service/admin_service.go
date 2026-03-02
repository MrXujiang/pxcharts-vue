package service

import (
	"mvtable/internal/app/team/model"
	"mvtable/internal/app/team/repo"
	"mvtable/internal/pkg/errorx"
	"mvtable/pkg/log"

	"go.uber.org/zap"
)

func (s *TeamService) AdminGetTeamList(req *model.AdminGetTeamListReq) (*model.AdminGetTeamListRes, error) {
	list, total, err := repo.GetTeamList(req.Page, req.Size, req.SearchWord)
	if err != nil {
		log.Error("get team list error", zap.Error(err))
		return nil, errorx.InternalServerError("获取失败")
	}
	return &model.AdminGetTeamListRes{
		List:  list,
		Total: total,
	}, nil
}
