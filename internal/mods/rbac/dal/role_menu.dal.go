package dal

import (
	"context"
	"gin-admin/internal/mods/rbac/schema"
	"gin-admin/pkg/errors"
	"gin-admin/pkg/util"

	"gorm.io/gorm"
)

func GetRoleMenuDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return util.GetDB(ctx, defDB).Model(new(schema.RoleMenu))
}

type RoleMenu struct {
	DB *gorm.DB
}

func (a *RoleMenu) Query(ctx context.Context, params schema.RoleMenuQueryParam, opts ...schema.RoleMenuQueryOptions) (*schema.RoleMenuQueryResult, error) {
	var opt schema.RoleMenuQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	db := GetRoleMenuDB(ctx, a.DB)
	if v := params.RoleID; len(v) > 0 {
		db = db.Where("role_id = ?", v)
	}

	var list schema.RoleMenus
	pageResult, err := util.WrapPageQuery(ctx, db, params.PaginationParam, opt.QueryOptions, &list)
	if err != nil {
		return nil, err
	}
	queryResult := &schema.RoleMenuQueryResult{
		PageResult: pageResult,
		Data:       list,
	}
	return queryResult, nil
}

func (a *RoleMenu) Get(ctx context.Context, id string, opts ...schema.RoleMenuQueryOptions) (*schema.RoleMenu, error) {
	var opt schema.RoleMenuQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	item := new(schema.RoleMenu)
	ok, err := util.FindOne(ctx, GetRoleMenuDB(ctx, a.DB).Where("id=?", id), opt.QueryOptions, item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}
	return item, nil
}

func (a *RoleMenu) Create(ctx context.Context, item *schema.RoleMenu) error {
	result := GetRoleMenuDB(ctx, a.DB).Create(item)
	return errors.WithStack(result.Error)
}

func (a *RoleMenu) Update(ctx context.Context, id string, item *schema.RoleMenu) error {
	result := GetRoleMenuDB(ctx, a.DB).Where("id=?", id).Updates(item)
	return errors.WithStack(result.Error)
}

func (a *RoleMenu) Delete(ctx context.Context, id string) error {
	result := GetRoleMenuDB(ctx, a.DB).Where("id=?", id).Delete(new(schema.RoleMenu))
	return errors.WithStack(result.Error)
}

func (a *RoleMenu) DeleteByRoleID(ctx context.Context, roleID string) error {
	result := GetRoleMenuDB(ctx, a.DB).Where("role_id=?", roleID).Delete(new(schema.RoleMenu))
	return errors.WithStack(result.Error)
}

func (a *RoleMenu) DeleteByMenuID(ctx context.Context, menuID string) error {
	result := GetRoleMenuDB(ctx, a.DB).Where("menu_id=?", menuID).Delete(new(schema.RoleMenu))
	return errors.WithStack(result.Error)
}
