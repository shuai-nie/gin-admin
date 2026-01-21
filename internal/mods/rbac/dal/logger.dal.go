package dal

import (
	"context"
	"gin-admin/internal/mods/rbac/schema"
	"gin-admin/pkg/util"

	"gorm.io/gorm"
)

func GetLoggerDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return util.GetDB(ctx, defDB).Model(new(schema.Logger))
}
