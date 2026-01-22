package schema

import (
	"encoding/json"
	"gin-admin/internal/config"
	"gin-admin/pkg/errors"
	"gin-admin/pkg/util"
	"strings"
	"time"
)

const (
	MenuStatusDisabled = "disabled"
	MenuStatusEnabled  = "enabled"
)

var (
	MenusOrderParams = []util.OrderByParam{
		{Field: "sequence", Direction: util.DESC},
		{Field: "created_at", Direction: util.DESC},
	}
)

type Menu struct {
	ID          string
	Code        string
	Name        string
	Description string
	Sequence    int
	Type        string
	Path        string
	Properties  string
	Status      string
	ParentID    string
	ParentPath  string
	Children    *Menus
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Resources   MenuResources
}

func (a *Menu) TableName() string {
	return config.C.FormatTableName("menu")
}

type MenuQueryParam struct {
	util.PaginationParam
	CodePath         string   `form:"code"`
	LikeName         string   `form:"name"`
	IncludeResources bool     `form:"includeResources"`
	InIDs            []string `form:"-"`
	Status           string   `form:"-"`
	ParentID         string   `form:"-"`
	ParentPathPrefix string   `form:"-"`
	UserID           string   `form:"-"`
	RoleID           string   `form:"-"`
}

type MenuQueryOptions struct {
	util.QueryOptions
}

type MenuQueryResult struct {
	Data       Menus
	PageResult *util.PaginationResult
}

type Menus []*Menu

func (a Menus) Len() int {
	return len(a)
}

func (a Menus) Less(i, j int) bool {
	if a[i].Sequence == a[j].Sequence {
		return a[i].CreatedAt.Unix() > a[j].CreatedAt.Unix()
	}
	return a[i].Sequence > a[j].Sequence
}

func (a Menus) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a Menus) ToMap() map[string]*Menu {
	m := make(map[string]*Menu)
	for _, item := range a {
		m[item.ID] = item
	}
	return m
}

func (a Menus) SplitParentIDs() []string {
	parentIDs := make([]string, 0, len(a))
	idMapper := make(map[string]struct{})
	for _, item := range a {
		if _, ok := idMapper[item.ID]; ok {
			continue
		}
		idMapper[item.ID] = struct{}{}
		if pp := item.ParentPath; pp != "" {
			for _, pid := range strings.Split(pp, util.TreePathDelimiter) {
				if pid == "" {
					continue
				}
				if _, ok := idMapper[pid]; ok {
					continue
				}
				parentIDs = append(parentIDs, pid)
				idMapper[pid] = struct{}{}
			}
		}
	}
	return parentIDs
}

func (a Menus) ToTree() Menus {
	var list Menus
	m := a.ToMap()
	for _, item := range a {
		if item.ParentID == "" {
			list = append(list, item)
			continue
		}
		if parent, ok := m[item.ParentID]; ok {
			if parent.Children == nil {
				children := Menus{item}
				parent.Children = &children
				continue
			}
			*parent.Children = append(*parent.Children, item)
		}
	}
	return list
}

type MenuForm struct {
	Code        string
	Name        string
	Description string
	Sequence    int
	Type        string
	Path        string
	Properties  string
	Status      string
	ParentID    string
	Resources   MenuResources
}

func (a *MenuForm) Validate() error {
	if v := a.Properties; v != "" {
		if !json.Valid([]byte(v)) {
			return errors.BadRequest("", "Invalid")
		}
	}
	return nil
}

func (a *MenuForm) FillTo(menu *Menu) error {
	menu.Code = a.Code
	menu.Name = a.Name
	menu.Description = a.Description
	menu.Sequence = a.Sequence
	menu.Type = a.Type
	menu.Path = a.Path
	menu.Properties = a.Properties
	menu.Status = a.Status
	menu.ParentID = a.ParentID
	return nil
}
