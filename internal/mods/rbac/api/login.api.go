package api

import "gin-admin/internal/mods/rbac/biz"

type Login struct {
	LoginBIZ *biz.Login
}
