package handler

import (
	"mvtable/internal/app/mv_view_board/model"
	"mvtable/internal/app/mv_view_board/service"
	"mvtable/internal/middleware"
	"mvtable/internal/pkg/errorx"

	"github.com/gin-gonic/gin"
)

type MvViewBoardHandler struct {
	service *service.MvViewBoardService
}

func NewMvViewBoardHandler(service *service.MvViewBoardService) *MvViewBoardHandler {
	return &MvViewBoardHandler{service: service}
}

// UpdateMvViewBoard 更新视图看板
// @Summary 更新视图看板
// @Description 更新MV视图看板信息
// @Description showFieldTitle: 是否显示字段标题
// @Description filterConfig: [{fieldId: "字段id", operator: "equal=等于/notEqual=不等于/contains=包含/notContains=不包含/null=为空/notNull=不为空/before=早于/after=晚于/beforeOrEqual=早于或等于/afterOrEqual=晚于或等于", value: "值"}]
// @Description groupConfig: {fieldId: ""}
// @Description sortConfig: [{fieldId: "", order: "asc/desc"}]
// @Tags 视图看板
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.UpdateMvViewBoardReq true "更新视图看板请求"
// @Success 200 {object} middleware.Response "更新成功"
// @Router /v1/view/board/update [post]
func (h *MvViewBoardHandler) UpdateMvViewBoard(c *gin.Context) {
	var req model.UpdateMvViewBoardReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	err := h.service.UpdateMvViewBoard(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, nil, "")
}
