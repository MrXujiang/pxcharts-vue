package handler

import (
	"mvtable/internal/app/mv_project_advanced_perm/model"
	"mvtable/internal/app/mv_project_advanced_perm/service"
	"mvtable/internal/middleware"
	"mvtable/internal/pkg/errorx"

	"github.com/gin-gonic/gin"
)

type MvProjectAdvancedPermHandler struct {
	service *service.MvProjectAdvancedPermService
}

func NewMvProjectAdvancedPermHandler(service *service.MvProjectAdvancedPermService) *MvProjectAdvancedPermHandler {
	return &MvProjectAdvancedPermHandler{service: service}
}

// EnableMvProjectAdvancedPerm 开启项目高级权限
// @Summary 开启项目高级权限
// @Description 开启项目高级权限
// @Tags 项目高级权限
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.EnableMvProjectAdvancedPermReq true "开启项目高级权限请求"
// @Success 200 {object} middleware.Response "开启成功"
// @Router /v1/project/advanced-perm/enable [post]
func (h *MvProjectAdvancedPermHandler) EnableMvProjectAdvancedPerm(c *gin.Context) {
	var req model.EnableMvProjectAdvancedPermReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	err := h.service.EnableMvProjectAdvancedPermReq(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, nil, "")
}

// DisableMvProjectAdvancedPerm 关闭项目高级权限
// @Summary 关闭项目高级权限
// @Description 关闭项目高级权限
// @Tags 项目高级权限
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.DisableMvProjectAdvancedPermReq true "关闭项目高级权限请求"
// @Success 200 {object} middleware.Response "关闭成功"
// @Router /v1/project/advanced-perm/disable [post]
func (h *MvProjectAdvancedPermHandler) DisableMvProjectAdvancedPerm(c *gin.Context) {
	var req model.DisableMvProjectAdvancedPermReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	err := h.service.DisableMvProjectAdvancedPermReq(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, nil, "")
}

// GetMvProjectAdvancedPerm 获取项目高级权限
// @Summary 获取项目高级权限
// @Description 获取项目高级权限
// @Tags 项目高级权限
// @Accept json
// @Produce json
// @Security Bearer
// @Param request query model.GetMvProjectAdvancedPermReq true "获取项目高级权限请求"
// @Success 200 {object} middleware.Response{data=model.GetMvProjectAdvancedPermRes} "获取成功"
// @Router /v1/project/advanced-perm/get [get]
func (h *MvProjectAdvancedPermHandler) GetMvProjectAdvancedPerm(c *gin.Context) {
	var req model.GetMvProjectAdvancedPermReq
	if err := c.ShouldBindQuery(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	resp, err := h.service.GetMvProjectAdvancedPerm(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, resp, "")
}

// UpdateMvProjectAdvancedPerm 更新项目高级权限
// @Summary 更新项目高级权限
// @Description 更新项目高级权限
// @Tags 项目高级权限
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.UpdateMvProjectAdvancedPermReq true "更新项目高级权限请求"
// @Success 200 {object} middleware.Response "更新成功"
// @Router /v1/project/advanced-perm/update [put]
func (h *MvProjectAdvancedPermHandler) UpdateMvProjectAdvancedPerm(c *gin.Context) {
	var req model.UpdateMvProjectAdvancedPermReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	err := h.service.UpdateMvProjectAdvancedPermReq(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, nil, "")
}
