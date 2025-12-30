package util

import "gin-admin/pkg/errors"

const (
	ReqVBodyKey       = "req-body"
	ResBodyKey        = "req-body"
	TreePathDelimiter = "."
)

type ResponseResult struct {
	Success bool
	Data    interface{}
	Total   int64
	Error   *errors.Error
}

type PaginationResult struct {
	Total    int64
	Current  int
	PageSize int
}

type PaginationParam struct {
	Pagination bool
	OnlyCount  bool
	Current    int
	PageSize   int
}

type QueryOptions struct {
	SelectFields []string
	OmitFields   []string
	OrderFields  OrderByParams
}

type Direction string

const (
	ASC  Direction = "ASC"
	DESC Direction = "DESC"
)

type OrderByParam struct {
	Field     string
	Direction Direction
}

type OrderByParams []OrderByParam

func (a OrderByParams) ToSQL() string {
	if len(a) == 0 {
		return ""
	}
	var sql string
	for _, v := range a {
		sql += v.Field + " " + string(v.Direction) + ","
	}
	return sql[:len(sql)-1]
}
