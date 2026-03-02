package handler

import (
	"mvtable/internal/app/mv_dashboard/model"
	"mvtable/internal/app/mv_dashboard/service"
	"mvtable/internal/middleware"
	"mvtable/internal/pkg/errorx"

	"github.com/gin-gonic/gin"
)

type MvDashboardHandler struct {
	service *service.MvDashboardService
}

func NewMvDashboardHandler(service *service.MvDashboardService) *MvDashboardHandler {
	return &MvDashboardHandler{service: service}
}

// CreateMvDashboard 创建仪表盘
// @Summary 创建仪表盘
// @Description 在项目下创建新的仪表盘
// @Tags 仪表盘
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.CreateMvDashboardReq true "创建仪表盘请求"
// @Success 200 {object} middleware.Response{data=string} "创建成功，data为仪表盘ID"
// @Router /v1/dashboard/create [post]
func (h *MvDashboardHandler) CreateMvDashboard(c *gin.Context) {
	var req model.CreateMvDashboardReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	id, err := h.service.CreateMvDashboard(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, id, "")
}

// UpdateMvDashboard 更新仪表盘
// @Summary 更新仪表盘
// @Description 更新仪表盘基础信息
// @Tags 仪表盘
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.UpdateMvDashboardReq true "更新仪表盘请求"
// @Success 200 {object} middleware.Response "更新成功"
// @Router /v1/dashboard/update [post]
func (h *MvDashboardHandler) UpdateMvDashboard(c *gin.Context) {
	var req model.UpdateMvDashboardReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	if err := h.service.UpdateMvDashboard(&req); err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, nil, "")
}

// DeleteMvDashboard 删除仪表盘
// @Summary 删除仪表盘
// @Description 删除指定的仪表盘
// @Tags 仪表盘
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.DeleteMvDashboardReq true "删除仪表盘请求"
// @Success 200 {object} middleware.Response "删除成功"
// @Router /v1/dashboard/delete [post]
func (h *MvDashboardHandler) DeleteMvDashboard(c *gin.Context) {
	var req model.DeleteMvDashboardReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	if err := h.service.DeleteMvDashboard(&req); err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, nil, "")
}

// GetDashboard 获取仪表盘详情
// @Summary 获取仪表盘详情
// @Description 获取仪表盘基础信息及其下所有图表
// @Tags 仪表盘
// @Accept json
// @Produce json
// @Security Bearer
// @Param request query model.GetDashboardReq true "获取仪表盘请求"
// @Success 200 {object} middleware.Response{data=model.GetDashboardRes} "获取成功"
// @Router /v1/dashboard/get [get]
func (h *MvDashboardHandler) GetDashboard(c *gin.Context) {
	var req model.GetDashboardReq
	if err := c.ShouldBindQuery(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	resp, err := h.service.GetDashboard(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, resp, "")
}

// CopyMvDashboard 复制仪表盘
// @Summary 复制仪表盘
// @Description 复制指定仪表盘及其下所有图表，返回新创建的仪表盘ID
// @Tags 仪表盘
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.CopyMvDashboardReq true "复制仪表盘请求"
// @Success 200 {object} middleware.Response{data=model.CopyMvDashboardRes} "复制成功"
// @Router /v1/dashboard/copy [post]
func (h *MvDashboardHandler) CopyMvDashboard(c *gin.Context) {
	var req model.CopyMvDashboardReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	resp, err := h.service.CopyMvDashboard(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, resp, "")
}
