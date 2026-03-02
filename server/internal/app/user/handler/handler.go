package handler

import (
	"mvtable/internal/app/user/model"
	"mvtable/internal/app/user/service"
	"mvtable/internal/pkg/errorx"

	"github.com/gin-gonic/gin"

	"mvtable/internal/middleware"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// AdminLogin 管理员登录
// @Summary 管理员登录
// @Description 管理员登录接口
// @Tags 管理后台
// @Accept json
// @Produce json
// @Param request body model.AdminLoginReq true "登录请求"
// @Success 200 {object} middleware.Response{data=object} "登录成功，data字段包含token和refreshToken"
// @Router /v1/admin/login [post]
func (h *UserHandler) AdminLogin(c *gin.Context) {
	var req model.AdminLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	response, err := h.userService.AdminLogin(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, response, "登录成功")
}

// AdminGetUserList 获取用户列表
// @Summary 获取用户列表
// @Description 管理员获取用户列表
// @Tags 管理后台
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Param email query string false "邮箱筛选"
// @Param nickname query string false "昵称筛选"
// @Param status query int false "状态筛选"
// @Success 200 {object} middleware.Response{data=model.AdminGetUserListRes} "获取成功"
// @Router /v1/admin/user/list [get]
func (h *UserHandler) AdminGetUserList(c *gin.Context) {
	var req model.AdminGetUserListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	resp, err := h.userService.AdminGetUserList(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, resp, "操作成功")
}

// AdminUpdateInfo 更新用户信息
// @Summary 更新用户信息
// @Description 管理员更新用户信息
// @Tags 管理后台
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.AdminUpdateInfoReq true "更新请求"
// @Success 200 {object} middleware.Response "更新成功"
// @Router /v1/admin/user/info/update [post]
func (h *UserHandler) AdminUpdateInfo(c *gin.Context) {
	var req model.AdminUpdateInfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	err := h.userService.AdminUpdateInfo(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, nil, "更新成功")
}

// AdminDeleteUser 删除用户
// @Summary 删除用户
// @Description 管理员删除用户
// @Tags 管理后台
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.AdminDeleteUserReq true "删除请求"
// @Success 200 {object} middleware.Response "删除成功"
// @Router /v1/admin/user/delete [post]
func (h *UserHandler) AdminDeleteUser(c *gin.Context) {
	var req model.AdminDeleteUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	err := h.userService.AdminDeleteUser(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, nil, "删除成功")
}

// EmailRegister 邮箱注册
// @Summary 邮箱注册
// @Description 通过邮箱验证码注册
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body model.RegisterReq true "注册请求"
// @Success 200 {object} middleware.Response{data=model.RegisterRes} "注册成功"
// @Router /v1/auth/register [post]
func (h *UserHandler) EmailRegister(c *gin.Context) {
	var req model.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	resp, err := h.userService.Register(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, resp, "")
}

// InviteCodeRegister 邀请码注册
// @Summary 邀请码注册
// @Description 通过邀请码注册
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body model.InviteCodeRegisterReq true "注册请求"
// @Success 200 {object} middleware.Response{data=model.RegisterRes} "注册成功"
// @Router /v1/auth/invite-code/register [post]
func (h *UserHandler) InviteCodeRegister(c *gin.Context) {
	var req model.InviteCodeRegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	resp, err := h.userService.InviteCodeRegister(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, resp, "")
}

// SendEmailVerifyCode 发送邮箱验证码
// @Summary 发送邮箱验证码
// @Description 发送邮箱验证码用于注册
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body model.SendEmailVerifyCodeReq true "发送验证码请求"
// @Success 200 {object} middleware.Response "发送成功"
// @Router /v1/auth/register/send-code [post]
func (h *UserHandler) SendEmailVerifyCode(c *gin.Context) {
	var req model.SendEmailVerifyCodeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	err := h.userService.SendEmailVerifyCode(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, nil, "发送验证码成功")
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录接口
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body model.LoginReq true "登录请求"
// @Success 200 {object} middleware.Response{data=model.LoginRes} "登录成功"
// @Router /v1/auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req model.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	response, err := h.userService.Login(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	// 设置 token 到 cookie
	// Cookie 过期时间设置为 24 小时（与 access token 过期时间一致）
	c.SetCookie(
		"token",        // cookie 名称
		response.Token, // cookie 值
		24*60*60,       // 过期时间（秒），24小时
		"/",            // 路径
		"",             // 域名（空字符串表示当前域名）
		false,          // Secure（生产环境建议设置为 true）
		true,           // HttpOnly（防止 XSS 攻击）
	)

	middleware.ResOK(c, response, "登录成功")
}

// RefreshToken 刷新令牌
// @Summary 刷新令牌
// @Description 刷新访问令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Success 200 {object} middleware.Response "刷新成功"
// @Router /v1/auth/refresh [post]
func (h *UserHandler) RefreshToken(c *gin.Context) {
	// 这里可以添加刷新令牌逻辑
	middleware.ResOK(c, nil, "令牌刷新成功")
}

// GetUser 获取用户信息
// @Summary 获取用户信息
// @Description 获取当前登录用户信息
// @Tags 用户
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} middleware.Response{data=model.UserResponse} "获取成功"
// @Router /v1/user/info [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	user, err := h.userService.GetUser(userID)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, user, "")
}

// UpdateInfo 更新用户信息
// @Summary 更新用户信息
// @Description 更新当前登录用户信息
// @Tags 用户
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.UpdateInfoReq true "更新请求"
// @Success 200 {object} middleware.Response{data=model.UserResponse} "更新成功"
// @Router /v1/user/update [post]
func (h *UserHandler) UpdateInfo(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	var req model.UpdateInfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	user, err := h.userService.UpdateInfo(userID, &req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, user, "用户信息更新成功")
}

// ResetPassword 重置密码
// @Summary 重置密码
// @Description 重置当前用户密码
// @Tags 用户
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.ResetPasswordReq true "重置密码请求"
// @Success 200 {object} middleware.Response "重置成功"
// @Router /v1/user/password/reset [post]
func (h *UserHandler) ResetPassword(c *gin.Context) {
	var req model.ResetPasswordReq

	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResErr(c, errorx.BadRequest(err.Error()))
		return
	}

	err := h.userService.ResetPassword(&req)
	if err != nil {
		middleware.ResErr(c, err)
		return
	}

	middleware.ResOK(c, nil, "密码修改成功")
}
