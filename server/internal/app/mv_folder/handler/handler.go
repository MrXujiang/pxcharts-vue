package handler

import (
	"mvtable/internal/app/mv_folder/model"
	"mvtable/internal/app/mv_folder/service"
	"mvtable/internal/middleware"
	"mvtable/internal/pkg/errorx"

	"github.com/gin-gonic/gin"
)

type MvFolderHandler struct {
	service *service.MvFolderService
}

func NewMvFolderHandler(service *service.MvFolderService) *MvFolderHandler {
	return &MvFolderHandler{service: service}
}

// CreateMvFolder 创建文件夹
// @Summary 创建文件夹
// @Description 创建新的文件夹
// @Tags 文件夹
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.CreateMvFolderReq true "创建文件夹请求"
// @Success 200 {object} middleware.Response "创建成功"
// @Router /v1/folder/create [post]
func (h *MvFolderHandler) CreateMvFolder(c *gin.Context) {
	var req model.CreateMvFolderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	err := h.service.CreateMvFolder(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, nil, "")
}

// UpdateMvFolder 更新文件夹
// @Summary 更新文件夹
// @Description 更新文件夹信息
// @Tags 文件夹
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.UpdateMvFolderReq true "更新文件夹请求"
// @Success 200 {object} middleware.Response "更新成功"
// @Router /v1/folder/update [post]
func (h *MvFolderHandler) UpdateMvFolder(c *gin.Context) {
	var req model.UpdateMvFolderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	err := h.service.UpdateMvFolder(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, nil, "")
}

// DeleteMvFolder 删除文件夹
// @Summary 删除文件夹
// @Description 删除文件夹
// @Tags 文件夹
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.DeleteMvFolderReq true "删除文件夹请求"
// @Success 200 {object} middleware.Response "删除成功"
// @Router /v1/folder/delete [post]
func (h *MvFolderHandler) DeleteMvFolder(c *gin.Context) {
	var req model.DeleteMvFolderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	err := h.service.DeleteMvFolder(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, nil, "")
}

// Subquery 查询子节点
// @Summary 查询子节点
// @Description 查询指定父节点下的子节点（文件夹、表格、仪表盘和文档）
// @Tags 文件夹
// @Accept json
// @Produce json
// @Security Bearer
// @Param request query model.QuerySubNodeReq true "查询子节点请求"
// @Success 200 {object} middleware.Response{data=model.QuerySubNodeRes} "查询成功"
// @Router /v1/folder/subquery [get]
func (h *MvFolderHandler) Subquery(c *gin.Context) {
	var req model.QuerySubNodeReq
	if err := c.ShouldBindQuery(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	resp, err := h.service.Subquery(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}
	middleware.ResOK(c, resp, "")
}

// QueryAllNodes 查询项目下所有节点（树状）
// @Summary 查询所有节点
// @Description 查询指定项目下的所有节点（文件夹、表格、仪表盘和文档），以树状结构返回
// @Tags 文件夹
// @Accept json
// @Produce json
// @Security Bearer
// @Param request query model.QueryAllNodeReq true "查询所有节点请求"
// @Success 200 {object} middleware.Response{data=model.QuerySubNodeRes} "查询成功"
// @Router /v1/folder/all-nodes [get]
func (h *MvFolderHandler) QueryAllNodes(c *gin.Context) {
	var req model.QueryAllNodeReq
	if err := c.ShouldBindQuery(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	resp, err := h.service.QueryAllNodes(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}
	middleware.ResOK(c, resp, "")
}

// SearchNodes 搜索文件夹、表格、仪表盘与文档
// @Summary 搜索文件夹、表格、仪表盘与文档
// @Description 根据关键字模糊搜索文件夹、表格、仪表盘与文档，返回树状结构
// @Tags 文件夹
// @Accept json
// @Produce json
// @Security Bearer
// @Param request query model.SearchReq true "搜索请求"
// @Success 200 {object} middleware.Response{data=model.SearchRes} "搜索成功，返回树状结构"
// @Router /v1/folder/search [get]
func (h *MvFolderHandler) Search(c *gin.Context) {
	var req model.SearchReq
	if err := c.ShouldBindQuery(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	resp, err := h.service.Search(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}
	middleware.ResOK(c, resp, "")
}

// ListProjectFolders 查询项目下所有文件夹（扁平列表）
// @Summary 查询项目下所有文件夹（扁平列表）
// @Description 根据项目ID查询该项目下所有文件夹，返回扁平结构列表
// @Tags 文件夹
// @Accept json
// @Produce json
// @Security Bearer
// @Param request query model.QueryFolderListReq true "查询项目文件夹请求"
// @Success 200 {object} middleware.Response{data=model.QueryFolderListRes} "查询成功，返回扁平文件夹列表"
// @Router /v1/folder/list [get]
func (h *MvFolderHandler) ListProjectFolders(c *gin.Context) {
	var req model.QueryFolderListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	resp, err := h.service.ListProjectFolders(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, resp, "")
}

// MoveNode 移动节点到其他文件夹
// @Summary 移动节点
// @Description 将文件夹、数据表、仪表盘或文档移动到目标文件夹，targetFolderId为空字符串表示移动到根目录
// @Tags 文件夹
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.MoveNodeReq true "移动节点请求"
// @Success 200 {object} middleware.Response "移动成功"
// @Router /v1/folder/move [post]
func (h *MvFolderHandler) MoveNode(c *gin.Context) {
	var req model.MoveNodeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	if err := h.service.MoveNode(&req); err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, nil, "")
}
