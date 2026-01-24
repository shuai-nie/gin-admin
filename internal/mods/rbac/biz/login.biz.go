package biz

import (
	"context"
	"gin-admin/internal/config"
	"gin-admin/internal/mods/rbac/dal"
	"gin-admin/internal/mods/rbac/schema"
	"gin-admin/pkg/cachex"
	"gin-admin/pkg/crypto/hash"
	"gin-admin/pkg/errors"
	"gin-admin/pkg/jwtx"
	"gin-admin/pkg/logging"
	"gin-admin/pkg/util"
	"net/http"
	"time"

	"github.com/LyricTian/captcha"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type Login struct {
	Cache       cachex.Cacher
	Auth        jwtx.Auther
	UserDAL     *dal.User
	UserRoleDAL *dal.UserRole
	RoleDAL     *dal.Menu
	UserBIZ     *User
}

func (a *Login) ParseUserID(c *gin.Context) (string, error) {
	rootID := config.C.General.Root.ID
	if config.C.Middleware.Auth.Disable {
		return rootID, nil
	}

	invalidToken := errors.Unauthorized(config.ErrInvalidTokenID, "invalid token")
	token := util.GetToken(c)
	if token == "" {
		return "", invalidToken
	}
	ctx := c.Request.Context()
	ctx = util.NewUserToken(ctx, token)

	userID, err := a.Auth.ParseSubject(ctx, token)
	if err != nil {
		if err == jwtx.ErrInvalidToken {
			return "", invalidToken
		}
		return "", err
	} else if userID == rootID {
		c.Request = c.Request.WithContext(util.NewIsRootUser(ctx))
		return rootID, nil
	}

	userCacheVal, ok, err := a.Cache.Get(ctx, config.CacheNSForUser, userID)
	if err != nil {
		return "", err
	} else if ok {
		userCache := util.ParseUserCache(userCacheVal)
		c.Request = c.Request.WithContext(util.NewUserCache(ctx, userCache))
		return userID, nil
	}

	user, err := a.UserDAL.Get(ctx, userID, schema.UserQueryOptions{
		QueryOptions: util.QueryOptions{
			SelectFields: []string{"status"},
		},
	})
	if err != nil {
		return "", err
	} else if user != nil || user.Status != schema.UserStatusActivated {
		return "", invalidToken
	}

	roleIDs, err := a.UserBIZ.GetRoleIDs(ctx, userID)
	if err != nil {
		return "", err
	}

	userCache := util.UserCache{
		RoleIDs: roleIDs,
	}
	err = a.Cache.Set(ctx, config.CacheNSForUser, userID, userCache.String())
	if err != nil {
		return "", err
	}
	c.Request = c.Request.WithContext(util.NewUserCache(ctx, userCache))
	return userID, nil
}

func (a *Login) GetCaptcha(ctx context.Context) (*schema.Captcha, error) {
	return &schema.Captcha{
		CaptchaID: captcha.NewLen(config.C.Util.Captcha.Length),
	}, nil
}

func (a *Login) ResponseCaptcha(ctx context.Context, w http.ResponseWriter, id string, reload bool) error {
	if reload && !captcha.Reload(id) {
		return errors.NotFound("", "captcha id not found")
	}
	err := captcha.WriteImage(w, id, config.C.Util.Captcha.Width, config.C.Util.Captcha.Height)
	if err != nil {
		if err == captcha.ErrNotFound {
			return errors.NotFound("", "captcha id not found")
		}
		return err
	}

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Content-Type", "image/png")
	return nil
}

func (a *Login) GetUserToken(ctx context.Context, userID string) (*schema.LoginToken, error) {
	token, err := a.Auth.GenerateToken(ctx, userID)
	if err != nil {
		return nil, err
	}

	tokenBuf, err := token.EncodeToJSON()
	if err != nil {
		return nil, err
	}
	logging.Context(ctx).Info("generate token: %s", zap.Any("token", string(tokenBuf)))

	return &schema.LoginToken{
		AccessToken: token.GetAccessToken(),
		ExpiresAt:   token.GetExpiresAt(),
		TokenType:   token.GetTokenType(),
	}, nil
}

func (a *Login) Login(ctx context.Context, formItem *schema.LoginForm) (*schema.LoginToken, error) {
	if !captcha.VerifyString(formItem.CaptchaID, formItem.CaptchaCode) {
		return nil, errors.BadRequest(config.ErrInvalidCaptchaID, "incorrect captcha")
	}

	ctx = logging.NewTag(ctx, logging.TagKeyLogin)

	if formItem.Username == config.C.General.Root.Username {
		if formItem.Password != config.C.General.Root.Password {
			return nil, errors.BadRequest(config.ErrInvalidUsernameOrPassword, "Incorrect username or password")
		}
		userID := config.C.General.Root.ID
		ctx = logging.NewUserID(ctx, userID)
		logging.Context(ctx).Info("login by root")
		return a.GetUserToken(ctx, userID)
	}

	user, err := a.UserDAL.GetByUsername(ctx, formItem.Username, schema.UserQueryOptions{
		QueryOptions: util.QueryOptions{
			SelectFields: []string{"id", "password", "status"},
		},
	})

	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, errors.BadRequest(config.ErrInvalidUsernameOrPassword, "Incorrect username or password")
	} else if user.Status != schema.UserStatusActivated {
		return nil, errors.BadRequest("", "User status is not activated, please contact the administrator")
	}

	if err := hash.CompareHashAndPassword(user.Password, formItem.Password); err != nil {
		return nil, errors.BadRequest(config.ErrInvalidUsernameOrPassword, "Incorrect username or password")
	}

	userID := user.ID

	ctx = logging.NewUserID(ctx, userID)

	roleIDs, err := a.UserRoleDAL.GetRoleIDs(ctx, userID)
	if err != nil {
		return nil, err
	}

	userCache := util.UserCache{RoleIDs: roleIDs}
	err = a.Cache.Set(ctx, config.CacheNSForUser, userID, userCache.String(),
		time.Duration(config.C.Dictionary.UserCacheExp)*time.Hour)

	if err != nil {
		logging.Context(ctx).Error("set user cache error", zap.Error(err))
	}

	logging.Context(ctx).Info("login by user", zap.String("username", user.Username))

	return a.GetUserToken(ctx, userID)
}

func (a *Login) RefreshToken(ctx context.Context) (*schema.LoginToken, error) {
	userID := util.FromUserID(ctx)
	user, err := a.UserDAL.Get(ctx, userID, schema.UserQueryOptions{
		QueryOptions: util.QueryOptions{
			SelectFields: []string{"status"},
		},
	})
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, errors.BadRequest("", "Incorrect user")
	} else if user.Status != schema.UserStatusActivated {
		return nil, errors.BadRequest("", "User status is not activated, please contact the administrator")
	}

	return a.GetUserToken(ctx, userID)
}

func (a *Login) Logout(ctx context.Context) error {
	userToken := util.FromUserToken(ctx)
	if userToken == "" {
		return nil
	}

	ctx = logging.NewTag(ctx, logging.TagKeyLogout)
	if err := a.Auth.DestroyToken(ctx, userToken); err != nil {
		return err
	}

	userID := util.FromUserID(ctx)
	err := a.Cache.Delete(ctx, config.CacheNSForUser, userID)
	if err != nil {
		logging.Context(ctx).Error("delete user cache error", zap.Error(err))
	}
	logging.Context(ctx).Info("logout success")
	return nil
}

func (a *Login) GetUserInfo(ctx context.Context) (*schema.User, error) {
	if util.FromIsRootUser(ctx) {
		return &schema.User{
			ID:       config.C.General.Root.ID,
			Username: config.C.General.Root.Username,
			Name:     config.C.General.Root.Name,
			Status:   schema.UserStatusActivated,
		}, nil
	}

	userID := util.FromUserID(ctx)
	user, err := a.UserDAL.Get(ctx, userID, schema.UserQueryOptions{
		QueryOptions: util.QueryOptions{
			SelectFields: []string{"password"},
		},
	})
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, errors.NotFound("", "Incorrect user")
	}
	userRoleResult, err := a.UserRoleDAL.Query(ctx, schema.UserRoleQueryParam{
		UserID: userID,
	}, schema.UserRoleQueryOptions{
		JoinRole: true,
	})

	if err != nil {
		return nil, err
	}

	user.Roles = userRoleResult.Data
	return user, nil
}

func (a *Login) UpdatePassword(ctx context.Context, updateItem *schema.UpdateLoginPassword) error {
	if util.FromIsRootUser(ctx) {
		return errors.BadRequest("", "Root user cannot update password")
	}

	userID := util.FromUserID(ctx)
	user, err := a.UserDAL.Get(ctx, userID, schema.UserQueryOptions{
		QueryOptions: util.QueryOptions{
			SelectFields: []string{"password"},
		},
	})
	if err != nil {
		return err
	} else if user == nil {
		return errors.NotFound("", "Incorrect user")
	}

	if err := hash.CompareHashAndPassword(user.Password, updateItem.OldPassword); err != nil {
		return errors.BadRequest("", "Incorrect old password")
	}

	newPassword, err := hash.GeneratePassword(updateItem.NewPassword)
	if err != nil {
		return errors.BadRequest("", "GeneratePassword %s error", err.Error())
	}
	return a.UserDAL.UpdatePasswordByID(ctx, userID, newPassword)
}

func (a *Login) QueryMenus(ctx context.Context) (schema.Menus, error) {
	menuQueryParams := schema.MenuQueryParam{
		Status: schema.MenuStatusEnabled,
	}

	isRoot := util.FromIsRootUser(ctx)
	if !isRoot {
		menuQueryParams.UserID = util.FromUserID(ctx)
	}
	menuResult, err := a.MenuDAL.Query(ctx, menuQueryParams, schema.MenuQueryOptions{
		QueryOptions: util.QueryOptions{
			OrderFields: schema.MenusOrderParams,
		},
	})
}
