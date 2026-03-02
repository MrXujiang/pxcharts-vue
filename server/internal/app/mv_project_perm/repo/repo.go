package repo

import (
	"mvtable/internal/app/mv_project_perm/model"

	"gorm.io/gorm"
)

func QueryAvailableMembers(db *gorm.DB, userId, keyword, projectId string) ([]*model.MemberListItem, error) {
	var result []*model.MemberListItem
	sql := `
		WITH current_user_teams AS (
			SELECT team_id
			FROM user_team
			WHERE user_id = ? AND deleted_at IS NULL
		),
		team_members AS (
			SELECT DISTINCT ut.user_id
			FROM user_team ut
			JOIN current_user_teams cut ON ut.team_id = cut.team_id
			WHERE ut.deleted_at IS NULL
		),
		excluded_targets AS (
			SELECT target_id, target
			FROM mv_project_perm
			WHERE project_id = ? AND deleted_at IS NULL
		)
		-- 查询用户
		SELECT
			u.id AS target_id,
			'user' AS target,
			u.nickname AS name,
			u.avatar AS avatar
		FROM "user" u
		WHERE u.deleted_at IS NULL
		  AND u.status = 1
		  AND (
			  -- 同团队的：模糊匹配
			  (u.id IN (SELECT user_id FROM team_members) AND u.nickname ILIKE '%' || ? || '%')
			  OR
			  -- 不同团队的：完全匹配
			  (u.id NOT IN (SELECT user_id FROM team_members) AND u.nickname = ?)
		  )
		  AND NOT EXISTS (
			  SELECT 1 FROM excluded_targets et
			  WHERE et.target_id = u.id AND et.target = 'user'
		  )
		
		UNION ALL
		
		-- 查询团队
		SELECT
			t.id AS target_id,
			'team' AS target,
			t.name AS name,
			t.logo AS avatar
		FROM team t
		JOIN user_team ut ON t.id = ut.team_id
		WHERE ut.user_id = ?
		  AND t.deleted_at IS NULL
		  AND ut.deleted_at IS NULL
		  AND t.name ILIKE '%' || ? || '%'
		  AND NOT EXISTS (
			  SELECT 1 FROM excluded_targets et
			  WHERE et.target_id = t.id AND et.target = 'team'
		  )
	`
	err := db.Raw(sql, userId, projectId, keyword, keyword, userId, keyword).Scan(&result).Error
	return result, err
}

// QueryProjectMembers 查询项目已添加的成员（用户或团队）
func QueryProjectMembers(db *gorm.DB, projectId string) ([]*model.ProjectMemberItem, error) {
	var result []*model.ProjectMemberItem
	sql := `
		SELECT
			p.target_id,
			p.target,
			p.role,
			u.nickname AS name,
			u.avatar   AS avatar
		FROM mv_project_perm p
		JOIN "user" u ON p.target_id = u.id
		WHERE p.project_id = ?
		  AND p.deleted_at IS NULL
		  AND p.target = 'user'
		  AND u.deleted_at IS NULL
		  AND u.status = 1

		UNION ALL

		SELECT
			p.target_id,
			p.target,
			p.role,
			t.name  AS name,
			t.logo  AS avatar
		FROM mv_project_perm p
		JOIN team t ON p.target_id = t.id
		WHERE p.project_id = ?
		  AND p.deleted_at IS NULL
		  AND p.target = 'team'
		  AND t.deleted_at IS NULL

		ORDER BY name ASC
	`

	if err := db.Raw(sql, projectId, projectId).Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
