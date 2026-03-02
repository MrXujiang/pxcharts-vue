package handler

import (
	"mvtable/internal/app/mv_project_state/model"
	"mvtable/internal/app/mv_project_state/service"
	"mvtable/internal/middleware"
	"mvtable/internal/pkg/errorx"

	"github.com/gin-gonic/gin"
)

type MvProjectStateHandler struct {
	service *service.MvProjectStateService
}

func NewMvProjectStateHandler(service *service.MvProjectStateService) *MvProjectStateHandler {
	return &MvProjectStateHandler{service: service}
}

// SetShareRange 设置分享范围
// @Summary 设置分享范围
// @Description 设置项目分享范围
// @Tags 项目
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.SetShareRangeReq true "设置分享范围请求"
// @Success 200 {object} middleware.Response "设置成功"
// @Router /v1/project/state/set [post]
func (h *MvProjectStateHandler) SetShareRange(c *gin.Context) {
	var req model.SetShareRangeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	userId := middleware.GetCurrentUserID(c)
	err := h.service.SetShareRange(userId, &req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, nil, "")
}
