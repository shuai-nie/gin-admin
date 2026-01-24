package biz

import (
	"context"
	"gin-admin/internal/mods/rbac/dal"
	"gin-admin/internal/mods/rbac/schema"
	"gin-admin/pkg/util"
)

type Logger struct {
	LoggerDAL *dal.Logger
}

func (a *Logger) Query(ctx context.Context, params schema.LoggerQueryParam) (*schema.LoggerQueryResult, error) {
	params.Pagination = true
	result, err := a.LoggerDAL.Query(ctx, params, schema.LoggerQueryOptions{
		QueryOptions: util.QueryOptions{
			OrderFields: []util.OrderByParam{
				{Field: "created_at", Direction: util.DESC},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
