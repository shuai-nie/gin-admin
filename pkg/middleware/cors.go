package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type CORSConfig struct {
	Enable                 bool
	AllowAllOrigin         bool
	AllowOrigins           []string
	AllowMethods           []string
	AllowHeaders           []string
	AllowCredentials       bool
	ExposeHeaders          []string
	MaxAge                 int
	AllowWildcard          bool
	AllowBrowserExtensions bool
	AllowWebSockets        bool
	AllowFiles             bool
}

var DefaultCORSConfig = CORSConfig{
	AllowOrigins: []string{"*"},
	AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"},
}

func CORSWithConfig(cfg CORSConfig) gin.HandlerFunc {
	if !cfg.Enable {
		return Empty()
	}

	return cors.New(cors.Config{
		AllowAllOrigins:        cfg.AllowAllOrigin,
		AllowOrigins:           cfg.AllowOrigins,
		AllowMethods:           cfg.AllowMethods,
		AllowHeaders:           cfg.AllowHeaders,
		AllowCredentials:       cfg.AllowCredentials,
		ExposeHeaders:          cfg.ExposeHeaders,
		MaxAge:                 time.Second * time.Duration(cfg.MaxAge),
		AllowWildcard:          cfg.AllowWildcard,
		AllowBrowserExtensions: cfg.AllowBrowserExtensions,
		AllowWebSockets:        cfg.AllowWebSockets,
		AllowFiles:             cfg.AllowFiles,
	})
}
