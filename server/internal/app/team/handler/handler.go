package handler

import (
	"mvtable/internal/app/team/model"
	"mvtable/internal/app/team/service"
	"mvtable/internal/middleware"
	"mvtable/internal/pkg/errorx"

	"github.com/gin-gonic/gin"
)

type TeamHandler struct {
	teamService *service.TeamService
}

func NewTeamHandler(teamService *service.TeamService) *TeamHandler {
	return &TeamHandler{
		teamService: teamService,
	}
}

// Create 创建团队
// @Summary 创建团队
// @Description 创建新的团队
// @Tags 团队
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.CreateTeamReq true "创建团队请求"
// @Success 200 {object} middleware.Response "创建成功"
// @Router /v1/team/create [post]
func (h *TeamHandler) Create(c *gin.Context) {
	userId := middleware.GetCurrentUserID(c)
	var req model.CreateTeamReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	err := h.teamService.Create(userId, &req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, nil, "")
}

// Update 更新团队
// @Summary 更新团队
// @Description 更新团队信息
// @Tags 团队
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.UpdateTeamReq true "更新团队请求"
// @Success 200 {object} middleware.Response "更新成功"
// @Router /v1/team/update [post]
func (h *TeamHandler) Update(c *gin.Context) {
	userId := middleware.GetCurrentUserID(c)
	var req model.UpdateTeamReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	err := h.teamService.Update(userId, &req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, nil, "")
}

// Delete 删除团队
// @Summary 删除团队
// @Description 删除团队
// @Tags 团队
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.DeleteTeamReq true "删除团队请求"
// @Success 200 {object} middleware.Response "删除成功"
// @Router /v1/team/delete [post]
func (h *TeamHandler) Delete(c *gin.Context) {
	userId := middleware.GetCurrentUserID(c)
	var req model.DeleteTeamReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	err := h.teamService.Delete(userId, &req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, nil, "")
}

// List 获取团队列表
// @Summary 获取团队列表
// @Description 获取当前用户的团队列表
// @Tags 团队
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} middleware.Response{data=model.ListTeamReq} "获取成功"
// @Router /v1/team/list [get]
func (h *TeamHandler) List(c *gin.Context) {
	userId := middleware.GetCurrentUserID(c)
	resp, err := h.teamService.List(userId)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, resp, "")
}

// GetMemberList 获取团队成员列表
// @Summary 获取团队成员列表
// @Description 获取团队成员列表
// @Tags 团队
// @Accept json
// @Produce json
// @Security Bearer
// @Param teamId query string true "团队ID"
// @Param searchWord query string false "搜索关键词"
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} middleware.Response{data=model.GetMemberListRes} "获取成功"
// @Router /v1/team/member-list [get]
func (h *TeamHandler) GetMemberList(c *gin.Context) {
	userId := middleware.GetCurrentUserID(c)
	var req model.GetMemberListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	resp, err := h.teamService.GetMemberList(userId, &req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, resp, "")
}

// SearchUser 搜索用户
// @Summary 搜索用户
// @Description 根据邮箱或昵称搜索用户
// @Tags 团队
// @Accept json
// @Produce json
// @Security Bearer
// @Param searchWord query string true "搜索关键词"
// @Success 200 {object} middleware.Response{data=model.SearchUserRes} "搜索成功"
// @Router /v1/team/search-user [get]
func (h *TeamHandler) SearchUser(c *gin.Context) {
	userId := middleware.GetCurrentUserID(c)
	var req model.SearchUserReq
	if err := c.ShouldBindQuery(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	resp, err := h.teamService.SearchUser(userId, &req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, resp, "")
}

// BatchGetUsers 批量查询用户
// @Summary 批量查询用户
// @Description 根据用户ID列表批量查询用户信息
// @Tags 团队
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.BatchGetUsersReq true "批量查询用户请求"
// @Success 200 {object} middleware.Response{data=model.BatchGetUsersRes} "查询成功"
// @Router /v1/team/batch-get-users [post]
func (h *TeamHandler) BatchGetUsers(c *gin.Context) {
	var req model.BatchGetUsersReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	resp, err := h.teamService.BatchGetUsers(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, resp, "")
}

// AddMember 添加团队成员
// @Summary 添加团队成员
// @Description 向团队添加新成员
// @Tags 团队
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.AddMemberReq true "添加成员请求"
// @Success 200 {object} middleware.Response "添加成功"
// @Router /v1/team/add-member [post]
func (h *TeamHandler) AddMember(c *gin.Context) {
	userId := middleware.GetCurrentUserID(c)
	var req model.AddMemberReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	err := h.teamService.AddMember(userId, &req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, nil, "")
}

// UpdateMember 更新团队成员
// @Summary 更新团队成员
// @Description 更新团队成员信息
// @Tags 团队
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.UpdateMemberReq true "更新成员请求"
// @Success 200 {object} middleware.Response "更新成功"
// @Router /v1/team/update-member [post]
func (h *TeamHandler) UpdateMember(c *gin.Context) {
	userId := middleware.GetCurrentUserID(c)
	var req model.UpdateMemberReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	err := h.teamService.UpdateMember(userId, &req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, nil, "")
}

// DeleteMember 删除团队成员
// @Summary 删除团队成员
// @Description 从团队中删除成员
// @Tags 团队
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.DeleteMemberReq true "删除成员请求"
// @Success 200 {object} middleware.Response "删除成功"
// @Router /v1/team/delete-member [post]
func (h *TeamHandler) DeleteMember(c *gin.Context) {
	userId := middleware.GetCurrentUserID(c)
	var req model.DeleteMemberReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	err := h.teamService.DeleteMember(userId, &req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, nil, "")
}
