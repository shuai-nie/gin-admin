package api

import (
	"gin-admin/internal/mods/rbac/biz"
	"gin-admin/internal/mods/rbac/schema"
	"gin-admin/pkg/util"

	"github.com/gin-gonic/gin"
)

type Role struct {
	RoleBIZ *biz.Role
}

func (a *Role) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.RoleQueryParam
	if err := util.ParseQuery(c, &params); err != nil {
		util.ResError(c, err)
		return
	}

	result, err := a.RoleBIZ.Query(ctx, params)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResPage(c, result.Data, result.PageResult)
}

func (a *Role) Get(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.RoleBIZ.Get(ctx, c.Param("id"))
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, item)
}

func (a *Role) Create(c *gin.Context) {
	ctx := c.Request.Context()
	item := new(schema.RoleForm)
	if err := util.ParseJSON(c, item); err != nil {
		util.ResError(c, err)
		return
	} else if err := item.Validate(); err != nil {
		util.ResError(c, err)
		return
	}

	result, err := a.RoleBIZ.Create(ctx, item)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, result)
}

func (a *Role) Update(c *gin.Context) {
	ctx := c.Request.Context()
	item := new(schema.RoleForm)
	if err := util.ParseJSON(c, item); err != nil {
		util.ResError(c, err)
		return
	} else if err := item.Validate(); err != nil {
		util.ResError(c, err)
		return
	}

	err := a.RoleBIZ.Update(ctx, c.Param("id"), item)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResOk(c)
}

func (a *Role) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.RoleBIZ.Delete(ctx, c.Param("id"))
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResOk(c)
}
