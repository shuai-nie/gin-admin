package schema

import (
	"gin-admin/internal/config"
	"gin-admin/pkg/util"
	"time"
)

type RoleMenu struct {
	ID        string
	RoleID    string
	MenuID    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (a *RoleMenu) TableName() string {
	return config.C.FormatTableName("role_menu")
}

type RoleMenuQueryParam struct {
	util.PaginationParam
	RoleID string
}

type RoleMenuQueryOptions struct {
	util.QueryOptions
}

type RoleMenuQueryResult struct {
	Data       RoleMenus
	PageResult *util.PaginationResult
}

type RoleMenus []*RoleMenu

type RoleMenuForm struct {
}

func (a *RoleMenuForm) Validate() error {
	return nil
}

func (a *RoleMenuForm) FillTo(roleMenu *RoleMenu) error {
	return nil
}
