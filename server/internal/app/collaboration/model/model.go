package model

import "mvtable/internal/pkg/constants"

// CollaborationMessage 协同编辑消息
type CollaborationMessage struct {
	Type      string `json:"type"`      // 消息类型
	TableID   string `json:"tableId"`   // 表格ID
	UserID    string `json:"userId"`    // 用户ID
	Timestamp int64  `json:"timestamp"` // 时间戳
	Data      any    `json:"data"`      // 消息数据
}

// LockRequest 申请锁请求
type LockRequest struct {
	RecordID string `json:"recordId"` // 记录ID（行ID）
	FieldID  string `json:"fieldId"`  // 字段ID（列ID）
}

// LockRelease 释放锁请求
type LockRelease struct {
	RecordID string `json:"recordId"` // 记录ID
	FieldID  string `json:"fieldId"`  // 字段ID
}

// LockResponse 锁响应
type LockResponse struct {
	RecordID string `json:"recordId"` // 记录ID
	FieldID  string `json:"fieldId"`  // 字段ID
	UserID   string `json:"userId"`   // 持有锁的用户ID
	Success  bool   `json:"success"`  // 是否成功
	Message  string `json:"message"`  // 消息
}

// CellUpdate 单元格更新
type CellUpdate struct {
	RecordID string      `json:"recordId"` // 记录ID
	FieldID  string      `json:"fieldId"`  // 字段ID
	Value    interface{} `json:"value"`    // 值
}

// RowInsert 插入行
type RowInsert struct {
	PrevRecordID string `json:"prevRecordId"` // 前一行记录ID
	RecordID     string `json:"recordId"`     // 新插入的记录ID
}

// RowDelete 删除行
type RowDelete struct {
	RecordIDs []string `json:"recordIds"` // 要删除的记录ID列表
}

// ColumnInsert 插入列
type ColumnInsert struct {
	PrevFieldID string              `json:"prevFieldId"` // 前一列字段ID（用于定位插入位置）
	FieldID     string              `json:"fieldId"`     // 新插入的字段ID
	Title       string              `json:"title"`       // 字段标题
	Type        constants.FieldType `json:"type"`        // 字段类型
	Config      any                 `json:"config"`      // 字段配置
	OrderIndex  string              `json:"orderIndex"`  // 排序索引（可选）
}

// ColumnUpdate 更新列
type ColumnUpdate struct {
	FieldID string               `json:"fieldId"` // 字段ID
	Title   *string              `json:"title"`   // 字段标题（可选）
	Type    *constants.FieldType `json:"type"`    // 字段类型（可选）
	Config  any                  `json:"config"`  // 字段配置（可选）
}

// ColumnDelete 删除列
type ColumnDelete struct {
	FieldID string `json:"fieldId"` // 要删除的字段ID
}

// LockBroadcast 锁状态广播
type LockBroadcast struct {
	Locks []LockInfo `json:"locks"` // 锁列表
}

// LockInfo 锁信息
type LockInfo struct {
	RecordID string `json:"recordId"` // 记录ID
	FieldID  string `json:"fieldId"`  // 字段ID
	UserID   string `json:"userId"`   // 持有锁的用户ID
}

// DataBroadcast 数据更新广播
type DataBroadcast struct {
	Updates []CellUpdate `json:"updates"` // 更新列表
}

// JoinRoomRequest 加入房间请求
type JoinRoomRequest struct {
	TableID string `json:"tableId"` // 表格ID
}

// LeaveRoomRequest 离开房间请求
type LeaveRoomRequest struct {
	TableID string `json:"tableId"` // 表格ID
}

// UserJoinBroadcast 用户加入房间广播
type UserJoinBroadcast struct {
	UserID   string `json:"userId"`   // 用户ID
	Nickname string `json:"nickname"` // 用户昵称
	Avatar   string `json:"avatar"`   // 用户头像
	Color    string `json:"color"`    // 用户颜色
}

// UserLeaveBroadcast 用户离开房间广播
type UserLeaveBroadcast struct {
	UserID string `json:"userId"` // 用户ID
}

// RoomUsersBroadcast 房间用户列表广播
type RoomUsersBroadcast struct {
	Users []UserJoinBroadcast `json:"users"` // 用户列表
}

// ErrorMessage 错误消息
type ErrorMessage struct {
	Code    int    `json:"code"`    // 错误码
	Message string `json:"message"` // 错误消息
}

// 消息类型常量
const (
	MsgTypeLockRequest        = "lock_request"         // 申请锁
	MsgTypeLockRelease        = "lock_release"         // 释放锁
	MsgTypeLockResponse       = "lock_response"        // 锁响应
	MsgTypeCellUpdate         = "cell_update"          // 单元格更新
	MsgTypeRowInsert          = "row_insert"           // 插入行
	MsgTypeRowDelete          = "row_delete"           // 删除行
	MsgTypeColumnInsert       = "column_insert"        // 插入列
	MsgTypeColumnUpdate       = "column_update"        // 更新列
	MsgTypeColumnDelete       = "column_delete"        // 删除列
	MsgTypeLockBroadcast      = "lock_broadcast"       // 锁状态广播
	MsgTypeDataBroadcast      = "data_broadcast"       // 数据更新广播
	MsgTypeJoinRoom           = "join_room"            // 加入房间
	MsgTypeLeaveRoom          = "leave_room"           // 离开房间
	MsgTypeUserJoinBroadcast  = "user_join_broadcast"  // 用户加入广播
	MsgTypeUserLeaveBroadcast = "user_leave_broadcast" // 用户离开广播
	MsgTypeRoomUsersBroadcast = "room_users_broadcast" // 房间用户列表广播
	MsgTypeError              = "error"                // 错误消息
	MsgTypePing               = "ping"                 // 心跳
	MsgTypePong               = "pong"                 // 心跳响应
)
