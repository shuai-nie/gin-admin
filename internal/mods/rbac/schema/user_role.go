package schema

import (
	"gin-admin/internal/config"
	"gin-admin/pkg/util"
	"time"
)

type UserRole struct {
	ID        string
	UserID    string
	RoleID    string
	CreatedAt time.Time
	UpdatedAt time.Time
	RoleName  string
}

func (a *UserRole) TableName() string {
	return config.C.FormatTableName("user_role")
}

type UserRoleQueryParam struct {
	util.PaginationParam
	InUserIDs []string
	UserID    string
	RoleID    string
}

type UserRoleQueryOptions struct {
	util.QueryOptions
	JoinRole bool
}

type UserRoleQueryResult struct {
	Data       UserRoles
	PageResult *util.PaginationResult
}

type UserRoles []*UserRole

func (a UserRoles) ToUserIDMap() map[string]UserRoles {
	m := make(map[string]UserRoles)
	for _, userRole := range a {
		m[userRole.UserID] = append(m[userRole.UserID], userRole)
	}
	return m
}

func (a UserRoles) ToRoleIDs() []string {
	var ids []string
	for _, item := range a {
		ids = append(ids, item.RoleID)
	}
	return ids
}

type UserRoleForm struct{}

func (a *UserRoleForm) Validate() error {
	return nil
}

func (a *UserRoleForm) FillTo(UserRole *UserRole) error {
	return nil
}
