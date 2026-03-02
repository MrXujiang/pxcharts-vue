package handler

import (
	"mvtable/internal/app/team/model"
	"mvtable/internal/middleware"
	"mvtable/internal/pkg/errorx"

	"github.com/gin-gonic/gin"
)

// AdminGetTeamList 获取团队列表
// @Summary 获取团队列表
// @Description 管理员获取所有团队列表
// @Tags 管理后台
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Param searchWord query string false "搜索关键词"
// @Success 200 {object} middleware.Response{data=model.AdminGetTeamListRes} "获取成功"
// @Router /v1/admin/team/list [get]
func (h *TeamHandler) AdminGetTeamList(c *gin.Context) {
	var req model.AdminGetTeamListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	resp, err := h.teamService.AdminGetTeamList(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, resp, "")
}
