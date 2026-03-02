package handler

import (
	"mvtable/internal/app/mv_table_schema/model"
	"mvtable/internal/app/mv_table_schema/service"
	"mvtable/internal/middleware"
	"mvtable/internal/pkg/errorx"

	"github.com/gin-gonic/gin"
)

type MvTableSchemaHandler struct {
	service *service.MvTableSchemaService
}

func NewMvTableSchemaHandler(service *service.MvTableSchemaService) *MvTableSchemaHandler {
	return &MvTableSchemaHandler{service: service}
}

// CreateMvTableSchema 创建表格
// @Summary 创建表格
// @Description 创建新的表格schema
// @Tags 表格
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.CreateMvTableSchemaReq true "创建表格请求"
// @Success 200 {object} middleware.Response{data=string} "创建成功，data字段为表格ID"
// @Router /v1/table/create [post]
func (h *MvTableSchemaHandler) CreateMvTableSchema(c *gin.Context) {
	var req model.CreateMvTableSchemaReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	userId := middleware.GetCurrentUserID(c)
	resp, err := h.service.CreateMvTableSchema(userId, &req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}
	middleware.ResOK(c, resp, "")
}

// UpdateMvTableSchema 更新表格
// @Summary 更新表格
// @Description 更新表格schema信息
// @Tags 表格
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.UpdateMvTableSchemaReq true "更新表格请求"
// @Success 200 {object} middleware.Response "更新成功"
// @Router /v1/table/update [post]
func (h *MvTableSchemaHandler) UpdateMvTableSchema(c *gin.Context) {
	var req model.UpdateMvTableSchemaReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	err := h.service.UpdateMvTableSchema(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}
	middleware.ResOK(c, nil, "")
}

// DeleteMvTableSchema 删除表格
// @Summary 删除表格
// @Description 删除表格schema
// @Tags 表格
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.DeleteMvTableSchemaReq true "删除表格请求"
// @Success 200 {object} middleware.Response "删除成功"
// @Router /v1/table/delete [post]
func (h *MvTableSchemaHandler) DeleteMvTableSchema(c *gin.Context) {
	var req model.DeleteMvTableSchemaReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	err := h.service.DeleteMvTableSchema(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}
	middleware.ResOK(c, nil, "")
}

// CopyMvTableSchema 复制表格
// @Summary 复制表格
// @Description 根据源表ID和复制范围，复制表格结构，支持选择是否复制所有记录
// @Tags 表格
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.CopyMvTableSchemaReq true "复制表格请求"
// @Success 200 {object} middleware.Response{data=model.CopyMvTableSchemaRes} "复制成功"
// @Router /v1/table/copy [post]
func (h *MvTableSchemaHandler) CopyMvTableSchema(c *gin.Context) {
	var req model.CopyMvTableSchemaReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	userId := middleware.GetCurrentUserID(c)
	resp, err := h.service.CopyMvTableSchema(userId, &req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, resp, "")
}

// GetTableData 获取表格数据
// @Summary 获取表格数据
// @Description 获取表格的完整数据，包括字段和记录
// @Tags 表格
// @Accept json
// @Produce json
// @Security Bearer
// @Param request query model.GetTableDataReq true "获取表格数据请求"
// @Success 200 {object} middleware.Response{data=model.GetTableDataRes} "获取成功"
// @Router /v1/table/get [get]
func (h *MvTableSchemaHandler) GetTableData(c *gin.Context) {
	var req model.GetTableDataReq
	if err := c.ShouldBindQuery(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	resp, err := h.service.GetTableData(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}
	middleware.ResOK(c, resp, "")
}
