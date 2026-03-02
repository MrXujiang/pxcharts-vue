package handler

import (
	"mvtable/internal/app/mv_dashboard_chart/model"
	"mvtable/internal/app/mv_dashboard_chart/service"
	"mvtable/internal/middleware"
	"mvtable/internal/pkg/errorx"

	"github.com/gin-gonic/gin"
)

type MvDashboardChartHandler struct {
	service *service.MvDashboardChartService
}

func NewMvDashboardChartHandler(service *service.MvDashboardChartService) *MvDashboardChartHandler {
	return &MvDashboardChartHandler{service: service}
}

// CreateMvDashboardChart 创建仪表盘图表
// @Summary 创建仪表盘图表
// @Description 在指定仪表盘下创建新的图表
// @Tags 仪表盘
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.CreateMvDashboardChartReq true "创建仪表盘图表请求"
// @Success 200 {object} middleware.Response{data=string} "创建成功，data为图表ID"
// @Router /v1/dashboard/chart/create [post]
func (h *MvDashboardChartHandler) CreateMvDashboardChart(c *gin.Context) {
	var req model.CreateMvDashboardChartReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	userId := middleware.GetCurrentUserID(c)
	resp, err := h.service.CreateMvDashboardChart(userId, &req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, resp, "")
}

// UpdateMvDashboardChart 更新仪表盘图表
// @Summary 更新仪表盘图表
// @Description 更新仪表盘图表的基础信息
// @Tags 仪表盘
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.UpdateMvDashboardChartReq true "更新仪表盘图表请求"
// @Success 200 {object} middleware.Response "更新成功"
// @Router /v1/dashboard/chart/update [post]
func (h *MvDashboardChartHandler) UpdateMvDashboardChart(c *gin.Context) {
	var req model.UpdateMvDashboardChartReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	if err := h.service.UpdateMvDashboardChart(&req); err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, nil, "")
}

// DeleteMvDashboardChart 删除仪表盘图表
// @Summary 删除仪表盘图表
// @Description 删除指定的仪表盘图表
// @Tags 仪表盘
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.DeleteMvDashboardChartReq true "删除仪表盘图表请求"
// @Success 200 {object} middleware.Response "删除成功"
// @Router /v1/dashboard/chart/delete [post]
func (h *MvDashboardChartHandler) DeleteMvDashboardChart(c *gin.Context) {
	var req model.DeleteMvDashboardChartReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	if err := h.service.DeleteMvDashboardChart(&req); err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, nil, "")
}
