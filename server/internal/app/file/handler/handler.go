package handler

import (
	"mime/multipart"
	"mvtable/internal/app/file/model"
	"mvtable/internal/app/file/service"
	"mvtable/internal/middleware"
	"mvtable/internal/pkg/errorx"

	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	fileService *service.FileService
}

func NewFileHandler(fileService *service.FileService) *FileHandler {
	return &FileHandler{
		fileService: fileService,
	}
}

// Upload 文件上传
// @Summary 文件上传
// @Description 上传文件到OSS
// @Tags 文件
// @Accept multipart/form-data
// @Produce json
// @Security Bearer
// @Param file formData file true "文件"
// @Param remark formData string false "备注"
// @Success 200 {object} middleware.Response{data=model.UploadRes} "上传成功"
// @Router /v1/file/upload [post]
func (h *FileHandler) Upload(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	var (
		file   *multipart.FileHeader
		remark string
	)
	if len(form.File["file"]) == 0 {
		middleware.ResErr(c, errorx.BadRequest("file cannot be empty"))
		return
	}

	file = form.File["file"][0]
	if len(form.Value["remark"]) > 0 {
		remark = form.Value["remark"][0]
	}

	userID := middleware.GetCurrentUserID(c)
	resp, err := h.fileService.Upload(userID, file, remark)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, resp, "")
}

// AdminGetList 获取文件列表
// @Summary 获取文件列表
// @Description 管理员获取文件列表
// @Tags 管理后台
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Param filename query string false "文件名筛选"
// @Param oss query string false "OSS类型筛选"
// @Param remark query string false "备注筛选"
// @Success 200 {object} middleware.Response{data=model.AdminGetListRes} "获取成功"
// @Router /v1/admin/file/list [get]
func (h *FileHandler) AdminGetList(c *gin.Context) {
	var req model.AdminGetListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	resp, err := h.fileService.AdminGetList(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, resp, "操作成功")
}
