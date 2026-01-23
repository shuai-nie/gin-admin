package schema

import (
	"gin-admin/internal/config"
	"gin-admin/pkg/util"
	"time"
)

const (
	RoleStatusEnabled    = "enabled"
	RoleStatusDisabled   = "disabled"
	RoleResultTypeSelect = "select"
)

type Role struct {
	ID          string
	Code        string
	Name        string
	Description string
	Sequence    int
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Menus       RoleMenus
}

func (a *Role) TableName() string {
	return config.C.FormatTableName("role")
}

type RoleQueryParam struct {
	util.PaginationParam
	LikeName    string
	Status      string
	ResultType  string
	InIDs       []string
	GtUpdatedAt *time.Time
}

type RoleQueryOptions struct {
	util.QueryOptions
}

type RoleQueryResult struct {
	Data       Roles
	PageResult *util.PaginationResult
}

type Roles []*Role

type RoleForm struct {
	Code        string
	Name        string
	Description string
	Sequence    int
	Status      string
	Menus       RoleMenus
}

func (a *RoleForm) Validate() error {
	return nil
}

func (a *RoleForm) FillTo(role *Role) error {
	role.Code = a.Code
	role.Name = a.Name
	role.Description = a.Description
	role.Sequence = a.Sequence
	role.Status = a.Status
	role.Menus = a.Menus
	return nil
}
