package api

import (
	"gin-admin/internal/mods/rbac/biz"
	"gin-admin/internal/mods/rbac/schema"
	"gin-admin/pkg/util"

	"github.com/gin-gonic/gin"
)

type Menu struct {
	MenuBIZ *biz.Menu
}

func (a *Menu) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.MenuQueryParam
	if err := util.ParseQuery(c, &params); err != nil {
		util.ResError(c, err)
		return
	}

	result, err := a.MenuBIZ.Query(ctx, params)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResPage(c, result.Data, result.PageResult)
}

func (a *Menu) Get(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.MenuBIZ.Get(ctx, c.Param("id"))
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, item)
}

func (a *Menu) Create(c *gin.Context) {
	ctx := c.Request.Context()
	item := new(schema.MenuForm)
	if err := util.ParseJSON(c, item); err != nil {
		util.ResError(c, err)
		return
	}

	result, err := a.MenuBIZ.Create(ctx, item)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, result)
}

func (a *Menu) Update(c *gin.Context) {
	ctx := c.Request.Context()
	item := new(schema.MenuForm)
	if err := util.ParseJSON(c, item); err != nil {
		util.ResError(c, err)
		return
	} else if err := item.Validate(); err != nil {
		util.ResError(c, err)
		return
	}

	err := a.MenuBIZ.Update(ctx, c.Param("id"), item)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResOk(c)
}

func (a *Menu) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.MenuBIZ.Delete(ctx, c.Param("id"))
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResOk(c)
}
