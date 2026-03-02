package handler

import (
	"mvtable/internal/app/mv_view_table/model"
	"mvtable/internal/app/mv_view_table/service"
	"mvtable/internal/middleware"
	"mvtable/internal/pkg/errorx"

	"github.com/gin-gonic/gin"
)

type MvViewTableHandler struct {
	service *service.MvViewTableService
}

func NewMvViewTableHandler(service *service.MvViewTableService) *MvViewTableHandler {
	return &MvViewTableHandler{service: service}
}

// UpdateMvViewTable 更新视图表格
// @Summary 更新视图表格
// @Description 更新MV视图表格信息
// @Description filterConfig: [{fieldId: "字段id", operator: "equal=等于/notEqual=不等于/contains=包含/notContains=不包含/null=为空/notNull=不为空/before=早于/after=晚于/beforeOrEqual=早于或等于/afterOrEqual=晚于或等于", value: "值"}]
// @Description groupConfig: {fieldId: "", order: "asc/desc"}
// @Description sortConfig: [{fieldId: "", order: "asc/desc"}]
// @Description rowHeight: 1=常规, 2=中等, 3=高, 4=超高
// @Description colorConfig: [{color: "#cccccc", scope: "cell/row/rowHeader/column", fieldId: "字段id", operator: "equal=等于/notEqual=不等于/greaterThan=大于/lessThan=小于/greaterThanOrEqual=大于等于/lessThanOrEqual=小于等于/contains=包含/notContains=不包含/null=为空/notNull=不为空/before=早于/after=晚于/beforeOrEqual=早于或等于/afterOrEqual=晚于或等于", value: "值"}]
// @Tags 视图表格
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.UpdateMvViewTableReq true "更新视图表格请求"
// @Success 200 {object} middleware.Response "更新成功"
// @Router /v1/view/table/update [post]
func (h *MvViewTableHandler) UpdateMvViewTable(c *gin.Context) {
	var req model.UpdateMvViewTableReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	err := h.service.UpdateMvViewTable(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, nil, "")
}
