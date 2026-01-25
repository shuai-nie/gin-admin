package api

import (
	"gin-admin/internal/mods/rbac/biz"
	"gin-admin/internal/mods/rbac/schema"
	"gin-admin/pkg/util"

	"github.com/gin-gonic/gin"
)

type Login struct {
	LoginBIZ *biz.Login
}

func (a *Login) GetCaptcha(c *gin.Context) {
	ctx := c.Request.Context()
	data, err := a.LoginBIZ.GetCaptcha(ctx)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, data)
}

func (a *Login) ResponseCaptcha(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.LoginBIZ.ResponseCaptcha(ctx, c.Writer, c.Query("id"), c.Query("reload") == "1")
	if err != nil {
		util.ResError(c, err)
	}
}

func (a *Login) Login(c *gin.Context) {
	ctx := c.Request.Context()
	item := new(schema.LoginForm)
	if err := util.ParseJSON(c, item); err != nil {
		util.ResError(c, err)
		return
	}
	data, err := a.LoginBIZ.Login(ctx, item.Trim())
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, data)
}

func (a *Login) Logout(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.LoginBIZ.Logout(ctx)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResOk(c)
}

func (a *Login) RefreshToken(c *gin.Context) {
	ctx := c.Request.Context()
	data, err := a.LoginBIZ.RefreshToken(ctx)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, data)
}

func (a *Login) GetUserInfo(c *gin.Context) {
	ctx := c.Request.Context()
	data, err := a.LoginBIZ.GetUserInfo(ctx)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, data)
}

func (a *Login) UpdatePassword(c *gin.Context) {
	ctx := c.Request.Context()
	item := new(schema.UpdateLoginPassword)
	if err := util.ParseJSON(c, item); err != nil {
		util.ResError(c, err)
		return
	}
	err := a.LoginBIZ.UpdatePassword(ctx, item)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResOk(c)
}

func (a *Login) QueryMenus(c *gin.Context) {
	ctx := c.Request.Context()
	data, err := a.LoginBIZ.QueryMenus(ctx)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, data)
}

func (a *Login) UpdateUser(c *gin.Context) {
	ctx := c.Request.Context()
	item := new(schema.UpdateCurrentUser)
	if err := util.ParseJSON(c, item); err != nil {
		util.ResError(c, err)
		return
	}
	err := a.LoginBIZ.UpdateUser(ctx, item)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResOk(c)
}
