package server

import (
	collabHandler "mvtable/internal/app/collaboration/handler"
	handler3 "mvtable/internal/app/file/handler"
	handler2 "mvtable/internal/app/invite_code/handler"
	handler21 "mvtable/internal/app/mv_dashboard/handler"
	handler20 "mvtable/internal/app/mv_dashboard_chart/handler"
	handler24 "mvtable/internal/app/mv_doc/handler"
	handler8 "mvtable/internal/app/mv_field/handler"
	handler9 "mvtable/internal/app/mv_folder/handler"
	handler17 "mvtable/internal/app/mv_form_submit/handler"
	handler6 "mvtable/internal/app/mv_project/handler"
	handler13 "mvtable/internal/app/mv_project_advanced_perm/handler"
	handler12 "mvtable/internal/app/mv_project_perm/handler"
	handler11 "mvtable/internal/app/mv_project_state/handler"
	handler10 "mvtable/internal/app/mv_record/handler"
	handler22 "mvtable/internal/app/mv_rich_text_content/handler"
	handler7 "mvtable/internal/app/mv_table_schema/handler"
	handler19 "mvtable/internal/app/mv_template_project/handler"
	handler18 "mvtable/internal/app/mv_template_tag/handler"
	handler14 "mvtable/internal/app/mv_view/handler"
	handler23 "mvtable/internal/app/mv_view_board/handler"
	handler15 "mvtable/internal/app/mv_view_form/handler"
	handler16 "mvtable/internal/app/mv_view_table/handler"
	handler5 "mvtable/internal/app/team/handler"
	"mvtable/internal/app/user/handler"
	handler4 "mvtable/internal/app/website_config/handler"
	handler26 "mvtable/internal/app/website_visit/handler"
	"mvtable/internal/middleware"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Router struct {
	userHandler                  *handler.UserHandler
	inviteCodeHandler            *handler2.InviteCodeHandler
	fileHandler                  *handler3.FileHandler
	websiteConfigHandler         *handler4.WebsiteConfigHandler
	teamHandler                  *handler5.TeamHandler
	mvProjectHandler             *handler6.MvProjectHandler
	mvTableSchemaHandler         *handler7.MvTableSchemaHandler
	mvFieldHandler               *handler8.MvFieldHandler
	mvFolderHandler              *handler9.MvFolderHandler
	mvRecordHandler              *handler10.MvRecordHandler
	mvProjectStateHandler        *handler11.MvProjectStateHandler
	mvProjectPermissionHandler   *handler12.MvProjectPermissionHandler
	mvProjectAdvancedPermHandler *handler13.MvProjectAdvancedPermHandler
	mvViewHandler                *handler14.MvViewHandler
	mvViewFormHandler            *handler15.MvViewFormHandler
	mvViewTableHandler           *handler16.MvViewTableHandler
	mvViewBoardHandler           *handler23.MvViewBoardHandler
	mvDocHandler                 *handler24.MvDocHandler
	mvFormSubmitHandler          *handler17.MvFormSubmitHandler
	mvTemplateTagHandler         *handler18.MvTemplateTagHandler
	mvTemplateProjectHandler     *handler19.MvTemplateProjectHandler
	mvDashboardChartHandler      *handler20.MvDashboardChartHandler
	mvDashboardHandler           *handler21.MvDashboardHandler
	mvRichTextContentHandler     *handler22.MvRichTextContentHandler
	collaborationHandler         *collabHandler.CollaborationHandler
	websiteVisitHandler          *handler26.WebsiteVisitHandler
}

func NewRouter(
	userHandler *handler.UserHandler,
	inviteCodeHandler *handler2.InviteCodeHandler,
	fileHandler *handler3.FileHandler,
	websiteConfigHandler *handler4.WebsiteConfigHandler,
	teamHandler *handler5.TeamHandler,
	projectHandler *handler6.MvProjectHandler,
	tableSchemaHandler *handler7.MvTableSchemaHandler,
	fieldHandler *handler8.MvFieldHandler,
	folderHandler *handler9.MvFolderHandler,
	recordHandler *handler10.MvRecordHandler,
	mvProjectStateHandler *handler11.MvProjectStateHandler,
	mvProjectPermissionHandler *handler12.MvProjectPermissionHandler,
	mvProjectAdvancedPermHandler *handler13.MvProjectAdvancedPermHandler,
	mvViewHandler *handler14.MvViewHandler,
	mvViewFormHandler *handler15.MvViewFormHandler,
	mvViewTableHandler *handler16.MvViewTableHandler,
	mvViewBoardHandler *handler23.MvViewBoardHandler,
	mvDocHandler *handler24.MvDocHandler,
	mvFormSubmitHandler *handler17.MvFormSubmitHandler,
	mvTemplateTagHandler *handler18.MvTemplateTagHandler,
	mvTemplateProjectHandler *handler19.MvTemplateProjectHandler,
	mvDashboardChartHandler *handler20.MvDashboardChartHandler,
	mvDashboardHandler *handler21.MvDashboardHandler,
	mvRichTextContentHandler *handler22.MvRichTextContentHandler,
	collaborationHandler *collabHandler.CollaborationHandler,
	websiteVisitHandler *handler26.WebsiteVisitHandler,
) *Router {
	return &Router{
		userHandler:                  userHandler,
		inviteCodeHandler:            inviteCodeHandler,
		fileHandler:                  fileHandler,
		websiteConfigHandler:         websiteConfigHandler,
		teamHandler:                  teamHandler,
		mvProjectHandler:             projectHandler,
		mvTableSchemaHandler:         tableSchemaHandler,
		mvFieldHandler:               fieldHandler,
		mvFolderHandler:              folderHandler,
		mvRecordHandler:              recordHandler,
		mvProjectStateHandler:        mvProjectStateHandler,
		mvProjectPermissionHandler:   mvProjectPermissionHandler,
		mvProjectAdvancedPermHandler: mvProjectAdvancedPermHandler,
		mvViewHandler:                mvViewHandler,
		mvViewFormHandler:            mvViewFormHandler,
		mvViewTableHandler:           mvViewTableHandler,
		mvViewBoardHandler:           mvViewBoardHandler,
		mvDocHandler:                 mvDocHandler,
		mvFormSubmitHandler:          mvFormSubmitHandler,
		mvTemplateTagHandler:         mvTemplateTagHandler,
		mvTemplateProjectHandler:     mvTemplateProjectHandler,
		mvDashboardChartHandler:      mvDashboardChartHandler,
		mvDashboardHandler:           mvDashboardHandler,
		mvRichTextContentHandler:     mvRichTextContentHandler,
		collaborationHandler:         collaborationHandler,
		websiteVisitHandler:          websiteVisitHandler,
	}
}

func (r *Router) SetupRoutes(engine *gin.Engine) {
	// WebSocket 路由需要在 timeout 中间件之前注册
	// 因为 WebSocket 会 hijack HTTP 连接，timeout 中间件不兼容
	engine.GET("/api/ws/collaboration", middleware.AuthMiddleware(), r.collaborationHandler.HandleWebSocket)

	api := engine.Group("/api")

	// 应用访问记录中间件（记录所有API访问）
	api.Use(middleware.VisitRecordMiddleware())

	// 健康检查
	// @Summary 健康检查
	// @Description 检查服务器运行状态
	// @Tags 系统
	// @Accept json
	// @Produce json
	// @Success 200 {object} middleware.Response "服务器运行正常"
	// @Router /health [get]
	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"msg":    "Server is running",
			"data":   time.Now().Format(time.RFC3339),
		})
	})

	{
		// 管理后台
		admin := api.Group("/v1/admin")
		{
			admin.POST("/login", r.userHandler.AdminLogin)

			admin.GET("/user/list", r.userHandler.AdminGetUserList)
			admin.POST("/user/info/update", r.userHandler.AdminUpdateInfo)
			admin.POST("/user/delete", middleware.AuthMiddleware(), r.userHandler.AdminDeleteUser)

			admin.POST("/invite-code/batch-create", middleware.AuthMiddleware(), r.inviteCodeHandler.AdminBatchCreate)
			admin.GET("/invite-code/list", middleware.AuthMiddleware(), r.inviteCodeHandler.AdminGetList)

			// 文件列表
			admin.GET("/file/list", middleware.AuthMiddleware(), r.fileHandler.AdminGetList)

			// 保存网站配置
			admin.POST("/website-config/save", middleware.AuthMiddleware(), r.websiteConfigHandler.AdminSave)

			// 团队列表
			admin.GET("/team/list", middleware.AuthMiddleware(), r.teamHandler.AdminGetTeamList)

			// 模板标签管理
			admin.POST("/template-tag/create", middleware.AuthMiddleware(), r.mvTemplateTagHandler.CreateMvTemplateTag)
			admin.POST("/template-tag/update", middleware.AuthMiddleware(), r.mvTemplateTagHandler.UpdateMvTemplateTag)
			admin.POST("/template-tag/delete", middleware.AuthMiddleware(), r.mvTemplateTagHandler.DeleteMvTemplateTag)
			admin.POST("/template-tag/sort", middleware.AuthMiddleware(), r.mvTemplateTagHandler.UpdateTagSort)

			// 模板项目管理
			admin.POST("/template-project/save-as-template", middleware.AuthMiddleware(), r.mvTemplateProjectHandler.SaveAsTemplate)
			admin.POST("/template-project/update", middleware.AuthMiddleware(), r.mvTemplateProjectHandler.UpdateMvTemplateProject)
			admin.POST("/template-project/delete", middleware.AuthMiddleware(), r.mvTemplateProjectHandler.DeleteMvTemplateProject)

			// 网站统计数据
			admin.GET("/website/statistics", middleware.AuthMiddleware(), r.websiteVisitHandler.GetStatistics)
			admin.GET("/website/new-users", middleware.AuthMiddleware(), r.websiteVisitHandler.GetNewUsers)
		}

		// 公开路由
		public := api.Group("/v1")
		{

			public.POST("/auth/register", r.userHandler.EmailRegister)
			public.POST("/auth/invite-code/register", r.userHandler.InviteCodeRegister)
			public.POST("/auth/register/send-code", r.userHandler.SendEmailVerifyCode)
			public.POST("/auth/login", r.userHandler.Login)
			public.POST("/auth/refresh", r.userHandler.RefreshToken)

			// 文件上传接口
			public.POST("/file/upload", middleware.AuthMiddleware(), r.fileHandler.Upload)

			// 获取网站配置
			public.GET("/website-config", r.websiteConfigHandler.Get)
		}

		// 用户路由
		user := api.Group("/v1/user")
		{
			user.GET("/info", middleware.AuthMiddleware(), r.userHandler.GetUser)
			user.POST("/update", middleware.AuthMiddleware(), r.userHandler.UpdateInfo)
			user.POST("/password/reset", r.userHandler.ResetPassword)
		}

		// 团队路由
		team := api.Group("/v1/team")
		team.Use(middleware.AuthMiddleware())
		{
			team.POST("/create", middleware.AuthMiddleware(), r.teamHandler.Create)
			team.POST("/update", middleware.AuthMiddleware(), r.teamHandler.Update)
			team.POST("/delete", middleware.AuthMiddleware(), r.teamHandler.Delete)
			team.GET("/list", middleware.AuthMiddleware(), r.teamHandler.List)

			team.GET("/member-list", middleware.AuthMiddleware(), r.teamHandler.GetMemberList)
			team.GET("/search-user", middleware.AuthMiddleware(), r.teamHandler.SearchUser)
			team.POST("/batch-get-users", middleware.AuthMiddleware(), r.teamHandler.BatchGetUsers)
			team.POST("/add-member", middleware.AuthMiddleware(), r.teamHandler.AddMember)
			team.POST("/update-member", middleware.AuthMiddleware(), r.teamHandler.UpdateMember)
			team.POST("/delete-member", middleware.AuthMiddleware(), r.teamHandler.DeleteMember)
		}

		// 项目路由
		project := api.Group("/v1/project")
		project.Use(middleware.AuthMiddleware())
		{
			project.POST("/create", middleware.AuthMiddleware(), r.mvProjectHandler.CreateMvProject)
			project.POST("/update", middleware.AuthMiddleware(), r.mvProjectHandler.UpdateMvProject)
			project.POST("/delete", middleware.AuthMiddleware(), r.mvProjectHandler.DeleteMvProject)
			project.POST("/favorite/set", middleware.AuthMiddleware(), r.mvProjectHandler.SetFavoriteProject)
			project.GET("/get", middleware.AuthMiddleware(), r.mvProjectHandler.GetProject)
			project.GET("/query", middleware.AuthMiddleware(), r.mvProjectHandler.QueryProject)
			project.POST("/node/create", middleware.AuthMiddleware(), r.mvProjectHandler.CreateProjectNode)
			project.POST("/node/rename", middleware.AuthMiddleware(), r.mvProjectHandler.RenameProjectNode)
			project.GET("/tables", middleware.AuthMiddleware(), r.mvProjectHandler.GetProjectTables)
		}

		// 文件夹路由
		folder := api.Group("/v1/folder")
		folder.Use(middleware.AuthMiddleware())
		{
			folder.POST("/create", middleware.AuthMiddleware(), r.mvFolderHandler.CreateMvFolder)
			folder.POST("/update", middleware.AuthMiddleware(), r.mvFolderHandler.UpdateMvFolder)
			folder.POST("/delete", middleware.AuthMiddleware(), r.mvFolderHandler.DeleteMvFolder)
			folder.POST("/move", middleware.AuthMiddleware(), r.mvFolderHandler.MoveNode)
			folder.GET("/subquery", middleware.AuthMiddleware(), r.mvFolderHandler.Subquery)
			folder.GET("/all-nodes", middleware.AuthMiddleware(), r.mvFolderHandler.QueryAllNodes)
			folder.GET("/search", middleware.AuthMiddleware(), r.mvFolderHandler.Search)
			folder.GET("/list", middleware.AuthMiddleware(), r.mvFolderHandler.ListProjectFolders)
		}

		// 表格路由
		table := api.Group("/v1/table")
		table.Use(middleware.AuthMiddleware())
		{
			table.POST("/create", middleware.AuthMiddleware(), r.mvTableSchemaHandler.CreateMvTableSchema)
			table.POST("/update", middleware.AuthMiddleware(), r.mvTableSchemaHandler.UpdateMvTableSchema)
			table.POST("/delete", middleware.AuthMiddleware(), r.mvTableSchemaHandler.DeleteMvTableSchema)
			table.POST("/copy", middleware.AuthMiddleware(), r.mvTableSchemaHandler.CopyMvTableSchema)
			table.GET("/get",
				middleware.AuthMiddleware(),
				middleware.AdvancedPermMiddleware(middleware.AdvancedPermConfig{
					Action:         "read",
					CheckFieldPerm: true,
				}),
				r.mvTableSchemaHandler.GetTableData)
		}

		// 字段路由
		field := api.Group("/v1/field")
		field.Use(middleware.AuthMiddleware())
		{
			field.GET("/list",
				middleware.AuthMiddleware(),
				middleware.AdvancedPermMiddleware(middleware.AdvancedPermConfig{
					Action: "read",
				}),
				r.mvFieldHandler.GetMvFieldList)
			field.POST("/create",
				middleware.AuthMiddleware(),
				middleware.AdvancedPermMiddleware(middleware.AdvancedPermConfig{
					Action: "field_manage",
				}),
				r.mvFieldHandler.CreateMvField)
			field.POST("/update",
				middleware.AuthMiddleware(),
				middleware.AdvancedPermMiddleware(middleware.AdvancedPermConfig{
					Action: "field_manage",
				}),
				r.mvFieldHandler.UpdateMvField)
			field.POST("/sort",
				middleware.AuthMiddleware(),
				middleware.AdvancedPermMiddleware(middleware.AdvancedPermConfig{
					Action: "field_manage",
				}),
				r.mvFieldHandler.SortMvFields)
			field.POST("/delete",
				middleware.AuthMiddleware(),
				middleware.AdvancedPermMiddleware(middleware.AdvancedPermConfig{
					Action: "field_manage",
				}),
				r.mvFieldHandler.DeleteMvField)
		}

		// 记录路由
		record := api.Group("/v1/record")
		record.Use(middleware.AuthMiddleware())
		{
			record.POST("/insert",
				middleware.AuthMiddleware(),
				middleware.AdvancedPermMiddleware(middleware.AdvancedPermConfig{
					Action: "add",
				}),
				r.mvRecordHandler.InsertMvRecord)
			record.GET("/get",
				middleware.AuthMiddleware(),
				middleware.AdvancedPermMiddleware(middleware.AdvancedPermConfig{
					Action: "read",
				}),
				r.mvRecordHandler.GetRecord)
			record.GET("/list",
				middleware.AuthMiddleware(),
				middleware.AdvancedPermMiddleware(middleware.AdvancedPermConfig{
					Action: "read",
				}),
				r.mvRecordHandler.GetRecords)
			record.GET("/options",
				middleware.AuthMiddleware(),
				middleware.AdvancedPermMiddleware(middleware.AdvancedPermConfig{
					Action: "read",
				}),
				r.mvRecordHandler.GetRecordOptions)
			record.POST("/sort",
				middleware.AuthMiddleware(),
				middleware.AdvancedPermMiddleware(middleware.AdvancedPermConfig{
					Action:            "edit",
					CheckOperateRange: true,
				}),
				r.mvRecordHandler.SaveRowSort)
			record.POST("/cell/update",
				middleware.AuthMiddleware(),
				middleware.AdvancedPermMiddleware(middleware.AdvancedPermConfig{
					Action:            "edit",
					CheckFieldPerm:    true,
					CheckOperateRange: true,
				}),
				r.mvRecordHandler.UpdateCell)
			record.POST("/row/update",
				middleware.AuthMiddleware(),
				middleware.AdvancedPermMiddleware(middleware.AdvancedPermConfig{
					Action:            "edit",
					CheckOperateRange: true,
				}),
				r.mvRecordHandler.UpdateRow)
			record.POST("/delete",
				middleware.AuthMiddleware(),
				middleware.AdvancedPermMiddleware(middleware.AdvancedPermConfig{
					Action:            "delete",
					CheckOperateRange: true,
				}),
				r.mvRecordHandler.DeleteRecords)
			record.POST("/import",
				middleware.AuthMiddleware(),
				middleware.AdvancedPermMiddleware(middleware.AdvancedPermConfig{
					Action: "add",
				}),
				r.mvRecordHandler.ImportRecords)
		}

		// 富文本路由
		richText := api.Group("/v1/rich-text")
		richText.Use(middleware.AuthMiddleware())
		{
			richText.GET("/get",
				middleware.AuthMiddleware(),
				middleware.AdvancedPermMiddleware(middleware.AdvancedPermConfig{
					Action: "read",
				}),
				r.mvRichTextContentHandler.GetRichTextContent)
		}

		// 项目分享范围路由
		projectState := api.Group("/v1/project/state")
		projectState.Use(middleware.AuthMiddleware())
		{
			projectState.POST("/set", middleware.AuthMiddleware(), r.mvProjectStateHandler.SetShareRange)
		}

		// 项目权限路由
		projectPermission := api.Group("/v1/project/member")
		projectPermission.Use(middleware.AuthMiddleware())
		{
			projectPermission.POST("/batch-create", middleware.AuthMiddleware(), r.mvProjectPermissionHandler.BatchCreateMember)
			projectPermission.POST("/update", middleware.AuthMiddleware(), r.mvProjectPermissionHandler.UpdateMember)
			projectPermission.POST("/delete", middleware.AuthMiddleware(), r.mvProjectPermissionHandler.DeleteMember)
			projectPermission.GET("/list", middleware.AuthMiddleware(), r.mvProjectPermissionHandler.ListMembers)
			projectPermission.GET("/search", middleware.AuthMiddleware(), r.mvProjectPermissionHandler.SearchMember)
		}

		// 项目高级权限路由
		projectAdvancedPerm := api.Group("/v1/project/advanced-perm")
		projectAdvancedPerm.Use(middleware.AuthMiddleware())
		{
			projectAdvancedPerm.POST("/enable", middleware.AuthMiddleware(), r.mvProjectAdvancedPermHandler.EnableMvProjectAdvancedPerm)
			projectAdvancedPerm.POST("/disable", middleware.AuthMiddleware(), r.mvProjectAdvancedPermHandler.DisableMvProjectAdvancedPerm)
			projectAdvancedPerm.GET("/get", middleware.AuthMiddleware(), r.mvProjectAdvancedPermHandler.GetMvProjectAdvancedPerm)
			projectAdvancedPerm.PUT("/update", middleware.AuthMiddleware(), r.mvProjectAdvancedPermHandler.UpdateMvProjectAdvancedPerm)
		}

		// 视图路由
		view := api.Group("/v1/view")
		view.Use(middleware.AuthMiddleware())
		{
			view.GET("/query",
				middleware.AuthMiddleware(),
				middleware.AdvancedPermMiddleware(middleware.AdvancedPermConfig{
					Action: "read",
				}),
				r.mvViewHandler.QueryMvView)
			view.GET("/get",
				middleware.AuthMiddleware(),
				middleware.AdvancedPermMiddleware(middleware.AdvancedPermConfig{
					Action:        "view_check",
					CheckViewPerm: true,
				}),
				r.mvViewHandler.GetMvView)
			view.POST("/create",
				middleware.AuthMiddleware(),
				middleware.AdvancedPermMiddleware(middleware.AdvancedPermConfig{
					Action: "view_operate",
				}),
				r.mvViewHandler.CreateMvView)
			view.POST("/update",
				middleware.AuthMiddleware(),
				middleware.AdvancedPermMiddleware(middleware.AdvancedPermConfig{
					Action: "view_operate",
				}),
				r.mvViewHandler.UpdateMvView)
			view.POST("/delete",
				middleware.AuthMiddleware(),
				middleware.AdvancedPermMiddleware(middleware.AdvancedPermConfig{
					Action: "view_operate",
				}),
				r.mvViewHandler.DeleteMvView)
			view.POST("/switch-active",
				middleware.AuthMiddleware(),
				middleware.AdvancedPermMiddleware(middleware.AdvancedPermConfig{
					Action: "read",
				}),
				r.mvViewHandler.SwitchActiveView)
		}

		// 视图表单路由
		viewForm := api.Group("/v1/view/form")
		viewForm.Use(middleware.AuthMiddleware())
		{
			viewForm.POST("/update", middleware.AuthMiddleware(), r.mvViewFormHandler.UpdateMvViewForm)
		}

		// 视图表格路由
		viewTable := api.Group("/v1/view/table")
		viewTable.Use(middleware.AuthMiddleware())
		{
			viewTable.POST("/update", middleware.AuthMiddleware(), r.mvViewTableHandler.UpdateMvViewTable)
		}

		// 视图看板路由
		viewBoard := api.Group("/v1/view/board")
		viewBoard.Use(middleware.AuthMiddleware())
		{
			viewBoard.POST("/update", middleware.AuthMiddleware(), r.mvViewBoardHandler.UpdateMvViewBoard)
		}

		// 仪表盘图表路由
		dashboardChart := api.Group("/v1/dashboard/chart")
		dashboardChart.Use(middleware.AuthMiddleware())
		{
			dashboardChart.POST("/create", middleware.AuthMiddleware(), r.mvDashboardChartHandler.CreateMvDashboardChart)
			dashboardChart.POST("/update", middleware.AuthMiddleware(), r.mvDashboardChartHandler.UpdateMvDashboardChart)
			dashboardChart.POST("/delete", middleware.AuthMiddleware(), r.mvDashboardChartHandler.DeleteMvDashboardChart)
		}

		// 表单提交路由
		formSubmit := api.Group("/v1/form")
		formSubmit.Use(middleware.AuthMiddleware())
		{
			formSubmit.POST("/submit", middleware.AuthMiddleware(), r.mvFormSubmitHandler.SubmitForm)
		}

		// 模板标签路由（查询接口，普通用户可访问）
		templateTag := api.Group("/v1/template-tag")
		templateTag.Use(middleware.AuthMiddleware())
		{
			templateTag.GET("/query", middleware.AuthMiddleware(), r.mvTemplateTagHandler.QueryMvTemplateTag)
		}

		// 模板项目路由
		templateProject := api.Group("/v1/template-project")
		templateProject.Use(middleware.AuthMiddleware())
		{
			templateProject.GET("/query", middleware.AuthMiddleware(), r.mvTemplateProjectHandler.QueryMvTemplateProject)
		}

		// 仪表盘路由
		dashboard := api.Group("/v1/dashboard")
		dashboard.Use(middleware.AuthMiddleware())
		{
			dashboard.POST("/create", middleware.AuthMiddleware(), r.mvDashboardHandler.CreateMvDashboard)
			dashboard.POST("/update", middleware.AuthMiddleware(), r.mvDashboardHandler.UpdateMvDashboard)
			dashboard.POST("/delete", middleware.AuthMiddleware(), r.mvDashboardHandler.DeleteMvDashboard)
			dashboard.POST("/copy", middleware.AuthMiddleware(), r.mvDashboardHandler.CopyMvDashboard)
			dashboard.GET("/get", middleware.AuthMiddleware(), r.mvDashboardHandler.GetDashboard)
		}

		// 文档路由
		doc := api.Group("/v1/doc")
		doc.Use(middleware.AuthMiddleware())
		{
			doc.POST("/create", middleware.AuthMiddleware(), r.mvDocHandler.CreateMvDoc)
			doc.POST("/update", middleware.AuthMiddleware(), r.mvDocHandler.UpdateMvDoc)
			doc.POST("/delete", middleware.AuthMiddleware(), r.mvDocHandler.DeleteMvDoc)
			doc.GET("/get", middleware.AuthMiddleware(), r.mvDocHandler.GetMvDoc)
			doc.GET("/list", middleware.AuthMiddleware(), r.mvDocHandler.ListMvDocs)
		}
	}
}
