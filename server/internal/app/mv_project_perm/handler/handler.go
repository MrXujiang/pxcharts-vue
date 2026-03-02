package handler

import (
	"mvtable/internal/app/mv_project_perm/model"
	"mvtable/internal/app/mv_project_perm/service"
	"mvtable/internal/middleware"
	"mvtable/internal/pkg/errorx"

	"github.com/gin-gonic/gin"
)

type MvProjectPermissionHandler struct {
	service *service.MvProjectPermissionService
}

func NewMvProjectPermissionHandler(service *service.MvProjectPermissionService) *MvProjectPermissionHandler {
	return &MvProjectPermissionHandler{service: service}
}

// BatchCreateMember 批量创建项目成员
// @Summary 批量创建项目成员
// @Description 批量向项目添加协作者，设置成员的权限（read/edit/manage）
// @Tags 项目权限
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.BatchCreateMemberReq true "批量创建成员请求"
// @Success 200 {object} middleware.Response "创建成功"
// @Router /v1/project/member/batch-create [post]
func (h *MvProjectPermissionHandler) BatchCreateMember(c *gin.Context) {
	var req model.BatchCreateMemberReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	userId := middleware.GetCurrentUserID(c)
	err := h.service.BatchCreateMember(userId, &req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, nil, "")
}

// UpdateMember 更新项目成员权限
// @Summary 更新项目成员权限
// @Description 更新项目协作者的权限类型（read/edit/manage）
// @Tags 项目权限
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.UpdateMemberReq true "更新成员请求"
// @Success 200 {object} middleware.Response "更新成功"
// @Router /v1/project/member/update [post]
func (h *MvProjectPermissionHandler) UpdateMember(c *gin.Context) {
	var req model.UpdateMemberReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	userId := middleware.GetCurrentUserID(c)
	err := h.service.UpdateMember(userId, &req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, nil, "")
}

// DeleteMember 删除项目成员
// @Summary 删除项目成员
// @Description 从项目中移除协作者
// @Tags 项目权限
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.DeleteMemberReq true "删除成员请求"
// @Success 200 {object} middleware.Response "删除成功"
// @Router /v1/project/member/delete [post]
func (h *MvProjectPermissionHandler) DeleteMember(c *gin.Context) {
	var req model.DeleteMemberReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	userId := middleware.GetCurrentUserID(c)
	err := h.service.DeleteMember(userId, &req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, nil, "")
}

// SearchMember 搜索非当前项目的成员
// @Summary 搜索非当前项目的成员
// @Description 根据昵称模糊匹配或者根据邮箱完全匹配
// @Tags 项目权限
// @Accept json
// @Produce json
// @Security Bearer
// @Param projectId query string true "项目ID"
// @Param keywords query string true "搜索关键词"
// @Success 200 {object} middleware.Response{data=model.SearchMemberRes} "搜索成功"
// @Router /v1/project/member/search [get]
func (h *MvProjectPermissionHandler) SearchMember(c *gin.Context) {
	var req model.SearchMemberReq
	if err := c.ShouldBindQuery(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	userId := middleware.GetCurrentUserID(c)
	resp, err := h.service.SearchMember(userId, &req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, resp, "")
}

// ListMembers 查询项目成员列表
// @Summary 查询项目成员列表
// @Description 获取项目中所有成员（用户/团队）及其角色
// @Tags 项目权限
// @Accept json
// @Produce json
// @Security Bearer
// @Param projectId query string true "项目ID"
// @Success 200 {object} middleware.Response{data=model.ProjectMemberListRes} "查询成功"
// @Router /v1/project/member/list [get]
func (h *MvProjectPermissionHandler) ListMembers(c *gin.Context) {
	var req model.ProjectMemberListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	userId := middleware.GetCurrentUserID(c)
	resp, err := h.service.ListProjectMembers(userId, &req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, resp, "")
}
