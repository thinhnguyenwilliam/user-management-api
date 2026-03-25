// user-management-api/internal/domain/user/role.go
package domain

type Role string

const (
	RoleAdmin     Role = "admin"
	RoleModerator Role = "moderator"
	RoleMember    Role = "member"
)

func MapRole(level int32) Role {
	switch level {
	case 1:
		return RoleAdmin
	case 2:
		return RoleModerator
	case 3:
		return RoleMember
	default:
		return "unknown"
	}
}
