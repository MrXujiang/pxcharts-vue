package handler

import (
	"mvtable/internal/app/mv_view_form/model"
	"mvtable/internal/app/mv_view_form/service"
	"mvtable/internal/middleware"
	"mvtable/internal/pkg/errorx"

	"github.com/gin-gonic/gin"
)

type MvViewFormHandler struct {
	service *service.MvViewFormService
}

func NewMvViewFormHandler(service *service.MvViewFormService) *MvViewFormHandler {
	return &MvViewFormHandler{service: service}
}

// UpdateMvViewForm 更新视图表单
// @Summary 更新视图表单
// @Description 更新MV视图表单信息
// @Tags 视图表单
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.UpdateMvViewFormReq true "更新视图表单请求"
// @Success 200 {object} middleware.Response "更新成功"
// @Router /v1/view/form/update [post]
func (h *MvViewFormHandler) UpdateMvViewForm(c *gin.Context) {
	var req model.UpdateMvViewFormReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	err := h.service.UpdateMvViewForm(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, nil, "")
}
