package repo

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

func GetChildrenIdsByParentId(db *gorm.DB, parentId string) ([]string, error) {
	query := `
		WITH RECURSIVE folder_tree AS (
			SELECT id
			FROM mv_folder
			WHERE parent_id = ?
			  AND deleted_at IS NULL

			UNION ALL

			SELECT f.id
			FROM mv_folder f
			         INNER JOIN folder_tree ft ON f.parent_id = ft.id
			WHERE f.deleted_at IS NULL
		)
		SELECT id
		FROM folder_tree;
	`

	rows, err := db.Raw(query, parentId).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

// DeleteDataByFolderIds 根据文件夹ID及其子文件夹，软删除所有相关的数据
// 包括：记录、字段、表结构、视图、表单、仪表盘等
func DeleteDataByFolderIds(db *gorm.DB, folderIds []string) error {
	if len(folderIds) == 0 {
		return nil
	}

	// 先查询出所有需要删除的文件夹ID（包括子文件夹）
	folderTreeSQL := `
		WITH RECURSIVE folder_tree AS (
				SELECT id
				FROM mv_folder
				WHERE id = ANY(?)
					AND deleted_at IS NULL
				UNION ALL
				SELECT f.id
				FROM mv_folder f
								INNER JOIN folder_tree ft ON f.parent_id = ft.id
				WHERE f.deleted_at IS NULL
		)
		SELECT id FROM folder_tree;
	`

	rows, err := db.Raw(folderTreeSQL, pq.Array(folderIds)).Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	var allFolderIds []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return err
		}
		allFolderIds = append(allFolderIds, id)
	}

	if len(allFolderIds) == 0 {
		return nil
	}

	// 查询所有相关的 table_schema_id
	tableSchemaSQL := `
		SELECT id
		FROM mv_table_schema
		WHERE folder_id = ANY(?)
			AND deleted_at IS NULL;
	`
	tableSchemaRows, err := db.Raw(tableSchemaSQL, pq.Array(allFolderIds)).Rows()
	if err != nil {
		return err
	}
	defer tableSchemaRows.Close()

	var tableSchemaIds []string
	for tableSchemaRows.Next() {
		var id string
		if err := tableSchemaRows.Scan(&id); err != nil {
			return err
		}
		tableSchemaIds = append(tableSchemaIds, id)
	}

	// 按照依赖关系逆序删除（从叶子节点开始删除）

	// 1. 删除表单提交记录
	if len(tableSchemaIds) > 0 {
		updateFormSubmitSQL := `
			UPDATE mv_form_submit
			SET deleted_at = NOW()
			WHERE view_form_id IN (
				SELECT id
				FROM mv_view_form
				WHERE view_id IN (
					SELECT id
					FROM mv_view
					WHERE table_schema_id = ANY(?)
						AND deleted_at IS NULL
				)
				AND deleted_at IS NULL
			)
			AND deleted_at IS NULL;
		`
		if err := db.Exec(updateFormSubmitSQL, pq.Array(tableSchemaIds)).Error; err != nil {
			return err
		}

		// 2. 删除富文本内容
		updateRichTextSQL := `
			UPDATE mv_rich_text_content
			SET deleted_at = NOW()
			WHERE record_id IN (
				SELECT id
				FROM mv_record
				WHERE table_schema_id = ANY(?)
					AND deleted_at IS NULL
			)
			AND deleted_at IS NULL;
		`
		if err := db.Exec(updateRichTextSQL, pq.Array(tableSchemaIds)).Error; err != nil {
			return err
		}
	}

	// 3. 删除视图活跃状态
	if len(tableSchemaIds) > 0 {
		updateViewActiveSQL := `
			UPDATE mv_view_active
			SET deleted_at = NOW()
			WHERE table_schema_id = ANY(?)
				AND deleted_at IS NULL;
		`
		if err := db.Exec(updateViewActiveSQL, pq.Array(tableSchemaIds)).Error; err != nil {
			return err
		}
	}

	// 4. 删除项目高级权限
	if len(tableSchemaIds) > 0 {
		updateAdvancedPermSQL := `
			UPDATE mv_project_advanced_perm
			SET deleted_at = NOW()
			WHERE table_schema_id = ANY(?)
				AND deleted_at IS NULL;
		`
		if err := db.Exec(updateAdvancedPermSQL, pq.Array(tableSchemaIds)).Error; err != nil {
			return err
		}
	}

	// 5. 删除记录
	updateRecordSQL := `
		UPDATE mv_record
		SET deleted_at = NOW()
		WHERE table_schema_id = ANY(?)
			AND deleted_at IS NULL;
	`
	if err := db.Exec(updateRecordSQL, pq.Array(tableSchemaIds)).Error; err != nil {
		return err
	}

	// 6. 删除字段
	updateFieldSQL := `
		UPDATE mv_field
		SET deleted_at = NOW()
		WHERE table_schema_id = ANY(?)
			AND deleted_at IS NULL;
	`
	if err := db.Exec(updateFieldSQL, pq.Array(tableSchemaIds)).Error; err != nil {
		return err
	}

	// 7. 删除视图配置和视图（在 DeleteViewAndDashboardByFolderIds 中处理）
	// 这里不需要重复处理，因为另一个函数会处理

	// 8. 删除表结构
	updateTableSchemaSQL := `
		UPDATE mv_table_schema
		SET deleted_at = NOW()
		WHERE id = ANY(?)
			AND deleted_at IS NULL;
	`
	if err := db.Exec(updateTableSchemaSQL, pq.Array(tableSchemaIds)).Error; err != nil {
		return err
	}

	// 9. 删除仪表盘（在 DeleteViewAndDashboardByFolderIds 中处理）
	// 这里不需要重复处理，因为另一个函数会处理

	// 10. 删除文件夹
	updateFolderSQL := `
		UPDATE mv_folder
		SET deleted_at = NOW()
		WHERE id = ANY(?)
			AND deleted_at IS NULL;
	`
	if err := db.Exec(updateFolderSQL, pq.Array(allFolderIds)).Error; err != nil {
		return err
	}

	return nil
}

// DeleteViewAndDashboardByFolderIds 根据文件夹及其子文件夹，删除（软删）对应的视图、视图配置以及仪表盘和仪表盘图表
func DeleteViewAndDashboardByFolderIds(db *gorm.DB, folderIds []string) error {
	if len(folderIds) == 0 {
		return nil
	}

	// 先查询出所有需要删除的文件夹ID（包括子文件夹）
	folderTreeSQL := `
		WITH RECURSIVE folder_tree AS (
				SELECT id
				FROM mv_folder
				WHERE id = ANY(?)
					AND deleted_at IS NULL
				UNION ALL
				SELECT f.id
				FROM mv_folder f
								INNER JOIN folder_tree ft ON f.parent_id = ft.id
				WHERE f.deleted_at IS NULL
		)
		SELECT id FROM folder_tree;
	`

	rows, err := db.Raw(folderTreeSQL, pq.Array(folderIds)).Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	var allFolderIds []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return err
		}
		allFolderIds = append(allFolderIds, id)
	}

	if len(allFolderIds) == 0 {
		return nil
	}

	// 查询相关的 table_schema_id
	tableSchemaSQL := `
		SELECT id
		FROM mv_table_schema
		WHERE folder_id = ANY(?)
			AND deleted_at IS NULL;
	`
	tableSchemaRows, err := db.Raw(tableSchemaSQL, pq.Array(allFolderIds)).Rows()
	if err != nil {
		return err
	}
	defer tableSchemaRows.Close()

	var tableSchemaIds []string
	for tableSchemaRows.Next() {
		var id string
		if err := tableSchemaRows.Scan(&id); err != nil {
			return err
		}
		tableSchemaIds = append(tableSchemaIds, id)
	}

	// 查询相关的 view_id
	var viewIds []string
	if len(tableSchemaIds) > 0 {
		viewSQL := `
			SELECT id
			FROM mv_view
			WHERE table_schema_id = ANY(?)
				AND deleted_at IS NULL;
		`
		viewRows, err := db.Raw(viewSQL, pq.Array(tableSchemaIds)).Rows()
		if err != nil {
			return err
		}
		defer viewRows.Close()

		for viewRows.Next() {
			var id string
			if err := viewRows.Scan(&id); err != nil {
				return err
			}
			viewIds = append(viewIds, id)
		}
	}

	// 查询相关的 dashboard_id
	dashboardSQL := `
		SELECT id
		FROM mv_dashboard
		WHERE folder_id = ANY(?)
			AND deleted_at IS NULL;
	`
	dashboardRows, err := db.Raw(dashboardSQL, pq.Array(allFolderIds)).Rows()
	if err != nil {
		return err
	}
	defer dashboardRows.Close()

	var dashboardIds []string
	for dashboardRows.Next() {
		var id string
		if err := dashboardRows.Scan(&id); err != nil {
			return err
		}
		dashboardIds = append(dashboardIds, id)
	}

	// 分别执行每个UPDATE语句
	// 1. 删除仪表盘下的图表
	if len(dashboardIds) > 0 {
		updateDashboardChartSQL := `
			UPDATE mv_dashboard_chart
			SET deleted_at = NOW()
			WHERE dashboard_id = ANY(?)
				AND deleted_at IS NULL;
		`
		if err := db.Exec(updateDashboardChartSQL, pq.Array(dashboardIds)).Error; err != nil {
			return err
		}
	}

	// 2. 删除仪表盘
	if len(dashboardIds) > 0 {
		updateDashboardSQL := `
			UPDATE mv_dashboard
			SET deleted_at = NOW()
			WHERE id = ANY(?)
				AND deleted_at IS NULL;
		`
		if err := db.Exec(updateDashboardSQL, pq.Array(dashboardIds)).Error; err != nil {
			return err
		}
	}

	// 3. 删除视图表格配置
	if len(viewIds) > 0 {
		updateViewTableSQL := `
			UPDATE mv_view_table
			SET deleted_at = NOW()
			WHERE view_id = ANY(?)
				AND deleted_at IS NULL;
		`
		if err := db.Exec(updateViewTableSQL, pq.Array(viewIds)).Error; err != nil {
			return err
		}
	}

	// 4. 删除视图表单配置
	if len(viewIds) > 0 {
		updateViewFormSQL := `
			UPDATE mv_view_form
			SET deleted_at = NOW()
			WHERE view_id = ANY(?)
				AND deleted_at IS NULL;
		`
		if err := db.Exec(updateViewFormSQL, pq.Array(viewIds)).Error; err != nil {
			return err
		}
	}

	// 5. 删除视图
	if len(viewIds) > 0 {
		updateViewSQL := `
			UPDATE mv_view
			SET deleted_at = NOW()
			WHERE id = ANY(?)
				AND deleted_at IS NULL;
		`
		if err := db.Exec(updateViewSQL, pq.Array(viewIds)).Error; err != nil {
			return err
		}
	}

	return nil
}
