package service

import (
	"mvtable/internal/app/team/model"
	"mvtable/internal/app/team/repo"
	model3 "mvtable/internal/app/user/model"
	model2 "mvtable/internal/app/user_team/model"
	"mvtable/internal/pkg/constants"
	"mvtable/internal/pkg/errorx"
	"mvtable/internal/storage/db"
	"mvtable/pkg/log"
	"slices"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type TeamService struct{}

func NewTeamService() *TeamService {
	return &TeamService{}
}

func (s *TeamService) Create(userId string, req *model.CreateTeamReq) error {
	// 判断用户创建团队数是否到达最大值
	_, total, err := db.List[model2.UserTeam](db.GetDB(), 0, 0, map[string]any{"user_id": userId, "identity": constants.TeamIdentityCreator}, nil)
	if err != nil {
		log.Error("query user team error", zap.Error(err))
		return errorx.InternalServerError("创建失败")
	}

	user, err := db.Get[model3.User](db.GetDB(), map[string]any{"id": userId})
	if err != nil {
		log.Error("get user error", zap.Error(err))
		return errorx.InternalServerError("创建失败")
	}

	if total >= int64(user.MaxTeamCount) {
		return errorx.New(errorx.ErrOperationFailed, "创建团队数已达上限")
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		team := &model.Team{
			UserID:      userId,
			Name:        req.Name,
			Description: req.Description,
			Logo:        req.Logo,
		}
		if err = db.Create(tx, team); err != nil {
			log.Error("create team error", zap.Error(err))
			return errorx.InternalServerError("创建失败")
		}

		if err = db.Create(tx, &model2.UserTeam{
			UserID:   userId,
			TeamID:   team.ID,
			Identity: constants.TeamIdentityCreator,
		}); err != nil {
			log.Error("create user team error", zap.Error(err))
			return errorx.InternalServerError("创建失败")
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *TeamService) Update(userId string, req *model.UpdateTeamReq) error {
	// 创建人和管理员都能更新团队信息
	if !hasPermission(userId, req.ID, constants.TeamIdentityCreator, constants.TeamIdentityManager) {
		return errorx.New(errorx.ErrNoPermission, "无权限")
	}

	if err := db.Update(db.GetDB(), &model.Team{
		Name:        req.Name,
		Description: req.Description,
		Logo:        req.Logo,
	}, map[string]any{
		"id": req.ID,
	}, "name", "description", "logo"); err != nil {
		log.Error("update team error", zap.Error(err))
		return errorx.InternalServerError("更新失败")
	}

	return nil
}

func (s *TeamService) Delete(userId string, req *model.DeleteTeamReq) error {
	// 创建人才能删除团队
	if !hasPermission(userId, req.ID, constants.TeamIdentityCreator) {
		return errorx.New(errorx.ErrNoPermission, "无权限")
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := db.Delete[model.Team](tx, map[string]any{
			"id": req.ID,
		}); err != nil {
			log.Error("delete team error", zap.Error(err))
			return errorx.InternalServerError("删除失败")
		}

		// 删除团队成员
		if err := db.Delete[model2.UserTeam](tx, map[string]any{
			"team_id": req.ID,
		}); err != nil {
			log.Error("delete user team error", zap.Error(err))
			return errorx.InternalServerError("删除失败")
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *TeamService) List(userId string) (*model.ListTeamReq, error) {
	// 获取用户团队列表
	userTeams, _, err := db.List[model2.UserTeam](db.GetDB(), 0, 0, map[string]any{"user_id": userId}, nil)
	if err != nil {
		log.Error("query user team error", zap.Error(err))
		return nil, errorx.InternalServerError("获取失败")
	}

	teamIds := make([]string, len(userTeams))
	identityMap := make(map[string]constants.TeamIdentity)

	for _, userTeam := range userTeams {
		teamIds = append(teamIds, userTeam.TeamID)
		identityMap[userTeam.TeamID] = userTeam.Identity
	}

	teams, _, err := db.List[model.Team](db.GetDB(), 0, 0, map[string]any{"id": []any{"IN", teamIds}}, nil)
	if err != nil {
		log.Error("query team error", zap.Error(err))
		return nil, errorx.InternalServerError("获取失败")
	}

	userIds := make([]string, 0, len(teams))
	for _, v := range teams {
		if !slices.Contains(userIds, v.UserID) {
			userIds = append(userIds, v.UserID)
		}
	}

	users, _, err := db.List[model3.User](db.GetDB(), 0, 0, map[string]any{"id": []any{"IN", userIds}}, nil)
	if err != nil {
		log.Error("query user error", zap.Error(err))
		return nil, errorx.InternalServerError("获取失败")
	}

	userMap := make(map[string]string, len(users))
	for _, v := range users {
		userMap[v.ID] = v.Nickname
	}

	list := make([]model.TeamListItem, len(teams))
	for i, v := range teams {
		list[i].FillFromTeam(v)
		list[i].Identity = identityMap[v.ID]
		list[i].Creator = userMap[v.UserID]

		// 根据team查询成员id
		userTeams, _, err := db.List[model2.UserTeam](db.GetDB(), 0, 0, map[string]any{"team_id": v.ID}, nil)
		if err != nil {
			log.Error("query user team error", zap.Error(err))
			return nil, errorx.InternalServerError("获取失败")
		}
		list[i].MemberCount = len(userTeams)

		// 取前四名成员，获取头像
		userIds := make([]string, 4)
		for _, v := range userTeams {
			userIds = append(userIds, v.UserID)
		}
		users, _, err = db.List[model3.User](db.GetDB(), 0, 0, map[string]any{"id": []any{"IN", userIds}}, nil)
		if err != nil {
			log.Error("query user error", zap.Error(err))
			return nil, errorx.InternalServerError("获取失败")
		}
		avatarList := make([]string, 0, len(users))
		for _, v := range users {
			avatarList = append(avatarList, v.Avatar)
		}
		list[i].DisplayAvatarList = avatarList
	}

	return &model.ListTeamReq{
		List: list,
	}, nil
}

func (s *TeamService) GetMemberList(userId string, req *model.GetMemberListReq) (*model.GetMemberListRes, error) {
	if !hasPermission(userId, req.TeamID, constants.TeamIdentityCreator, constants.TeamIdentityManager, constants.TeamIdentityMember) {
		return nil, errorx.New(errorx.ErrNoPermission, "无权限")
	}

	list, total, err := repo.GetMemberList(req.TeamID, req.SearchWord, req.Page, req.Size)
	if err != nil {
		log.Error("get member list error", zap.Error(err))
		return nil, errorx.InternalServerError("获取失败")
	}

	return &model.GetMemberListRes{
		List:  list,
		Total: total,
	}, nil
}

func (s *TeamService) SearchUser(userId string, req *model.SearchUserReq) (*model.SearchUserRes, error) {
	list, err := repo.SearchUser(userId, req.SearchWord)
	if err != nil {
		log.Error("search user error", zap.Error(err))
		return nil, errorx.InternalServerError("搜索失败")
	}

	return &model.SearchUserRes{
		List: list,
	}, nil
}

func (s *TeamService) BatchGetUsers(req *model.BatchGetUsersReq) (*model.BatchGetUsersRes, error) {
	list, err := repo.BatchGetUsers(req.UserIDs)
	if err != nil {
		log.Error("batch get users error", zap.Error(err))
		return nil, errorx.InternalServerError("查询失败")
	}

	return &model.BatchGetUsersRes{
		List: list,
	}, nil
}

func (s *TeamService) AddMember(userId string, req *model.AddMemberReq) error {
	// 创建人和管理员都能添加成员
	if !hasPermission(userId, req.TeamID, constants.TeamIdentityCreator, constants.TeamIdentityManager) {
		return errorx.New(errorx.ErrNoPermission, "无权限")
	}

	// 判断团队成员是否达到上限
	team, err := db.Get[model.Team](db.GetDB(), map[string]any{"id": req.TeamID})
	if err != nil {
		log.Error("get team error", zap.Error(err))
		return errorx.InternalServerError("添加失败")
	}
	if team == nil {
		return errorx.New(errorx.ErrNotFound, "团队不存在")
	}

	// 判断用户是否已添加到团队
	userTeam, err := db.Get[model2.UserTeam](db.GetDB(), map[string]any{"user_id": req.UserID, "team_id": req.TeamID})
	if err != nil {
		log.Error("get user team error", zap.Error(err))
		return errorx.InternalServerError("添加失败")
	}
	if userTeam != nil {
		return errorx.New(errorx.ErrOperationFailed, "用户已添加到团队")
	}

	creator, err := db.Get[model3.User](db.GetDB(), map[string]any{"id": team.UserID})
	if err != nil {
		log.Error("get creator error", zap.Error(err))
		return errorx.InternalServerError("添加失败")
	}
	if creator == nil {
		return errorx.New(errorx.ErrNotFound, "创建人不存在")
	}

	userTeams, _, err := db.List[model2.UserTeam](db.GetDB(), 0, 0, map[string]any{"team_id": req.TeamID}, nil)
	if err != nil {
		log.Error("query user team error", zap.Error(err))
		return errorx.InternalServerError("添加失败")
	}

	if len(userTeams) >= creator.Identity.LimitTeamMemberCount() {
		return errorx.New(errorx.ErrOperationFailed, "团队成员已达上限")
	}

	if err := db.Create(db.GetDB(), &model2.UserTeam{
		UserID:   req.UserID,
		TeamID:   req.TeamID,
		Identity: req.Identity,
	}); err != nil {
		log.Error("create user team error", zap.Error(err))
		return errorx.InternalServerError("添加失败")
	}

	return nil
}

func (s *TeamService) UpdateMember(userId string, req *model.UpdateMemberReq) error {
	// 创建人和管理员都能更新成员
	if !hasPermission(userId, req.TeamID, constants.TeamIdentityCreator, constants.TeamIdentityManager) {
		return errorx.New(errorx.ErrNoPermission, "无权限")
	}

	team, err := db.Get[model.Team](db.GetDB(), map[string]any{"id": req.TeamID})
	if err != nil {
		log.Error("get team error", zap.Error(err))
		return errorx.InternalServerError("更新失败")
	}
	if team == nil {
		return errorx.New(errorx.ErrNotFound, "团队不存在")
	}

	// 不允许修改创建者权限
	if team.UserID == req.UserID {
		return errorx.New(errorx.ErrOperationFailed, "无法修改创建者权限")
	}

	// 不允许修改自身权限
	if userId == req.UserID {
		return errorx.New(errorx.ErrOperationFailed, "无法修改自身权限")
	}

	userTeam, err := db.Get[model2.UserTeam](db.GetDB(), map[string]any{"user_id": req.UserID, "team_id": req.TeamID})
	if err != nil {
		log.Error("get user team error", zap.Error(err))
		return errorx.InternalServerError("更新失败")
	}
	if userTeam == nil {
		return errorx.New(errorx.ErrNotFound, "用户不属于该组织")
	}

	if !slices.Contains([]constants.TeamIdentity{constants.TeamIdentityManager, constants.TeamIdentityMember}, req.Identity) {
		return errorx.New(errorx.ErrOperationFailed, "设置的权限不合法")
	}

	if err := db.Update(db.GetDB(), &model2.UserTeam{
		Identity: req.Identity,
	}, map[string]any{
		"id": userTeam.ID,
	}, "identity"); err != nil {
		log.Error("update user team error", zap.Error(err))
		return errorx.InternalServerError("更新失败")
	}

	return nil
}

func (s *TeamService) DeleteMember(userId string, req *model.DeleteMemberReq) error {
	// 创建人和管理员都能删除成员
	if !hasPermission(userId, req.TeamID, constants.TeamIdentityCreator, constants.TeamIdentityManager) {
		return errorx.New(errorx.ErrNoPermission, "无权限")
	}

	team, err := db.Get[model.Team](db.GetDB(), map[string]any{"id": req.TeamID})
	if err != nil {
		log.Error("get team error", zap.Error(err))
		return errorx.InternalServerError("删除失败")
	}
	if team == nil {
		return errorx.New(errorx.ErrNotFound, "团队不存在")
	}

	// 不允许删除创建者
	if team.UserID == req.UserID {
		return errorx.New(errorx.ErrOperationFailed, "无法删除创建者")
	}

	// 不允许删除自身
	if userId == req.UserID {
		return errorx.New(errorx.ErrOperationFailed, "无法删除自身")
	}

	userTeam, err := db.Get[model2.UserTeam](db.GetDB(), map[string]any{"user_id": req.UserID, "team_id": req.TeamID})
	if err != nil {
		log.Error("get user team error", zap.Error(err))
		return errorx.InternalServerError("删除失败")
	}
	if userTeam == nil {
		return errorx.New(errorx.ErrNotFound, "用户已被移除")
	}

	if err := db.Delete[model2.UserTeam](db.GetDB(), map[string]any{
		"id": userTeam.ID,
	}); err != nil {
		log.Error("delete user team error", zap.Error(err))
		return errorx.InternalServerError("删除失败")
	}

	return nil
}
