package dal

import (
	"context"
	"gin-admin/internal/mods/rbac/schema"
	"gin-admin/pkg/errors"
	"gin-admin/pkg/util"

	"gorm.io/gorm"
)

func GetUserDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return util.GetDB(ctx, defDB).Model(new(schema.User))
}

type User struct {
	DB *gorm.DB
}

func (a *User) Query(ctx context.Context, params schema.UserQueryParam, opts ...schema.UserQueryOptions) (*schema.UserQueryResult, error) {
	var opt schema.UserQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	db := GetUserDB(ctx, a.DB)

	if v := params.LikeUsername; len(v) > 0 {
		db = db.Where("username like ?", "%"+v+"%")
	}
	if v := params.LikeName; len(v) > 0 {
		db = db.Where("name LIKE ?", "%"+v+"%")
	}

	if v := params.Status; len(v) > 0 {
		db = db.Where("status = ?", v)
	}

	var list schema.Users
	pageResult, err := util.WrapPageQuery(ctx, db, params.PaginationParam, opt.QueryOptions, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	queryResult := &schema.UserQueryResult{
		PageResult: pageResult,
		Data:       list,
	}
	return queryResult, nil
}

func (a *User) Get(ctx context.Context, id string, opts ...schema.UserQueryOptions) (*schema.User, error) {
	var opt schema.UserQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	item := new(schema.User)
	ok, err := util.FindOne(ctx, GetUserDB(ctx, a.DB), opt.QueryOptions, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}
	return item, nil
}

func (a *User) GetByUsername(ctx context.Context, username string, opts ...schema.UserQueryOptions) (*schema.User, error) {
	var opt schema.UserQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	item := new(schema.User)
	ok, err := util.FindOne(ctx, GetUserDB(ctx, a.DB).Where("username = ?", username), opt.QueryOptions, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}
	return item, nil
}

func (a *User) Exists(ctx context.Context, id string) (bool, error) {
	ok, err := util.Exists(ctx, GetUserDB(ctx, a.DB).Where("id=?", id))
	return ok, errors.WithStack(err)
}

func (a *User) ExistsUsername(ctx context.Context, username string) (bool, error) {
	ok, err := util.Exists(ctx, GetUserDB(ctx, a.DB).Where("username=?", username))
	return ok, errors.WithStack(err)
}

func (a *User) Create(ctx context.Context, item *schema.User) error {
	result := GetUserDB(ctx, a.DB).Create(item)
	return errors.WithStack(result.Error)
}

func (a *User) Update(ctx context.Context, item *schema.User, selectFields ...string) error {
	db := GetUserDB(ctx, a.DB).Where("id=?", item.ID)
	if len(selectFields) > 0 {
		db = db.Select(selectFields)
	} else {
		db = db.Select("*").Omit("create_at")
	}
	result := db.Updates(item)
	return errors.WithStack(result.Error)
}

func (a *User) Delete(ctx context.Context, id string) error {
	result := GetUserDB(ctx, a.DB).Where("id=?", id).Delete(new(schema.User))
	return errors.WithStack(result.Error)
}

func (a *User) UpdatePasswordByID(ctx context.Context, id string, password string) error {
	result := GetUserDB(ctx, a.DB).Where("id=?", id).Select("password").Updates(schema.User{Password: password})
	return errors.WithStack(result.Error)
}
