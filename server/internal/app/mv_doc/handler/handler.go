package handler

import (
	"mvtable/internal/app/mv_doc/model"
	"mvtable/internal/app/mv_doc/service"
	"mvtable/internal/middleware"
	"mvtable/internal/pkg/errorx"

	"github.com/gin-gonic/gin"
)

type MvDocHandler struct {
	service *service.MvDocService
}

func NewMvDocHandler(service *service.MvDocService) *MvDocHandler {
	return &MvDocHandler{service: service}
}

// CreateMvDoc 创建文档
// @Summary 创建文档
// @Description 在项目下创建新的文档
// @Tags 文档
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.CreateMvDocReq true "创建文档请求"
// @Success 200 {object} middleware.Response{data=string} "创建成功，data为文档ID"
// @Router /v1/doc/create [post]
func (h *MvDocHandler) CreateMvDoc(c *gin.Context) {
	var req model.CreateMvDocReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	userId := middleware.GetCurrentUserID(c)
	id, err := h.service.CreateMvDoc(userId, &req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, id, "")
}

// UpdateMvDoc 更新文档
// @Summary 更新文档
// @Description 更新文档信息
// @Tags 文档
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.UpdateMvDocReq true "更新文档请求"
// @Success 200 {object} middleware.Response "更新成功"
// @Router /v1/doc/update [post]
func (h *MvDocHandler) UpdateMvDoc(c *gin.Context) {
	var req model.UpdateMvDocReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	userId := middleware.GetCurrentUserID(c)
	if err := h.service.UpdateMvDoc(userId, &req); err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, nil, "")
}

// DeleteMvDoc 删除文档
// @Summary 删除文档
// @Description 删除指定的文档
// @Tags 文档
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.DeleteMvDocReq true "删除文档请求"
// @Success 200 {object} middleware.Response "删除成功"
// @Router /v1/doc/delete [post]
func (h *MvDocHandler) DeleteMvDoc(c *gin.Context) {
	var req model.DeleteMvDocReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	if err := h.service.DeleteMvDoc(&req); err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, nil, "")
}

// GetMvDoc 获取文档详情
// @Summary 获取文档详情
// @Description 获取文档的详细信息
// @Tags 文档
// @Accept json
// @Produce json
// @Security Bearer
// @Param request query model.GetMvDocReq true "获取文档请求"
// @Success 200 {object} middleware.Response{data=model.GetMvDocRes} "获取成功"
// @Router /v1/doc/get [get]
func (h *MvDocHandler) GetMvDoc(c *gin.Context) {
	var req model.GetMvDocReq
	if err := c.ShouldBindQuery(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	resp, err := h.service.GetMvDoc(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, resp, "")
}

// ListMvDocs 查询文档列表
// @Summary 查询文档列表
// @Description 分页查询文档列表，支持按项目、文件夹筛选
// @Tags 文档
// @Accept json
// @Produce json
// @Security Bearer
// @Param request query model.ListMvDocReq true "查询文档列表请求"
// @Success 200 {object} middleware.Response{data=model.ListMvDocRes} "查询成功"
// @Router /v1/doc/list [get]
func (h *MvDocHandler) ListMvDocs(c *gin.Context) {
	var req model.ListMvDocReq
	if err := c.ShouldBindQuery(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	resp, err := h.service.ListMvDocs(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, resp, "")
}
