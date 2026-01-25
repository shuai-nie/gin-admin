package api

import (
	"gin-admin/internal/mods/rbac/biz"
	"gin-admin/internal/mods/rbac/schema"
	"gin-admin/pkg/util"

	"github.com/gin-gonic/gin"
)

type Logger struct {
	LoggerBIZ *biz.Logger
}

func (a *Logger) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.LoggerQueryParam
	if err := util.ParseQuery(c, &params); err != nil {
		util.ResError(c, err)
		return
	}

	result, err := a.LoggerBIZ.Query(ctx, params)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResPage(c, result.Data, result.PageResult)
}
