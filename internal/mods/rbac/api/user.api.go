package api

import (
	"gin-admin/internal/mods/rbac/biz"
	"gin-admin/internal/mods/rbac/schema"
	"gin-admin/pkg/util"

	"github.com/gin-gonic/gin"
)

type User struct {
	UserBIZ *biz.User
}

func (a *User) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.UserQueryParam
	if err := util.ParseQuery(c, &params); err != nil {
		util.ResError(c, err)
		return
	}

	result, err := a.UserBIZ.Query(ctx, params)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResPage(c, result.Data, result.PageResult)
}

func (a *User) Get(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.UserBIZ.Get(ctx, c.Param("id"))
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, item)
}

func (a *User) Create(c *gin.Context) {
	ctx := c.Request.Context()
	item := new(schema.UserForm)
	if err := util.ParseJSON(c, item); err != nil {
		util.ResError(c, err)
		return
	} else if err := item.Validate(); err != nil {
		util.ResError(c, err)
		return
	}
	result, err := a.UserBIZ.Create(ctx, item)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, result)
}

func (a *User) Update(c *gin.Context) {
	ctx := c.Request.Context()
	item := new(schema.UserForm)
	if err := util.ParseJSON(c, item); err != nil {
		util.ResError(c, err)
		return
	} else if err := item.Validate(); err != nil {
		util.ResError(c, err)
		return
	}

	err := a.UserBIZ.Update(ctx, c.Param("id"), item)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResOk(c)
}

func (a *User) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.UserBIZ.Delete(ctx, c.Param("id"))
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResOk(c)
}

func (a *User) ResetPassword(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.UserBIZ.ResetPassword(ctx, c.Param("id"))
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResOk(c)
}
