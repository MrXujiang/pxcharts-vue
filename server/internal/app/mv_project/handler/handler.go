package handler

import (
	"mvtable/internal/app/mv_project/model"
	"mvtable/internal/app/mv_project/service"
	"mvtable/internal/middleware"
	"mvtable/internal/pkg/errorx"

	"github.com/gin-gonic/gin"
)

type MvProjectHandler struct {
	service *service.MvProjectService
}

func NewMvProjectHandler(service *service.MvProjectService) *MvProjectHandler {
	return &MvProjectHandler{service: service}
}

// CreateMvProject 创建项目
// @Summary 创建项目
// @Description 创建新的MV项目
// @Tags 项目
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.CreateMvProjectReq true "创建项目请求"
// @Success 200 {object} middleware.Response{data=model.CreateMvProjectRes} "创建成功"
// @Router /v1/project/create [post]
func (h *MvProjectHandler) CreateMvProject(c *gin.Context) {
	var req model.CreateMvProjectReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	userId := middleware.GetCurrentUserID(c)
	resp, err := h.service.CreateMvProject(userId, &req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, resp, "")
}

// UpdateMvProject 更新项目
// @Summary 更新项目
// @Description 更新MV项目信息
// @Tags 项目
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.UpdateMvProjectReq true "更新项目请求"
// @Success 200 {object} middleware.Response "更新成功"
// @Router /v1/project/update [post]
func (h *MvProjectHandler) UpdateMvProject(c *gin.Context) {
	var req model.UpdateMvProjectReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	userId := middleware.GetCurrentUserID(c)
	err := h.service.UpdateMvProject(userId, &req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, nil, "")
}

// DeleteMvProject 删除项目
// @Summary 删除项目
// @Description 删除MV项目
// @Tags 项目
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.DeleteMvProjectReq true "删除项目请求"
// @Success 200 {object} middleware.Response "删除成功"
// @Router /v1/project/delete [post]
func (h *MvProjectHandler) DeleteMvProject(c *gin.Context) {
	var req model.DeleteMvProjectReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	userId := middleware.GetCurrentUserID(c)
	err := h.service.DeleteMvProject(userId, &req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, nil, "")
}

// SetFavoriteProject 设置项目为收藏或取消收藏
// @Summary 设置项目为收藏或取消收藏
// @Description 设置MV项目为收藏或取消收藏
// @Tags 项目
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.SetFavoriteProjectReq true "设置收藏请求"
// @Success 200 {object} middleware.Response "设置成功"
// @Router /v1/project/favorite/set [post]
func (h *MvProjectHandler) SetFavoriteProject(c *gin.Context) {
	var req model.SetFavoriteProjectReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	userId := middleware.GetCurrentUserID(c)
	err := h.service.SetFavoriteProject(userId, &req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, nil, "")
}

// QueryProject 查询项目列表
// @Summary 查询项目列表
// @Description 根据类型查询MV项目列表
// @Tags 项目
// @Accept json
// @Produce json
// @Security Bearer
// @Param request query model.QueryProjectReq true "查询项目请求"
// @Success 200 {object} middleware.Response{data=model.QueryProjectRes} "查询成功"
// @Router /v1/project/query [get]
func (h *MvProjectHandler) QueryProject(c *gin.Context) {
	var req model.QueryProjectReq
	if err := c.ShouldBindQuery(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	userId := middleware.GetCurrentUserID(c)
	resp, err := h.service.QueryProject(userId, &req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, resp, "")
}

// GetProject 获取项目详情
// @Summary 获取项目详情
// @Description 获取MV项目的详细信息
// @Tags 项目
// @Accept json
// @Produce json
// @Security Bearer
// @Param request query model.GetProjectReq true "获取项目请求"
// @Success 200 {object} middleware.Response{data=model.GetProjectRes} "获取成功"
// @Router /v1/project/get [get]
func (h *MvProjectHandler) GetProject(c *gin.Context) {
	var req model.GetProjectReq
	if err := c.ShouldBindQuery(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	userId := middleware.GetCurrentUserID(c)
	resp, err := h.service.GetProject(userId, &req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, resp, "")
}

// RenameProjectNode 重命名项目下节点（文件夹 / 表格 / 表单 / 仪表盘 / 文档）
// @Summary 重命名项目节点
// @Description 根据类型对项目下的 folder / table / form / dashboard / doc 进行重命名
// @Tags 项目
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.RenameProjectNodeReq true "重命名项目节点请求"
// @Success 200 {object} middleware.Response "重命名成功"
// @Router /v1/project/node/rename [post]
func (h *MvProjectHandler) RenameProjectNode(c *gin.Context) {
	var req model.RenameProjectNodeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	userId := middleware.GetCurrentUserID(c)
	if err := h.service.RenameProjectNode(userId, &req); err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, nil, "")
}

// CreateProjectNode 在项目下创建节点（文件夹 / 表格 / 表单 / 仪表盘 / 文档）
// @Summary 在项目下创建节点
// @Description 在指定项目和文件夹下创建节点，type: table（创建数据表及表格视图）、form（创建数据表及表单视图）、dashboard（创建仪表盘）、folder（创建文件夹）、doc（创建文档）
// @Tags 项目
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.CreateProjectNodeReq true "创建项目节点请求"
// @Success 200 {object} middleware.Response{data=model.CreateProjectNodeRes} "创建成功"
// @Router /v1/project/node/create [post]
func (h *MvProjectHandler) CreateProjectNode(c *gin.Context) {
	var req model.CreateProjectNodeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	userId := middleware.GetCurrentUserID(c)
	resp, err := h.service.CreateProjectNode(userId, &req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, resp, "")
}

// GetProjectTables 获取项目下的所有数据表
// @Summary 获取项目下的所有数据表
// @Description 根据项目ID获取该项目下的所有数据表列表
// @Tags 项目
// @Accept json
// @Produce json
// @Security Bearer
// @Param request query model.GetProjectTablesReq true "获取项目数据表请求"
// @Success 200 {object} middleware.Response{data=model.GetProjectTablesRes} "获取成功"
// @Router /v1/project/tables [get]
func (h *MvProjectHandler) GetProjectTables(c *gin.Context) {
	var req model.GetProjectTablesReq
	if err := c.ShouldBindQuery(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	userId := middleware.GetCurrentUserID(c)
	resp, err := h.service.GetProjectTables(userId, &req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, resp, "")
}
