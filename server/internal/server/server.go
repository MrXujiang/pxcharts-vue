package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	_ "mvtable/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	"mvtable/internal/middleware"
	"mvtable/internal/pkg/config"
	"mvtable/pkg/log"
)

// HTTPServer HTTP服务器
type HTTPServer struct {
	config *config.Config
	engine *gin.Engine
	server *http.Server
}

// NewHTTPServer 创建HTTP服务器
func NewHTTPServer(cfg *config.Config, router *Router) *HTTPServer {
	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	// 创建Gin引擎
	engine := gin.New()

	// 添加Swagger
	if cfg.Server.Mode == "debug" {
		log.Info("Swagger enabled")
		log.Info(fmt.Sprintf("Swagger URL: http://%s:%d/swagger/index.html", cfg.Server.Host, cfg.Server.Port))
		// 配置Swagger UI，支持Authorization头设置
		engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
			ginSwagger.DeepLinking(true),
			ginSwagger.DefaultModelsExpandDepth(-1),
			ginSwagger.DocExpansion("none"),
		))
	}

	// 添加中间件
	engine.Use(middleware.GinLogger())
	engine.Use(middleware.GinRecovery())
	engine.Use(requestid.New())
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 添加超时中间件（跳过 WebSocket 升级请求）
	timeoutMiddleware := timeout.New(
		timeout.WithTimeout(cfg.Server.ReadTimeout),
		timeout.WithHandler(func(c *gin.Context) {
			c.Next()
		}),
		timeout.WithResponse(func(c *gin.Context) {
			c.JSON(http.StatusRequestTimeout, gin.H{
				"code": 408,
				"msg":  "Request timeout",
			})
		}),
	)

	engine.Use(func(c *gin.Context) {
		// 检查是否是 WebSocket 升级请求
		if c.GetHeader("Upgrade") == "websocket" {
			// WebSocket 连接会 hijack HTTP 连接，跳过 timeout 中间件
			c.Next()
			return
		}
		// 对于普通 HTTP 请求，应用 timeout 中间件
		timeoutMiddleware(c)
	})

	// 设置路由
	httpServer := &HTTPServer{
		config: cfg,
		engine: engine,
	}
	router.SetupRoutes(engine)

	// 创建HTTP服务器
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      engine,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	httpServer.server = server
	return httpServer
}

// Start 启动服务器
func (s *HTTPServer) Start() error {
	log.Info("Starting HTTP server",
		zap.String("host", s.config.Server.Host),
		zap.Int("port", s.config.Server.Port),
		zap.String("mode", s.config.Server.Mode),
	)

	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("启动HTTP服务器失败: %w", err)
	}

	return nil
}

// Stop 停止服务器
func (s *HTTPServer) Stop(ctx context.Context) error {
	log.Info("Stopping HTTP server")

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("停止HTTP服务器失败: %w", err)
	}

	log.Info("HTTP server stopped")
	return nil
}

// GetEngine 获取Gin引擎
func (s *HTTPServer) GetEngine() *gin.Engine {
	return s.engine
}
