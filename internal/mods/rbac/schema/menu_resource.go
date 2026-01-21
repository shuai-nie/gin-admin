package schema

import (
	"gin-admin/internal/config"
	"gin-admin/pkg/util"
	"time"
)

type MenuResource struct {
	ID        string
	MenuID    string
	Method    string
	Path      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (a *MenuResource) TableName() string {
	return config.C.FormatTableName("menu_resource")
}

type MenuResourceQueryParam struct {
	util.PaginationParam
	MenuID  string
	MenuIDs []string
}

type MenuResourceQueryOptions struct {
	util.QueryOptions
}

type MenuResourceQueryResult struct {
	Data       MenuResources
	PageResult *util.PaginationResult
}

type MenuResources []*MenuResource

type MenuResourceForm struct{}

func (a *MenuResourceForm) Validate() error {
	return nil
}

func (a *MenuResources) FillTo(menuResource *MenuResource) error {
	return nil
}
