package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	projectModel "mvtable/internal/app/mv_project/model"
	advancedPermModel "mvtable/internal/app/mv_project_advanced_perm/model"
	projectPermModel "mvtable/internal/app/mv_project_perm/model"
	tableSchemaModel "mvtable/internal/app/mv_table_schema/model"
	userTeamModel "mvtable/internal/app/user_team/model"
	"mvtable/internal/pkg/constants"
	"mvtable/internal/pkg/errorx"
	"mvtable/internal/storage/db"
	"mvtable/pkg/log"
	"slices"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// AdvancedPermConfig 高级权限中间件配置
type AdvancedPermConfig struct {
	// 需要校验的操作类型：add, delete, edit, read, view_operate, view_check, field_manage
	Action string
	// 是否需要校验字段权限（用于编辑操作）
	CheckFieldPerm bool
	// 字段ID（用于字段权限校验）
	FieldID string
	// 是否需要校验视图权限
	CheckViewPerm bool
	// 视图ID（用于视图权限校验）
	ViewID string
	// 是否需要校验操作范围（own/all）
	CheckOperateRange bool
	// 记录ID（用于操作范围校验）
	RecordID string
}

// AdvancedPermMiddleware 高级权限校验中间件
// 使用方式：router.Use(middleware.AdvancedPermMiddleware(middleware.AdvancedPermConfig{Action: "add"}))
func AdvancedPermMiddleware(config AdvancedPermConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := GetCurrentUserID(c)

		// 从请求中提取必要参数
		var projectID, tableSchemaID string

		// 尝试从不同位置获取 projectID
		if projectID = c.Query("projectId"); projectID == "" {
			if projectID = c.PostForm("projectId"); projectID == "" {
				// 尝试从 JSON body 中获取（不消耗请求体）
				if bodyData := getRequestBody(c); bodyData != nil {
					var req struct {
						ProjectID string `json:"projectId"`
					}
					if err := json.Unmarshal(bodyData, &req); err == nil {
						projectID = req.ProjectID
					}
				}
			}
		}

		// 尝试从不同位置获取 tableSchemaID
		if tableSchemaID = c.Query("tableSchemaId"); tableSchemaID == "" {
			if tableSchemaID = c.PostForm("tableSchemaId"); tableSchemaID == "" {
				// 尝试从 JSON body 中获取（不消耗请求体）
				if bodyData := getRequestBody(c); bodyData != nil {
					var req struct {
						TableSchemaID string `json:"tableSchemaId"`
					}
					if err := json.Unmarshal(bodyData, &req); err == nil {
						tableSchemaID = req.TableSchemaID
					}
				}
			}
		}

		// 如果无法获取 projectID 或 tableSchemaID，尝试从 tableSchemaID 获取 projectID
		if projectID == "" && tableSchemaID != "" {
			tableSchema, err := db.Get[tableSchemaModel.MvTableSchema](db.GetDB(), map[string]any{"id": tableSchemaID})
			if err == nil && tableSchema != nil {
				projectID = tableSchema.ProjectID
			}
		}

		// 如果仍然无法获取 projectID，跳过高级权限校验
		if projectID == "" {
			c.Next()
			return
		}

		// 检查项目是否开启高级权限
		project, err := db.Get[projectModel.MvProject](db.GetDB(), map[string]any{"id": projectID})
		if err != nil {
			log.Error("get project error", zap.Error(err))
			c.Next()
			return
		}
		if project == nil || !project.EnableAdvancedPerm {
			// 项目未开启高级权限，跳过校验
			c.Next()
			return
		}

		// 如果 tableSchemaID 为空，尝试从字段ID获取（用于字段更新和删除操作）
		if tableSchemaID == "" && (config.Action == "field_manage" || config.Action == "field_update" || config.Action == "field_delete") {
			// 尝试从请求中获取字段ID
			var fieldID string
			if fieldID = c.Query("id"); fieldID == "" {
				if fieldID = c.PostForm("id"); fieldID == "" {
					// 尝试从 JSON body 中获取（不消耗请求体）
					if bodyData := getRequestBody(c); bodyData != nil {
						var req struct {
							ID string `json:"id"`
						}
						if err := json.Unmarshal(bodyData, &req); err == nil {
							fieldID = req.ID
						}
					}
				}
			}
			if fieldID != "" {
				// 从字段ID获取 tableSchemaID（直接查询数据库）
				var fieldTableSchemaID string
				if err := db.GetDB().Table("mv_field").Where("id = ? AND deleted_at IS NULL", fieldID).Select("table_schema_id").Scan(&fieldTableSchemaID).Error; err == nil && fieldTableSchemaID != "" {
					tableSchemaID = fieldTableSchemaID
					// 如果还没有 projectID，从 tableSchemaID 获取
					if projectID == "" && tableSchemaID != "" {
						tableSchema, err := db.Get[tableSchemaModel.MvTableSchema](db.GetDB(), map[string]any{"id": tableSchemaID})
						if err == nil && tableSchema != nil {
							projectID = tableSchema.ProjectID
						}
					}
				}
			}
		}

		// 如果 tableSchemaID 为空，无法进行高级权限校验
		if tableSchemaID == "" {
			c.Next()
			return
		}

		// 获取用户在项目中的角色
		role, err := getUserProjectRole(db.GetDB(), userID, projectID)
		if err != nil {
			log.Error("get user project role error", zap.Error(err))
			ResErr(c, errorx.InternalServerError("获取权限失败"))
			return
		}
		if role == "" {
			ResErr(c, errorx.New(errorx.ErrNoPermission, "无权限访问该项目"))
			return
		}

		// 获取高级权限配置
		advancedPerm, err := db.Get[advancedPermModel.MvProjectAdvancedPerm](db.GetDB(), map[string]any{
			"table_schema_id": tableSchemaID,
			"role":            role,
		})
		if err != nil {
			log.Error("get advanced perm error", zap.Error(err))
			ResErr(c, errorx.InternalServerError("获取权限配置失败"))
			return
		}
		if advancedPerm == nil {
			ResErr(c, errorx.New(errorx.ErrNoPermission, "无权限配置"))
			return
		}

		// 根据操作类型校验权限
		switch config.Action {
		case "add":
			// 检查 data_action 是否允许新增（新增需要 edit 或 manage 权限）
			if advancedPerm.DataAction != constants.ActionEdit && advancedPerm.DataAction != constants.ActionMange {
				ResErr(c, errorx.New(errorx.ErrNoPermission, "无权限新增数据"))
				return
			}
			// 检查 can_add 标志
			if !advancedPerm.CanAdd {
				ResErr(c, errorx.New(errorx.ErrNoPermission, "无权限新增数据"))
				return
			}
			// 校验字段权限（新增记录时可能需要检查字段权限）
			if err := checkFieldPermission(c, advancedPerm, config, "add"); err != nil {
				ResErr(c, err)
				return
			}
		case "delete":
			// 检查 data_action 是否允许删除（删除需要 edit 或 manage 权限）
			if advancedPerm.DataAction != constants.ActionEdit && advancedPerm.DataAction != constants.ActionMange {
				ResErr(c, errorx.New(errorx.ErrNoPermission, "无权限删除数据"))
				return
			}
			// 检查 can_delete 标志
			if !advancedPerm.CanDelete {
				ResErr(c, errorx.New(errorx.ErrNoPermission, "无权限删除数据"))
				return
			}
			// 校验操作范围
			if config.CheckOperateRange && advancedPerm.OperateRange == constants.AdvancedOperateRangeOwn {
				recordID := config.RecordID
				if recordID == "" {
					// 尝试从请求中获取 recordID（不消耗请求体）
					if bodyData := getRequestBody(c); bodyData != nil {
						var req struct {
							RecordIDs []string `json:"recordIds"`
							RecordID  string   `json:"recordId"`
						}
						if err := json.Unmarshal(bodyData, &req); err == nil {
							if len(req.RecordIDs) > 0 {
								// 对于批量删除，需要检查所有记录是否都属于当前用户
								for _, rid := range req.RecordIDs {
									if !isRecordOwner(db.GetDB(), rid, userID) {
										ResErr(c, errorx.New(errorx.ErrNoPermission, "只能删除自己创建的记录"))
										return
									}
								}
								return // 所有记录都通过校验
							} else if req.RecordID != "" {
								recordID = req.RecordID
							}
						}
					}
				}
				if recordID != "" {
					if !isRecordOwner(db.GetDB(), recordID, userID) {
						ResErr(c, errorx.New(errorx.ErrNoPermission, "只能删除自己创建的记录"))
						return
					}
				}
			}
		case "edit":
			// 检查 data_action 是否允许编辑
			if advancedPerm.DataAction != constants.ActionEdit && advancedPerm.DataAction != constants.ActionMange {
				ResErr(c, errorx.New(errorx.ErrNoPermission, "无权限编辑数据"))
				return
			}
			// 校验操作范围
			if config.CheckOperateRange && advancedPerm.OperateRange == constants.AdvancedOperateRangeOwn {
				recordID := config.RecordID
				if recordID == "" {
					// 尝试从请求中获取 recordID（不消耗请求体）
					if bodyData := getRequestBody(c); bodyData != nil {
						var req struct {
							RecordID        string `json:"recordId"`
							CurrentRecordID string `json:"currentId"`
						}
						if err := json.Unmarshal(bodyData, &req); err == nil {
							if req.RecordID != "" {
								recordID = req.RecordID
							} else if req.CurrentRecordID != "" {
								recordID = req.CurrentRecordID
							}
						}
					}
				}
				if recordID != "" {
					if !isRecordOwner(db.GetDB(), recordID, userID) {
						ResErr(c, errorx.New(errorx.ErrNoPermission, "只能编辑自己创建的记录"))
						return
					}
				}
			}
			// 校验字段权限
			if err := checkFieldPermission(c, advancedPerm, config, "edit"); err != nil {
				ResErr(c, err)
				return
			}
		case "read":
			// 检查 data_action 是否允许查看
			if advancedPerm.DataAction == constants.ActionNone {
				ResErr(c, errorx.New(errorx.ErrNoPermission, "无权限查看数据"))
				return
			}
			// 校验字段权限
			if err := checkFieldPermission(c, advancedPerm, config, "read"); err != nil {
				ResErr(c, err)
				return
			}
		case "view_operate":
			// 检查是否可以操作视图
			if !advancedPerm.CanOperateView {
				ResErr(c, errorx.New(errorx.ErrNoPermission, "无权限操作视图"))
				return
			}
		case "view_check":
			// 检查视图访问权限
			if config.CheckViewPerm {
				viewID := config.ViewID
				if viewID == "" {
					// 尝试从请求中获取 viewID（不消耗请求体）
					if viewID = c.Query("viewId"); viewID == "" {
						if viewID = c.Query("id"); viewID == "" {
							if bodyData := getRequestBody(c); bodyData != nil {
								var req struct {
									ViewID string `json:"viewId"`
									ID     string `json:"id"`
								}
								if err := json.Unmarshal(bodyData, &req); err == nil {
									if req.ViewID != "" {
										viewID = req.ViewID
									} else if req.ID != "" {
										viewID = req.ID
									}
								}
							}
						}
					}
				}
				if viewID != "" {
					if advancedPerm.ViewAccess == constants.AdvancedViewAccessCustom {
						if !slices.Contains(advancedPerm.CanCheckViews, viewID) {
							ResErr(c, errorx.New(errorx.ErrNoPermission, "无权限查看该视图"))
							return
						}
					}
				}
			}
		case "field_manage":
			// 字段管理操作（创建、更新、删除字段）需要 manage 权限
			if advancedPerm.DataAction != constants.ActionMange {
				ResErr(c, errorx.New(errorx.ErrNoPermission, "无权限管理字段"))
				return
			}
		}

		// 将权限信息存储到上下文中，供后续使用
		c.Set("projectRole", role)
		c.Set("advancedPerm", advancedPerm)

		c.Next()
	}
}

// getRequestBody 安全地读取请求体，不消耗请求体内容
// 返回请求体的字节数组，如果读取失败返回 nil
func getRequestBody(c *gin.Context) []byte {
	// 检查是否已经读取过请求体
	if c.Request.Body == nil {
		return nil
	}

	// 读取请求体
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil
	}

	// 重新设置请求体，以便后续 handler 可以读取
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return bodyBytes
}

// getUserProjectRole 获取用户在项目中的角色
// 优先级：owner > admin > editor > reader
// 考虑用户直接权限和团队权限
func getUserProjectRole(gormDB *gorm.DB, userID, projectID string) (constants.ProjectRole, error) {
	// 首先检查用户是否是项目创建者
	project, err := db.Get[projectModel.MvProject](gormDB, map[string]any{"id": projectID})
	if err != nil {
		return "", err
	}
	if project != nil && project.UserID == userID {
		return constants.ProjectActionOwner, nil
	}

	// 查询用户直接权限
	userPerm, err := db.Get[projectPermModel.MvProjectPerm](gormDB, map[string]any{
		"project_id": projectID,
		"target_id":  userID,
		"target":     "user",
	})
	if err != nil {
		return "", err
	}

	var highestRole constants.ProjectRole
	if userPerm != nil {
		highestRole = userPerm.Role
	}

	// 查询用户所属的团队
	var userTeams []userTeamModel.UserTeam
	if err := gormDB.Where("user_id = ? AND deleted_at IS NULL", userID).Find(&userTeams).Error; err != nil {
		return "", err
	}

	// 查询团队权限
	for _, userTeam := range userTeams {
		teamPerm, err := db.Get[projectPermModel.MvProjectPerm](gormDB, map[string]any{
			"project_id": projectID,
			"target_id":  userTeam.TeamID,
			"target":     "team",
		})
		if err != nil {
			continue
		}
		if teamPerm != nil {
			// 比较权限等级，取最高权限
			if compareRole(teamPerm.Role, highestRole) > 0 {
				highestRole = teamPerm.Role
			}
		}
	}

	return highestRole, nil
}

// compareRole 比较角色权限等级，返回 1 表示 role1 更高，-1 表示 role2 更高，0 表示相等
func compareRole(role1, role2 constants.ProjectRole) int {
	roleLevel := map[constants.ProjectRole]int{
		constants.ProjectActionOwner:  4,
		constants.ProjectActionAdmin:  3,
		constants.ProjectActionEditor: 2,
		constants.ProjectActionReader: 1,
	}

	level1 := roleLevel[role1]
	level2 := roleLevel[role2]

	if level1 > level2 {
		return 1
	} else if level1 < level2 {
		return -1
	}
	return 0
}

// isRecordOwner 检查记录是否属于指定用户
func isRecordOwner(gormDB *gorm.DB, recordID, userID string) bool {
	var createdBy string
	if err := gormDB.Table("mv_record").Where("id = ? AND deleted_at IS NULL", recordID).Select("created_by").Scan(&createdBy).Error; err != nil {
		return false
	}
	return createdBy == userID
}

// hasFieldPermission 检查是否有字段权限
// 参数：
//   - perm: 高级权限配置
//   - fieldID: 字段ID
//   - action: 操作类型（read/add/edit）
//
// 返回：
//   - true: 有权限
//   - false: 无权限
func hasFieldPermission(perm *advancedPermModel.MvProjectAdvancedPerm, fieldID, action string) bool {
	// 如果字段访问权限为 all，则所有字段都有权限
	if perm.FieldAccess != constants.AdvancedFieldAccessCustom {
		return true
	}

	// 如果字段访问权限为 custom，需要检查自定义字段权限配置
	for _, fieldPerm := range perm.CustomFieldPerm {
		if fieldPerm.FieldID == fieldID {
			switch action {
			case "read":
				return fieldPerm.CanRead
			case "add":
				return fieldPerm.CanAdd
			case "edit":
				return fieldPerm.CanEdit
			}
		}
	}

	// 如果字段不在自定义权限列表中，默认无权限
	return false
}

// checkFieldPermission 统一的字段权限校验函数
// 参数：
//   - c: gin上下文
//   - advancedPerm: 高级权限配置
//   - config: 中间件配置
//   - action: 操作类型（read/add/edit）
//
// 返回：
//   - error: 如果有权限错误，返回错误；否则返回 nil
func checkFieldPermission(c *gin.Context, advancedPerm *advancedPermModel.MvProjectAdvancedPerm, config AdvancedPermConfig, action string) error {
	// 如果不需要校验字段权限，直接返回
	if !config.CheckFieldPerm {
		return nil
	}

	// 获取字段ID
	fieldID := config.FieldID
	if fieldID == "" {
		// 尝试从请求中获取 fieldID（不消耗请求体）
		if bodyData := getRequestBody(c); bodyData != nil {
			var req struct {
				FieldID string `json:"fieldId"`
			}
			if err := json.Unmarshal(bodyData, &req); err == nil {
				fieldID = req.FieldID
			}
		}
	}

	// 如果没有字段ID，无法进行字段权限校验
	if fieldID == "" {
		return nil // 没有字段ID时不进行校验，由业务逻辑处理
	}

	// 检查字段权限
	if advancedPerm.FieldAccess == constants.AdvancedFieldAccessCustom {
		if !hasFieldPermission(advancedPerm, fieldID, action) {
			var errMsg string
			switch action {
			case "read":
				errMsg = "无权限查看该字段"
			case "add":
				errMsg = "无权限新增该字段数据"
			case "edit":
				errMsg = "无权限编辑该字段"
			default:
				errMsg = "无权限操作该字段"
			}
			return errorx.New(errorx.ErrNoPermission, errMsg)
		}
	}

	return nil
}

// GetProjectRole 从上下文中获取项目角色
func GetProjectRole(c *gin.Context) constants.ProjectRole {
	if role, exists := c.Get("projectRole"); exists {
		if r, ok := role.(constants.ProjectRole); ok {
			return r
		}
	}
	return ""
}

// GetAdvancedPerm 从上下文中获取高级权限配置
func GetAdvancedPerm(c *gin.Context) *advancedPermModel.MvProjectAdvancedPerm {
	if perm, exists := c.Get("advancedPerm"); exists {
		if p, ok := perm.(*advancedPermModel.MvProjectAdvancedPerm); ok {
			return p
		}
	}
	return nil
}

// CheckAdvancedPerm 便捷函数：检查高级权限
// 用于在 handler 中手动调用权限检查
func CheckAdvancedPerm(c *gin.Context, config AdvancedPermConfig) error {
	userID := GetCurrentUserID(c)

	// 从请求中提取必要参数
	var projectID, tableSchemaID string

	// 尝试从不同位置获取 projectID
	if projectID = c.Query("projectId"); projectID == "" {
		if projectID = c.PostForm("projectId"); projectID == "" {
			// 尝试从 JSON body 中获取（不消耗请求体）
			if bodyData := getRequestBody(c); bodyData != nil {
				var req struct {
					ProjectID string `json:"projectId"`
				}
				if err := json.Unmarshal(bodyData, &req); err == nil {
					projectID = req.ProjectID
				}
			}
		}
	}

	// 尝试从不同位置获取 tableSchemaID
	if tableSchemaID = c.Query("tableSchemaId"); tableSchemaID == "" {
		if tableSchemaID = c.PostForm("tableSchemaId"); tableSchemaID == "" {
			// 尝试从 JSON body 中获取（不消耗请求体）
			if bodyData := getRequestBody(c); bodyData != nil {
				var req struct {
					TableSchemaID string `json:"tableSchemaId"`
				}
				if err := json.Unmarshal(bodyData, &req); err == nil {
					tableSchemaID = req.TableSchemaID
				}
			}
		}
	}

	// 如果 tableSchemaID 为空，尝试从字段ID获取（用于字段更新和删除操作）
	if tableSchemaID == "" && (config.Action == "field_manage" || config.Action == "field_update" || config.Action == "field_delete") {
		// 尝试从请求中获取字段ID
		var fieldID string
		if fieldID = c.Query("id"); fieldID == "" {
			if fieldID = c.PostForm("id"); fieldID == "" {
				// 尝试从 JSON body 中获取（不消耗请求体）
				if bodyData := getRequestBody(c); bodyData != nil {
					var req struct {
						ID string `json:"id"`
					}
					if err := json.Unmarshal(bodyData, &req); err == nil {
						fieldID = req.ID
					}
				}
			}
		}
		if fieldID != "" {
			// 从字段ID获取 tableSchemaID（直接查询数据库）
			var fieldTableSchemaID string
			if err := db.GetDB().Table("mv_field").Where("id = ? AND deleted_at IS NULL", fieldID).Select("table_schema_id").Scan(&fieldTableSchemaID).Error; err == nil && fieldTableSchemaID != "" {
				tableSchemaID = fieldTableSchemaID
			}
		}
	}

	// 如果无法获取 projectID 或 tableSchemaID，尝试从 tableSchemaID 获取 projectID
	if projectID == "" && tableSchemaID != "" {
		tableSchema, err := db.Get[tableSchemaModel.MvTableSchema](db.GetDB(), map[string]any{"id": tableSchemaID})
		if err == nil && tableSchema != nil {
			projectID = tableSchema.ProjectID
		}
	}

	if projectID == "" || tableSchemaID == "" {
		return errorx.BadRequest("缺少必要参数")
	}

	// 检查项目是否开启高级权限
	project, err := db.Get[projectModel.MvProject](db.GetDB(), map[string]any{"id": projectID})
	if err != nil {
		log.Error("get project error", zap.Error(err))
		return errorx.InternalServerError("获取项目信息失败")
	}
	if project == nil || !project.EnableAdvancedPerm {
		// 项目未开启高级权限，直接通过
		return nil
	}

	// 获取用户在项目中的角色
	role, err := getUserProjectRole(db.GetDB(), userID, projectID)
	if err != nil {
		log.Error("get user project role error", zap.Error(err))
		return errorx.InternalServerError("获取权限失败")
	}
	if role == "" {
		return errorx.New(errorx.ErrNoPermission, "无权限访问该项目")
	}

	// 获取高级权限配置
	advancedPerm, err := db.Get[advancedPermModel.MvProjectAdvancedPerm](db.GetDB(), map[string]any{
		"table_schema_id": tableSchemaID,
		"role":            role,
	})
	if err != nil {
		log.Error("get advanced perm error", zap.Error(err))
		return errorx.InternalServerError("获取权限配置失败")
	}
	if advancedPerm == nil {
		return errorx.New(errorx.ErrNoPermission, "无权限配置")
	}

	// 根据操作类型校验权限
	switch config.Action {
	case "add":
		// 检查 data_action 是否允许新增（新增需要 edit 或 manage 权限）
		if advancedPerm.DataAction != constants.ActionEdit && advancedPerm.DataAction != constants.ActionMange {
			return errorx.New(errorx.ErrNoPermission, "无权限新增数据")
		}
		// 检查 can_add 标志
		if !advancedPerm.CanAdd {
			return errorx.New(errorx.ErrNoPermission, "无权限新增数据")
		}
		// 校验字段权限（新增记录时可能需要检查字段权限）
		if err := checkFieldPermission(c, advancedPerm, config, "add"); err != nil {
			return err
		}
	case "delete":
		// 检查 data_action 是否允许删除（删除需要 edit 或 manage 权限）
		if advancedPerm.DataAction != constants.ActionEdit && advancedPerm.DataAction != constants.ActionMange {
			return errorx.New(errorx.ErrNoPermission, "无权限删除数据")
		}
		// 检查 can_delete 标志
		if !advancedPerm.CanDelete {
			return errorx.New(errorx.ErrNoPermission, "无权限删除数据")
		}
		if config.CheckOperateRange && advancedPerm.OperateRange == constants.AdvancedOperateRangeOwn {
			if config.RecordID != "" {
				if !isRecordOwner(db.GetDB(), config.RecordID, userID) {
					return errorx.New(errorx.ErrNoPermission, "只能删除自己创建的记录")
				}
			}
		}
	case "edit":
		if advancedPerm.DataAction != constants.ActionEdit && advancedPerm.DataAction != constants.ActionMange {
			return errorx.New(errorx.ErrNoPermission, "无权限编辑数据")
		}
		if config.CheckOperateRange && advancedPerm.OperateRange == constants.AdvancedOperateRangeOwn {
			if config.RecordID != "" {
				if !isRecordOwner(db.GetDB(), config.RecordID, userID) {
					return errorx.New(errorx.ErrNoPermission, "只能编辑自己创建的记录")
				}
			}
		}
		// 校验字段权限
		if err := checkFieldPermission(c, advancedPerm, config, "edit"); err != nil {
			return err
		}
	case "read":
		if advancedPerm.DataAction == constants.ActionNone {
			return errorx.New(errorx.ErrNoPermission, "无权限查看数据")
		}
		// 校验字段权限
		if err := checkFieldPermission(c, advancedPerm, config, "read"); err != nil {
			return err
		}
	case "view_operate":
		if !advancedPerm.CanOperateView {
			return errorx.New(errorx.ErrNoPermission, "无权限操作视图")
		}
	case "view_check":
		if config.CheckViewPerm && config.ViewID != "" {
			if advancedPerm.ViewAccess == constants.AdvancedViewAccessCustom {
				if !slices.Contains(advancedPerm.CanCheckViews, config.ViewID) {
					return errorx.New(errorx.ErrNoPermission, "无权限查看该视图")
				}
			}
		}
	case "field_manage":
		// 字段管理操作（创建、更新、删除字段）需要 manage 权限
		if advancedPerm.DataAction != constants.ActionMange {
			return errorx.New(errorx.ErrNoPermission, "无权限管理字段")
		}
	}

	return nil
}
