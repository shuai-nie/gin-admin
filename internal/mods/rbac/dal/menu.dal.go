package dal

import (
	"context"
	"gin-admin/internal/mods/rbac/schema"
	"gin-admin/pkg/errors"
	"gin-admin/pkg/util"

	"gorm.io/gorm"
)

func GetMenuDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return util.GetDB(ctx, defDB).Model(new(schema.Menu))
}

type Menu struct {
	DB *gorm.DB
}

func (a *Menu) Query(ctx context.Context, params schema.MenuQueryParam, opts ...schema.MenuQueryOptions) (*schema.MenuQueryResult, error) {
	var opt schema.MenuQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	db := GetMenuDB(ctx, a.DB)
	if v := params.InIDs; len(v) > 0 {
		db = db.Where("id in ?", v)
	}

	if v := params.LikeName; v != "" {
		db = db.Where("name like ?", "%"+v+"%")
	}

	if v := params.Status; len(v) > 0 {
		db = db.Where("status=?", v)
	}

	if v := params.ParentID; len(v) > 0 {
		db = db.Where("parent_id=?", v)
	}

	if v := params.ParentPathPrefix; len(v) > 0 {
		db = db.Where("parent_path like ?", v+"%")
	}

	if v := params.RoleID; len(v) > 0 {
		userRoleQuery := GetUserRoleDB(ctx, a.DB).Where("user_id = ?", v).Select("role_id")
		roleMenuQuery := GetRoleMenuDB(ctx, a.DB).Where("role_id in ?", userRoleQuery).Select("menu_id")
		db = db.Where("id in ?", roleMenuQuery)
	}

	if v := params.RoleID; len(v) > 0 {
		roleMenuQuery := GetRoleMenuDB(ctx, a.DB).Where("role_id = ?", v).Select("menu_id")
		db = db.Where("id in ?", roleMenuQuery)
	}

	var list schema.Menus
	pageResult, err := util.WrapPageQuery(ctx, db, params.PaginationParam, opt.QueryOptions, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	queryResult := schema.MenuQueryResult{
		PageResult: pageResult,
		Data:       list,
	}
	return &queryResult, nil
}

func (a *Menu) Get(ctx context.Context, id string, opts ...schema.MenuQueryOptions) (*schema.Menu, error) {
	var opt schema.MenuQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	item := new(schema.Menu)
	ok, err := util.FindOne(ctx, GetMenuDB(ctx, a.DB).Where("id=?", id), opt.QueryOptions, item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}
	return item, nil
}

func (a *Menu) GetByCodeAndParentID(ctx context.Context, code, parentID string, opts ...schema.MenuQueryOptions) (*schema.Menu, error) {
	var opt schema.MenuQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	item := new(schema.Menu)
	ok, err := util.FindOne(ctx, GetMenuDB(ctx, a.DB).Where("code=? and parent_id=?", code, parentID), opt.QueryOptions, item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item, nil
}

func (a *Menu) GetByNameAndParentID(ctx context.Context, name, parentID string, opts ...schema.MenuQueryOptions) (*schema.Menu, error) {
	var opt schema.MenuQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	item := new(schema.Menu)
	ok, err := util.FindOne(ctx, GetMenuDB(ctx, a.DB).Where("name=? and parent_id=?", name, parentID), opt.QueryOptions, item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}
	return item, nil
}

func (a *Menu) Exists(ctx context.Context, id string) (bool, error) {
	ok, err := util.Exists(ctx, GetMenuDB(ctx, a.DB).Where("id=?", id))
	return ok, errors.WithStack(err)
}

func (a *Menu) ExistsCodeByParentID(ctx context.Context, code, parentID string) (bool, error) {
	ok, err := util.Exists(ctx, GetMenuDB(ctx, a.DB).Where("code=? and parent_id=?", code, parentID))
	return ok, errors.WithStack(err)
}

func (a *Menu) ExistsNameByParentID(ctx context.Context, name, parentID string) (bool, error) {
	ok, err := util.Exists(ctx, GetMenuDB(ctx, a.DB).Where("name=? and parent_id=?", name, parentID))
	return ok, errors.WithStack(err)
}

func (a *Menu) Create(ctx context.Context, item *schema.Menu) error {
	result := GetMenuDB(ctx, a.DB).Create(item)
	return errors.WithStack(result.Error)
}

func (a *Menu) Update(ctx context.Context, item *schema.Menu) error {
	result := GetMenuDB(ctx, a.DB).Where("id=?", item.ID).Select("*").Omit("created_at").Updates(item)
	return errors.WithStack(result.Error)
}

func (a *Menu) Delete(ctx context.Context, id string) error {
	result := GetMenuDB(ctx, a.DB).Where("id=?", id).Delete(new(schema.Menu))
	return errors.WithStack(result.Error)
}

func (a *Menu) UpdateParentPath(ctx context.Context, id, parentPath string) error {
	result := GetMenuDB(ctx, a.DB).Where("id=?", id).Update("parent_path", parentPath)
	return errors.WithStack(result.Error)
}

func (a *Menu) UpdateStatusByParentPath(ctx context.Context, parentPath string, status string) error {
	result := GetMenuDB(ctx, a.DB).Where("parent_path like ?", parentPath+"%").Update("status", status)
	return errors.WithStack(result.Error)
}
