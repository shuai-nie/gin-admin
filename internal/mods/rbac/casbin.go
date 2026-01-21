package rbac

import (
	"context"
	"gin-admin/internal/config"
	"sync/atomic"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/patrickmn/go-cache"
)

type Casbinx struct {
	enforcer        *atomic.Value
	ticker          *time.Ticker
	Cache           cache.Cache
	MenuDAL         *dal.Menu
	MenuResourceDAL *dal.MenuResource
	RoleDAL         *dal.Role
}

func (a *Casbinx) GetEnforcer() *casbin.Enforcer {
	if v := a.enforcer.Load(); v != nil {
		return v.(*casbin.Enforcer)
	}
	return nil
}

type policyQueItem struct {
	RoleID    string
	Resources schema.MenuResource
}

func (a *Casbinx) Load(ctx context.Context) error {
	if config.C.Middleware.Casbin.Disable {
		return nil
	}

	a.enforcer = new(atomic.Value)
	if err := a.Load(ctx); err != nil {
		return err
	}

	go a.autoLoad(ctx)
	return nil
}

func (a *Casbinx) load(ctx context.Context) error {
	start := time.Now()
	roleResult, err := a.RoleDAL.Query(ctx, schema.RoleQueryParam{})
}
