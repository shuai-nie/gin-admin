package prom

import (
	"gin-admin/pkg/promx"

	"github.com/gin-gonic/gin"
)

var (
	Ins *promx.PrometheusWrapper
	GinMiddleware gin.HandlerFunc
)

func Init() {
	logMethod := make(map[string]struct{})
	logAPI := make(map[string]struct{})
	for _, m := range
}