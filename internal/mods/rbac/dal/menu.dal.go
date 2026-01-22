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
