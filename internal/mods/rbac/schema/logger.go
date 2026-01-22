package schema

import (
	"gin-admin/internal/config"
	"gin-admin/pkg/util"
	"time"
)

type Logger struct {
	ID        string
	Level     string
	TraceID   string
	UserID    string
	Tag       string
	Message   string
	Stack     string
	Data      string
	CreatedAt time.Time
	LoginName string
	UserName  string
}

func (a *Logger) TableName() string {
	return config.C.FormatTableName("logger")
}

type LoggerQueryParam struct {
	util.PaginationParam
	Level        string
	TraceID      string
	LikeUserName string
	Tag          string
	LikeMessage  string
	StartTime    string
	EndTime      string
}

type LoggerQueryOptions struct {
	util.QueryOptions
}

type LoggerQueryResult struct {
	Data       Loggers
	PageResult *util.PaginationResult
}

type Loggers []*Logger
