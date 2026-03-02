package handler

import (
	"mvtable/internal/app/mv_template_project/model"
	"mvtable/internal/app/mv_template_project/service"
	"mvtable/internal/middleware"
	"mvtable/internal/pkg/errorx"

	"github.com/gin-gonic/gin"
)

type MvTemplateProjectHandler struct {
	service *service.MvTemplateProjectService
}

func NewMvTemplateProjectHandler(service *service.MvTemplateProjectService) *MvTemplateProjectHandler {
	return &MvTemplateProjectHandler{service: service}
}

// QueryMvTemplateProject 查询模板项目列表
// @Summary 查询模板项目列表
// @Description 获取模板项目列表，支持分页和按标签筛选
// @Tags 项目模版
// @Accept json
// @Produce json
// @Security Bearer
// @Param request query model.QueryMvTemplateProjectReq true "查询模板项目请求"
// @Success 200 {object} middleware.Response{data=model.QueryMvTemplateProjectRes} "查询成功"
// @Router /v1/template-project/query [get]
func (h *MvTemplateProjectHandler) QueryMvTemplateProject(c *gin.Context) {
	var req model.QueryMvTemplateProjectReq
	if err := c.ShouldBindQuery(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	resp, err := h.service.QueryMvTemplateProject(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, resp, "")
}

// SaveAsTemplate 保存项目为模板
// @Summary 保存项目为模板
// @Description 将项目保存为模板，复制项目的所有相关数据（文件夹、表格、字段、记录、视图等）
// @Tags 管理后台
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.SaveAsTemplateReq true "保存为模板请求"
// @Success 200 {object} middleware.Response{data=model.SaveAsTemplateRes} "保存成功"
// @Router /v1/admin/template-project/save-as-template [post]
func (h *MvTemplateProjectHandler) SaveAsTemplate(c *gin.Context) {
	userId := middleware.GetCurrentUserID(c)

	var req model.SaveAsTemplateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	resp, err := h.service.SaveAsTemplate(userId, &req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, resp, "")
}

// UpdateMvTemplateProject 更新模板项目
// @Summary 更新模板项目
// @Description 更新模板项目信息（名称、描述、封面、标签）
// @Tags 管理后台
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.UpdateMvTemplateProjectReq true "更新模板项目请求"
// @Success 200 {object} middleware.Response "更新成功"
// @Router /v1/admin/template-project/update [post]
func (h *MvTemplateProjectHandler) UpdateMvTemplateProject(c *gin.Context) {
	var req model.UpdateMvTemplateProjectReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	err := h.service.UpdateMvTemplateProject(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, nil, "")
}

// DeleteMvTemplateProject 删除模板项目
// @Summary 删除模板项目
// @Description 删除模板项目及其所有相关数据（表格、字段、记录、视图等）
// @Tags 管理后台
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.DeleteMvTemplateProjectReq true "删除模板项目请求"
// @Success 200 {object} middleware.Response "删除成功"
// @Router /v1/admin/template-project/delete [post]
func (h *MvTemplateProjectHandler) DeleteMvTemplateProject(c *gin.Context) {
	var req model.DeleteMvTemplateProjectReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	err := h.service.DeleteMvTemplateProject(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, nil, "")
}
