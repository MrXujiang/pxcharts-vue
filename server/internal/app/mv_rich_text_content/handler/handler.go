package handler

import (
	"mvtable/internal/app/mv_rich_text_content/model"
	"mvtable/internal/app/mv_rich_text_content/service"
	"mvtable/internal/middleware"
	"mvtable/internal/pkg/errorx"

	"github.com/gin-gonic/gin"
)

type MvRichTextContentHandler struct {
	service *service.MvRichTextContentService
}

func NewMvRichTextContentHandler(service *service.MvRichTextContentService) *MvRichTextContentHandler {
	return &MvRichTextContentHandler{service: service}
}

// GetRichTextContent 获取富文本内容
// @Summary 获取富文本内容
// @Description 根据记录ID和字段ID获取富文本内容
// @Tags 记录
// @Accept json
// @Produce json
// @Security Bearer
// @Param recordId query string true "记录ID"
// @Param fieldId query string true "字段ID"
// @Success 200 {object} middleware.Response{data=model.GetRichTextContentRes} "获取成功"
// @Router /v1/rich-text/get [get]
func (h *MvRichTextContentHandler) GetRichTextContent(c *gin.Context) {
	var req model.GetRichTextContentReq
	if err := c.ShouldBindQuery(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	content, err := h.service.GetRichTextContent(req.RecordID, req.FieldID)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	res := &model.GetRichTextContentRes{
		Content: content,
	}
	middleware.ResOK(c, res, "")
}
