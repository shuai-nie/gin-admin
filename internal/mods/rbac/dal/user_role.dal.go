package dal

import (
	"context"
	"fmt"
	"gin-admin/internal/mods/rbac/schema"
	"gin-admin/pkg/errors"
	"gin-admin/pkg/util"

	"gorm.io/gorm"
)

func GetUserRoleDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return util.GetDB(ctx, defDB).Model(new(schema.UserRole))
}

type UserRole struct {
	DB *gorm.DB
}

func (a *UserRole) Query(ctx context.Context, params schema.UserRoleQueryParam, opts ...schema.UserRoleQueryOptions) (*schema.UserRoleQueryResult, error) {
	var opt schema.UserRoleQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	db := a.DB.Table(fmt.Sprintf("%s as a", new(schema.UserRole).TableName()))
	if opt.JoinRole {
		db = db.Joins(fmt.Sprintf("left join %s b on a.role_id=b.id", new(schema.Role).TableName()))
		db = db.Select("a.*,b.name as role_name")
	}

	if v := params.InUserIDs; len(v) > 0 {
		db = db.Where("a.user_id in (?)", v)
	}

	if v := params.UserID; len(v) > 0 {
		db = db.Where("a.user_id=?", v)
	}
	if v := params.RoleID; len(v) > 0 {
		db = db.Where("a.role_id=?", v)
	}

	var list schema.UserRoles
	pageResult, err := util.WrapPageQuery(ctx, db, params.PaginationParam, opt.QueryOptions, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	queryResult := &schema.UserRoleQueryResult{
		Data:       list,
		PageResult: pageResult,
	}
	return queryResult, nil
}

func (a *UserRole) Get(ctx context.Context, id string, opts ...schema.UserRoleQueryOptions) (*schema.UserRole, error) {
	var opt schema.UserRoleQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	item := new(schema.UserRole)
	ok, err := util.FindOne(ctx, GetUserRoleDB(ctx, a.DB).Where("id=?", id), opt.QueryOptions, item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}
	return item, nil
}

func (a *UserRole) Exists(ctx context.Context, id string) (bool, error) {
	ok, err := util.Exists(ctx, GetUserRoleDB(ctx, a.DB).Where("id=?", id))
	return ok, errors.WithStack(err)
}

func (a *UserRole) Create(ctx context.Context, item schema.UserRole) error {
	result := GetUserRoleDB(ctx, a.DB).Create(item)
	return errors.WithStack(result.Error)
}

func (a *UserRole) Update(ctx context.Context, id string, item schema.UserRole) error {
	result := GetUserRoleDB(ctx, a.DB).Where("id=?", id).Updates(item)
	return errors.WithStack(result.Error)
}

func (a *UserRole) Delete(ctx context.Context, id string) error {
	result := GetUserRoleDB(ctx, a.DB).Where("id=?", id).Delete(new(schema.UserRole))
	return errors.WithStack(result.Error)
}

func (a *UserRole) DeleteByUserID(ctx context.Context, userID string) error {
	result := GetUserRoleDB(ctx, a.DB).Where("user_id=?", userID).Delete(new(schema.UserRole))
	return errors.WithStack(result.Error)
}

func (a *UserRole) DeleteByRoleID(ctx context.Context, roleID string) error {
	result := GetUserRoleDB(ctx, a.DB).Where("role_id=?", roleID).Delete(new(schema.UserRole))
	return errors.WithStack(result.Error)
}
