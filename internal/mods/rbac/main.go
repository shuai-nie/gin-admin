package rbac

import "gorm.io/gorm"

type RBAC struct {
	DB      *gorm.DB
	MenuAPI *api.Menu
}
