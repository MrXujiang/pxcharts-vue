package repo

import (
	"mvtable/internal/app/team/model"
	userModel "mvtable/internal/app/user/model"
	"mvtable/internal/pkg/utils"
	"mvtable/internal/storage/db"
	"mvtable/pkg/log"

	"go.uber.org/zap"
)

func GetMemberList(teamId string, searchWord string, page, size int) (list []*model.TeamMember, total int64, err error) {
	query := db.GetDB().Table("user_team ut").
		Select("u.id AS user_id, u.nickname, u.avatar, u.email, ut.identity, ut.updated_at, ut.created_at").
		Joins(`LEFT JOIN "user" u ON u.id = ut.user_id`).
		Where("ut.team_id = ?", teamId).
		Where("ut.deleted_at IS NULL")

	if utils.ValidateEmail(searchWord) {
		query = query.Where("u.email = ?", searchWord)
	} else if searchWord != "" {
		query = query.Where("u.nickname LIKE ?", "%"+searchWord+"%")
	}

	err = query.Count(&total).
		Order("ut.created_at").
		Limit(size).
		Offset((page - 1) * size).
		Scan(&list).Error

	if list == nil {
		list = []*model.TeamMember{}
	}

	return list, total, err
}

func GetTeamList(page, size int, searchWord string) (list []*model.AdminGetTeamListItem, total int64, err error) {
	query := db.GetDB().Table("team t").
		Select(`
			t.id, 
			u.nickname AS creator_nickname, 
			u.email AS creator_email, 
			t.name, 
			t.description, 
			t.logo, 
			COALESCE(member_count.count, 0) AS member_count, 
			t.created_at, 
			t.updated_at
		`).
		Joins(`LEFT JOIN "user" u ON t.user_id = u.id`).
		Joins(`LEFT JOIN (
				SELECT team_id, COUNT(*) AS count
				FROM user_team
				WHERE deleted_at IS NULL
				GROUP BY team_id
			) AS member_count ON t.id = member_count.team_id`).
		Where("t.deleted_at IS NULL")

	if searchWord != "" {
		query = query.Where("t.name LIKE ?", "%"+searchWord+"%")
	}

	err = query.Count(&total).
		Order("t.created_at DESC").
		Limit(size).
		Offset((page - 1) * size).
		Scan(&list).Error

	return list, total, err
}

func SearchUser(userId, searchWord string) (list []*model.UserBasicInfo, err error) {
	if utils.ValidateEmail(searchWord) {
		// 如果是邮箱格式，则在所有用户中进行完全匹配
		err = db.GetDB().Table("\"user\"").
			Select("id, email, nickname, avatar").
			Where("email = ? AND deleted_at IS NULL", searchWord).
			Scan(&list).Error
	} else {
		// 如果不是邮箱，则在当前用户加入的团队成员范围内，根据昵称进行模糊匹配
		query := db.GetDB().Table(`"user" u`).
			Select("DISTINCT u.id, u.email, u.nickname, u.avatar").
			Joins("INNER JOIN user_team ut ON u.id = ut.user_id").
			Where("ut.deleted_at IS NULL AND u.deleted_at IS NULL").
			Where("ut.team_id IN (SELECT team_id FROM user_team WHERE user_id = ? AND deleted_at IS NULL)", userId)

		if searchWord != "" {
			query = query.Where("u.nickname LIKE ?", "%"+searchWord+"%")
		}

		err = query.Scan(&list).Error
	}

	if err != nil {
		log.Error("search user error", zap.Error(err))
		return nil, err
	}

	if list == nil {
		list = []*model.UserBasicInfo{}
	}

	return list, nil
}

func BatchGetUsers(userIDs []string) (list []*model.UserBasicInfo, err error) {
	if len(userIDs) == 0 {
		return []*model.UserBasicInfo{}, nil
	}

	// 查询用户信息
	users, _, err := db.List[userModel.User](db.GetDB(), 0, 0, map[string]any{"id": []any{"IN", userIDs}}, nil)
	if err != nil {
		log.Error("batch get users error", zap.Error(err))
		return nil, err
	}

	// 转换为 UserBasicInfo 格式
	list = make([]*model.UserBasicInfo, 0, len(users))
	for _, user := range users {
		list = append(list, &model.UserBasicInfo{
			ID:       user.ID,
			Email:    user.Email,
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
		})
	}

	return list, nil
}
