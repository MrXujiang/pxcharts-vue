package handler

import (
	"mvtable/internal/app/mv_template_tag/model"
	"mvtable/internal/app/mv_template_tag/service"
	"mvtable/internal/middleware"
	"mvtable/internal/pkg/errorx"

	"github.com/gin-gonic/gin"
)

type MvTemplateTagHandler struct {
	service *service.MvTemplateTagService
}

func NewMvTemplateTagHandler(service *service.MvTemplateTagService) *MvTemplateTagHandler {
	return &MvTemplateTagHandler{service: service}
}

// QueryMvTemplateTag 查询模板标签列表
// @Summary 查询模板标签列表
// @Description 获取模板标签列表，支持搜索
// @Tags 项目模版
// @Accept json
// @Produce json
// @Security Bearer
// @Param searchWord query string false "搜索关键词"
// @Success 200 {object} middleware.Response{data=model.QueryMvTemplateTagRes} "查询成功"
// @Router /v1/template-tag/query [get]
func (h *MvTemplateTagHandler) QueryMvTemplateTag(c *gin.Context) {
	var req model.QueryMvTemplateTagReq
	if err := c.ShouldBindQuery(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	resp, err := h.service.QueryMvTemplateTag(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, resp, "")
}

// CreateMvTemplateTag 创建模板标签
// @Summary 创建模板标签
// @Description 创建新的模板标签
// @Tags 管理后台
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.CreateMvTemplateTagReq true "创建模板标签请求"
// @Success 200 {object} middleware.Response "创建成功"
// @Router /v1/admin/template-tag/create [post]
func (h *MvTemplateTagHandler) CreateMvTemplateTag(c *gin.Context) {
	var req model.CreateMvTemplateTagReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	err := h.service.CreateMvTemplateTag(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, nil, "")
}

// UpdateMvTemplateTag 更新模板标签
// @Summary 更新模板标签
// @Description 更新模板标签信息
// @Tags 管理后台
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.UpdateMvTemplateTagReq true "更新模板标签请求"
// @Success 200 {object} middleware.Response "更新成功"
// @Router /v1/admin/template-tag/update [post]
func (h *MvTemplateTagHandler) UpdateMvTemplateTag(c *gin.Context) {
	var req model.UpdateMvTemplateTagReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	err := h.service.UpdateMvTemplateTag(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, nil, "")
}

// DeleteMvTemplateTag 删除模板标签
// @Summary 删除模板标签
// @Description 删除模板标签（软删除）
// @Tags 管理后台
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.DeleteMvTemplateTagReq true "删除模板标签请求"
// @Success 200 {object} middleware.Response "删除成功"
// @Router /v1/admin/template-tag/delete [post]
func (h *MvTemplateTagHandler) DeleteMvTemplateTag(c *gin.Context) {
	var req model.DeleteMvTemplateTagReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	err := h.service.DeleteMvTemplateTag(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, nil, "")
}

// UpdateTagSort 调整标签排序
// @Summary 调整标签排序
// @Description 更新模板标签的展示顺序，根据前后标签的位置计算新的排序索引
// @Tags 管理后台
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.UpdateTagSortReq true "调整排序请求"
// @Success 200 {object} middleware.Response "更新成功"
// @Router /v1/admin/template-tag/sort [post]
func (h *MvTemplateTagHandler) UpdateTagSort(c *gin.Context) {
	var req model.UpdateTagSortReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	err := h.service.UpdateTagSort(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, nil, "")
}
