package handler

import (
	"mvtable/internal/app/collaboration/service"
	"mvtable/internal/middleware"
	"mvtable/pkg/log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// 允许所有来源，生产环境应该检查 Origin
		return true
	},
}

// CollaborationHandler 协同编辑处理器
type CollaborationHandler struct {
	hub     *service.Hub
	service *service.CollaborationService
}

// NewCollaborationHandler 创建协同编辑处理器
func NewCollaborationHandler(hub *service.Hub, collaborationService *service.CollaborationService) *CollaborationHandler {
	// 启动 Hub（在后台运行）
	go hub.Run()

	return &CollaborationHandler{
		hub:     hub,
		service: collaborationService,
	}
}

// HandleWebSocket 处理 WebSocket 连接
// @Summary WebSocket 协同编辑连接
// @Description 建立 WebSocket 连接进行协同编辑
// @Tags 协同编辑
// @Router /ws/collaboration [get]
func (h *CollaborationHandler) HandleWebSocket(c *gin.Context) {
	// 升级 HTTP 连接为 WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Error("websocket upgrade error", zap.Error(err))
		return
	}

	// 获取用户ID
	userID := middleware.GetCurrentUserID(c)
	if userID == "" {
		log.Error("user ID not found")
		conn.Close()
		return
	}

	// 创建客户端
	client := service.NewClient(h.hub, conn, userID)

	// 注册客户端
	h.hub.RegisterClient(client)

	// 启动读写协程
	go client.WritePump()
	go client.ReadPump(h.service)
}
