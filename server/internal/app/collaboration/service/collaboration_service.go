package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"mvtable/internal/app/collaboration/model"
	fieldModel "mvtable/internal/app/mv_field/model"
	fieldService "mvtable/internal/app/mv_field/service"
	projectModel "mvtable/internal/app/mv_project/model"
	advancedPermModel "mvtable/internal/app/mv_project_advanced_perm/model"
	projectPermModel "mvtable/internal/app/mv_project_perm/model"
	recordModel "mvtable/internal/app/mv_record/model"
	"mvtable/internal/app/mv_record/service"
	tableSchemaModel "mvtable/internal/app/mv_table_schema/model"
	viewModel "mvtable/internal/app/mv_view/model"
	userService "mvtable/internal/app/user/service"
	userTeamModel "mvtable/internal/app/user_team/model"
	"mvtable/internal/pkg/constants"
	"mvtable/internal/storage/db"
	"mvtable/internal/storage/redis"
	"mvtable/pkg/log"

	"strings"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// CollaborationService 协同编辑服务
type CollaborationService struct {
	lockManager   *LockManager
	recordService *service.MvRecordService
	fieldService  *fieldService.MvFieldService
	userService   *userService.UserService
	hub           *Hub
}

// NewCollaborationService 创建协同编辑服务
func NewCollaborationService(lockManager *LockManager, recordService *service.MvRecordService, fieldService *fieldService.MvFieldService, userService *userService.UserService, hub *Hub) *CollaborationService {
	return &CollaborationService{
		lockManager:   lockManager,
		recordService: recordService,
		fieldService:  fieldService,
		userService:   userService,
		hub:           hub,
	}
}

// HandleMessage 处理消息
func (s *CollaborationService) HandleMessage(ctx context.Context, client *Client, msg *model.CollaborationMessage) error {
	switch msg.Type {
	case model.MsgTypeLockRequest:
		return s.handleLockRequest(ctx, client, msg)
	case model.MsgTypeLockRelease:
		return s.handleLockRelease(ctx, client, msg)
	case model.MsgTypeCellUpdate:
		return s.handleCellUpdate(ctx, client, msg)
	case model.MsgTypeRowInsert:
		return s.handleRowInsert(ctx, client, msg)
	case model.MsgTypeRowDelete:
		return s.handleRowDelete(ctx, client, msg)
	case model.MsgTypeColumnInsert:
		return s.handleColumnInsert(ctx, client, msg)
	case model.MsgTypeColumnUpdate:
		return s.handleColumnUpdate(ctx, client, msg)
	case model.MsgTypeColumnDelete:
		return s.handleColumnDelete(ctx, client, msg)
	case model.MsgTypeJoinRoom:
		return s.handleJoinRoom(ctx, client, msg)
	case model.MsgTypeLeaveRoom:
		return s.handleLeaveRoom(ctx, client, msg)
	case model.MsgTypePing:
		return s.handlePing(ctx, client, msg)
	// 以下消息类型是服务端到客户端的广播消息，客户端不应该发送
	case model.MsgTypeUserJoinBroadcast, model.MsgTypeUserLeaveBroadcast, model.MsgTypeRoomUsersBroadcast,
		model.MsgTypeLockBroadcast, model.MsgTypeDataBroadcast, model.MsgTypeLockResponse, model.MsgTypePong:
		// 忽略这些只读的广播消息，不返回错误
		return nil
	default:
		return s.sendError(client, "unknown message type: "+msg.Type)
	}
}

// resolveTableSchemaID: accept either tableSchemaID or viewID, return tableSchemaID if found
func (s *CollaborationService) resolveTableSchemaID(id string) (string, error) {
	if id == "" {
		return "", nil
	}

	// try table schema first
	tableSchema, err := db.Get[tableSchemaModel.MvTableSchema](db.GetDB(), map[string]any{"id": id})
	if err != nil {
		return "", err
	}
	if tableSchema != nil {
		return tableSchema.ID, nil
	}

	// try view id -> map to table schema id
	view, err := db.Get[viewModel.MvView](db.GetDB(), map[string]any{"id": id})
	if err != nil {
		return "", err
	}
	if view != nil {
		return view.TableSchemaID, nil
	}

	return "", nil
}

// handleLockRequest 处理申请锁请求
func (s *CollaborationService) handleLockRequest(ctx context.Context, client *Client, msg *model.CollaborationMessage) error {
	var req model.LockRequest
	if err := s.unmarshalData(msg.Data, &req); err != nil {
		return s.sendError(client, "invalid lock request: "+err.Error())
	}

	// resolve table id: accept either tableSchemaID or viewID
	resolvedTableID, err := s.resolveTableSchemaID(msg.TableID)
	if err != nil {
		log.Error("resolve table id error", zap.Error(err))
		return s.sendError(client, "获取表格信息失败")
	}
	if resolvedTableID == "" {
		return s.sendError(client, "表格不存在")
	}
	// override to canonical table schema id
	msg.TableID = resolvedTableID

	// 校验记录是否存在，并且属于该表格
	record, err := db.Get[recordModel.MvRecord](db.GetDB(), map[string]any{"id": req.RecordID})
	if err != nil {
		log.Error("get record error", zap.Error(err))
		return s.sendError(client, "获取记录信息失败")
	}
	if record == nil {
		return s.sendError(client, "记录不存在")
	}
	if record.TableSchemaID != msg.TableID {
		return s.sendError(client, "记录不属于该表格")
	}

	// 校验字段是否存在，并且属于该表格
	field, err := db.Get[fieldModel.MvField](db.GetDB(), map[string]any{"id": req.FieldID})
	if err != nil {
		log.Error("get field error", zap.Error(err))
		return s.sendError(client, "获取字段信息失败")
	}
	if field == nil {
		return s.sendError(client, "字段不存在")
	}
	if field.TableSchemaID != msg.TableID {
		return s.sendError(client, "字段不属于该表格")
	}

	// 所有校验通过，申请锁
	success, err := s.lockManager.AcquireLock(ctx, msg.TableID, req.RecordID, req.FieldID, msg.UserID)
	if err != nil {
		log.Error("acquire lock error", zap.Error(err))
		return s.sendError(client, "failed to acquire lock: "+err.Error())
	}

	response := model.CollaborationMessage{
		Type:      model.MsgTypeLockResponse,
		TableID:   msg.TableID,
		UserID:    msg.UserID,
		Timestamp: time.Now().Unix(),
		Data: model.LockResponse{
			RecordID: req.RecordID,
			FieldID:  req.FieldID,
			UserID:   msg.UserID,
			Success:  success,
			Message:  map[bool]string{true: "锁申请成功", false: "锁已被其他用户持有"}[success],
		},
	}

	if err := client.Send(&response); err != nil {
		log.Error("send lock response error", zap.Error(err))
		return err
	}

	// 如果成功，广播锁状态更新
	if success {
		s.broadcastLockState(ctx, msg.TableID)
	}

	return nil
}

// handleLockRelease 处理释放锁请求
func (s *CollaborationService) handleLockRelease(ctx context.Context, client *Client, msg *model.CollaborationMessage) error {
	var req model.LockRelease
	if err := s.unmarshalData(msg.Data, &req); err != nil {
		return s.sendError(client, "invalid lock release: "+err.Error())
	}

	// resolve table id if viewId was provided
	resolvedTableID, err := s.resolveTableSchemaID(msg.TableID)
	if err != nil {
		log.Error("resolve table id error", zap.Error(err))
		return s.sendError(client, "获取表格信息失败")
	}
	if resolvedTableID == "" {
		return s.sendError(client, "表格不存在")
	}
	msg.TableID = resolvedTableID

	err = s.lockManager.ReleaseLock(ctx, msg.TableID, req.RecordID, req.FieldID, msg.UserID)
	if err != nil {
		log.Error("release lock error", zap.Error(err))
		return s.sendError(client, "failed to release lock: "+err.Error())
	}

	// 广播锁状态更新
	s.broadcastLockState(ctx, msg.TableID)

	return nil
}

// handleCellUpdate 处理单元格更新
func (s *CollaborationService) handleCellUpdate(ctx context.Context, client *Client, msg *model.CollaborationMessage) error {
	var update model.CellUpdate
	if err := s.unmarshalData(msg.Data, &update); err != nil {
		return s.sendError(client, "invalid cell update: "+err.Error())
	}

	// resolve table id if viewId was provided
	resolvedTableID, err := s.resolveTableSchemaID(msg.TableID)
	if err != nil {
		log.Error("resolve table id error", zap.Error(err))
		return s.sendError(client, "获取表格信息失败")
	}
	if resolvedTableID == "" {
		return s.sendError(client, "表格不存在")
	}
	msg.TableID = resolvedTableID

	// 检查锁是否由当前用户持有
	owner, err := s.lockManager.GetLockOwner(ctx, msg.TableID, update.RecordID, update.FieldID)
	if err != nil {
		log.Error("get lock owner error", zap.Error(err))
		return s.sendError(client, "failed to check lock: "+err.Error())
	}

	if owner != msg.UserID {
		return s.sendError(client, "you don't have the lock for this cell")
	}

	// 更新数据库
	updateReq := &recordModel.UpdateCellReq{
		RecordID: update.RecordID,
		FieldID:  update.FieldID,
		Value:    update.Value,
	}
	if err := s.recordService.UpdateCell(msg.UserID, updateReq); err != nil {
		log.Error("update cell error", zap.Error(err))
		return s.sendError(client, "failed to update cell: "+err.Error())
	}

	// 广播数据更新
	s.broadcastDataUpdate(msg.TableID, []model.CellUpdate{update}, msg.UserID)

	return nil
}

// handleRowInsert 处理插入行
func (s *CollaborationService) handleRowInsert(ctx context.Context, client *Client, msg *model.CollaborationMessage) error {
	var req model.RowInsert
	if err := s.unmarshalData(msg.Data, &req); err != nil {
		return s.sendError(client, "invalid row insert: "+err.Error())
	}

	// resolve table id if viewId was provided
	resolvedTableID, err := s.resolveTableSchemaID(msg.TableID)
	if err != nil {
		log.Error("resolve table id error", zap.Error(err))
		return s.sendError(client, "获取表格信息失败")
	}
	if resolvedTableID == "" {
		return s.sendError(client, "表格不存在")
	}
	msg.TableID = resolvedTableID

	// 插入行到数据库
	insertReq := &recordModel.InsertMvRecordReq{
		PrevID:        req.PrevRecordID,
		TableSchemaID: msg.TableID,
	}
	recordData, err := s.recordService.InsertMvRecord(msg.UserID, insertReq)
	if err != nil {
		log.Error("insert row error", zap.Error(err))
		return s.sendError(client, "failed to insert row: "+err.Error())
	}

	// 从返回的记录数据中提取记录ID
	if recordMap, ok := recordData.(map[string]any); ok {
		if rowId, exists := recordMap["rowId"]; exists {
			if recordIdStr, ok := rowId.(string); ok {
				req.RecordID = recordIdStr
			} else {
				log.Error("invalid recordId type", zap.Any("rowId", rowId))
				return s.sendError(client, "获取记录ID失败")
			}
		} else {
			log.Error("rowId not found in record data", zap.Any("recordData", recordData))
			return s.sendError(client, "获取记录ID失败")
		}
	} else {
		log.Error("invalid record data type", zap.Any("recordData", recordData))
		return s.sendError(client, "获取记录数据失败")
	}

	// 广播行插入（向所有客户端广播，包括发送者，这样发送者也能收到真实ID）
	broadcastMsg := model.CollaborationMessage{
		Type:      model.MsgTypeRowInsert,
		TableID:   msg.TableID,
		UserID:    msg.UserID,
		Timestamp: time.Now().Unix(),
		Data:      req, // req 已经包含了真实的 recordID
	}
	s.hub.BroadcastToRoom(msg.TableID, &broadcastMsg, "") // 空字符串表示向所有客户端广播

	return nil
}

// handleRowDelete 处理删除行
func (s *CollaborationService) handleRowDelete(ctx context.Context, client *Client, msg *model.CollaborationMessage) error {
	var req model.RowDelete
	if err := s.unmarshalData(msg.Data, &req); err != nil {
		return s.sendError(client, "invalid row delete: "+err.Error())
	}

	// resolve table id if viewId was provided
	resolvedTableID, err := s.resolveTableSchemaID(msg.TableID)
	if err != nil {
		log.Error("resolve table id error", zap.Error(err))
		return s.sendError(client, "获取表格信息失败")
	}
	if resolvedTableID == "" {
		return s.sendError(client, "表格不存在")
	}
	msg.TableID = resolvedTableID

	// 删除行
	deleteReq := &recordModel.DeleteRecordsReq{
		RecordIDs: req.RecordIDs,
	}
	if err := s.recordService.DeleteRecords(msg.UserID, deleteReq); err != nil {
		log.Error("delete rows error", zap.Error(err))
		return s.sendError(client, "failed to delete rows: "+err.Error())
	}

	// 释放相关锁
	// 这里需要知道所有字段ID，简化处理：释放该记录的所有锁
	// 实际应该查询表格的所有字段
	// 暂时不处理，锁会自动过期
	_ = req.RecordIDs

	// 广播行删除
	broadcastMsg := model.CollaborationMessage{
		Type:      model.MsgTypeRowDelete,
		TableID:   msg.TableID,
		UserID:    msg.UserID,
		Timestamp: time.Now().Unix(),
		Data:      req,
	}
	s.hub.BroadcastToRoom(msg.TableID, &broadcastMsg, msg.UserID)

	return nil
}

// handleColumnInsert 处理插入列
func (s *CollaborationService) handleColumnInsert(ctx context.Context, client *Client, msg *model.CollaborationMessage) error {
	var req model.ColumnInsert
	if err := s.unmarshalData(msg.Data, &req); err != nil {
		return s.sendError(client, "invalid column insert: "+err.Error())
	}

	// resolve table id if viewId was provided
	resolvedTableID, err := s.resolveTableSchemaID(msg.TableID)
	if err != nil {
		log.Error("resolve table id error", zap.Error(err))
		return s.sendError(client, "获取表格信息失败")
	}
	if resolvedTableID == "" {
		return s.sendError(client, "表格不存在")
	}
	msg.TableID = resolvedTableID

	// 校验字段管理权限
	if err := s.checkFieldManagePermission(ctx, msg.TableID, msg.UserID); err != nil {
		log.Warn("field manage permission denied",
			zap.String("userID", msg.UserID),
			zap.String("tableID", msg.TableID),
			zap.Error(err))
		return s.sendError(client, "无权限管理字段: "+err.Error())
	}

	// 创建字段
	createReq := &fieldModel.CreateMvFieldReq{
		TableSchemaID: msg.TableID,
		Title:         req.Title,
		Type:          req.Type,
		Config:        req.Config,
	}

	var createFieldID string
	if createFieldID, err = s.fieldService.CreateMvField(createReq); err != nil {
		log.Error("create field error", zap.Error(err))
		return s.sendError(client, "创建字段失败: "+err.Error())
	}

	// 获取创建后的字段信息（包含生成的ID）
	var createdField fieldModel.MvField
	if err := db.GetDB().Where("table_schema_id = ? AND id = ?", msg.TableID, createFieldID).
		Order("created_at DESC").First(&createdField).Error; err == nil {
		req.FieldID = createdField.ID
	}

	// 广播字段插入（向所有客户端广播，包括发送者，这样发送者也能收到真实ID）
	broadcastMsg := model.CollaborationMessage{
		Type:      model.MsgTypeColumnInsert,
		TableID:   msg.TableID,
		UserID:    msg.UserID,
		Timestamp: time.Now().Unix(),
		Data:      req, // req 已经包含了真实的 fieldID
	}
	s.hub.BroadcastToRoom(msg.TableID, &broadcastMsg, "") // 空字符串表示向所有客户端广播

	return nil
}

// handleColumnUpdate 处理更新列
func (s *CollaborationService) handleColumnUpdate(ctx context.Context, client *Client, msg *model.CollaborationMessage) error {
	var req model.ColumnUpdate
	if err := s.unmarshalData(msg.Data, &req); err != nil {
		return s.sendError(client, "invalid column update: "+err.Error())
	}

	// resolve table id if viewId was provided
	resolvedTableID, err := s.resolveTableSchemaID(msg.TableID)
	if err != nil {
		log.Error("resolve table id error", zap.Error(err))
		return s.sendError(client, "获取表格信息失败")
	}
	if resolvedTableID == "" {
		return s.sendError(client, "表格不存在")
	}
	msg.TableID = resolvedTableID

	// 校验字段是否存在，并且属于该表格
	field, err := db.Get[fieldModel.MvField](db.GetDB(), map[string]any{"id": req.FieldID})
	if err != nil {
		log.Error("get field error", zap.Error(err))
		return s.sendError(client, "获取字段信息失败")
	}
	if field == nil {
		return s.sendError(client, "字段不存在")
	}
	if field.TableSchemaID != msg.TableID {
		return s.sendError(client, "字段不属于该表格")
	}

	// 校验字段管理权限
	if err := s.checkFieldManagePermission(ctx, msg.TableID, msg.UserID); err != nil {
		log.Warn("field manage permission denied",
			zap.String("userID", msg.UserID),
			zap.String("tableID", msg.TableID),
			zap.Error(err))
		return s.sendError(client, "无权限管理字段: "+err.Error())
	}

	// 构建更新请求
	updateReq := &fieldModel.UpdateMvFieldReq{
		ID: req.FieldID,
	}
	if req.Title != nil {
		updateReq.Title = *req.Title
	} else {
		updateReq.Title = field.Title // 保持原值
	}
	if req.Type != nil {
		updateReq.Type = *req.Type
	} else {
		updateReq.Type = field.Type // 保持原值
	}
	if req.Config != nil {
		updateReq.Config = req.Config
	} else {
		updateReq.Config = field.Config // 保持原值
	}

	// 更新字段
	if err := s.fieldService.UpdateMvField(updateReq); err != nil {
		log.Error("update field error", zap.Error(err))
		return s.sendError(client, "更新字段失败: "+err.Error())
	}

	// 释放该字段上的所有锁（字段更新时建议释放锁）
	releasedLocks, err := s.lockManager.ReleaseAllLocksByField(ctx, msg.TableID, req.FieldID)
	if err != nil {
		log.Error("release locks by field error", zap.Error(err))
	} else if len(releasedLocks) > 0 {
		log.Info("released locks due to field update",
			zap.String("fieldID", req.FieldID),
			zap.Int("lockCount", len(releasedLocks)))
		// 通知持有锁的用户
		s.notifyLocksReleased(ctx, msg.TableID, releasedLocks)
	}

	// 广播字段更新
	broadcastMsg := model.CollaborationMessage{
		Type:      model.MsgTypeColumnUpdate,
		TableID:   msg.TableID,
		UserID:    msg.UserID,
		Timestamp: time.Now().Unix(),
		Data:      req,
	}
	s.hub.BroadcastToRoom(msg.TableID, &broadcastMsg, msg.UserID)

	return nil
}

// handleColumnDelete 处理删除列
func (s *CollaborationService) handleColumnDelete(ctx context.Context, client *Client, msg *model.CollaborationMessage) error {
	var req model.ColumnDelete
	if err := s.unmarshalData(msg.Data, &req); err != nil {
		return s.sendError(client, "invalid column delete: "+err.Error())
	}

	// resolve table id if viewId was provided
	resolvedTableID, err := s.resolveTableSchemaID(msg.TableID)
	if err != nil {
		log.Error("resolve table id error", zap.Error(err))
		return s.sendError(client, "获取表格信息失败")
	}
	if resolvedTableID == "" {
		return s.sendError(client, "表格不存在")
	}
	msg.TableID = resolvedTableID

	// 校验字段是否存在，并且属于该表格
	field, err := db.Get[fieldModel.MvField](db.GetDB(), map[string]any{"id": req.FieldID})
	if err != nil {
		log.Error("get field error", zap.Error(err))
		return s.sendError(client, "获取字段信息失败")
	}
	if field == nil {
		return s.sendError(client, "字段不存在")
	}
	if field.TableSchemaID != msg.TableID {
		return s.sendError(client, "字段不属于该表格")
	}

	// 校验字段管理权限
	if err := s.checkFieldManagePermission(ctx, msg.TableID, msg.UserID); err != nil {
		log.Warn("field manage permission denied",
			zap.String("userID", msg.UserID),
			zap.String("tableID", msg.TableID),
			zap.Error(err))
		return s.sendError(client, "无权限管理字段: "+err.Error())
	}

	// 释放该字段上的所有锁
	releasedLocks, err := s.lockManager.ReleaseAllLocksByField(ctx, msg.TableID, req.FieldID)
	if err != nil {
		log.Error("release locks by field error", zap.Error(err))
		// 不阻止删除操作，只记录错误
	} else if len(releasedLocks) > 0 {
		log.Info("released locks due to field deletion",
			zap.String("fieldID", req.FieldID),
			zap.Int("lockCount", len(releasedLocks)))
		// 通知持有锁的用户
		s.notifyLocksReleased(ctx, msg.TableID, releasedLocks)
	}

	// 删除字段
	deleteReq := &fieldModel.DeleteMvFieldReq{
		ID: req.FieldID,
	}
	if err := s.fieldService.DeleteMvField(deleteReq); err != nil {
		log.Error("delete field error", zap.Error(err))
		return s.sendError(client, "删除字段失败: "+err.Error())
	}

	// 广播列删除
	broadcastMsg := model.CollaborationMessage{
		Type:      model.MsgTypeColumnDelete,
		TableID:   msg.TableID,
		UserID:    msg.UserID,
		Timestamp: time.Now().Unix(),
		Data:      req,
	}
	s.hub.BroadcastToRoom(msg.TableID, &broadcastMsg, msg.UserID)

	return nil
}

// handleJoinRoom 处理加入房间
func (s *CollaborationService) handleJoinRoom(ctx context.Context, client *Client, msg *model.CollaborationMessage) error {
	log.Info("处理join_room消息", zap.String("clientID", client.UserID), zap.Any("msg", msg))
	var req model.JoinRoomRequest
	if err := s.unmarshalData(msg.Data, &req); err != nil {
		return s.sendError(client, "invalid join room request: "+err.Error())
	}
	log.Info("解析join_room请求", zap.String("tableID", req.TableID))

	// 如果用户信息未加载，先加载用户信息
	if client.Nickname == "" {
		user, err := s.userService.GetUser(client.UserID)
		if err != nil {
			log.Error("get user error", zap.Error(err))
			// 如果获取用户失败，使用默认值
			client.Nickname = "用户" + client.UserID[:8]
			client.Avatar = ""
		} else {
			client.Nickname = user.Nickname
			if client.Nickname == "" {
				client.Nickname = user.Email
			}
			client.Avatar = user.Avatar
		}
		// 生成用户颜色（基于用户ID）
		client.Color = generateUserColor(client.UserID)
	}

	// 加入房间
	// resolve viewId -> tableSchemaId if needed
	resolvedTableID, err := s.resolveTableSchemaID(req.TableID)
	if err != nil {
		log.Error("resolve table id error", zap.Error(err))
		return s.sendError(client, "获取表格信息失败")
	}
	if resolvedTableID == "" {
		return s.sendError(client, "表格不存在")
	}
	req.TableID = resolvedTableID
	s.hub.JoinRoom(req.TableID, client)

	// 广播用户加入消息给房间内其他用户
	s.broadcastUserJoin(req.TableID, client)

	// 发送当前房间用户列表给新加入的用户
	log.Info("调用sendRoomUsers", zap.String("tableID", req.TableID), zap.String("clientID", client.UserID))
	s.sendRoomUsers(req.TableID, client)

	// 发送当前锁状态
	s.broadcastLockState(ctx, req.TableID)

	return nil
}

// handleLeaveRoom 处理离开房间
func (s *CollaborationService) handleLeaveRoom(ctx context.Context, client *Client, msg *model.CollaborationMessage) error {
	var req model.LeaveRoomRequest
	if err := s.unmarshalData(msg.Data, &req); err != nil {
		return s.sendError(client, "invalid leave room request: "+err.Error())
	}

	// resolve viewId -> tableSchemaId if needed
	resolvedTableID, err := s.resolveTableSchemaID(req.TableID)
	if err != nil {
		log.Error("resolve table id error", zap.Error(err))
		return s.sendError(client, "获取表格信息失败")
	}
	if resolvedTableID == "" {
		return s.sendError(client, "表格不存在")
	}
	req.TableID = resolvedTableID

	// 广播用户离开消息
	s.broadcastUserLeave(req.TableID, client)

	// 离开房间
	s.hub.LeaveRoom(req.TableID, client)

	return nil
}

// handlePing 处理心跳
func (s *CollaborationService) handlePing(ctx context.Context, client *Client, msg *model.CollaborationMessage) error {
	response := model.CollaborationMessage{
		Type:      model.MsgTypePong,
		TableID:   msg.TableID,
		UserID:    msg.UserID,
		Timestamp: time.Now().Unix(),
		Data:      nil,
	}
	return client.Send(&response)
}

// broadcastLockState 广播锁状态
func (s *CollaborationService) broadcastLockState(ctx context.Context, tableID string) {
	// 查询 Redis 中该表格的所有锁键：collab:lock:{tableID}:*:{fieldID}
	pattern := fmt.Sprintf("collab:lock:%s:*:*", tableID)
	var cursor uint64 = 0
	var count int64 = 100
	var locks []model.LockInfo

	for {
		keys, nextCursor, err := redis.GetClient().Scan(ctx, cursor, pattern, count).Result()
		if err != nil {
			log.Error("scan locks error", zap.Error(err))
			// on error, fall back to empty list
			locks = []model.LockInfo{}
			break
		}

		for _, key := range keys {
			parts := strings.Split(key, ":")
			// expect format collab:lock:{tableID}:{recordID}:{fieldID}
			if len(parts) != 5 {
				continue
			}
			recordID := parts[3]
			fieldID := parts[4]
			userID, err := redis.Get(ctx, key)
			if err != nil {
				if err == redis.Nil {
					// key expired between scan and get
					continue
				}
				log.Warn("get lock owner error", zap.String("key", key), zap.Error(err))
				continue
			}
			if userID == "" {
				continue
			}
			locks = append(locks, model.LockInfo{
				RecordID: recordID,
				FieldID:  fieldID,
				UserID:   userID,
			})
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	broadcastMsg := model.CollaborationMessage{
		Type:      model.MsgTypeLockBroadcast,
		TableID:   tableID,
		Timestamp: time.Now().Unix(),
		Data: model.LockBroadcast{
			Locks: locks,
		},
	}
	s.hub.BroadcastToRoom(tableID, &broadcastMsg, "")
}

// broadcastDataUpdate 广播数据更新
func (s *CollaborationService) broadcastDataUpdate(tableID string, updates []model.CellUpdate, excludeUserID string) {
	broadcastMsg := model.CollaborationMessage{
		Type:      model.MsgTypeDataBroadcast,
		TableID:   tableID,
		Timestamp: time.Now().Unix(),
		Data: model.DataBroadcast{
			Updates: updates,
		},
	}
	s.hub.BroadcastToRoom(tableID, &broadcastMsg, excludeUserID)
}

// sendError 发送错误消息
func (s *CollaborationService) sendError(client *Client, message string) error {
	errorMsg := model.CollaborationMessage{
		Type:      model.MsgTypeError,
		Timestamp: time.Now().Unix(),
		Data: model.ErrorMessage{
			Code:    400,
			Message: message,
		},
	}
	return client.Send(&errorMsg)
}

// unmarshalData 解析消息数据
func (s *CollaborationService) unmarshalData(data interface{}, v interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, v)
}

// generateUserColor 基于用户ID生成颜色
func generateUserColor(userID string) string {
	// 使用 MD5 哈希用户ID，然后取前6位作为颜色
	hash := md5.Sum([]byte(userID))
	hexHash := hex.EncodeToString(hash[:])
	// 取前6位作为颜色值
	return "#" + hexHash[:6]
}

// broadcastUserJoin 广播用户加入
func (s *CollaborationService) broadcastUserJoin(tableID string, client *Client) {
	// If the user already has other active connections in this room, skip broadcasting duplicate join
	users := s.hub.GetRoomUsers(tableID)
	count := 0
	for _, u := range users {
		if u.UserID == client.UserID {
			count++
		}
	}
	if count > 1 {
		// not the first connection for this user in the room, skip broadcast
		return
	}

	broadcastMsg := model.CollaborationMessage{
		Type:      model.MsgTypeUserJoinBroadcast,
		TableID:   tableID,
		UserID:    client.UserID,
		Timestamp: time.Now().Unix(),
		Data: model.UserJoinBroadcast{
			UserID:   client.UserID,
			Nickname: client.Nickname,
			Avatar:   client.Avatar,
			Color:    client.Color,
		},
	}
	// 广播给房间内其他用户（不包括自己）
	s.hub.BroadcastToRoom(tableID, &broadcastMsg, client.UserID)
}

// broadcastUserLeave 广播用户离开
func (s *CollaborationService) broadcastUserLeave(tableID string, client *Client) {
	// Only broadcast leave when this was the last active connection for the user in the room
	users := s.hub.GetRoomUsers(tableID)
	remaining := 0
	for _, u := range users {
		if u.UserID == client.UserID && u != client {
			remaining++
		}
	}
	if remaining > 0 {
		// there are still other connections for this user in the room, do not broadcast leave
		return
	}

	broadcastMsg := model.CollaborationMessage{
		Type:      model.MsgTypeUserLeaveBroadcast,
		TableID:   tableID,
		UserID:    client.UserID,
		Timestamp: time.Now().Unix(),
		Data: model.UserLeaveBroadcast{
			UserID: client.UserID,
		},
	}
	// 广播给房间内其他用户
	s.hub.BroadcastToRoom(tableID, &broadcastMsg, client.UserID)
}

// sendRoomUsers 发送房间用户列表给指定客户端
func (s *CollaborationService) sendRoomUsers(tableID string, client *Client) {
	users := s.hub.GetRoomUsers(tableID)
	log.Info("发送房间用户列表", zap.String("tableID", tableID), zap.String("clientID", client.UserID), zap.Int("userCount", len(users)))

	// 转换为广播格式，去重同一 userId（可能存在多连接）
	userList := make([]model.UserJoinBroadcast, 0, len(users))
	seen := make(map[string]bool, len(users))
	for _, user := range users {
		// 排除自己
		if user.UserID == client.UserID {
			continue
		}
		if seen[user.UserID] {
			continue
		}
		seen[user.UserID] = true
		userList = append(userList, model.UserJoinBroadcast{
			UserID:   user.UserID,
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
			Color:    user.Color,
		})
	}

	broadcastMsg := model.CollaborationMessage{
		Type:      model.MsgTypeRoomUsersBroadcast,
		TableID:   tableID,
		Timestamp: time.Now().Unix(),
		Data: model.RoomUsersBroadcast{
			Users: userList,
		},
	}

	// 只发送给当前客户端
	log.Info("发送room_users_broadcast消息", zap.String("tableID", tableID), zap.String("clientID", client.UserID), zap.Int("usersCount", len(userList)))
	client.Send(&broadcastMsg)
}

// checkFieldManagePermission 检查字段管理权限
func (s *CollaborationService) checkFieldManagePermission(ctx context.Context, tableID, userID string) error {
	// 获取表格信息
	tableSchema, err := db.Get[tableSchemaModel.MvTableSchema](db.GetDB(), map[string]any{"id": tableID})
	if err != nil {
		log.Error("get table schema error", zap.Error(err))
		return fmt.Errorf("获取表格信息失败")
	}
	if tableSchema == nil {
		return fmt.Errorf("表格不存在")
	}

	projectID := tableSchema.ProjectID

	// 检查项目是否开启高级权限
	project, err := db.Get[projectModel.MvProject](db.GetDB(), map[string]any{"id": projectID})
	if err != nil {
		log.Error("get project error", zap.Error(err))
		return fmt.Errorf("获取项目信息失败")
	}
	if project == nil {
		return fmt.Errorf("项目不存在")
	}

	// 如果项目未开启高级权限，直接通过
	if !project.EnableAdvancedPerm {
		return nil
	}

	// 获取用户在项目中的角色
	role, err := s.getUserProjectRole(db.GetDB(), userID, projectID)
	if err != nil {
		log.Error("get user project role error", zap.Error(err))
		return fmt.Errorf("获取权限失败")
	}
	if role == "" {
		return fmt.Errorf("无权限访问该项目")
	}

	// 获取高级权限配置
	advancedPerm, err := db.Get[advancedPermModel.MvProjectAdvancedPerm](db.GetDB(), map[string]any{
		"table_schema_id": tableID,
		"role":            role,
	})
	if err != nil {
		log.Error("get advanced perm error", zap.Error(err))
		return fmt.Errorf("获取权限配置失败")
	}
	if advancedPerm == nil {
		return fmt.Errorf("无权限配置")
	}

	// 字段管理操作需要 manage 权限
	if advancedPerm.DataAction != constants.ActionMange {
		return fmt.Errorf("无权限管理字段")
	}

	return nil
}

// getUserProjectRole 获取用户在项目中的角色
func (s *CollaborationService) getUserProjectRole(gormDB *gorm.DB, userID, projectID string) (constants.ProjectRole, error) {
	// 首先检查用户是否是项目创建者
	project, err := db.Get[projectModel.MvProject](gormDB, map[string]any{"id": projectID})
	if err != nil {
		return "", err
	}
	if project != nil && project.UserID == userID {
		return constants.ProjectActionOwner, nil
	}

	// 查询用户直接权限
	userPerm, err := db.Get[projectPermModel.MvProjectPerm](gormDB, map[string]any{
		"project_id": projectID,
		"target_id":  userID,
		"target":     "user",
	})
	if err != nil {
		return "", err
	}

	var highestRole constants.ProjectRole
	if userPerm != nil {
		highestRole = userPerm.Role
	}

	// 查询用户所属的团队
	var userTeams []userTeamModel.UserTeam
	if err := gormDB.Where("user_id = ? AND deleted_at IS NULL", userID).Find(&userTeams).Error; err != nil {
		return "", err
	}

	// 查询团队权限
	for _, userTeam := range userTeams {
		teamPerm, err := db.Get[projectPermModel.MvProjectPerm](gormDB, map[string]any{
			"project_id": projectID,
			"target_id":  userTeam.TeamID,
			"target":     "team",
		})
		if err != nil {
			continue
		}
		if teamPerm != nil {
			// 比较权限等级，取最高权限
			if s.compareRole(teamPerm.Role, highestRole) > 0 {
				highestRole = teamPerm.Role
			}
		}
	}

	return highestRole, nil
}

// compareRole 比较角色权限等级
func (s *CollaborationService) compareRole(role1, role2 constants.ProjectRole) int {
	roleLevel := map[constants.ProjectRole]int{
		constants.ProjectActionOwner:  4,
		constants.ProjectActionAdmin:  3,
		constants.ProjectActionEditor: 2,
		constants.ProjectActionReader: 1,
	}

	level1 := roleLevel[role1]
	level2 := roleLevel[role2]

	if level1 > level2 {
		return 1
	} else if level1 < level2 {
		return -1
	}
	return 0
}

// notifyLocksReleased 通知用户锁已被释放
func (s *CollaborationService) notifyLocksReleased(ctx context.Context, tableID string, locks []model.LockInfo) {
	// 按用户分组锁
	userLocks := make(map[string][]model.LockInfo)
	for _, lock := range locks {
		userLocks[lock.UserID] = append(userLocks[lock.UserID], lock)
	}

	// 向每个用户发送锁释放通知
	for userID, userLockList := range userLocks {
		// 查找该用户在房间内的客户端
		roomUsers := s.hub.GetRoomUsers(tableID)
		for _, client := range roomUsers {
			if client.UserID == userID {
				// 发送锁状态更新（这些锁已被释放）
				notification := model.CollaborationMessage{
					Type:      model.MsgTypeLockBroadcast,
					TableID:   tableID,
					UserID:    userID,
					Timestamp: time.Now().Unix(),
					Data: model.LockBroadcast{
						Locks: userLockList,
					},
				}
				// 客户端收到后应该释放这些锁
				if err := client.Send(&notification); err != nil {
					log.Warn("send lock release notification error",
						zap.String("userID", userID),
						zap.Error(err))
				}
				break
			}
		}
	}
}
