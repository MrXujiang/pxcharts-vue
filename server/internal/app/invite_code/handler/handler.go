package handler

import (
	"mvtable/internal/app/invite_code/model"
	"mvtable/internal/app/invite_code/service"
	"mvtable/internal/middleware"
	"mvtable/internal/pkg/errorx"

	"github.com/gin-gonic/gin"
)

type InviteCodeHandler struct {
	inviteCodeService *service.InviteCodeService
}

func NewInviteCodeHandler(inviteCodeService *service.InviteCodeService) *InviteCodeHandler {
	return &InviteCodeHandler{inviteCodeService: inviteCodeService}
}

// AdminBatchCreate 批量创建邀请码
// @Summary 批量创建邀请码
// @Description 管理员批量创建邀请码
// @Tags 管理后台
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.BatchCreateInviteCodeReq true "批量创建邀请码请求"
// @Success 200 {object} middleware.Response{data=[]string} "创建成功，data字段为邀请码列表"
// @Router /v1/admin/invite-code/batch-create [post]
func (h *InviteCodeHandler) AdminBatchCreate(c *gin.Context) {
	userId := middleware.GetCurrentUserID(c)

	var req model.BatchCreateInviteCodeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	resp, err := h.inviteCodeService.AdminBatchCreat(userId, &req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, resp, "操作成功")
}

// AdminGetList 获取邀请码列表
// @Summary 获取邀请码列表
// @Description 管理员获取邀请码列表
// @Tags 管理后台
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Param isUsed query bool false "是否已使用"
// @Success 200 {object} middleware.Response{data=model.GetInviteCodeListRes} "获取成功"
// @Router /v1/admin/invite-code/list [get]
func (h *InviteCodeHandler) AdminGetList(c *gin.Context) {
	var req model.GetInviteCodeListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	resp, err := h.inviteCodeService.AdminGetList(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, resp, "操作成功")
}
