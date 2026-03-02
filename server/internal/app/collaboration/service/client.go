package service

import (
	"context"
	"encoding/json"
	"mvtable/internal/app/collaboration/model"
	"mvtable/pkg/log"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	// 写超时时间
	writeWait = 10 * time.Second

	// 读超时时间
	pongWait = 60 * time.Second

	// ping 间隔（必须小于 pongWait）
	pingPeriod = (pongWait * 9) / 10

	// 最大消息大小
	maxMessageSize = 512 * 1024 // 512KB
)

// Client 表示一个 WebSocket 客户端连接
type Client struct {
	Hub      *Hub
	Conn     *websocket.Conn
	SendChan chan *model.CollaborationMessage
	UserID   string
	Nickname string
	Avatar   string
	Color    string
	Rooms    map[string]bool // 当前所在的房间（tableID -> true）
}

// NewClient 创建客户端
func NewClient(hub *Hub, conn *websocket.Conn, userID string) *Client {
	return &Client{
		Hub:      hub,
		Conn:     conn,
		SendChan: make(chan *model.CollaborationMessage, 256),
		UserID:   userID,
		Rooms:    make(map[string]bool),
	}
}

// ReadPump 读取消息
func (c *Client) ReadPump(service *CollaborationService) {
	defer func() {
		// 在断开连接前，先释放该用户在 Redis 上的所有锁，并通知相关房间用户
		releasedMap, err := service.lockManager.ReleaseAllLocksByUser(context.Background(), c.UserID)
		if err != nil {
			log.Warn("release all locks by user error", zap.Error(err), zap.String("userID", c.UserID))
		} else {
			for tableID, locks := range releasedMap {
				if len(locks) > 0 {
					// 通知持有锁的用户这些锁已经被释放
					service.notifyLocksReleased(context.Background(), tableID, locks)
					// 广播最新锁状态到房间
					service.broadcastLockState(context.Background(), tableID)
				}
			}
		}

		// 在断开连接前，广播用户离开所有房间
		for tableID := range c.Rooms {
			service.broadcastUserLeave(tableID, c)
		}
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, messageBytes, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Error("websocket read error", zap.Error(err))
			}
			break
		}

		var msg model.CollaborationMessage
		if err := json.Unmarshal(messageBytes, &msg); err != nil {
			log.Error("unmarshal message error", zap.Error(err))
			service.sendError(c, "invalid message format")
			continue
		}

		// 设置用户ID
		msg.UserID = c.UserID

		// 处理消息
		ctx := context.Background()
		if err := service.HandleMessage(ctx, c, &msg); err != nil {
			log.Error("handle message error", zap.Error(err))
		}
	}
}

// WritePump 写入消息
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.SendChan:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub 关闭了通道
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			jsonBytes, err := json.Marshal(message)
			if err != nil {
				log.Error("marshal message error", zap.Error(err))
				continue
			}

			w.Write(jsonBytes)

			// 批量发送队列中的其他消息
			n := len(c.SendChan)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				msg := <-c.SendChan
				jsonBytes, err := json.Marshal(msg)
				if err != nil {
					log.Error("marshal message error", zap.Error(err))
					continue
				}
				w.Write(jsonBytes)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// Send 发送消息
func (c *Client) Send(message *model.CollaborationMessage) error {
	select {
	case c.SendChan <- message:
		return nil
	default:
		return ErrSendBufferFull
	}
}

// ErrSendBufferFull 发送缓冲区已满错误
var ErrSendBufferFull = &SendBufferFullError{}

type SendBufferFullError struct{}

func (e *SendBufferFullError) Error() string {
	return "send buffer is full"
}
