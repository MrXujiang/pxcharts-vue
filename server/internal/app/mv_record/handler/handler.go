package handler

import (
	"mvtable/internal/app/mv_record/model"
	"mvtable/internal/app/mv_record/service"
	"mvtable/internal/middleware"
	"mvtable/internal/pkg/errorx"

	"github.com/gin-gonic/gin"
)

type MvRecordHandler struct {
	service *service.MvRecordService
}

func NewMvRecordHandler(service *service.MvRecordService) *MvRecordHandler {
	return &MvRecordHandler{service: service}
}

// InsertMvRecord 插入记录
// @Summary 插入记录
// @Description 在指定表格中插入新的记录行，支持指定在前一条记录之后插入
// @Tags 记录
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.InsertMvRecordReq true "插入记录请求"
// @Success 200 {object} middleware.Response "插入成功"
// @Router /v1/record/insert [post]
func (h *MvRecordHandler) InsertMvRecord(c *gin.Context) {
	var req model.InsertMvRecordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	userId := middleware.GetCurrentUserID(c)
	resp, err := h.service.InsertMvRecord(userId, &req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}
	middleware.ResOK(c, resp, "")
}

// GetRecord 获取单条记录
// @Summary 获取单条记录
// @Description 根据表格ID和记录ID获取单条完整记录数据
// @Tags 记录
// @Accept json
// @Produce json
// @Security Bearer
// @Param tableSchemaId query string true "表格ID"
// @Param recordId query string true "记录ID"
// @Success 200 {object} middleware.Response "获取成功"
// @Router /v1/record/get [get]
func (h *MvRecordHandler) GetRecord(c *gin.Context) {
	var req model.GetRecordReq
	if err := c.ShouldBindQuery(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	resp, err := h.service.GetRecord(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}
	middleware.ResOK(c, resp, "")
}

// GetRecords 获取表格所有记录
// @Summary 获取表格所有记录（不进行筛选，只返回记录列表）
// @Description 根据表格ID获取该表格的所有完整记录数据
// @Tags 记录
// @Accept json
// @Produce json
// @Security Bearer
// @Param tableSchemaId query string true "表格ID"
// @Success 200 {object} middleware.Response "获取成功"
// @Router /v1/record/list [get]
func (h *MvRecordHandler) GetRecords(c *gin.Context) {
	var req model.GetRecordsReq
	if err := c.ShouldBindQuery(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	resp, err := h.service.GetRecords(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}
	middleware.ResOK(c, resp, "")
}

// GetRecordOptions 获取表格记录选项列表
// @Summary 获取表格记录选项列表
// @Description 根据表格ID获取该表格的所有记录，返回id和第一个字段值的选项列表
// @Tags 记录
// @Accept json
// @Produce json
// @Security Bearer
// @Param tableSchemaId query string true "表格ID"
// @Success 200 {object} middleware.Response{data=model.GetRecordOptionsRes} "获取成功"
// @Router /v1/record/options [get]
func (h *MvRecordHandler) GetRecordOptions(c *gin.Context) {
	var req model.GetRecordOptionsReq
	if err := c.ShouldBindQuery(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	resp, err := h.service.GetRecordOptions(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}
	middleware.ResOK(c, resp, "")
}

// SaveRowSort 保存行排序
// @Summary 保存行排序
// @Description 更新表格记录行的排序顺序，根据前后记录的位置计算新的排序索引
// @Tags 记录
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.SaveRowSortReq true "保存排序请求"
// @Success 200 {object} middleware.Response "保存成功"
// @Router /v1/record/sort [post]
func (h *MvRecordHandler) SaveRowSort(c *gin.Context) {
	var req model.SaveRowSortReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	userId := middleware.GetCurrentUserID(c)
	err := h.service.SaveRowSort(userId, &req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}
	middleware.ResOK(c, nil, "")
}

// UpdateCell 更新单元格
// @Summary 更新单元格
// @Description 更新表格记录中指定字段的单元格值
// @Tags 记录
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.UpdateCellReq true "更新单元格请求"
// @Success 200 {object} middleware.Response "更新成功"
// @Router /v1/record/cell/update [post]
func (h *MvRecordHandler) UpdateCell(c *gin.Context) {
	var req model.UpdateCellReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	userId := middleware.GetCurrentUserID(c)

	err := h.service.UpdateCell(userId, &req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}
	middleware.ResOK(c, nil, "")
}

// UpdateRow 更新整条记录
// @Summary 更新整条记录
// @Description 更新表格记录中多个字段的值 {"rowId": "", "field1Id": "", "field2Id": ""}
// @Tags 记录
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.UpdateRowReq true "更新记录请求"
// @Success 200 {object} middleware.Response "更新成功"
// @Router /v1/record/row/update [post]
func (h *MvRecordHandler) UpdateRow(c *gin.Context) {
	var req model.UpdateRowReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	userId := middleware.GetCurrentUserID(c)

	err := h.service.UpdateRow(userId, req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}
	middleware.ResOK(c, nil, "")
}

// DeleteRecords 删除多行记录
// @Summary 删除多行记录
// @Description 批量删除指定的记录行
// @Tags 记录
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.DeleteRecordsReq true "删除记录请求"
// @Success 200 {object} middleware.Response "删除成功"
// @Router /v1/record/delete [post]
func (h *MvRecordHandler) DeleteRecords(c *gin.Context) {
	var req model.DeleteRecordsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	userId := middleware.GetCurrentUserID(c)
	err := h.service.DeleteRecords(userId, &req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, nil, "")
}

// ImportRecords 导入记录
// @Summary 导入记录
// @Description 从 CSV 或 Excel 文件导入数据到表格，文件第一行为标题，根据标题名称匹配字段
// @Tags 记录
// @Accept multipart/form-data
// @Produce json
// @Security Bearer
// @Param tableSchemaId formData string true "表格ID"
// @Param file formData file true "导入文件（CSV 或 Excel）"
// @Success 200 {object} middleware.Response{data=model.ImportRecordsRes} "导入成功"
// @Router /v1/record/import [post]
func (h *MvRecordHandler) ImportRecords(c *gin.Context) {
	tableSchemaID := c.PostForm("tableSchemaId")
	if tableSchemaID == "" {
		middleware.ResErr(c, errorx.BadRequest("tableSchemaId 不能为空"))
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		middleware.ResErr(c, errorx.BadRequest("请选择要导入的文件"))
		return
	}

	userId := middleware.GetCurrentUserID(c)
	res, err := h.service.ImportRecords(userId, tableSchemaID, file)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, res, "")
}
