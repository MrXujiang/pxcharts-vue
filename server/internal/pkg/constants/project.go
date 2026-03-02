package constants

type QueryProjectType int

const (
	// QueryProjectRecently 最近访问
	QueryProjectRecently QueryProjectType = 1
	// QueryProjectCreated 我创建的
	QueryProjectCreated QueryProjectType = 2
	// QueryProjectShared 我参与的
	QueryProjectShared QueryProjectType = 3
	// QueryProjectFavorite 我收藏的
	QueryProjectFavorite QueryProjectType = 4
)

type ShareRange int

const (
	ShareRangePrivate ShareRange = 1 // 仅协作者可见
	ShareRangePublic  ShareRange = 2 // 互联网公开
	ShareRangeTeam    ShareRange = 3 // 团队内公开
)

type ProjectRole string

// 项目角色
const (
	ProjectActionOwner  ProjectRole = "owner"
	ProjectActionAdmin  ProjectRole = "admin"
	ProjectActionReader ProjectRole = "reader"
	ProjectActionEditor ProjectRole = "editor"
)

func (p *ProjectRole) String() ProjectRole {
	switch *p {
	case ProjectActionOwner:
		return "owner"
	case ProjectActionAdmin:
		return "admin"
	case ProjectActionEditor:
		return "editor"
	case ProjectActionReader:
		return "reader"
	default:
		return ""
	}
}

// 项目权限等级
const (
	ActionRead  = "read"
	ActionEdit  = "edit"
	ActionMange = "manage"
	ActionNone  = "none"
)

type PermissionLevel int

const (
	PermissionNone PermissionLevel = iota
	PermissionRead
	PermissionEdit
	PermissionManage
)

func (p PermissionLevel) String() string {
	switch p {
	case PermissionRead:
		return ActionRead
	case PermissionEdit:
		return ActionEdit
	case PermissionManage:
		return ActionMange
	default:
		return ActionNone
	}
}

// 项目高级权限
const (
	AdvancedOperateRangeAll = "all"
	AdvancedOperateRangeOwn = "own"

	AdvancedFieldAccessAll    = "all"
	AdvancedFieldAccessCustom = "custom"

	AdvancedViewAccessAll    = "all"
	AdvancedViewAccessCustom = "custom"
)
