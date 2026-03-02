package service

import (
	"context"
	"errors"
	"fmt"
	model2 "mvtable/internal/app/invite_code/model"
	"mvtable/internal/app/user/model"
	"mvtable/internal/pkg/constants"
	"mvtable/internal/pkg/errorx"
	"mvtable/internal/pkg/utils"
	"mvtable/internal/storage/db"
	"mvtable/internal/storage/redis"
	"mvtable/pkg/mail"
	"strings"
	"time"

	"gorm.io/gorm"

	"golang.org/x/crypto/bcrypt"

	"mvtable/pkg/log"

	"go.uber.org/zap"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

// AdminUpdateInfo 管理员更新用户资料
func (s *UserService) AdminUpdateInfo(req *model.AdminUpdateInfoReq) error {
	var (
		user         *model.User
		updateFields = make([]string, 0)
		err          error
	)
	err = db.Transaction(func(tx *gorm.DB) error {
		user, err = db.Get[model.User](tx, map[string]any{"id": req.ID})
		if err != nil {
			log.Error("query user error", zap.Error(err))
			return errorx.InternalServerError("更新失败")
		}
		if user == nil {
			return errorx.New(errorx.ErrNotFound, "用户不存在")
		}

		if req.Nickname != nil && *req.Nickname != "" {
			user.Nickname = *req.Nickname
			updateFields = append(updateFields, "nickname")
		}
		if req.Avatar != nil && *req.Avatar != "" {
			user.Avatar = *req.Avatar
			updateFields = append(updateFields, "avatar")
		}
		if req.Tags != nil && len(*req.Tags) > 0 {
			user.Tags = strings.Join(*req.Tags, ",")
			updateFields = append(updateFields, "tags")
		}
		if req.Status != nil && *req.Status != 0 {
			user.Status = *req.Status
			updateFields = append(updateFields, "status")
		}
		if req.Role != nil && *req.Role != "" {
			user.Role = *req.Role
			updateFields = append(updateFields, "role")
		}
		if req.Identity != nil && *req.Identity != 0 {
			user.Identity = *req.Identity
			user.MaxProjectCount = req.Identity.MaxProjectCount()
			user.MaxTeamCount = req.Identity.MaxTeamCount()
			updateFields = append(updateFields, "identity", "max_project_count", "max_team_count")
		}

		if err = db.Update(tx, user, map[string]any{"id": req.ID}, updateFields...); err != nil {
			log.Error("update user error", zap.Error(err))
			return errorx.InternalServerError("更新失败")
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// AdminDeleteUser 管理员删除用户
func (s *UserService) AdminDeleteUser(req *model.AdminDeleteUserReq) error {
	var (
		user *model.User
		err  error
	)
	err = db.Transaction(func(tx *gorm.DB) error {
		user, err = db.Get[model.User](tx, map[string]any{"id": req.ID})
		if err != nil {
			log.Error("query user error", zap.Error(err))
			return errorx.InternalServerError("删除失败")
		}
		if user == nil {
			return errorx.New(errorx.ErrNotFound, "用户不存在")
		}

		// 删除用户
		if err = db.Delete[model.User](tx, map[string]any{"id": req.ID}); err != nil {
			log.Error("delete user error", zap.Error(err))
			return errorx.InternalServerError("删除失败")
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// Register 用户注册
func (s *UserService) Register(req *model.RegisterReq) (*model.RegisterRes, error) {
	var (
		user           *model.User
		token          string
		refreshToken   string
		err            error
		hashedPassword []byte
	)
	err = db.Transaction(func(tx *gorm.DB) error {
		if user, err = db.Get[model.User](tx, map[string]any{"email": req.Email}); err != nil {
			log.Error("query user by email error", zap.Error(err))
			return errorx.InternalServerError("注册失败")
		}
		if user != nil {
			return errorx.New(errorx.ErrRegister, "邮箱已存在")
		}

		// 校验邮箱验证码
		if err = judgeVerifyCode(req.Email, req.VerifyCode); err != nil {
			log.Error("verify code error", zap.Error(err))
			return err
		}

		// 加密密码
		hashedPassword, err = bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Error("encrypt password error", zap.Error(err))
			return errorx.InternalServerError("注册失败")
		}

		// 创建用户
		user = &model.User{
			Email:           req.Email,
			Password:        string(hashedPassword),
			Nickname:        genNickname(req.Email),
			Identity:        constants.UserIdentityBasic,
			MaxProjectCount: constants.UserIdentityBasic.MaxProjectCount(),
			MaxTeamCount:    constants.UserIdentityBasic.MaxTeamCount(),
		}

		if err = db.Create(tx, user); err != nil {
			log.Error("create user error", zap.Error(err))
			return errorx.InternalServerError("注册失败")
		}

		if err = createTeam(tx, user.ID); err != nil {
			return err
		}

		token, refreshToken, err = genTokens(user.ID)
		if err != nil {
			log.Error("generate token error", zap.Error(err))
			return errorx.InternalServerError("注册失败")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &model.RegisterRes{
		Token:        token,
		RefreshToken: refreshToken,
		UserInfo:     user.ToResponse(),
	}, nil
}

func (s *UserService) InviteCodeRegister(req *model.InviteCodeRegisterReq) (*model.RegisterRes, error) {
	var (
		user           *model.User
		inviteCode     *model2.InviteCode
		token          string
		refreshToken   string
		err            error
		hashedPassword []byte
	)

	err = db.Transaction(func(tx *gorm.DB) error {
		if user, err = db.Get[model.User](tx, map[string]any{"email": req.Email}); err != nil {
			log.Error("query user by email error", zap.Error(err))
			return errorx.InternalServerError("注册失败")
		}

		if user != nil {
			return errorx.New(errorx.ErrRegister, "邮箱已存在")
		}

		// 判断邀请码是否有效
		if inviteCode, err = db.Get[model2.InviteCode](tx, map[string]any{"value": req.InviteCode}); err != nil {
			log.Error("get invite code error", zap.Error(err))
			return errorx.InternalServerError("注册失败")
		}

		if inviteCode == nil || inviteCode.IsUsed {
			return errorx.New(errorx.ErrOperationFailed, "邀请码无效")
		}

		// 加密密码
		hashedPassword, err = bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Error("encrypt password error", zap.Error(err))
			return errorx.InternalServerError("注册失败")
		}

		// 创建用户
		user = &model.User{
			Email:           req.Email,
			Password:        string(hashedPassword),
			Nickname:        genNickname(req.Email),
			Identity:        constants.UserIdentityBasic,
			MaxProjectCount: constants.UserIdentityBasic.MaxProjectCount(),
			MaxTeamCount:    constants.UserIdentityBasic.MaxTeamCount(),
		}
		if err = db.Create(tx, user); err != nil {
			log.Error("create user error", zap.Error(err))
			return errorx.InternalServerError("注册失败")
		}

		if err = createTeam(tx, user.ID); err != nil {
			return err
		}

		token, refreshToken, err = genTokens(user.ID)
		if err != nil {
			log.Error("generate token error", zap.Error(err))
			return errorx.InternalServerError("注册失败")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &model.RegisterRes{
		Token:        token,
		RefreshToken: refreshToken,
		UserInfo:     user.ToResponse(),
	}, nil
}

// SendEmailVerifyCode 发送注册邮件
func (s *UserService) SendEmailVerifyCode(req *model.SendEmailVerifyCodeReq) error {
	if !utils.ValidateEmail(req.Email) {
		return errorx.New(errorx.ErrInvalidEmail, "邮箱格式不正确")
	}

	key := constants.RedisRegisterVerifyCodePrefix + req.Email

	// 判断是否已经发送过验证码
	cacheCode, err := redis.Get(context.Background(), key)
	if err != nil && !errors.Is(err, redis.Nil) {
		log.Error("get verify code error", zap.Error(err))
		return errorx.New(errorx.ErrOperationFailed, "发送验证码失败")
	}

	if cacheCode != "" {
		return errorx.New(errorx.ErrOperationFailed, "验证码已发送")
	}

	// 生成验证码
	code := utils.GenRandomNumber(6)

	// 发送邮件
	if err = mail.SendPlanText(req.Email, "邮箱验证码", fmt.Sprintf("您的邮箱验证码是：%s", code)); err != nil {
		log.Error("send register email error", zap.Error(err))
		return errorx.New(errorx.ErrOperationFailed, "发送验证码失败")
	}

	// 存入Redis
	if err = redis.Set(context.Background(), key, code, constants.RedisRegisterVerifyCodeTTL); err != nil {
		log.Error("set verify code error", zap.Error(err))
		return errorx.New(errorx.ErrOperationFailed, "发送验证码失败")
	}

	return nil
}

// AdminLogin 管理员登录
func (s *UserService) AdminLogin(req *model.AdminLoginReq) (any, error) {
	var (
		user         *model.User
		err          error
		token        string
		refreshToken string
	)
	err = db.Transaction(func(tx *gorm.DB) error {
		if user, err = db.Get[model.User](tx, map[string]any{"email": req.Email}); err != nil {
			log.Error("query user by username error", zap.Error(err))
			return errorx.InternalServerError("登录失败")
		}
		if user == nil {
			return errorx.New(errorx.ErrOperationFailed, "邮箱或密码错误")
		}

		if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			return errorx.New(errorx.ErrOperationFailed, "邮箱或密码错误")
		}

		if user.Status == constants.UserStatusDisabled {
			return errorx.New(errorx.ErrNoPermission, "用户已被禁用")
		}

		token, refreshToken, err = genTokens(user.ID)
		if err != nil {
			log.Error("generate token error", zap.Error(err))
			return errorx.InternalServerError("登录失败")
		}

		// 更新最后登录时间
		now := time.Now()
		user.LastLoginAt = &now
		if err = db.Update(tx, user, map[string]any{"id": user.ID}, "last_login_at"); err != nil {
			log.Error("update last login time error", zap.Error(err))
			return errorx.InternalServerError("登录失败")
		}

		return nil
	})

	if err != nil {
		log.Error("用户登录失败", zap.Error(err))
		return nil, err
	}

	return map[string]string{
		"token":        token,
		"refreshToken": refreshToken,
	}, nil
}

func (s *UserService) AdminGetUserList(req *model.AdminGetUserListReq) (*model.AdminGetUserListRes, error) {
	var (
		users []*model.User
		total int64
		err   error
	)

	// 构建查询条件
	conditions := make(map[string]any)

	// 邮箱筛选
	if req.Email != "" {
		conditions["email"] = []any{"LIKE", "%" + req.Email + "%"}
	}

	// 昵称筛选
	if req.Nickname != "" {
		conditions["nickname"] = []any{"LIKE", "%" + req.Nickname + "%"}
	}

	// 状态筛选
	if req.Status != 0 {
		conditions["status"] = req.Status
	}

	// 设置默认分页参数
	page := req.Page
	size := req.Size
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}

	// 查询用户列表
	users, total, err = db.List[model.User](db.GetDB(), page, size, conditions, []string{"created_at DESC"})
	if err != nil {
		log.Error("查询用户列表失败", zap.Error(err))
		return nil, errorx.InternalServerError("查询用户列表失败")
	}

	// 转换为响应格式
	userResponses := make([]*model.UserResponse, 0, len(users))
	for _, user := range users {
		userResponses = append(userResponses, user.ToResponse())
	}

	return &model.AdminGetUserListRes{
		List:  userResponses,
		Total: total,
	}, nil
}

// Login 用户登录
func (s *UserService) Login(req *model.LoginReq) (*model.LoginRes, error) {
	var (
		user         *model.User
		err          error
		token        string
		refreshToken string
	)

	err = db.Transaction(func(tx *gorm.DB) error {
		if user, err = db.Get[model.User](tx, map[string]any{"email": req.Email}); err != nil {
			log.Error("query user by username error", zap.Error(err))
			return errorx.InternalServerError("登录失败")
		}
		if user == nil {
			return errorx.New(errorx.ErrOperationFailed, "用户名或密码错误")
		}

		if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			return errorx.New(errorx.ErrOperationFailed, "用户名或密码错误")
		}

		if user.Status == constants.UserStatusDisabled {
			return errorx.New(errorx.ErrNoPermission, "用户已被禁用")
		}

		// 生成JWT令牌
		token, refreshToken, err = genTokens(user.ID)
		if err != nil {
			log.Error("generate token error", zap.Error(err))
			return errorx.InternalServerError("登录失败")
		}

		// 更新最后登录时间
		now := time.Now()
		user.LastLoginAt = &now
		if err = db.Update(tx, user, map[string]any{"id": user.ID}, "last_login_at"); err != nil {
			log.Error("update last login time error", zap.Error(err))
			return errorx.InternalServerError("登录失败")
		}

		return nil
	})

	if err != nil {
		log.Error("用户登录失败", zap.Error(err))
		return nil, err
	}

	response := &model.LoginRes{
		Token:        token,
		RefreshToken: refreshToken,
		UserInfo:     user.ToResponse(),
	}

	return response, nil
}

// GetUser 获取用户资料
func (s *UserService) GetUser(userID string) (*model.UserResponse, error) {
	user, err := db.Get[model.User](db.GetDB(), map[string]any{"id": userID})
	if err != nil {
		log.Error("获取用户失败", zap.Error(err))
		return nil, errors.New("服务错误")
	}
	if user == nil {
		return nil, errors.New("用户不存在")
	}

	response := user.ToResponse()
	return response, nil
}

// UpdateInfo 更新用户资料
func (s *UserService) UpdateInfo(userID string, req *model.UpdateInfoReq) (*model.UserResponse, error) {
	var (
		updateFields []string
		user         *model.User
		err          error
	)
	err = db.Transaction(func(tx *gorm.DB) error {
		user, err = db.Get[model.User](tx, map[string]any{"id": userID})
		if err != nil {
			log.Error("query user error", zap.Error(err))
			return errorx.InternalServerError("更新失败")
		}
		if user == nil {
			return errorx.New(errorx.ErrNotFound, "用户不存在")
		}

		if req.Nickname != nil && *req.Nickname != "" {
			user.Nickname = *req.Nickname
			updateFields = append(updateFields, "nickname")
		}
		if req.Avatar != nil && *req.Avatar != "" {
			user.Avatar = *req.Avatar
			updateFields = append(updateFields, "avatar")
		}

		if req.Tags != nil && len(*req.Tags) > 0 {
			user.Tags = strings.Join(*req.Tags, ",")
			updateFields = append(updateFields, "tags")
		}

		if err = db.Update(tx, user, map[string]any{"id": userID}, updateFields...); err != nil {
			log.Error("update user error", zap.Error(err))
			return errorx.InternalServerError("更新失败")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return response, nil
}

// ResetPassword 重置密码
func (s *UserService) ResetPassword(req *model.ResetPasswordReq) error {
	var (
		user           *model.User
		hashedPassword []byte
		err            error
	)

	// 校验邮箱验证码
	if err = judgeVerifyCode(req.Email, req.VerifyCode); err != nil {
		log.Error("verify code error", zap.Error(err))
		return err
	}

	if err = db.Transaction(func(tx *gorm.DB) error {
		user, err = db.Get[model.User](tx, map[string]any{"email": req.Email})
		if err != nil {
			log.Error("query user error", zap.Error(err))
			return errorx.InternalServerError("修改密码失败")
		}
		if user == nil {
			return errorx.New(errorx.ErrNotFound, "用户不存在")
		}

		// 加密新密码
		hashedPassword, err = bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Error("encrypt password error", zap.Error(err))
			return errorx.InternalServerError("修改密码失败")
		}

		user.Password = string(hashedPassword)

		if err = db.Update(tx, user, map[string]any{"id": user.ID}, "password"); err != nil {
			log.Error("update password error", zap.Error(err))
			return errorx.InternalServerError("修改密码失败")
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}
