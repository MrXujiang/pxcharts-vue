package handler

import (
	"mvtable/internal/app/mv_view/model"
	"mvtable/internal/app/mv_view/service"
	"mvtable/internal/middleware"
	"mvtable/internal/pkg/errorx"

	"github.com/gin-gonic/gin"
)

type MvViewHandler struct {
	service *service.MvViewFormService
}

func NewMvViewHandler(service *service.MvViewFormService) *MvViewHandler {
	return &MvViewHandler{service: service}
}

// QueryMvView 查询视图列表
// @Summary 查询视图列表
// @Description 根据表格ID查询视图列表
// @Tags 视图
// @Accept json
// @Produce json
// @Security Bearer
// @Param request query model.QueryMvViewReq true "查询视图请求"
// @Success 200 {object} middleware.Response{data=model.QueryMvViewRes} "查询成功"
// @Router /v1/view/query [get]
func (h *MvViewHandler) QueryMvView(c *gin.Context) {
	var req model.QueryMvViewReq
	if err := c.ShouldBindQuery(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	userId := middleware.GetCurrentUserID(c)
	resp, err := h.service.QueryMvView(userId, &req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, resp, "")
}

// CreateMvView 创建视图
// @Summary 创建视图
// @Description 创建新的MV视图
// @Tags 视图
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.CreateMvViewReq true "创建视图请求"
// @Success 200 {object} middleware.Response{data=string} "创建成功"
// @Router /v1/view/create [post]
func (h *MvViewHandler) CreateMvView(c *gin.Context) {
	var req model.CreateMvViewReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	resp, err := h.service.CreateMvView(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, resp, "")
}

// UpdateMvView 更新视图
// @Summary 更新视图
// @Description 更新MV视图信息
// @Tags 视图
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.UpdateMvViewReq true "更新视图请求"
// @Success 200 {object} middleware.Response "更新成功"
// @Router /v1/view/update [post]
func (h *MvViewHandler) UpdateMvView(c *gin.Context) {
	var req model.UpdateMvViewReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	err := h.service.UpdateMvView(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, nil, "")
}

// DeleteMvView 删除视图
// @Summary 删除视图
// @Description 删除MV视图
// @Tags 视图
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.DeleteMvViewReq true "删除视图请求"
// @Success 200 {object} middleware.Response "删除成功"
// @Router /v1/view/delete [post]
func (h *MvViewHandler) DeleteMvView(c *gin.Context) {
	var req model.DeleteMvViewReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	err := h.service.DeleteMvView(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, nil, "")
}

// GetMvView 获取视图详情
// @Summary 获取视图详情
// @Description 根据视图ID获取视图详情信息
// @Tags 视图
// @Accept json
// @Produce json
// @Security Bearer
// @Param request query model.GetMvViewReq true "获取视图请求"
// @Success 200 {object} middleware.Response{data=model.GetMvViewRes} "获取成功"
// @Router /v1/view/get [get]
func (h *MvViewHandler) GetMvView(c *gin.Context) {
	var req model.GetMvViewReq
	if err := c.ShouldBindQuery(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	resp, err := h.service.GetMvView(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, resp, "")
}

// SwitchActiveView 切换激活视图
// @Summary 切换激活视图
// @Description 切换指定表格的激活视图
// @Tags 视图
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.SwitchActiveViewReq true "切换激活视图请求"
// @Success 200 {object} middleware.Response "切换成功"
// @Router /v1/view/switch-active [post]
func (h *MvViewHandler) SwitchActiveView(c *gin.Context) {
	var req model.SwitchActiveViewReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	userId := middleware.GetCurrentUserID(c)
	err := h.service.SwitchActiveView(userId, &req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, nil, "")
}
