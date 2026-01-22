package dal

import (
	"context"
	"gin-admin/internal/mods/rbac/schema"
	"gin-admin/pkg/errors"
	"gin-admin/pkg/util"

	"gorm.io/gorm"
)

func GetMenuResourceDB(ctx context.Context, defDb *gorm.DB) *gorm.DB {
	return util.GetDB(ctx, defDb).Model(new(schema.MenuResource))
}

type MenuResource struct {
	DB *gorm.DB
}

func (a *MenuResource) Query(ctx context.Context, params schema.MenuResourceQueryParam, opts ...schema.MenuResourceQueryOptions) (*schema.MenuResourceQueryResult, error) {
	var opt schema.MenuResourceQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	db := GetMenuResourceDB(ctx, a.DB)
	if v := params.MenuID; len(v) > 0 {
		db = db.Where("menu_id = ?", v)
	}

	if v := params.MenuIDs; len(v) > 0 {
		db = db.Where("menu_id in ?", v)
	}

	var list schema.MenuResources
	pageResult, err := util.WrapPageQuery(ctx, db, params.PaginationParam, opt.QueryOptions, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	queryResult := &schema.MenuResourceQueryResult{
		PageResult: pageResult,
		Data:       list,
	}
	return queryResult, nil
}

func (a *MenuResource) Get(ctx context.Context, id string, opts ...schema.MenuResourceQueryOptions) (*schema.MenuResource, error) {
	var opt schema.MenuResourceQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	item := new(schema.MenuResource)
	ok, err := util.FindOne(ctx, GetMenuResourceDB(ctx, a.DB), opt.QueryOptions, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}
	return item, nil
}

func (a *MenuResource) Exists(ctx context.Context, id string) (bool, error) {
	ok, err := util.Exists(ctx, GetMenuResourceDB(ctx, a.DB).Where("id=?", id))
	return ok, errors.WithStack(err)
}

func (a *MenuResource) ExistsMethodPathByMenuID(ctx context.Context, menuID, method, path string) (bool, error) {
	ok, err := util.Exists(ctx, GetMenuResourceDB(ctx, a.DB).Where("menu_id=? and method=? and path=?", menuID, method, path))
	return ok, errors.WithStack(err)
}

func (a *MenuResource) Create(ctx context.Context, item schema.MenuResource) error {
	result := GetMenuResourceDB(ctx, a.DB).Create(item)
	return errors.WithStack(result.Error)
}

func (a *MenuResource) Update(ctx context.Context, id string, item schema.MenuResource) error {
	result := GetMenuResourceDB(ctx, a.DB).Where("id=?", id).Delete(new(schema.MenuResource))
	return errors.WithStack(result.Error)
}

func (a *MenuResource) DeleteByMenuID(ctx context.Context, menuID string) error {
	result := GetMenuResourceDB(ctx, a.DB).Where("menu_id=?", menuID).Delete(new(schema.MenuResource))
	return errors.WithStack(result.Error)
}
