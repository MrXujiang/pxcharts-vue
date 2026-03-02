package middleware

import (
	"mvtable/internal/app/user/model"
	"mvtable/internal/pkg/constants"
	"mvtable/internal/storage/db"
	"mvtable/pkg/jwt"
	"mvtable/pkg/log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func isAdminURL(path string) bool {
	return strings.HasPrefix(path, "/api/v1/admin")
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractToken(c)
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "missing token",
			})
			c.Abort()
			return
		}

		// 验证token
		claims, err := jwt.GetInstance().ParseToken(tokenString)
		if err != nil {
			log.Warn("parse token warning", zap.Error(err))
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "invalid token",
			})
			c.Abort()
			return
		}

		user, err := db.Get[model.User](db.GetDB(), map[string]any{"id": claims.UserID})
		if err != nil {
			log.Error("get user error", zap.Error(err))
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "获取用户信息失败",
			})
			c.Abort()
			return
		}
		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "用户不存在",
			})
			c.Abort()
			return
		}

		// 校验用户是否被禁用
		if user.Status != constants.UserStatusActive {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "用户已被禁用",
			})
			c.Abort()
			return
		}

		// 管理后台只有管理员能够访问
		if isAdminURL(c.Request.URL.Path) {
			if user.Role != constants.RoleAdmin {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code": 401,
					"msg":  "无权限",
				})
				c.Abort()
				return
			}
		}

		// 将用户信息存储到上下文中
		c.Set("userId", claims.UserID)
		c.Set("user", user)

		c.Next()
	}
}

// extractToken 尝试从 Header、Cookie 或 URL 参数中获取 JWT（优先 Header）
// 支持：
// - Authorization: Bearer <token>
// - Cookie: token=<token> 或 Authorization=<token>（可带或不带 Bearer 前缀）
// - URL 参数: ?token=<token>（用于 WebSocket 连接）
func extractToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		if after, ok := strings.CutPrefix(authHeader, "Bearer "); ok {
			return after
		}
		// 非 Bearer 前缀时，直接返回原值作为 token
		return authHeader
	}

	// 尝试从 URL 参数获取（用于 WebSocket 连接）
	if token := c.Query("token"); token != "" {
		return token
	}

	// 尝试从 Cookie 获取
	cookieCandidates := []string{"token", "Authorization"}
	for _, key := range cookieCandidates {
		if val, err := c.Cookie(key); err == nil && val != "" {
			if after, ok := strings.CutPrefix(val, "Bearer "); ok {
				return after
			}
			return val
		}
	}
	return ""
}

func GetCurrentUserID(c *gin.Context) string {
	userID := c.MustGet("userId").(string)
	return userID
}

func GetCurrentUser(c *gin.Context) *model.User {
	user := c.MustGet("user").(*model.User)
	return user
}
