package constants

type TeamIdentity string

const (
	// TeamIdentityCreator 创建者
	TeamIdentityCreator TeamIdentity = "creator"
	// TeamIdentityManager 管理员
	TeamIdentityManager TeamIdentity = "manager"
	// TeamIdentityMember 成员
	TeamIdentityMember TeamIdentity = "member"
)
