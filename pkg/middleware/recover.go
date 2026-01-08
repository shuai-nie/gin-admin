package middleware

import (
	"fmt"
	"gin-admin/pkg/errors"
	"gin-admin/pkg/logging"
	"gin-admin/pkg/util"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RecoverConfig struct {
	Skip int
}

var DefaultRecoverConfig = RecoverConfig{
	Skip: 3,
}

func Recover() gin.HandlerFunc {
	return RecoveryWithConfig(DefaultRecoverConfig)
}

func RecoveryWithConfig(config RecoverConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rv := recover(); rv != nil {
				ctx := c.Request.Context()
				ctx = logging.NewTag(ctx, logging.TagKeyRecovery)

				var fileds []zap.Field
				fileds = append(fileds, zap.Strings("error", []string{fmt.Sprintf("%v", rv)}))
				fileds = append(fileds, zap.Int("skip", config.Skip))

				if gin.IsDebugging() {
					httpRequest, _ := httputil.DumpRequest(c.Request, false)
					headers := strings.Split(string(httpRequest), "\r\n")
					for idx, header := range headers {
						current := strings.Split(header, ":")
						if current[0] == "Authorization" {
							headers[idx] = current[0] + ":"
						}
					}
					fileds = append(fileds, zap.Strings("request", headers))
				}

				logging.Context(ctx).Error(fmt.Sprintf("[Recovery] %s panic recovered", time.Now().Format("2006/01/02 - 15:04:05")), fileds...)
				util.ResError(c, errors.InternalServerError("", "Internal server error, please try again later"))

			}
		}()
		c.Next()
	}
}
