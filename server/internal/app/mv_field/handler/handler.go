package handler

import (
	collabModel "mvtable/internal/app/collaboration/model"
	collabService "mvtable/internal/app/collaboration/service"
	"mvtable/internal/app/mv_field/model"
	"mvtable/internal/app/mv_field/service"
	"mvtable/internal/middleware"
	"mvtable/internal/pkg/constants"
	"mvtable/internal/pkg/errorx"
	"mvtable/internal/storage/db"
	"time"

	"mvtable/pkg/log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type MvFieldHandler struct {
	service *service.MvFieldService
	hub     *collabService.Hub
}

func NewMvFieldHandler(service *service.MvFieldService, hub *collabService.Hub) *MvFieldHandler {
	return &MvFieldHandler{
		service: service,
		hub:     hub,
	}
}

// GetMvFieldList 获取字段列表
// @Summary 获取字段列表
// @Description 根据表格ID获取字段列表
// @Tags 字段
// @Accept json
// @Produce json
// @Security Bearer
// @Param request query model.GetMvFieldListReq true "获取字段列表请求"
// @Success 200 {object} middleware.Response{data=model.GetMvFieldListRes} "获取成功"
// @Router /v1/field/list [get]
func (h *MvFieldHandler) GetMvFieldList(c *gin.Context) {
	var req model.GetMvFieldListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	resp, err := h.service.GetMvFieldList(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, resp, "")
}

// CreateMvField 创建字段
// @Summary 创建字段
// @Description 创建新的表格字段。必填字段：tableSchemaId（表格ID）、title（字段标题）、type（字段类型）。其他字段会自动放入config中。
// @Tags 字段
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body map[string]any true "创建字段请求"
// @Success 200 {object} middleware.Response{data=string} "创建成功 返回字段"
// @Router /v1/field/create [post]
func (h *MvFieldHandler) CreateMvField(c *gin.Context) {
	var originReq map[string]any
	if err := c.ShouldBindJSON(&originReq); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	config := make(map[string]any)
	for k, v := range originReq {
		if k == "tableSchemaId" || k == "title" || k == "type" {
			continue
		}
		config[k] = v
	}

	req := model.CreateMvFieldReq{
		TableSchemaID: originReq["tableSchemaId"].(string),
		Title:         originReq["title"].(string),
		Type:          constants.FieldType(originReq["type"].(string)),
		Config:        config,
	}

	resp, err := h.service.CreateMvField(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	// 获取创建后的字段信息（包含生成的ID）
	var createdField model.MvField
	if err := db.GetDB().Where("table_schema_id = ? AND title = ?", req.TableSchemaID, req.Title).
		Order("created_at DESC").First(&createdField).Error; err == nil {

		// 广播字段插入消息
		userID := middleware.GetCurrentUserID(c)
		broadcastMsg := &collabModel.CollaborationMessage{
			Type:      collabModel.MsgTypeColumnInsert,
			TableID:   req.TableSchemaID,
			UserID:    userID,
			Timestamp: time.Now().Unix(),
			Data: collabModel.ColumnInsert{
				FieldID: createdField.ID,
				Title:   createdField.Title,
				Type:    createdField.Type,
				Config:  createdField.Config,
			},
		}
		h.hub.BroadcastToRoom(req.TableSchemaID, broadcastMsg, userID)
	} else {
		log.Warn("failed to get created field for broadcast", zap.Error(err))
	}

	middleware.ResOK(c, resp, "")
}

// UpdateMvField 更新字段
// @Summary 更新字段
// @Description 更新表格字段信息。必填字段：id（字段ID）、title（字段标题）、type（字段类型）。其他字段会自动放入config中。
// @Tags 字段
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body map[string]any true "更新字段请求"
// @Success 200 {object} middleware.Response "更新成功"
// @Router /v1/field/update [post]
func (h *MvFieldHandler) UpdateMvField(c *gin.Context) {
	var originReq map[string]any
	if err := c.ShouldBindJSON(&originReq); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	config := make(map[string]any)
	for k, v := range originReq {
		if k == "tableSchemaId" || k == "title" || k == "type" {
			continue
		}
		config[k] = v
	}

	req := model.UpdateMvFieldReq{
		ID:     originReq["id"].(string),
		Title:  originReq["title"].(string),
		Type:   constants.FieldType(originReq["type"].(string)),
		Config: config,
	}

	// 先获取字段信息
	var field model.MvField
	if err := db.GetDB().Where("id = ?", req.ID).First(&field).Error; err != nil {
		middleware.ResErr(c, errorx.InternalServerError("字段不存在"))
		return
	}

	err := h.service.UpdateMvField(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	// 广播字段更新消息
	userID := middleware.GetCurrentUserID(c)
	title := req.Title
	typeStr := req.Type
	broadcastMsg := &collabModel.CollaborationMessage{
		Type:      collabModel.MsgTypeColumnUpdate,
		TableID:   field.TableSchemaID,
		UserID:    userID,
		Timestamp: time.Now().Unix(),
		Data: collabModel.ColumnUpdate{
			FieldID: req.ID,
			Title:   &title,
			Type:    &typeStr,
			Config:  req.Config,
		},
	}
	h.hub.BroadcastToRoom(field.TableSchemaID, broadcastMsg, userID)

	middleware.ResOK(c, nil, "")
}

// SortMvFields 字段排序
// @Summary 字段排序
// @Description 按字段ID列表重排表格字段顺序
// @Tags 字段
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.SortMvFieldReq true "排序请求"
// @Success 200 {object} middleware.Response "排序成功"
// @Router /v1/field/sort [post]
func (h *MvFieldHandler) SortMvFields(c *gin.Context) {
	var req model.SortMvFieldReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	if err := h.service.SortMvFields(&req); err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, nil, "")
}

// DeleteMvField 删除字段
// @Summary 删除字段
// @Description 删除表格字段
// @Tags 字段
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.DeleteMvFieldReq true "删除字段请求"
// @Success 200 {object} middleware.Response "删除成功"
// @Router /v1/field/delete [post]
func (h *MvFieldHandler) DeleteMvField(c *gin.Context) {
	var req model.DeleteMvFieldReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	// 先获取字段信息
	var field model.MvField
	if err := db.GetDB().Where("id = ?", req.ID).First(&field).Error; err != nil {
		middleware.ResErr(c, errorx.InternalServerError("字段不存在"))
		return
	}

	tableSchemaID := field.TableSchemaID

	err := h.service.DeleteMvField(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	// 广播字段删除消息
	userID := middleware.GetCurrentUserID(c)
	broadcastMsg := &collabModel.CollaborationMessage{
		Type:      collabModel.MsgTypeColumnDelete,
		TableID:   tableSchemaID,
		UserID:    userID,
		Timestamp: time.Now().Unix(),
		Data: collabModel.ColumnDelete{
			FieldID: req.ID,
		},
	}
	h.hub.BroadcastToRoom(tableSchemaID, broadcastMsg, userID)

	middleware.ResOK(c, nil, "")
}
