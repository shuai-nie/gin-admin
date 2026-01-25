package biz

import (
	"context"
	"gin-admin/internal/config"
	"gin-admin/internal/mods/rbac/dal"
	"gin-admin/internal/mods/rbac/schema"
	"gin-admin/pkg/cachex"
	"gin-admin/pkg/crypto/hash"
	"gin-admin/pkg/errors"
	"gin-admin/pkg/util"
	"time"
)

type User struct {
	Cache       cachex.Cacher
	Trans       *util.Trans
	UserDAL     *dal.User
	UserRoleDAL *dal.UserRole
}

func (a *User) Query(ctx context.Context, params schema.UserQueryParam) (*schema.UserQueryResult, error) {
	params.Pagination = true

	result, err := a.UserDAL.Query(ctx, params, schema.UserQueryOptions{
		QueryOptions: util.QueryOptions{
			OrderFields: []util.OrderByParam{
				{Field: "id", Direction: util.DESC},
			},
			OmitFields: []string{"password"},
		},
	})
	if err != nil {
		return nil, err
	}

	if userIDs := result.Data.ToIDs(); len(userIDs) > 0 {
		userRoleResult, err := a.UserRoleDAL.Query(ctx, schema.UserRoleQueryParam{
			InUserIDs: userIDs,
		}, schema.UserRoleQueryOptions{
			JoinRole: true,
		})
		if err != nil {
			return nil, err
		}
		userRolesMap := userRoleResult.Data.ToUserIDMap()
		for _, user := range result.Data {
			user.Roles = userRolesMap[user.ID]
		}
	}
	return result, nil
}

func (a *User) Get(ctx context.Context, id string) (*schema.User, error) {
	user, err := a.UserDAL.Get(ctx, id, schema.UserQueryOptions{
		QueryOptions: util.QueryOptions{
			OmitFields: []string{"password"},
		},
	})
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, errors.NotFound("", "Incorrect user")
	}

	userRoleResult, err := a.UserRoleDAL.Query(ctx, schema.UserRoleQueryParam{
		UserID: id,
	})
	if err != nil {
		return nil, err
	}

	user.Roles = userRoleResult.Data
	return user, nil
}

func (a *User) Create(ctx context.Context, formItem *schema.UserForm) (*schema.User, error) {
	existsUsername, err := a.UserDAL.ExistsUsername(ctx, formItem.Username)
	if err != nil {
		return nil, err
	} else if existsUsername {
		return nil, errors.BadRequest("", "Username already exists")
	}

	user := &schema.User{
		ID:        util.NewXID(),
		CreatedAt: time.Now(),
	}

	if formItem.Password == "" {
		formItem.Password = config.C.General.DefaultLoginPwd
	}

	if err := formItem.FillTo(user); err != nil {
		return nil, err
	}

	err = a.Trans.Exec(ctx, func(ctx context.Context) error {
		if err := a.UserDAL.Create(ctx, user); err != nil {
			return err
		}

		for _, userRole := range formItem.Roles {
			userRole.ID = util.NewXID()
			userRole.UserID = user.ID
			userRole.CreatedAt = time.Now()
			if err := a.UserRoleDAL.Create(ctx, userRole); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	user.Roles = formItem.Roles
	return user, nil
}

func (a *User) Update(ctx context.Context, id string, formItem *schema.UserForm) error {
	user, err := a.UserDAL.Get(ctx, id)
	if err != nil {
		return err
	} else if user == nil {
		return errors.NotFound("", "Incorrect user")
	} else if user.Username != formItem.Username {
		existsUsername, err := a.UserDAL.ExistsUsername(ctx, formItem.Username)
		if err != nil {
			return err
		} else if existsUsername {
			return errors.BadRequest("", "Username already exists")
		}
	}

	if err := formItem.FillTo(user); err != nil {
		return err
	}
	user.UpdatedAt = time.Now()
	return a.Trans.Exec(ctx, func(ctx context.Context) error {
		if err := a.UserDAL.Update(ctx, user); err != nil {
			return err
		}

		if err := a.UserRoleDAL.DeleteByUserID(ctx, id); err != nil {
			return err
		}

		for _, userRole := range formItem.Roles {
			if userRole.ID == "" {
				userRole.ID = util.NewXID()
			}
			userRole.UserID = user.ID
			if userRole.CreatedAt.IsZero() {
				userRole.CreatedAt = time.Now()
			}
			userRole.UpdatedAt = time.Now()
			if err := a.UserRoleDAL.Create(ctx, userRole); err != nil {
				return err
			}
		}
		return a.Cache.Delete(ctx, config.CacheNSForUser, id)
	})
}

func (a *User) Delete(ctx context.Context, id string) error {
	exists, err := a.UserDAL.Exists(ctx, id)
	if err != nil {
		return err
	} else if !exists {
		return errors.NotFound("", "User not found")
	}

	return a.Trans.Exec(ctx, func(ctx context.Context) error {
		if err := a.UserDAL.Delete(ctx, id); err != nil {
			return err
		}
		if err := a.UserRoleDAL.DeleteByUserID(ctx, id); err != nil {
			return err
		}
		return a.Cache.Delete(ctx, config.CacheNSForUser, id)
	})
}

func (a *User) ResetPassword(ctx context.Context, id string) error {
	exists, err := a.UserDAL.Exists(ctx, id)
	if err != nil {
		return err
	} else if !exists {
		return errors.NotFound("", "User not found")
	}

	hashPass, err := hash.GeneratePassword(config.C.General.DefaultLoginPwd)
	if err != nil {
		return errors.BadRequest("", "GeneratePassword %s error", err.Error())
	}
	return a.Trans.Exec(ctx, func(ctx context.Context) error {
		if err := a.UserDAL.UpdatePasswordByID(ctx, id, hashPass); err != nil {
			return err
		}
		return nil
	})
}

func (a *User) GetRoleIDs(ctx context.Context, id string) ([]string, error) {
	userRoleResult, err := a.UserRoleDAL.Query(ctx, schema.UserRoleQueryParam{
		UserID: id,
	}, schema.UserRoleQueryOptions{
		QueryOptions: util.QueryOptions{
			SelectFields: []string{"role_id"},
		},
	})
	if err != nil {
		return nil, err
	}
	return userRoleResult.Data.ToRoleIDs(), nil
}
