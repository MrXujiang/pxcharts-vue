package constants

const (
	UserStatusActive   = 1
	UserStatusDisabled = 2
)

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

// UserIdentity 用户权益
type UserIdentity int8

const (
	UserIdentityBasic    UserIdentity = 1
	UserIdentityPro      UserIdentity = 2
	UserIdentityFlagship UserIdentity = 3
)

func (i UserIdentity) MaxProjectCount() int {
	switch i {
	case UserIdentityBasic:
		return 10
	case UserIdentityPro:
		return 100
	case UserIdentityFlagship:
		return 1000
	default:
		return 0
	}
}

func (i UserIdentity) MaxTeamCount() int {
	switch i {
	case UserIdentityBasic:
		return 5
	case UserIdentityPro:
		return 20
	case UserIdentityFlagship:
		return 100
	default:
		return 0
	}
}

func (i UserIdentity) LimitTeamMemberCount() int {
	switch i {
	case UserIdentityBasic:
		return 10
	case UserIdentityPro:
		return 100
	case UserIdentityFlagship:
		return 1000
	default:
		return 0
	}
}
