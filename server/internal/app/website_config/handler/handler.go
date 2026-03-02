package handler

import (
	"mvtable/internal/app/website_config/model"
	"mvtable/internal/app/website_config/service"
	"mvtable/internal/middleware"

	"github.com/gin-gonic/gin"
)

type WebsiteConfigHandler struct {
	websiteConfigService *service.WebsiteConfigService
}

func NewWebsiteConfigHandler() *WebsiteConfigHandler {
	return &WebsiteConfigHandler{
		websiteConfigService: service.NewWebsiteConfigService(),
	}
}

// Get 获取网站配置
// @Summary 获取网站配置
// @Description 获取网站配置信息
// @Tags 网站配置
// @Accept json
// @Produce json
// @Success 200 {object} middleware.Response{data=model.WebsiteConfigSchema} "获取成功"
// @Router /v1/website-config [get]
func (h *WebsiteConfigHandler) Get(c *gin.Context) {
	resp, err := h.websiteConfigService.Get()
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, resp, "")
}

// AdminSave 保存网站配置
// @Summary 保存网站配置
// @Description 管理员保存网站配置
// @Tags 管理后台
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.WebsiteConfigSchema true "网站配置"
// @Success 200 {object} middleware.Response "保存成功"
// @Router /v1/admin/website-config/save [post]
func (h *WebsiteConfigHandler) AdminSave(c *gin.Context) {
	var req model.WebsiteConfigSchema
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, err)
		return
	}

	if err := h.websiteConfigService.AdminSave(&req); err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, nil, "保存成功")
}
