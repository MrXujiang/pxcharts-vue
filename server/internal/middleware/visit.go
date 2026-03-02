package middleware

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	websiteVisitService "mvtable/internal/app/website_visit/service"
)

// generateSessionID 生成会话ID
func generateSessionID(c *gin.Context) string {
	// 使用 IP + User-Agent + 时间戳生成会话ID
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	timestamp := time.Now().Unix() / 3600 // 按小时分组会话

	data := fmt.Sprintf("%s-%s-%d", ip, userAgent, timestamp)
	hash := md5.Sum([]byte(data))
	return fmt.Sprintf("%x", hash)[:16]
}

// shouldSkipVisitRecord 判断是否应该跳过记录访问
func shouldSkipVisitRecord(c *gin.Context) bool {
	path := c.Request.URL.Path
	method := c.Request.Method

	// 跳过静态资源
	if strings.HasPrefix(path, "/static/") ||
		strings.HasPrefix(path, "/assets/") ||
		strings.HasPrefix(path, "/public/") ||
		strings.Contains(path, ".") {
		return true
	}

	// 跳过健康检查
	if path == "/api/health" {
		return true
	}

	// 只记录 GET 请求，避免记录敏感的 POST/PUT/DELETE 数据
	if method != http.MethodGet {
		return true
	}

	// 可以根据需要添加更多过滤条件
	return false
}

// VisitRecordMiddleware 访问记录中间件
func VisitRecordMiddleware() gin.HandlerFunc {
	service := websiteVisitService.NewWebsiteVisitService()

	return func(c *gin.Context) {
		// 检查是否需要跳过记录
		if shouldSkipVisitRecord(c) {
			c.Next()
			return
		}

		// 获取请求信息
		ip := c.ClientIP()
		userAgent := c.GetHeader("User-Agent")
		referer := c.GetHeader("Referer")
		page := c.Request.URL.Path
		method := c.Request.Method

		// 生成会话ID
		sessionID := generateSessionID(c)

		// 获取用户ID（如果已登录）
		userID := ""
		if userIDValue, exists := c.Get("userId"); exists {
			if uid, ok := userIDValue.(string); ok {
				userID = uid
			}
		}

		// 记录访问（异步）
		service.RecordVisit(ip, userAgent, referer, page, method, sessionID, userID)

		c.Next()
	}
}