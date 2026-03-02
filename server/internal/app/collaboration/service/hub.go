package service

import (
	"mvtable/internal/app/collaboration/model"
	"mvtable/pkg/log"

	"go.uber.org/zap"
)

// Hub 维护所有活跃的连接和房间
type Hub struct {
	// 房间映射：tableID -> clients
	rooms map[string]map[*Client]bool

	// 客户端注册通道
	register chan *Client

	// 客户端注销通道
	unregister chan *Client

	// 广播通道
	broadcast chan *BroadcastMessage
}

// BroadcastMessage 广播消息
type BroadcastMessage struct {
	TableID       string                      // 表格ID
	Message       *model.CollaborationMessage // 消息
	ExcludeUserID string                      // 排除的用户ID（不发送给该用户）
}

// NewHub 创建 Hub
func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *BroadcastMessage, 256),
	}
}

// RegisterClient 注册客户端
func (h *Hub) RegisterClient(client *Client) {
	h.register <- client
}

// UnregisterClient 注销客户端
func (h *Hub) UnregisterClient(client *Client) {
	h.unregister <- client
}

// Run 运行 Hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			// 客户端注册，记录日志
			log.Info("client registered", zap.String("userID", client.UserID))

		case client := <-h.unregister:
			// 客户端注销
			h.removeClient(client)
			log.Info("client unregistered", zap.String("userID", client.UserID))

		case message := <-h.broadcast:
			// 广播消息
			h.broadcastToRoom(message.TableID, message.Message, message.ExcludeUserID)
		}
	}
}

// JoinRoom 加入房间
func (h *Hub) JoinRoom(tableID string, client *Client) {
	// 初始化房间（如果不存在）
	if h.rooms[tableID] == nil {
		h.rooms[tableID] = make(map[*Client]bool)
	}

	// 添加客户端到房间
	h.rooms[tableID][client] = true
	client.Rooms[tableID] = true

	log.Info("client joined room",
		zap.String("userID", client.UserID),
		zap.String("tableID", tableID),
	)
}

// LeaveRoom 离开房间
func (h *Hub) LeaveRoom(tableID string, client *Client) {
	if h.rooms[tableID] != nil {
		delete(h.rooms[tableID], client)
		delete(client.Rooms, tableID)

		// 如果房间为空，删除房间
		if len(h.rooms[tableID]) == 0 {
			delete(h.rooms, tableID)
		}

		log.Info("client left room",
			zap.String("userID", client.UserID),
			zap.String("tableID", tableID),
		)
	}
}

// BroadcastToRoom 广播消息到房间
func (h *Hub) BroadcastToRoom(tableID string, message *model.CollaborationMessage, excludeUserID string) {
	h.broadcast <- &BroadcastMessage{
		TableID:       tableID,
		Message:       message,
		ExcludeUserID: excludeUserID,
	}
}

// broadcastToRoom 实际执行广播
func (h *Hub) broadcastToRoom(tableID string, message *model.CollaborationMessage, excludeUserID string) {
	room := h.rooms[tableID]
	if room == nil {
		return
	}

	for client := range room {
		// 排除指定用户
		if excludeUserID != "" && client.UserID == excludeUserID {
			continue
		}

		select {
		case client.SendChan <- message:
		default:
			// 发送失败，关闭连接
			close(client.SendChan)
			delete(room, client)
			delete(client.Rooms, tableID)
		}
	}
}

// removeClient 移除客户端
func (h *Hub) removeClient(client *Client) {
	// 从所有房间中移除客户端
	for tableID := range client.Rooms {
		if h.rooms[tableID] != nil {
			delete(h.rooms[tableID], client)
			if len(h.rooms[tableID]) == 0 {
				delete(h.rooms, tableID)
			}
		}
	}

	// 关闭发送通道
	close(client.SendChan)
}

// GetRoomUsers 获取房间内所有用户
func (h *Hub) GetRoomUsers(tableID string) []*Client {
	room := h.rooms[tableID]
	if room == nil {
		return []*Client{}
	}

	users := make([]*Client, 0, len(room))
	for client := range room {
		users = append(users, client)
	}
	return users
}
