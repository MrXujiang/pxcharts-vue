package handler

import (
	"mvtable/internal/app/website_visit/model"
	"mvtable/internal/app/website_visit/service"
	"mvtable/internal/middleware"

	"github.com/gin-gonic/gin"
)

// WebsiteVisitHandler 网站访问处理器
type WebsiteVisitHandler struct {
	websiteVisitService *service.WebsiteVisitService
}

// NewWebsiteVisitHandler 创建处理器
func NewWebsiteVisitHandler(websiteVisitService *service.WebsiteVisitService) *WebsiteVisitHandler {
	return &WebsiteVisitHandler{
		websiteVisitService: websiteVisitService,
	}
}

// GetStatistics 获取网站统计数据
// @Summary 获取网站统计数据
// @Description 获取访问量、用户量、团队数、项目数等统计信息
// @Tags 管理后台
// @Accept json
// @Produce json
// @Success 200 {object} middleware.Response{data=object} "获取成功"
// @Router /v1/admin/website/statistics [get]
func (h *WebsiteVisitHandler) GetStatistics(c *gin.Context) {
	response, err := h.websiteVisitService.GetStatistics()
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, response, "获取统计数据成功")
}

// GetNewUsers 获取新增用户数据
// @Summary 获取新增用户数据
// @Description 根据时间范围获取新增用户统计数据，支持按日、周、月、年维度查看
// @Tags 管理后台
// @Accept json
// @Produce json
// @Param timeRange query string true "时间范围(day/week/month/year)" Enums(day,week,month,year)
// @Success 200 {object} middleware.Response{data=model.GetNewUsersRes} "获取成功"
// @Router /v1/admin/website/new-users [get]
func (h *WebsiteVisitHandler) GetNewUsers(c *gin.Context) {
	var req model.GetNewUsersReq
	if err := c.ShouldBindQuery(&req); err != nil {
		middleware.ResErr(c, err)
		return
	}

	response, err := h.websiteVisitService.GetNewUsers(req.TimeRange)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, response, "获取新增用户数据成功")
}
