package handler

import (
	"mvtable/internal/app/mv_form_submit/model"
	"mvtable/internal/app/mv_form_submit/service"
	"mvtable/internal/middleware"
	"mvtable/internal/pkg/errorx"

	"github.com/gin-gonic/gin"
)

type MvFormSubmitHandler struct {
	service *service.MvFormSubmitService
}

func NewMvFormSubmitHandler(service *service.MvFormSubmitService) *MvFormSubmitHandler {
	return &MvFormSubmitHandler{service: service}
}

// SubmitForm 提交表单
// @Summary 提交表单
// @Description 提交MV视图表单数据
// @Tags 表单提交
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.SubmitFormReq true "提交表单请求"
// @Success 200 {object} middleware.Response "提交成功"
// @Router /v1/form/submit [post]
func (h *MvFormSubmitHandler) SubmitForm(c *gin.Context) {
	var req model.SubmitFormReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	userId := middleware.GetCurrentUserID(c)
	err := h.service.SubmitForm(userId, &req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, nil, "提交成功")
}
