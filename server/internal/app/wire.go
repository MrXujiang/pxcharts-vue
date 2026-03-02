package app

import (
	collabHandler "mvtable/internal/app/collaboration/handler"
	collabService "mvtable/internal/app/collaboration/service"
	handler3 "mvtable/internal/app/file/handler"
	service3 "mvtable/internal/app/file/service"
	handler2 "mvtable/internal/app/invite_code/handler"
	service2 "mvtable/internal/app/invite_code/service"
	handler21 "mvtable/internal/app/mv_dashboard/handler"
	service21 "mvtable/internal/app/mv_dashboard/service"
	handler20 "mvtable/internal/app/mv_dashboard_chart/handler"
	service20 "mvtable/internal/app/mv_dashboard_chart/service"
	handler8 "mvtable/internal/app/mv_field/handler"
	service8 "mvtable/internal/app/mv_field/service"
	handler9 "mvtable/internal/app/mv_folder/handler"
	service9 "mvtable/internal/app/mv_folder/service"
	handler17 "mvtable/internal/app/mv_form_submit/handler"
	service17 "mvtable/internal/app/mv_form_submit/service"
	handler6 "mvtable/internal/app/mv_project/handler"
	service6 "mvtable/internal/app/mv_project/service"
	handler13 "mvtable/internal/app/mv_project_advanced_perm/handler"
	service13 "mvtable/internal/app/mv_project_advanced_perm/service"
	handler12 "mvtable/internal/app/mv_project_perm/handler"
	service12 "mvtable/internal/app/mv_project_perm/service"
	handler11 "mvtable/internal/app/mv_project_state/handler"
	service11 "mvtable/internal/app/mv_project_state/service"
	handler10 "mvtable/internal/app/mv_record/handler"
	service10 "mvtable/internal/app/mv_record/service"
	handler22 "mvtable/internal/app/mv_rich_text_content/handler"
	service22 "mvtable/internal/app/mv_rich_text_content/service"
	handler7 "mvtable/internal/app/mv_table_schema/handler"
	service7 "mvtable/internal/app/mv_table_schema/service"
	handler19 "mvtable/internal/app/mv_template_project/handler"
	service19 "mvtable/internal/app/mv_template_project/service"
	handler18 "mvtable/internal/app/mv_template_tag/handler"
	service18 "mvtable/internal/app/mv_template_tag/service"
	handler14 "mvtable/internal/app/mv_view/handler"
	service14 "mvtable/internal/app/mv_view/service"
	handler23 "mvtable/internal/app/mv_view_board/handler"
	service23 "mvtable/internal/app/mv_view_board/service"
	handler24 "mvtable/internal/app/mv_doc/handler"
	service24 "mvtable/internal/app/mv_doc/service"
	handler15 "mvtable/internal/app/mv_view_form/handler"
	service15 "mvtable/internal/app/mv_view_form/service"
	handler16 "mvtable/internal/app/mv_view_table/handler"
	service16 "mvtable/internal/app/mv_view_table/service"
	handler5 "mvtable/internal/app/team/handler"
	service5 "mvtable/internal/app/team/service"
	"mvtable/internal/app/user/handler"
	"mvtable/internal/app/user/service"
	handler4 "mvtable/internal/app/website_config/handler"
	service4 "mvtable/internal/app/website_config/service"
	handler25 "mvtable/internal/app/website_visit/handler"
	service25 "mvtable/internal/app/website_visit/service"

	"github.com/google/wire"
)

var Set = wire.NewSet(
	handler.NewUserHandler,
	service.NewUserService,

	handler2.NewInviteCodeHandler,
	service2.NewInviteCodeService,

	handler3.NewFileHandler,
	service3.NewFileService,

	handler4.NewWebsiteConfigHandler,
	service4.NewWebsiteConfigService,

	handler25.NewWebsiteVisitHandler,
	service25.NewWebsiteVisitService,

	handler5.NewTeamHandler,
	service5.NewTeamService,

	handler6.NewMvProjectHandler,
	service6.NewMvProjectService,

	handler7.NewMvTableSchemaHandler,
	service7.NewMvTableSchemaService,

	handler8.NewMvFieldHandler,
	service8.NewMvFieldService,

	handler9.NewMvFolderHandler,
	service9.NewMvFolderService,

	handler10.NewMvRecordHandler,
	service10.NewMvRecordService,

	handler11.NewMvProjectStateHandler,
	service11.NewMvProjectStateService,

	handler12.NewMvProjectPermissionHandler,
	service12.NewMvProjectPermissionService,

	handler13.NewMvProjectAdvancedPermHandler,
	service13.NewMvProjectAdvancedPermService,

	handler14.NewMvViewHandler,
	service14.NewMvViewFormService,

	handler15.NewMvViewFormHandler,
	service15.NewMvViewFormService,

	handler16.NewMvViewTableHandler,
	service16.NewMvViewTableService,

	handler17.NewMvFormSubmitHandler,
	service17.NewMvFormSubmitService,

	handler18.NewMvTemplateTagHandler,
	service18.NewMvTemplateTagService,

	handler19.NewMvTemplateProjectHandler,
	service19.NewMvTemplateProjectService,

	handler20.NewMvDashboardChartHandler,
	service20.NewMvDashboardChartService,

	handler21.NewMvDashboardHandler,
	service21.NewMvDashboardService,

	handler22.NewMvRichTextContentHandler,
	service22.NewMvRichTextContentService,

	handler23.NewMvViewBoardHandler,
	service23.NewMvViewBoardService,

	handler24.NewMvDocHandler,
	service24.NewMvDocService,

	// 协同编辑
	collabService.NewHub,
	collabService.NewLockManager,
	collabHandler.NewCollaborationHandler,
	collabService.NewCollaborationService,
)
